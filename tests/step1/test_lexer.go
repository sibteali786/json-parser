package main

import (
	"fmt"
	"json-parser/pkg/lexer"
	"json-parser/pkg/token"
)

func main() {
	input := `{}`

	l := lexer.New(input)

	for {
		tok := l.NextToken()
		fmt.Printf("%s\n", tok)

		if tok.Type == token.EOF {
			break
		}
	}
}
