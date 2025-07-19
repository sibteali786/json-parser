package token

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	// Special Tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// literals
	LEFT_BRACE    = "{"
	RIGHT_BRACE   = "}"
	LEFT_BRACKET  = "["
	RIGHT_BRACKET = "]"
	COLON         = ":"
	COMMA         = ","

	// values
	STRING = "STRING"
	NUMBER = "NUMBER"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NULL   = "NULL"

	// whitespace
	WHITESPACE = "WHITESPACE"
)

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %s, Line: %d, Column: %d}", t.Type, t.Literal, t.Line, t.Column)
}
