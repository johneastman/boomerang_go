package tokens

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Tokenizer struct {
	source     string
	currentPos int
}

func New(source string) Tokenizer {
	return Tokenizer{source: source, currentPos: 0}
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

func (t *Tokenizer) Next() Token {
	t.skipWhitespace()

	if t.current() == 0 {
		return EOF_TOKEN
	}

	if t.isIdentifier(false) {
		literal := t.readIdentifier()
		return GetKeywordToken(literal)

	} else if t.isNumber() {
		literal := t.readNumber()

		r, _ := regexp.Compile("^([0-9]*[.])?[0-9]+$")
		if !r.MatchString(literal) {
			panic(fmt.Sprintf("Invalid number literal: %s", literal))
		}
		return Token{Literal: literal, Type: NUMBER}
	}

	token, err := t.getMatchingTokens()
	if err != nil {
		panic(err.Error())
	}
	return *token
}

func (t *Tokenizer) getMatchingTokens() (*Token, error) {
	// Get list of tokens where the token literal start with current character
	var matchingTokens []Token
	for _, token := range tokenData {
		if strings.HasPrefix(token.Literal, string(t.current())) {
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
	return nil, fmt.Errorf("invalid symbol: %c", t.current())
}
