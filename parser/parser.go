package parser

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
)

const (
	LOWEST int = iota
	FUNC_CALL
	SUM
	PRODUCT
)

var precedenceLevels = map[string]int{
	tokens.OPEN_PAREN_TOKEN.Type:    FUNC_CALL,
	tokens.PLUS_TOKEN.Type:          SUM,
	tokens.MINUS_TOKEN.Type:         SUM,
	tokens.ASTERISK_TOKEN.Type:      PRODUCT,
	tokens.FORWARD_SLASH_TOKEN.Type: PRODUCT,
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
	for p.current != tokens.EOF_TOKEN {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *parser) parseStatement() node.Node {

	var statement node.Node

	if p.current.Type == tokens.IDENTIFIER_TOKEN.Type && p.peek.Type == tokens.ASSIGN_TOKEN.Type {
		variableName := p.current
		p.advance()
		p.advance()
		variableExpression := p.parseExpression(LOWEST)
		statement = node.CreateAssignmentStatement(variableName, variableExpression)

	} else if p.current.Type == tokens.PRINT_TOKEN.Type {
		p.advance()
		p.expectedToken(tokens.OPEN_PAREN_TOKEN)

		parameters := p.parseParameters()
		statement = node.CreatePrintStatement(parameters.Params)

	} else {
		statement = p.parseExpression(LOWEST)
	}

	p.expectedToken(tokens.SEMICOLON_TOKEN)
	return statement
}

func (p *parser) parseExpression(precedenceLevel int) node.Node {
	left := p.parsePrefix()

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left = p.parseInfix(left)
	}

	return left
}

func (p *parser) parsePrefix() node.Node {
	if p.current.Type == tokens.NUMBER_TOKEN.Type {
		value := p.current.Literal
		p.advance()
		return node.CreateNumber(value)

	} else if p.current.Type == tokens.MINUS_TOKEN.Type {
		op := p.current
		p.advance()
		expression := p.parsePrefix()
		return node.CreateUnaryExpression(op, expression)

	} else if p.current.Type == tokens.OPEN_PAREN_TOKEN.Type {
		p.advance()
		expression := p.parseExpression(LOWEST)
		if p.current.Type == tokens.CLOSED_PAREN_TOKEN.Type {
			p.advance()
			return expression
		}
		p.expectedToken(tokens.CLOSED_PAREN_TOKEN)

	} else if p.current.Type == tokens.FUNCTION_TOKEN.Type {
		return p.parseFunction()

	} else if p.current.Type == tokens.IDENTIFIER_TOKEN.Type {
		identifier := p.current
		p.advance()

		return node.CreateIdentifier(identifier.Literal)
	}

	panic(fmt.Sprintf("Invalid prefix: %s", p.current.Type))
}

func (p *parser) parseInfix(left node.Node) node.Node {
	var right node.Node

	op := p.current
	p.advance()
	if op.Type == tokens.OPEN_PAREN {
		// If the operator is an open parenthesis, then the operation is a function call
		right = p.parseParameters()
	} else {
		right = p.parseExpression(p.getPrecedenceLevel(op))
	}
	return node.CreateBinaryExpression(left, op, right)
}

func (p *parser) parseParameters() node.Node {
	params := []node.Node{}
	for {
		if p.current.Type == tokens.CLOSED_PAREN_TOKEN.Type {
			p.advance()
			break
		}

		expression := p.parseExpression(LOWEST)
		params = append(params, expression)

		if p.current.Type == tokens.COMMA_TOKEN.Type {
			p.advance()
			continue
		}
	}
	return node.Node{Type: node.PARAMETER, Params: params}
}

func (p *parser) parseFunction() node.Node {
	p.advance()
	p.expectedToken(tokens.OPEN_PAREN_TOKEN)

	parameters := p.parseParameters()
	p.expectedToken(tokens.OPEN_CURLY_BRACKET_TOKEN)

	statements := []node.Node{}
	for p.current != tokens.CLOSED_CURLY_BRACKET_TOKEN {
		statement := p.parseStatement()
		statements = append(statements, statement)
	}

	p.expectedToken(tokens.CLOSED_CURLY_BRACKET_TOKEN)
	return node.CreateFunction(parameters.Params, statements)
}

func (p *parser) getPrecedenceLevel(operator tokens.Token) int {
	level, ok := precedenceLevels[operator.Type]
	if !ok {
		return LOWEST
	}
	return level
}

func (p *parser) expectedToken(token tokens.Token) {
	// Check if the current token's type is the same as the expected token type. If not, throw an error; otherwise, advance to
	// the next token.
	if p.current.Type != token.Type {
		panic(fmt.Sprintf("Expected token type %s (%#v), got %s (%#v)", token.Type, token.Literal, p.current.Type, p.current.Literal))
	}
	p.advance()
}
