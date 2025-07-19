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
func (p *Parser) ParseJSON() bool {
	// For Step 1, we only handle empty objects: {}
	if p.curToken.Type != token.LEFT_BRACE {
		p.errors = append(p.errors, fmt.Sprintf("expected '{' at start, got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return false
	}

	p.nextToken() // consume '{'

	// Expect closing brace
	if p.curToken.Type != token.RIGHT_BRACE {
		p.errors = append(p.errors, fmt.Sprintf("expected '}', got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return false
	}

	p.nextToken() // consume '}'

	// Should be EOF
	if p.curToken.Type != token.EOF {
		p.errors = append(p.errors, fmt.Sprintf("expected EOF, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return false
	}

	return true
}
