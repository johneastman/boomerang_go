package parser

import (
	"fmt"
	"log"
	"my_lang/tokens"
)

type Expression interface {
	expressionNode()
	String() string
}

type Statement struct {
	Expr Expression
}

func (s *Statement) expressionNode() {}

func (s Statement) String() string {
	return s.Expr.String()
}

type BinaryOperation struct {
	Left  Expression
	OP    tokens.Token
	Right Expression
}

func (bo *BinaryOperation) expressionNode() {}

func (bo BinaryOperation) String() string {
	return fmt.Sprintf("BinaryOperation(left=%s, op=%s, right=%s)", bo.Left.String(), bo.OP.Type, bo.Right.String())
}

type Number struct {
	Value string
}

func (n *Number) expressionNode() {}

func (n Number) String() string {
	return n.Value
}

type parser struct {
	tokenizer tokens.Tokenizer
	current   tokens.Token
}

func New(tokenizer tokens.Tokenizer) parser {
	currentToken := tokenizer.Next()
	return parser{tokenizer: tokenizer, current: currentToken}
}

func (p *parser) advance() {
	p.current = p.tokenizer.Next()
}

func (p parser) Parse() []Statement {
	statements := []Statement{}
	for p.current.Type != tokens.EOF {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *parser) parseStatement() Statement {
	expr := p.parseExpression()

	if p.current.Type != tokens.SEMICOLON {
		p.expectedToken(tokens.SEMICOLON)
	}
	p.advance()

	return Statement{Expr: expr}
}

func (p *parser) parseExpression() Expression {
	left := p.parseTerm()
	for {
		if p.current.Type == tokens.EOF {
			return left
		} else if p.current.Type == tokens.PLUS || p.current.Type == tokens.MINUS {
			op := p.current
			p.advance()
			right := p.parseTerm()
			left = &BinaryOperation{Left: left, OP: op, Right: right}
		} else {
			break
		}
	}
	return left
}

func (p *parser) parseTerm() Expression {
	left := p.parseFactor()
	for {
		if p.current.Type == tokens.EOF {
			return left
		} else if p.current.Type == tokens.ASTERISK || p.current.Type == tokens.FORWARD_SLASH {
			op := p.current
			p.advance()
			right := p.parseFactor()
			left = &BinaryOperation{Left: left, OP: op, Right: right}
		} else {
			break
		}
	}
	return left

}

func (p *parser) parseFactor() Expression {

	switch p.current.Type {
	case tokens.NUMBER:
		val := p.current.Literal
		p.advance()
		return &Number{Value: val}
	case tokens.OPEN_PAREN:
		p.advance()
		expr := p.parseExpression()
		if p.current.Type == tokens.CLOSED_PAREN {
			p.advance()
			return expr
		}
		p.expectedToken(tokens.CLOSED_PAREN)
	}
	log.Fatalf("Unexpected factor %s", p.current.Type)
	return nil
}

func (p *parser) expectedToken(expectedType string) {
	log.Fatalf("Expected token type %s, got %s", expectedType, p.current.Type)
}
