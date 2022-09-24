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
		return utils.CreateError(err)
	}

	p.peek = *nextToken
	return nil
}

func (p Parser) Parse() (ast *[]node.Node, err error) {
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

func (p *Parser) parseStatement() (stmt *node.Node, stmtErr error) {
	defer func() {
		if err := p.expectToken(tokens.SEMICOLON_TOKEN); err != nil {
			stmtErr = err
			stmt = nil
		}
	}()

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
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, utils.CreateError(err)
	}
	returnNode := node.CreateReturnStatement(*expression)
	return &returnNode, nil
}

func (p *Parser) parseAssignmentStatement() (*node.Node, error) {
	variableName := p.current.Literal

	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	variableExpression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, utils.CreateError(err)
	}

	assignmentNode := node.CreateAssignmentStatement(variableName, *variableExpression)
	return &assignmentNode, nil
}

func (p *Parser) parsePrintStatement() (*node.Node, error) {
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	if err := p.expectToken(tokens.OPEN_PAREN_TOKEN); err != nil {
		return nil, err
	}

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
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}
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
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	identifierNode := node.CreateIdentifier(identifier.Literal)
	return &identifierNode, nil
}

func (p *Parser) parseNumber() (*node.Node, error) {
	value := p.current.Literal
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	numberNode := node.CreateNumber(value)
	return &numberNode, nil
}

func (p *Parser) parseBoolean() (*node.Node, error) {
	value := p.current.Literal
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

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

	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	stringNode := node.CreateString(stringLiteral, params)
	return &stringNode, nil
}

func (p *Parser) parseUnaryExpression() (*node.Node, error) {
	op := p.current
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}
	expression, err := p.parsePrefix()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	unaryNode := node.CreateUnaryExpression(op, *expression)
	return &unaryNode, nil
}

func (p *Parser) parseGroupedExpression() (*node.Node, error) {

	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
		if err := p.advance(); err != nil {
			return nil, utils.CreateError(err)
		}
		listNode := node.CreateList([]node.Node{})
		return &listNode, nil
	}

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, utils.CreateError(err)
	}

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN_TOKEN) {
		if err := p.advance(); err != nil {
			return nil, utils.CreateError(err)
		}
		return expression, nil

	} else if tokens.TokenTypesEqual(p.current, tokens.COMMA_TOKEN) {
		if err := p.advance(); err != nil {
			return nil, utils.CreateError(err)
		}
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
			if err := p.advance(); err != nil {
				return nil, utils.CreateError(err)
			}
			break
		}

		expression, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, utils.CreateError(err)
		}

		params = append(params, *expression)

		if tokens.TokenTypesEqual(p.current, tokens.COMMA_TOKEN) {
			if err := p.advance(); err != nil {
				return nil, utils.CreateError(err)
			}
			continue
		}
	}

	paramNode := node.CreateList(params)
	return &paramNode, nil
}

func (p *Parser) parseFunction() (*node.Node, error) {
	if err := p.advance(); err != nil {
		return nil, utils.CreateError(err)
	}

	if err := p.expectToken(tokens.OPEN_PAREN_TOKEN); err != nil {
		return nil, err
	}

	parameters, err := p.parseParameters()
	if err != nil {
		return nil, utils.CreateError(err)
	}

	if err := p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	statements := []node.Node{}
	for p.current.Type != tokens.CLOSED_CURLY_BRACKET_TOKEN.Type {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, utils.CreateError(err)
		}

		statements = append(statements, *statement)
	}

	if err := p.expectToken(tokens.CLOSED_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	functionNode := node.CreateFunction(parameters.Params, statements)
	return &functionNode, nil
}

func (p *Parser) expectToken(token tokens.Token) error {
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
		return err
	}

	return p.advance()
}
