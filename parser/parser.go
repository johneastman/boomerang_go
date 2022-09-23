package parser

import (
	"boomerang/node"
	"boomerang/tokens"
	"fmt"
	"regexp"
	"strings"
)

const (
	LOWEST int = iota
	FUNC_CALL
	SUM
	PRODUCT
)

var precedenceLevels = map[string]int{
	tokens.LEFT_PTR_TOKEN.Type:      FUNC_CALL,
	tokens.RIGHT_PTR_TOKEN.Type:     FUNC_CALL,
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
	for !tokens.TokenTypesEqual(p.current, tokens.EOF_TOKEN) {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *parser) parseStatement() node.Node {

	var statement node.Node

	if tokens.TokenTypesEqual(p.current, tokens.IDENTIFIER_TOKEN) && tokens.TokenTypesEqual(p.peek, tokens.ASSIGN_TOKEN) {
		statement = p.parseAssignmentStatement()

	} else if tokens.TokenTypesEqual(p.current, tokens.PRINT_TOKEN) {
		statement = p.parsePrintStatement()

	} else {
		statement = p.parseExpression(LOWEST)
	}

	p.expectToken(tokens.SEMICOLON_TOKEN)
	return statement
}

func (p *parser) parseAssignmentStatement() node.Node {
	variableName := p.current.Literal
	p.advance()
	p.advance()
	variableExpression := p.parseExpression(LOWEST)
	return node.CreateAssignmentStatement(variableName, variableExpression)
}

func (p *parser) parsePrintStatement() node.Node {
	p.advance()
	p.expectToken(tokens.OPEN_PAREN_TOKEN)

	parameters := p.parseParameters()
	return node.CreatePrintStatement(parameters.Params)
}

func (p *parser) parseExpression(precedenceLevel int) node.Node {
	left := p.parsePrefix()

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left = p.parseInfix(left)
	}

	return left
}

func (p *parser) parsePrefix() node.Node {
	switch p.current.Type {

	case tokens.NUMBER_TOKEN.Type:
		return p.parseNumber()

	case tokens.STRING_TOKEN.Type:
		return p.parseStrings()

	case tokens.MINUS_TOKEN.Type:
		return p.parseUnaryExpression()

	case tokens.OPEN_PAREN_TOKEN.Type:
		return p.parseGroupedExpression()

	case tokens.FUNCTION_TOKEN.Type:
		return p.parseFunction()

	case tokens.IDENTIFIER_TOKEN.Type:
		return p.parseIdentifier()
	}
	panic(fmt.Sprintf("Unexpected token: %s (%#v)", p.current.Type, p.current.Literal))
}

func (p *parser) parseInfix(left node.Node) node.Node {
	op := p.current
	p.advance()
	right := p.parseExpression(p.getPrecedenceLevel(op))
	return node.CreateBinaryExpression(left, op, right)
}

func (p *parser) getPrecedenceLevel(operator tokens.Token) int {
	level, ok := precedenceLevels[operator.Type]
	if !ok {
		return LOWEST
	}
	return level
}

func (p *parser) parseIdentifier() node.Node {
	identifier := p.current
	p.advance()
	return node.CreateIdentifier(identifier.Literal)
}

func (p *parser) parseNumber() node.Node {
	value := p.current.Literal
	p.advance()
	return node.CreateNumber(value)
}

func (p *parser) parseStrings() node.Node {
	stringLiteral := p.current.Literal
	params := []node.Node{}
	expressionIndex := 0

	r := regexp.MustCompile(`{[^{}]*}`)
	for {
		match := r.FindStringIndex(stringLiteral)
		if len(match) == 0 {
			break
		}

		startPos := match[0]
		endPos := match[1]

		expressionInString := stringLiteral[startPos+1 : endPos-1]

		tokenizer := tokens.New(expressionInString)
		parserObj := New(tokenizer)
		expression := parserObj.parseExpression(LOWEST)
		params = append(params, expression)

		stringLiteral = strings.Replace(stringLiteral, stringLiteral[startPos:endPos], fmt.Sprintf("<%d>", expressionIndex), 1)
		expressionIndex += 1
	}

	p.advance()
	return node.CreateString(stringLiteral, params)
}

func (p *parser) parseUnaryExpression() node.Node {
	op := p.current
	p.advance()
	expression := p.parsePrefix()
	return node.CreateUnaryExpression(op, expression)
}

func (p *parser) parseGroupedExpression() node.Node {

	p.advance()

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
		p.advance()
		return node.CreateParameters([]node.Node{})
	}

	expression := p.parseExpression(LOWEST)
	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
		p.advance()
		return expression

	} else if tokens.TokenTypesEqual(p.current, tokens.COMMA_TOKEN) {
		p.advance()
		stmts := []node.Node{expression}

		additionalParams := p.parseParameters()
		stmts = append(stmts, additionalParams.Params...)
		return node.CreateParameters(stmts)
	}

	panic(fmt.Sprintf(
		"Expected %s or %s, got %s",
		tokens.CLOSED_PAREN_TOKEN.Type,
		tokens.COMMA_TOKEN.Type,
		p.current.Type,
	))
}

func (p *parser) parseParameters() node.Node {
	params := []node.Node{}
	for {
		if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
			p.advance()
			break
		}

		expression := p.parseExpression(LOWEST)
		params = append(params, expression)

		if tokens.TokenTypesEqual(p.current, tokens.COMMA_TOKEN) {
			p.advance()
			continue
		}
	}
	return node.Node{Type: node.PARAMETER, Params: params}
}

func (p *parser) parseFunction() node.Node {
	p.advance()
	p.expectToken(tokens.OPEN_PAREN_TOKEN)

	parameters := p.parseParameters()
	p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN)

	statements := []node.Node{}
	for p.current != tokens.CLOSED_CURLY_BRACKET_TOKEN {
		statement := p.parseStatement()
		statements = append(statements, statement)
	}

	p.expectToken(tokens.CLOSED_CURLY_BRACKET_TOKEN)
	return node.CreateFunction(parameters.Params, statements)
}

func (p *parser) expectToken(token tokens.Token) {
	// Check if the current token's type is the same as the expected token type. If not, throw an error; otherwise, advance to
	// the next token.
	if !(tokens.TokenTypesEqual(p.current, token)) {
		panic(fmt.Sprintf("Expected token type %s (%#v), got %s (%#v)", token.Type, token.Literal, p.current.Type, p.current.Literal))
	}
	p.advance()
}
