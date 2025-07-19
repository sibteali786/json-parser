package main

import (
	"fmt"
	"io"
	"json-parser/pkg/lexer"
	"json-parser/pkg/parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <json-file>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read File
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n")
		os.Exit(1)
	}

	defer file.Close()

	input, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n")
		os.Exit(1)
	}

	l := lexer.New(string(input))
	p := parser.New(l)

	// Parse JSON
	if p.ParseJSON() {
		fmt.Println("Valid JSON")
		os.Exit(0)
	} else {
		fmt.Println("Invalid JSON")
		for _, err := range p.Errors() {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		os.Exit(1)
	}
}
