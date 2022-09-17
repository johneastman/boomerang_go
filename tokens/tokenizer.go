package tokens

type Tokenizer struct {
	source     string
	currentPos int
}

type Token struct {
	Literal string
	Type    string
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

func (t *Tokenizer) isIdentifier() bool {
	char := t.current()
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (t *Tokenizer) readIdentifier() string {
	startPos := t.currentPos
	endPos := startPos
	for t.isIdentifier() {
		endPos += 1
		t.advance()
	}
	return t.source[startPos:endPos]
}

func (t *Tokenizer) isNumber() bool {
	char := t.current()
	return '0' <= char && char <= '9'
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
		return Token{Literal: "", Type: EOF}
	}

	if t.isIdentifier() {
		literal := t.readIdentifier()
		tokenType := getTokenType(literal)
		return Token{Literal: literal, Type: tokenType}
	} else if t.isNumber() {
		literal := t.readNumber()
		return Token{Literal: literal, Type: NUMBER}
	}

	literal := t.current()
	tokenType := getSymbolType(literal)
	t.advance()
	return Token{Literal: string(literal), Type: tokenType}
}
