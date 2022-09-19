package tokens

import (
	"fmt"
	"regexp"
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

	literal := t.current()
	t.advance()
	return getSymbolType(literal)
}
