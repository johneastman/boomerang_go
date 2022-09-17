package parser

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
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

type parser struct {
	tokenizer tokens.Tokenizer
	current   tokens.Token
	peek      tokens.Token
}

func New(tokenizer tokens.Tokenizer) parser {
	currentToken := tokenizer.Next()
	peekToken := tokenizer.Next()
	return parser{tokenizer: tokenizer, current: currentToken, peek: peekToken}
}

func (p *parser) advance() {
	p.current = p.peek
	p.peek = p.tokenizer.Next()
}

func (p parser) Parse() []node.Node {
	statements := []node.Node{}
	for p.current.Type != tokens.EOF {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *parser) parseStatement() node.Node {

	expression := node.Node{}

	if p.current.Type == tokens.IDENTIFIER && p.peek.Type == tokens.ASSIGN {
		variableName := p.current
		p.advance()
		p.advance()
		variableExpression := p.parseExpression(LOWEST)
		expression = node.Node{
			Type: node.ASSIGN_STMT,
			Params: map[string]node.Node{
				node.ASSIGN_STMT_IDENTIFIER: {Type: variableName.Type, Value: variableName.Literal},
				node.EXPR:                   variableExpression,
			},
		}
	} else {
		expression = p.parseExpression(LOWEST)
	}

	if p.current.Type != tokens.SEMICOLON {
		p.expectedToken(tokens.SEMICOLON)
	}
	p.advance()

	return expression
}

func (p *parser) parseExpression(precedenceLevel int) node.Node {
	left := p.parsePrefix()

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left = p.parseInfix(left)
	}

	return left
}

func (p *parser) parsePrefix() node.Node {
	if p.current.Type == tokens.NUMBER {
		val := p.current.Literal
		p.advance()
		return node.Node{Type: node.NUMBER, Value: val}

	} else if p.current.Type == tokens.MINUS {
		op := p.current
		p.advance()
		expression := p.parsePrefix()
		return node.Node{
			Type: node.UNARY_EXPR,
			Params: map[string]node.Node{
				node.EXPR:     expression,
				node.OPERATOR: {Type: op.Type, Value: op.Literal},
			},
		}

	} else if p.current.Type == tokens.OPEN_PAREN {
		p.advance()
		expression := p.parseExpression(LOWEST)
		if p.current.Type == tokens.CLOSED_PAREN {
			p.advance()
			return expression
		}
		p.expectedToken(tokens.CLOSED_PAREN)

	} else if p.current.Type == tokens.IDENTIFIER {
		identifier := p.current
		p.advance()
		return node.Node{Type: node.IDENTIFIER, Value: identifier.Literal}
	}

	panic(fmt.Sprintf("Invalid prefix: %s", p.current.Type))
}

func (p *parser) parseInfix(left node.Node) node.Node {
	op := p.current
	p.advance()
	right := p.parseExpression(p.getPrecedenceLevel(op))
	return node.Node{
		Type: node.BIN_EXPR,
		Params: map[string]node.Node{
			node.BIN_EXPR_LEFT:  left,
			node.OPERATOR:       {Type: op.Type, Value: op.Literal},
			node.BIN_EXPR_RIGHT: right,
		},
	}
}

func (p *parser) getPrecedenceLevel(operator tokens.Token) int {
	level, ok := precedenceLevels[operator.Type]
	if !ok {
		return LOWEST
	}
	return level
}

func (p *parser) expectedToken(expectedType string) {
	panic(fmt.Sprintf("Expected token type %s, got %s", expectedType, p.current.Type))
}
