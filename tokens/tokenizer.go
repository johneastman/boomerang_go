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
	if t.currentPos < len(t.source) {
		return t.source[t.currentPos]
	}
	return 0
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
	token := t.Next()
	for token.Type != EOF {
		tokens = append(tokens, token)
		token = t.Next()
	}

	tokens = append(tokens, token)
	return tokens

}

func (t *tokenizer) Next() Token {
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
