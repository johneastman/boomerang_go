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
	defer p.expectedToken(tokens.SEMICOLON)

	if p.current.Type == tokens.IDENTIFIER && p.peek.Type == tokens.ASSIGN {
		variableName := p.current
		p.advance()
		p.advance()
		variableExpression := p.parseExpression(LOWEST)
		return node.CreateAssignmentStatement(variableName, variableExpression)

	} else if p.current.Type == tokens.PRINT {
		p.advance()
		p.expectedToken(tokens.OPEN_PAREN)

		parameters := p.parseParameters()
		return node.CreatePrintStatement(parameters)

	}
	return p.parseExpression(LOWEST)
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
		value := p.current.Literal
		p.advance()
		return node.CreateNumber(value)

	} else if p.current.Type == tokens.MINUS {
		op := p.current
		p.advance()
		expression := p.parsePrefix()
		return node.CreateUnaryExpression(op, expression)

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
		return node.CreateIdentifier(identifier.Literal)
	}

	panic(fmt.Sprintf("Invalid prefix: %s", p.current.Type))
}

func (p *parser) parseInfix(left node.Node) node.Node {
	op := p.current
	p.advance()
	right := p.parseExpression(p.getPrecedenceLevel(op))
	return node.CreateBinaryExpression(left, op, right)
}

func (p *parser) parseParameters() []node.Node {
	params := []node.Node{}
	for p.current.Type != tokens.CLOSED_PAREN {
		expression := p.parseExpression(LOWEST)
		params = append(params, expression)

		if p.current.Type == tokens.COMMA {
			p.advance()
			continue
		}

		if p.current.Type == tokens.CLOSED_PAREN {
			p.advance()
			break
		}
	}
	return params
}

func (p *parser) getPrecedenceLevel(operator tokens.Token) int {
	level, ok := precedenceLevels[operator.Type]
	if !ok {
		return LOWEST
	}
	return level
}

func (p *parser) expectedToken(expectedType string) {
	// Check if the current token's type is the same as the expected token type. If not, throw an error; otherwise, advance to
	// the next token.
	if p.current.Type != expectedType {
		panic(fmt.Sprintf("Expected token type %s, got %s", expectedType, p.current.Type))
	}
	p.advance()
}
