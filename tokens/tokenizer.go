package tokens

import (
	"boomerang/utils"
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

func (t *Tokenizer) isNumber() bool {
	char := t.current()
	return '0' <= char && char <= '9' || char == '.'
}

func (t *Tokenizer) readNumber() string {
	startPos := t.currentPos
	endPos := startPos
	for t.isNumber() {
		endPos += 1
		t.advance()
	}
	return t.source[startPos:endPos]
}

func (t *Tokenizer) isString() bool {
	return t.current() == DOUBLE_QUOTE_TOKEN.Literal[0]
}

func (t *Tokenizer) skipInlineComment() {
	if t.current() == '#' {
		for t.current() != '\n' && t.current() != EOF_CHAR {
			t.advance()
		}
	}
	// There might be whitespace after the comment, so that needs to be skipped as well
	t.skipWhitespace()
}

func (t *Tokenizer) readString() string {
	startPos := t.currentPos
	endPos := startPos
	for !t.isString() {
		endPos += 1
		t.advance()
	}
	return t.source[startPos:endPos]
}

func (t *Tokenizer) Next() (*Token, error) {
	t.skipWhitespace()
	t.skipInlineComment()

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

	} else if t.isNumber() {
		literal := t.readNumber()

		r, _ := regexp.Compile("^([0-9]*[.])?[0-9]+$")
		if !r.MatchString(literal) {
			return nil, utils.CreateError(
				t.currentLineNumber,
				"error at line %d: invalid number literal: %s",
				t.currentLineNumber,
				literal,
			)
		}

		token := t.createToken(NUMBER, literal)
		return &token, nil

	} else if t.isString() {
		t.advance()
		literal := t.readString()
		t.advance()

		token := t.createToken(STRING, literal)
		return &token, nil
	}

	token, err := t.getMatchingTokens()
	if err != nil {
		return nil, err
	}

	token.LineNumber = t.currentLineNumber
	return token, nil
}

func (t *Tokenizer) createToken(tokenType string, tokenLiteral string) Token {
	return Token{Type: tokenType, Literal: tokenLiteral, LineNumber: t.currentLineNumber}
}

func (t *Tokenizer) getMatchingTokens() (*Token, error) {
	/*
		Get list of tokens where the token literal start with current character. The check for
		literal values having a length of greater-than 0 ensures tokens without an initial value
		are not matched (like NUMBER and IDENTIFIER).
	*/
	var matchingTokens []Token
	for _, token := range tokenData {
		if strings.HasPrefix(token.Literal, string(t.current())) && len(token.Literal) > 0 {
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

			// Advance past n characters, where n is the length of the token literal
			for i := 0; i < len(source); i++ {
				t.advance()
			}

			return &matchingToken, nil
		}
	}
	return nil, utils.CreateError(t.currentLineNumber, "invalid symbol %c", t.current())
}
