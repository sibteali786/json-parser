package parser

import (
	"fmt"
	"json-parser/pkg/lexer"
	"json-parser/pkg/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d, column %d", t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.errors = append(p.errors, msg)
}

// ParseJSON parses the JSON input and returns true if valid
func (p *Parser) ParseJSON() (JSONValue, bool) {
	value := p.parseValue()
	if len(p.errors) > 0 {
		return nil, false
	}

	// Should be at EOF
	if p.curToken.Type != token.EOF {
		p.errors = append(p.errors, fmt.Sprintf("expected EOF, got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil, false
	}
	return value, true
}

// parseValue parses any JSON value
func (p *Parser) parseValue() JSONValue {
	switch p.curToken.Type {
	case token.LEFT_BRACE:
		return p.parseObject()
	case token.STRING:
		return p.parseString()
	case token.NUMBER:
		return p.parseNumber()
	case token.TRUE, token.FALSE:
		return p.parseBoolean()
	case token.NULL:
		return p.parseNull()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}
}

// parseObject parses a JSON object
func (p *Parser) parseObject() *JSONObject {
	obj := &JSONObject{Pairs: make(map[string]JSONValue)}

	if !p.expectToken(token.LEFT_BRACE) {
		return nil
	}

	// handle empty object
	if p.curToken.Type == token.RIGHT_BRACE {
		p.nextToken() // consume right bracket
		return obj
	}

	// Parse first pair
	if !p.parseObjectPair(obj) {
		return nil
	}

	// Parse additional pairs
	for p.curToken.Type == token.COMMA {
		p.nextToken() // consume comma
		if !p.parseObjectPair(obj) {
			return nil
		}
	}
	if !p.expectToken(token.RIGHT_BRACE) {
		return nil
	}

	return obj
}

// expectToken checks if current token matches expected type and advances
func (p *Parser) expectToken(t token.TokenType) bool {
	if p.curToken.Type != t {
		p.peekError(t)
		return false
	}
	p.nextToken()
	return true
}

// Helper parsing methods - YOUR TASKS TO IMPLEMENT
func (p *Parser) parseObjectPair(obj *JSONObject) bool {
	// 1. Expect current token to be STRING
	if p.curToken.Type != token.STRING {
		p.errors = append(p.errors, fmt.Sprintf("expected string key, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return false
	}

	// 2. Save the key
	key := p.curToken.Literal

	// 3. Advance and expect COLON
	p.nextToken()
	if p.curToken.Type != token.COLON {
		p.errors = append(p.errors, fmt.Sprintf("expected colon ':', got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return false
	}

	// 4. Advance to value
	p.nextToken()

	// 5. Parse the value
	value := p.parseValue()
	if value == nil {
		// parseValue() failed, errors already added
		return false
	}

	// 6. Store the key-value pair
	obj.Pairs[key] = value

	return true
}
func (p *Parser) parseString() *JSONString {
	if p.curToken.Type != token.STRING {
		p.errors = append(p.errors, fmt.Sprintf("expected string, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	str := &JSONString{Value: p.curToken.Literal}
	p.nextToken()
	return str
}

func (p *Parser) parseNumber() *JSONNumber {
	if p.curToken.Type != token.NUMBER {
		p.errors = append(p.errors, fmt.Sprintf("expected number, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	num := &JSONNumber{Value: p.curToken.Literal}
	p.nextToken()
	return num
}

func (p *Parser) parseBoolean() *JSONBoolean {
	// TODO: Parse true/false tokens
	// Check if token is TRUE or FALSE and create appropriate JSONBoolean
	switch p.curToken.Type {
	case token.FALSE:
		p.nextToken()
		return &JSONBoolean{
			Value: false,
		}
	case token.TRUE:
		p.nextToken()
		return &JSONBoolean{
			Value: true,
		}
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected boolean token, got %s", p.curToken.Type))
		return nil
	}
}

func (p *Parser) parseNull() *JSONNull {
	if p.curToken.Type != token.NULL {
		p.errors = append(p.errors, fmt.Sprintf("expected null, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	p.nextToken()
	return &JSONNull{}
}
