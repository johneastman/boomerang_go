package parser

import (
	"boomerang/node"
	"boomerang/tokens"
	"boomerang/utils"
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
	tokens.PTR_TOKEN.Type:           FUNC_CALL,
	tokens.PLUS_TOKEN.Type:          SUM,
	tokens.MINUS_TOKEN.Type:         SUM,
	tokens.ASTERISK_TOKEN.Type:      PRODUCT,
	tokens.FORWARD_SLASH_TOKEN.Type: PRODUCT,
}

type Parser struct {
	tokenizer tokens.Tokenizer
	current   tokens.Token
	peek      tokens.Token
}

func New(tokenizer tokens.Tokenizer) (*Parser, error) {
	currentToken, err := tokenizer.Next()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	peekToken, err := tokenizer.Next()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	return &Parser{tokenizer: tokenizer, current: *currentToken, peek: *peekToken}, nil
}

func (p *Parser) advance() error {
	p.current = p.peek
	nextToken, err := p.tokenizer.Next()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	p.peek = *nextToken
	return nil
}

func (p Parser) Parse() (*[]node.Node, error) {
	statements := []node.Node{}
	for !tokens.TokenTypesEqual(p.current, tokens.EOF_TOKEN) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, utils.CreateError(err)
		}

		statements = append(statements, *stmt)
	}
	return &statements, nil
}

func (p *Parser) parseStatement() (*node.Node, error) {

	defer p.expectToken(tokens.SEMICOLON_TOKEN)

	if tokens.TokenTypesEqual(p.current, tokens.IDENTIFIER_TOKEN) && tokens.TokenTypesEqual(p.peek, tokens.ASSIGN_TOKEN) {
		return p.parseAssignmentStatement()

	} else if tokens.TokenTypesEqual(p.current, tokens.PRINT_TOKEN) {
		return p.parsePrintStatement()

	} else if tokens.TokenTypesEqual(p.current, tokens.RETURN_TOKEN) {
		return p.parseReturnStatement()

	}
	return p.parseExpression(LOWEST)
}

func (p *Parser) parseReturnStatement() (*node.Node, error) {
	p.advance()
	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, utils.CreateError(err)
	}
	returnNode := node.CreateReturnStatement(*expression)
	return &returnNode, nil
}

func (p *Parser) parseAssignmentStatement() (*node.Node, error) {
	variableName := p.current.Literal
	p.advance()
	p.advance()
	variableExpression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, utils.CreateError(err)
	}

	assignmentNode := node.CreateAssignmentStatement(variableName, *variableExpression)
	return &assignmentNode, nil
}

func (p *Parser) parsePrintStatement() (*node.Node, error) {
	p.advance()
	p.expectToken(tokens.OPEN_PAREN_TOKEN)

	parameters, err := p.parseParameters()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	printNode := node.CreatePrintStatement(parameters.Params)
	return &printNode, nil
}

func (p *Parser) parseExpression(precedenceLevel int) (*node.Node, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left, err = p.parseInfix(*left)
		if err != nil {
			return nil, utils.CreateError(err)
		}
	}

	return left, nil
}

func (p *Parser) parsePrefix() (*node.Node, error) {
	switch p.current.Type {

	case tokens.NUMBER_TOKEN.Type:
		return p.parseNumber()

	case tokens.BOOLEAN_TOKEN.Type:
		return p.parseBoolean()

	case tokens.STRING_TOKEN.Type:
		return p.parseString()

	case tokens.MINUS_TOKEN.Type:
		return p.parseUnaryExpression()

	case tokens.OPEN_PAREN_TOKEN.Type:
		return p.parseGroupedExpression()

	case tokens.FUNCTION_TOKEN.Type:
		return p.parseFunction()

	case tokens.IDENTIFIER_TOKEN.Type:
		return p.parseIdentifier()

	default:
		return nil, fmt.Errorf("error at line %d: unexpected token at line %d: %s (%#v)",
			p.current.LineNumber,
			p.current.LineNumber,
			p.current.Type,
			p.current.Literal,
		)
	}
}

func (p *Parser) parseInfix(left node.Node) (*node.Node, error) {
	op := p.current
	p.advance()
	right, err := p.parseExpression(p.getPrecedenceLevel(op))
	if err != nil {
		return nil, utils.CreateError(err)
	}

	binaryNode := node.CreateBinaryExpression(left, op, *right)
	return &binaryNode, nil
}

func (p *Parser) getPrecedenceLevel(operator tokens.Token) int {
	level, ok := precedenceLevels[operator.Type]
	if !ok {
		return LOWEST
	}
	return level
}

func (p *Parser) parseIdentifier() (*node.Node, error) {
	identifier := p.current
	p.advance()

	identifierNode := node.CreateIdentifier(identifier.Literal)
	return &identifierNode, nil
}

func (p *Parser) parseNumber() (*node.Node, error) {
	value := p.current.Literal
	p.advance()

	numberNode := node.CreateNumber(value)
	return &numberNode, nil
}

func (p *Parser) parseBoolean() (*node.Node, error) {
	value := p.current.Literal
	p.advance()

	booleanNode := node.CreateBoolean(value)
	return &booleanNode, nil
}

func (p *Parser) parseString() (*node.Node, error) {
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
		parserObj, err := New(tokenizer)
		if err != nil {
			return nil, utils.CreateError(err)
		}

		expression, err := parserObj.parseExpression(LOWEST)
		if err != nil {
			return nil, utils.CreateError(err)
		}
		params = append(params, *expression)

		stringLiteral = strings.Replace(stringLiteral, stringLiteral[startPos:endPos], fmt.Sprintf("<%d>", expressionIndex), 1)
		expressionIndex += 1
	}

	p.advance()

	stringNode := node.CreateString(stringLiteral, params)
	return &stringNode, nil
}

func (p *Parser) parseUnaryExpression() (*node.Node, error) {
	op := p.current
	p.advance()
	expression, err := p.parsePrefix()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	unaryNode := node.CreateUnaryExpression(op, *expression)
	return &unaryNode, nil
}

func (p *Parser) parseGroupedExpression() (*node.Node, error) {

	p.advance()

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
		p.advance()
		listNode := node.CreateList([]node.Node{})
		return &listNode, nil
	}

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, utils.CreateError(err)
	}

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
		p.advance()
		return expression, nil

	} else if tokens.TokenTypesEqual(p.current, tokens.COMMA_TOKEN) {
		p.advance()
		stmts := []node.Node{*expression}

		additionalParams, err := p.parseParameters()
		if err != nil {
			return nil, utils.CreateError(err)
		}

		stmts = append(stmts, additionalParams.Params...)
		listNode := node.CreateList(stmts)
		return &listNode, nil
	}

	return nil, fmt.Errorf("error at line %d: expected %s or %s, got %s",
		p.current.LineNumber,
		tokens.CLOSED_PAREN_TOKEN.Type,
		tokens.COMMA_TOKEN.Type,
		p.current.Type,
	)
}

func (p *Parser) parseParameters() (*node.Node, error) {
	params := []node.Node{}
	for {
		if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
			p.advance()
			break
		}

		expression, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, utils.CreateError(err)
		}

		params = append(params, *expression)

		if tokens.TokenTypesEqual(p.current, tokens.COMMA_TOKEN) {
			p.advance()
			continue
		}
	}

	paramNode := node.CreateList(params)
	return &paramNode, nil
}

func (p *Parser) parseFunction() (*node.Node, error) {
	p.advance()
	p.expectToken(tokens.OPEN_PAREN_TOKEN)

	parameters, err := p.parseParameters()
	if err != nil {
		return nil, utils.CreateError(err)
	}
	p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN)

	statements := []node.Node{}
	for p.current.Type != tokens.CLOSED_CURLY_BRACKET_TOKEN.Type {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, utils.CreateError(err)
		}

		statements = append(statements, *statement)
	}

	p.expectToken(tokens.CLOSED_CURLY_BRACKET_TOKEN)

	functionNode := node.CreateFunction(parameters.Params, statements)
	return &functionNode, nil
}

func (p *Parser) expectToken(token tokens.Token) *error {
	// Check if the current token's type is the same as the expected token type. If not, throw an error; otherwise, advance to
	// the next token.
	if !(tokens.TokenTypesEqual(p.current, token)) {
		err := fmt.Errorf("error at line %d: expected token type %s (%#v), got %s (%#v)",
			p.current.LineNumber,
			token.Type,
			token.Literal,
			p.current.Type,
			p.current.Literal,
		)
		return &err
	}
	p.advance()
	return nil
}
