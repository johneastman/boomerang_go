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

func (t *Tokenizer) Next() (*Token, error) {
	t.skipWhitespace()

	if t.current() == EOF_CHAR {
		token := EOF_TOKEN
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

			matchStart := -1
			matchEnd := -1
			literalStart := -1
			literalEnd := -1

			matchLocations := r.FindStringSubmatchIndex(source)
			/*
				When there are no grouped expressions in the regex, "matchLocation" is 2-elements long: the start and end
				positions in the searched string. However, if grouped expressions are found in the searched string, they
				will be added to "matchLocations."

				For example, if the string is '\"hello, world\"' and the regex is '\"(.*)\"', "matchLocations" will be
				[0, 15, 1, 14].

				However, if the string is 'true', and the regex is 'true|false', "matchLocations" will be [0, 5].

				Note that end indices will be the index of the character after the last matched character. See documentation
				for more details: https://pkg.go.dev/regexp#Regexp.FindSubmatchIndex
			*/
			if len(matchLocations) == 2 {
				// No grouped expressions in regex
				matchStart = matchLocations[0]
				matchEnd = matchLocations[1]

				literalStart = matchLocations[0]
				literalEnd = matchLocations[1]
			}

			if len(matchLocations) == 4 {
				// Grouped expressions in regex
				matchStart = matchLocations[0]
				matchEnd = matchLocations[1]

				literalStart = matchLocations[2]
				literalEnd = matchLocations[3]
			}

			if matchStart != -1 && matchEnd != -1 && literalStart != -1 && literalEnd != -1 {
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
