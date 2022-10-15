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
	COMPARE
	SUM
	PRODUCT
	FUNC_CALL
	INDEX
)

var precedenceLevels = map[string]int{
	tokens.PTR:           FUNC_CALL,
	tokens.AT:            INDEX,
	tokens.PLUS:          SUM,
	tokens.MINUS:         SUM,
	tokens.NOT:           SUM,
	tokens.OR:            SUM,
	tokens.AND:           SUM,
	tokens.ASTERISK:      PRODUCT,
	tokens.FORWARD_SLASH: PRODUCT,
	tokens.EQ:            COMPARE,
}

type Parser struct {
	tokenizer tokens.Tokenizer
	current   tokens.Token
	peek      tokens.Token
}

func NewParser(tokenizer tokens.Tokenizer) (*Parser, error) {
	currentToken, err := tokenizer.Next()
	if err != nil {
		return nil, err
	}

	peekToken, err := tokenizer.Next()
	if err != nil {
		return nil, err
	}

	return &Parser{tokenizer: tokenizer, current: *currentToken, peek: *peekToken}, nil
}

func (p *Parser) advance() error {
	p.current = p.peek
	nextToken, err := p.tokenizer.Next()
	if err != nil {
		return err
	}

	p.peek = *nextToken
	return nil
}

func (p Parser) Parse() (*[]node.Node, error) {
	statements, err := p.parseGlobalStatements()
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) parseStatements(terminatingToken tokens.Token) (*[]node.Node, error) {
	statements := []node.Node{}
	for p.current.Type != terminatingToken.Type {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, *statement)
	}

	if err := p.expectToken(terminatingToken); err != nil {
		return nil, err
	}
	return &statements, nil
}

func (p *Parser) parseGlobalStatements() (*[]node.Node, error) {
	return p.parseStatements(tokens.EOF_TOKEN)
}

func (p *Parser) parseBlockStatements() (*node.Node, error) {
	lineNum := p.current.LineNumber
	statements, err := p.parseStatements(tokens.CLOSED_CURLY_BRACKET_TOKEN)
	if err != nil {
		return nil, err
	}

	blockStatementsNode := node.CreateBlockStatements(lineNum, *statements)
	return &blockStatementsNode, nil
}

func (p *Parser) parseStatement() (*node.Node, error) {
	var returnNode *node.Node
	var err error

	if tokens.TokenTypesEqual(p.current, tokens.IDENTIFIER) && tokens.TokenTypesEqual(p.peek, tokens.ASSIGN) {
		returnNode, err = p.parseAssignmentStatement()

	} else if tokens.TokenTypesEqual(p.current, tokens.PRINT) {
		returnNode, err = p.parsePrintStatement()

	} else {
		returnNode, err = p.parseExpression(LOWEST)
	}

	if err != nil {
		// This error check needs to return so the below expected-token error does not overwrite this error
		return nil, err
	}

	// Check that token at end of statement is a semicolon
	if expectedTokenErr := p.expectToken(tokens.SEMICOLON_TOKEN); expectedTokenErr != nil {
		returnNode = nil
		err = expectedTokenErr
	}

	return returnNode, err
}

func (p *Parser) parseAssignmentStatement() (*node.Node, error) {
	identifierToken := p.current

	if err := p.advance(); err != nil {
		return nil, err
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	variableExpression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	assignmentNode := node.CreateAssignmentStatement(identifierToken.LineNumber, identifierToken.Literal, *variableExpression)
	return &assignmentNode, nil
}

func (p *Parser) parseForLoop() (*node.Node, error) {
	lineNumber := p.current.LineNumber

	if err := p.advance(); err != nil {
		return nil, err
	}

	// placeholder variable for each element in the list
	elementPlaceholder := p.current
	if elementPlaceholder.Type != tokens.IDENTIFIER {
		return nil, utils.CreateError(
			lineNumber,
			"invalid type for for-loop element placeholder: %s",
			elementPlaceholder.ErrorDisplay(),
		)
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	// "in" keyword
	if err := p.expectToken(tokens.IN_TOKEN); err != nil {
		return nil, err
	}

	list, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	blockStatements, err := p.parseBlockStatements()
	if err != nil {
		return nil, err
	}

	return node.CreateForLoop(
		lineNumber,
		node.CreateIdentifier(lineNumber, elementPlaceholder.Literal),
		*list,
		*blockStatements,
	).Ptr(), nil
}

func (p *Parser) parsePrintStatement() (*node.Node, error) {

	lineNumber := p.current.LineNumber

	if err := p.advance(); err != nil {
		return nil, err
	}

	if err := p.expectToken(tokens.OPEN_PAREN_TOKEN); err != nil {
		return nil, err
	}

	parameters, err := p.parseParameters()
	if err != nil {
		return nil, err
	}

	printNode := node.CreatePrintStatement(lineNumber, parameters.Params)
	return &printNode, nil
}

func (p *Parser) parseExpression(precedenceLevel int) (*node.Node, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left, err = p.parseInfix(*left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) parsePrefix() (*node.Node, error) {
	switch p.current.Type {

	case tokens.NUMBER:
		return p.parseNumber()

	case tokens.BOOLEAN:
		return p.parseBoolean()

	case tokens.STRING:
		return p.parseString()

	case tokens.MINUS, tokens.NOT:
		return p.parseUnaryExpression()

	case tokens.OPEN_PAREN:
		return p.parseGroupedExpression()

	case tokens.FUNCTION:
		return p.parseFunction()

	case tokens.IDENTIFIER:
		return p.parseIdentifier()

	case tokens.WHEN:
		return p.parseWhenExpression()

	case tokens.FOR:
		return p.parseForLoop()

	default:
		current := p.current

		/*
			Advance to the next token so the statement error verifying the last token is a
			semicolon does not overwrite this error.
		*/
		if err := p.advance(); err != nil {
			return nil, err
		}

		return nil, utils.CreateError(current.LineNumber, "invalid prefix: %s",
			current.ErrorDisplay(),
		)
	}
}

func (p *Parser) parseInfix(left node.Node) (*node.Node, error) {
	op := p.current
	if err := p.advance(); err != nil {
		return nil, err
	}
	right, err := p.parseExpression(p.getPrecedenceLevel(op))
	if err != nil {
		return nil, err
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
	identifierToken := p.current
	if err := p.advance(); err != nil {
		return nil, err
	}

	identifierNode := node.CreateIdentifier(identifierToken.LineNumber, identifierToken.Literal)
	return &identifierNode, nil
}

func (p *Parser) parseNumber() (*node.Node, error) {
	numberToken := p.current

	if err := p.advance(); err != nil {
		return nil, err
	}

	numberNode := node.CreateNumber(numberToken.LineNumber, numberToken.Literal)
	return &numberNode, nil
}

func (p *Parser) parseBoolean() (*node.Node, error) {
	booleanToken := p.current

	if err := p.advance(); err != nil {
		return nil, err
	}

	booleanNode := node.CreateBoolean(booleanToken.LineNumber, booleanToken.Literal)
	return &booleanNode, nil
}

func (p *Parser) parseString() (*node.Node, error) {
	stringLiteral := p.current.Literal
	lineNumber := p.current.LineNumber

	params := []node.Node{}
	expressionIndex := 0

	// Parse each expression block in the string interpolation and save the result to the string
	r := regexp.MustCompile(`{[^{}]*}`)
	for {
		match := r.FindStringIndex(stringLiteral)
		if len(match) == 0 {
			break
		}

		startPos := match[0]
		endPos := match[1]

		expressionInString := stringLiteral[startPos+1 : endPos-1]

		tokenizer := tokens.NewTokenizer(expressionInString)
		parserObj, err := NewParser(tokenizer)
		if err != nil {
			return nil, err
		}

		expression, err := parserObj.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		params = append(params, *expression)

		stringLiteral = strings.Replace(stringLiteral, stringLiteral[startPos:endPos], fmt.Sprintf("<%d>", expressionIndex), 1)
		expressionIndex += 1
	}

	if err := p.advance(); err != nil {
		return nil, err
	}

	stringNode := node.CreateString(lineNumber, stringLiteral, params)
	return &stringNode, nil
}

func (p *Parser) parseUnaryExpression() (*node.Node, error) {
	op := p.current
	if err := p.advance(); err != nil {
		return nil, err
	}
	expression, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	unaryNode := node.CreateUnaryExpression(op, *expression)
	return &unaryNode, nil
}

func (p *Parser) parseGroupedExpression() (*node.Node, error) {

	lineNumber := p.current.LineNumber

	if err := p.advance(); err != nil {
		return nil, err
	}

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN) {
		if err := p.advance(); err != nil {
			return nil, err
		}
		listNode := node.CreateList(lineNumber, []node.Node{})
		return &listNode, nil
	}

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN) {
		if err := p.advance(); err != nil {
			return nil, err
		}
		return expression, nil

	} else if tokens.TokenTypesEqual(p.current, tokens.COMMA) {
		if err := p.advance(); err != nil {
			return nil, err
		}
		stmts := []node.Node{*expression}

		additionalParams, err := p.parseParameters()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, additionalParams.Params...)
		listNode := node.CreateList(lineNumber, stmts)
		return &listNode, nil
	}

	return nil, expectedMultipleTokens(
		p.current.LineNumber,
		p.current,
		[]tokens.Token{
			tokens.CLOSED_PAREN_TOKEN,
			tokens.COMMA_TOKEN,
		},
	)
}

func (p *Parser) parseParameters() (*node.Node, error) {

	lineNumber := p.current.LineNumber

	params := []node.Node{}
	for {
		if tokens.TokenTypesEqual(p.current, tokens.CLOSED_PAREN) {
			if err := p.advance(); err != nil {
				return nil, err
			}
			break
		}

		expression, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		params = append(params, *expression)

		if tokens.TokenTypesEqual(p.current, tokens.COMMA) {
			if err := p.advance(); err != nil {
				return nil, err
			}
			continue
		}
	}

	paramNode := node.CreateList(lineNumber, params)
	return &paramNode, nil
}

func (p *Parser) parseFunction() (*node.Node, error) {

	lineNumber := p.current.LineNumber

	if err := p.advance(); err != nil {
		return nil, err
	}

	if err := p.expectToken(tokens.OPEN_PAREN_TOKEN); err != nil {
		return nil, err
	}

	parameters, err := p.parseParameters()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	statements, err := p.parseBlockStatements()
	if err != nil {
		return nil, err
	}

	functionNode := node.CreateFunction(lineNumber, parameters.Params, *statements)
	return &functionNode, nil
}

func (p *Parser) parseWhenExpression() (*node.Node, error) {
	lineNumber := p.current.LineNumber

	// Skip over "when"
	if err := p.advance(); err != nil {
		return nil, err
	}

	// parse expression after "when"
	whenExpression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	// expect open curly bracket
	if err := p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	// Parse is/case expressions
	caseNodes := []node.Node{}
	for p.current.Type != tokens.ELSE {

		if err := p.expectToken(tokens.IS_TOKEN); err != nil {
			return nil, err
		}

		caseExpression, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		if err := p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN); err != nil {
			return nil, err
		}

		caseStatements, err := p.parseBlockStatements()
		if err != nil {
			return nil, err
		}

		caseNode := node.CreateCaseNode(caseExpression.LineNum, *caseExpression, *caseStatements)
		caseNodes = append(caseNodes, caseNode)
	}

	// Parse else/default case
	if err := p.expectToken(tokens.ELSE_TOKEN); err != nil {
		return nil, err
	}

	if err := p.expectToken(tokens.OPEN_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	elseStatements, err := p.parseBlockStatements()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(tokens.CLOSED_CURLY_BRACKET_TOKEN); err != nil {
		return nil, err
	}

	return node.CreateWhenNode(lineNumber, *whenExpression, caseNodes, *elseStatements).Ptr(), nil
}

func (p *Parser) expectToken(token tokens.Token) error {
	// Check if the current token's type is the same as the expected token type. If not, throw an error; otherwise, advance to
	// the next token.
	if !(tokens.TokenTypesEqual(p.current, token.Type)) {
		err := utils.CreateError(p.current.LineNumber, "expected token type %s, got %s",
			token.ErrorDisplay(),
			p.current.ErrorDisplay(),
		)
		return err
	}

	return p.advance()
}

func expectedMultipleTokens(lineNum int, actualToken tokens.Token, expectedTokens []tokens.Token) error {
	errorMessage := "expected "

	expectedTokenStrings := []string{}
	for _, expectedToken := range expectedTokens {
		message := expectedToken.ErrorDisplay()
		expectedTokenStrings = append(expectedTokenStrings, message)
	}
	errorMessage += strings.Join(expectedTokenStrings, " or ")
	errorMessage += fmt.Sprintf(", got %s", actualToken.ErrorDisplay())

	return utils.CreateError(lineNum, errorMessage)
}
