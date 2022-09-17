package parser

import (
	"fmt"
	"log"
	"my_lang/tokens"
)

const (
	LOWEST  = 1
	SUM     = 2
	PRODUCT = 3
)

var precedenceLevels = map[string]int{
	tokens.PLUS:          SUM,
	tokens.MINUS:         SUM,
	tokens.ASTERISK:      PRODUCT,
	tokens.FORWARD_SLASH: PRODUCT,
}

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
	expr := p.parseExpression(LOWEST)

	if p.current.Type != tokens.SEMICOLON {
		p.expectedToken(tokens.SEMICOLON)
	}
	p.advance()

	return Statement{Expr: expr}
}

func (p *parser) parseExpression(precedenceLevel int) Expression {
	left := p.parsePrefix()

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left = p.parseInfix(left)
	}

	return left
}

func (p *parser) parsePrefix() Expression {
	if p.current.Type == tokens.NUMBER {
		val := p.current.Literal
		p.advance()
		return &Number{Value: val}
	} else if p.current.Type == tokens.OPEN_PAREN {
		p.advance()
		expression := p.parseExpression(LOWEST)
		if p.current.Type == tokens.CLOSED_PAREN {
			p.advance()
			return expression
		}
		p.expectedToken(tokens.CLOSED_PAREN)
	}
	panic(fmt.Sprintf("Invalid prefix: %s", p.current.Type))
}

func (p *parser) parseInfix(left Expression) Expression {
	op := p.current
	p.advance()
	right := p.parseExpression(p.getPrecedenceLevel(op))
	return &BinaryOperation{Left: left, OP: op, Right: right}
}

func (p *parser) getPrecedenceLevel(operator tokens.Token) int {
	level, ok := precedenceLevels[operator.Type]
	if !ok {
		return LOWEST
	}
	return level
}

func (p *parser) expectedToken(expectedType string) {
	log.Fatalf("Expected token type %s, got %s", expectedType, p.current.Type)
}
