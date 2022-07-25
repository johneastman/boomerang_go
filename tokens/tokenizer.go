package tokens

type tokenizer struct {
	source     string
	currentPos int
}

type Token struct {
	Literal string
	Type    string
}

func New(source string) tokenizer {
	return tokenizer{source: source, currentPos: 0}
}

func (t *tokenizer) current() byte {
	return t.source[t.currentPos]
}

func (t *tokenizer) advance() {
	t.currentPos += 1
}

func (t *tokenizer) skipWhitespace() {
	for t.current() == ' ' || t.current() == '\t' || t.current() == '\n' || t.current() == '\r' {
		t.advance()
	}
}

func (t *tokenizer) isIdentifier() bool {
	char := t.current()
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (t *tokenizer) readIdentifier() string {
	startPos := t.currentPos
	endPos := startPos
	for t.isIdentifier() {
		endPos += 1
		t.advance()
	}
	return t.source[startPos:endPos]
}

func (t *tokenizer) isNumber() bool {
	char := t.current()
	return '0' <= char && char <= '9'
}

func (t *tokenizer) readNumber() string {
	startPos := t.currentPos
	endPos := startPos
	for t.isNumber() {
		endPos += 1
		t.advance()
	}
	return t.source[startPos:endPos]
}

func (t *tokenizer) Tokenize() []Token {
	tokens := []Token{}

	for t.currentPos < len(t.source) {

		t.skipWhitespace()

		if t.isIdentifier() {
			literal := t.readIdentifier()
			tokenType := getTokenType(literal)
			tokens = append(tokens, Token{Literal: literal, Type: tokenType})
		} else if t.isNumber() {
			literal := t.readNumber()
			tokens = append(tokens, Token{Literal: literal, Type: NUMBER})
		} else {
			literal := t.current()
			tokenType := getSymbolType(literal)
			tokens = append(tokens, Token{Literal: string(literal), Type: tokenType})
			t.advance()
		}
	}

	tokens = append(tokens, Token{Literal: "", Type: EOF})
	return tokens
}
