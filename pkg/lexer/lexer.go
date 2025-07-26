package lexer

import (
	"json-parser/pkg/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
	hasError     bool
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	// Track line and column for error reporting
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken scans input and returns the next token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Skip whitespace
	l.skipWhitespace()
	// Create token based on current character
	switch l.ch {
	case '{':
		tok = l.newToken(token.LEFT_BRACE, l.ch)
	case '}':
		tok = l.newToken(token.RIGHT_BRACE, l.ch)
	case '[':
		tok = l.newToken(token.LEFT_BRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RIGHT_BRACKET, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '"':
		l.hasError = false
		str := l.readString()

		if l.hasError {
			tok.Type = token.ILLEGAL
			tok.Literal = str
		} else {
			tok.Type = token.STRING
			tok.Literal = str
		}
		tok.Line = l.line
		tok.Column = l.column
		return tok // Return early to avoid readChar() call
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		tok.Line = l.line
		tok.Column = l.column

	default:
		// Handle numbers, booleans, null, and illegal characters
		if l.isDigit(l.ch) || l.ch == '-' {
			numStr := l.readNumber()
			if l.hasLeadingZero(numStr) {
				tok.Type = token.ILLEGAL
				tok.Literal = numStr
			} else {
				tok.Type = token.NUMBER
				tok.Literal = numStr
			}

			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if l.isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // Advance to the next character
	return tok
}

// Helper methods (you'll need to implement these)
func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// TODO: Implement these helper methods
func (l *Lexer) readString() string {
	position := l.position + 1 // Start after opening quote

	for {
		l.readChar()
		if l.ch == 0 {
			break
		}

		if l.ch == '\\' {
			// TODO: Handle escape sequences
			// Your task: implement this logic
			// Hint: read the next character and handle \", \\, \n, etc.
			l.handleEscapeSequence()
			continue
		}
		if l.ch == '"' { // Found closing quote or EOF
			break
		}
	}

	// Extract string content (between quotes)
	str := l.input[position:l.position]

	// Advance past the closing quote
	if l.ch == '"' {
		l.readChar()
	}

	return str
}

func (l *Lexer) readNumber() string {
	position := l.position

	// Handle negative sign
	if l.ch == '-' {
		l.readChar()
		if !l.isDigit(l.ch) {
			return l.input[position:l.position]
		}
	}
	if l.ch == '0' {
		l.readChar()

		if l.isDigit(l.ch) {
			// This is invalid! Leading zeros not allowed (except for "0" itself)
			// Continue reading the invalid number for error reporting
			for l.isDigit(l.ch) {
				l.readChar()
			}
			return l.input[position:l.position] // Return the invalid number
		}
	} else if l.isDigit(l.ch) {
		// Normal number: read all digits
		for l.isDigit(l.ch) {
			l.readChar()
		}
	}

	// Handle the decimal part
	if l.ch == '.' {
		l.readChar()
		if !l.isDigit(l.ch) {
			return l.input[position:l.position]
		}

		for l.isDigit(l.ch) {
			l.readChar()
		}
	}

	// Handle exponent part (optional)
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		// Handle optional + or - in exponent
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		// Must have at least one digit in exponent
		if !l.isDigit(l.ch) {
			return l.input[position:l.position] // Invalid: "1e" without digits
		}
		for l.isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]

}
func (l *Lexer) readIdentifier() string {
	position := l.position
	for l.isLetter(l.ch) || l.isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) lookupIdent(ident string) token.TokenType {
	switch ident {
	case "true":
		return token.TRUE
	case "false":
		return token.FALSE
	case "null":
		return token.NULL
	default:
		return token.ILLEGAL
	}
}

func (l *Lexer) hasLeadingZero(numStr string) bool {
	// make sure its at least not single digit
	if len(numStr) < 2 {
		return false
	}

	start := 0
	if numStr[0] == '-' {
		start = 1
		if len(numStr) < 3 { // "-0" or "-1" etc. are valid
			return false
		}
	}

	return numStr[start] == '0' && numStr[start+1] >= '0' && numStr[start+1] <= '9'

}

func (l *Lexer) handleEscapeSequence() {
	// read the character after back slash
	l.readChar()

	switch l.ch {
	case '"':
		// Escaped quote: continue reading
	case '\\':
		// Escaped backslash: continue reading
	case '/':
		// Escaped slash: continue reading
	case 'b':
		// Backspace: continue reading
	case 'f':
		// Form feed: continue reading
	case 'n':
		// Newline: continue reading
	case 'r':
		// Carriage return: continue reading
	case 't':
		// Tab: continue reading
	case 'u':
		// TODO: Unicode escape sequence \uXXXX
		// Your task: read 4 hex digits
		l.handleUnicodeEscape()
	default:
		// Invalid escape sequence - for now, just continue
		// In a production parser, you'd want to report an error
		l.hasError = true
	}
}

func (l *Lexer) handleUnicodeEscape() {
	for i := 0; i < 4; i++ {
		l.readChar()
		if !l.isHexDigit(l.ch) {
			l.hasError = true
			return
		}
	}
}
func (l *Lexer) isHexDigit(ch byte) bool {
	// TODO: Check if ch is 0-9, a-f, or A-F
	if ch >= '0' && ch <= '9' {
		return true
	}
	if (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F') {
		return true
	}
	return false
}
