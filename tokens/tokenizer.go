package tokens

import (
	"boomerang/utils"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Tokenizer struct {
	source            string
	currentPos        int
	currentLineNumber int
}

const EOF_CHAR = 0

func New(source string) Tokenizer {
	return Tokenizer{source: source, currentPos: 0, currentLineNumber: 1}
}

func (t *Tokenizer) current() byte {
	if t.currentPos < len(t.source) {
		return t.source[t.currentPos]
	}
	return 0
}

func (t *Tokenizer) peek() byte {
	nextCharIndex := t.currentPos + 1
	if nextCharIndex < len(t.source) {
		return t.source[nextCharIndex]
	}
	return 0
}

func (t *Tokenizer) advance() {
	t.currentPos += 1
}

func (t *Tokenizer) skipWhitespace() {
	for t.current() == ' ' || t.current() == '\t' || t.current() == '\n' || t.current() == '\r' {
		if t.current() == '\n' {
			t.currentLineNumber += 1
		}
		t.advance()
	}
}

func (t *Tokenizer) skipBlockComment() (*Token, error) {
	t.advance()
	t.advance()

	for {
		if t.current() == '#' && t.peek() == '#' {
			break
		}

		if t.peek() == 0 {
			return nil, utils.CreateError(t.currentLineNumber, "did not find ending ## while parsing block comment")
		}

		if t.current() == '\n' {
			t.currentLineNumber += 1
		}

		t.advance()
	}

	t.advance()
	t.advance()
	return t.Next()
}

func (t *Tokenizer) skipInlineComment() (*Token, error) {
	for t.current() != '\n' && t.current() != EOF_CHAR {
		t.advance()
	}
	return t.Next()
}

func (t *Tokenizer) isIdentifier(allowDigits bool) bool {
	/* Identifiers (e.g., variables) can include digits in the name but can't start with digits. When 'allowDigits' is false,
	 * only letters and underscores are allowed. When 'allowDigits' is true, digits are allowed.
	 */
	char := t.current()

	isIdentifierWithoutDigits := 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
	if allowDigits {
		return isIdentifierWithoutDigits || '0' <= char && char <= '9'
	}
	return isIdentifierWithoutDigits
}

func (t *Tokenizer) readIdentifier() string {
	startPos := t.currentPos
	endPos := startPos
	for t.isIdentifier(true) {
		endPos += 1
		t.advance()
	}
	return t.source[startPos:endPos]
}

func (t *Tokenizer) Next() (*Token, error) {
	t.skipWhitespace()

	if t.current() == EOF_CHAR {
		token := EOF_TOKEN
		token.LineNumber = t.currentLineNumber
		return &token, nil
	}

	if t.isIdentifier(false) {
		literal := t.readIdentifier()
		token := GetKeywordToken(literal)
		token.LineNumber = t.currentLineNumber
		return &token, nil
	}

	token, err := t.getMatchingTokens()
	if err != nil {
		return nil, err
	}

	token.LineNumber = t.currentLineNumber
	return token, nil
}

func (t *Tokenizer) getMatchingTokens() (*Token, error) {
	/*
		Get list of tokens where the token literal start with current character. The check for
		literal values having a length of greater-than 0 ensures tokens without an initial value
		are not matched (like NUMBER and IDENTIFIER).
	*/
	var matchingTokens []Token
	for _, tokenData := range tokenData {

		if tokenData.IsRegex {
			source := t.source[t.currentPos:]
			pattern := fmt.Sprintf("^%s", tokenData.Literal)

			r, err := regexp.Compile(pattern)
			if err != nil {
				panic(err.Error())
			}

			location := r.FindStringSubmatchIndex(source)
			/*
				location is nil if no match is found.

				If a match is found, location[0] and location[1] contain the start and end indices for the full regex match.
				location[3] and location[4] containg the start and end indices for sub-matches (e.g., a capturing group).
			*/
			if location != nil {
				matchStart := location[0]
				matchEnd := location[1]

				literalStart := location[2]
				literalEnd := location[3]

				literal := t.source[t.currentPos+literalStart : t.currentPos+literalEnd]
				token := Token{Type: tokenData.Type, Literal: literal}

				/*
					To advance past all the characters matching the regex, skip over the number of characters captured
					by the full regex match. For example, this ensures the double quotes for strings are skipped. However,
					when creating string tokens, we only want the text between the quotes, which is why the "literalStart"
					and "literalEnd" are used for token creation.
				*/
				for i := 0; i < (matchEnd - matchStart); i++ {
					t.advance()
				}

				return &token, nil
			}

			// For non-regex tokens, check if the current character matches the first character of the token literal
		} else if strings.HasPrefix(tokenData.Literal, string(t.current())) && len(tokenData.Literal) > 0 {
			token := Token{Type: tokenData.Type, Literal: tokenData.Literal}
			matchingTokens = append(matchingTokens, token)
		}
	}

	/*
		Sort the tokens by the length of the token literals in descending order. Sorting in descending ensures shorter
		tokens with similar characters to longer tokens are not mistakenly matched (for example, with '==', two '='
		tokens might be returned if the smaller tokens are ordered first).
	*/
	sort.SliceStable(matchingTokens, func(i, j int) bool {
		first := matchingTokens[i]
		second := matchingTokens[j]
		return len(first.Literal) > len(second.Literal)
	})

	/*
			For every matching token, check that the source code at the current position plus the length of the
			matching token literal are equal. For example:

				source = "1 == 1"
			              ^
		           pos: 2
			 len of '==': 2

			Search source from 2 to 4 (source[2 : 4]), but the last value is exclusive, so source[2 : 3] is returned,
			which is  "==".
	*/
	for _, matchingToken := range matchingTokens {
		source := t.source[t.currentPos : t.currentPos+len(matchingToken.Literal)]
		if source == matchingToken.Literal {

			if matchingToken.Type == INLINE_COMMENT_TOKEN.Type {
				return t.skipInlineComment()
			}

			if matchingToken.Type == BLOCK_COMMENT_TOKEN.Type {
				return t.skipBlockComment()
			}

			// Advance past n characters, where n is the length of the token literal

			for i := 0; i < len(source); i++ {
				t.advance()
			}

			return &matchingToken, nil
		}
	}
	return nil, utils.CreateError(t.currentLineNumber, "invalid character %c", t.current())
}
