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

type Node struct {
	Type   string
	Value  string
	Params map[string]Node
}

func (n *Node) GetParam(key string) Node {
	node, ok := n.Params[key]
	if !ok {
		panic(fmt.Sprintf("Key not in node params: %s", key))
	}
	return node
}

func (n *Node) String() string {
	return fmt.Sprintf("Node(Type: %s, Value: %s)", n.Type, n.Value)
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

func (p parser) Parse() []Node {
	statements := []Node{}
	for p.current.Type != tokens.EOF {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
	}
	return statements
}

func (p *parser) parseStatement() Node {
	expression := p.parseExpression(LOWEST)

	if p.current.Type != tokens.SEMICOLON {
		p.expectedToken(tokens.SEMICOLON)
	}
	p.advance()

	return Node{
		Type: "Statement",
		Params: map[string]Node{
			"Expression": expression,
		},
	}
}

func (p *parser) parseExpression(precedenceLevel int) Node {
	left := p.parsePrefix()

	for precedenceLevel < p.getPrecedenceLevel(p.current) {
		left = p.parseInfix(left)
	}

	return left
}

func (p *parser) parsePrefix() Node {
	if p.current.Type == tokens.NUMBER {
		val := p.current.Literal
		p.advance()
		return Node{Type: "Number", Value: val}
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

func (p *parser) parseInfix(left Node) Node {
	op := p.current
	p.advance()
	right := p.parseExpression(p.getPrecedenceLevel(op))
	return Node{
		Type: "BinaryExpression",
		Params: map[string]Node{
			"left":     left,
			"operator": {Type: op.Type, Value: op.Literal},
			"right":    right,
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
	log.Fatalf("Expected token type %s, got %s", expectedType, p.current.Type)
}
