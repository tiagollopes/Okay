package main

import (
	"fmt"
	"os"

	"github.com/tiagollopes/okay/lexer"
	"github.com/tiagollopes/okay/parser"
)

func main() {
	fmt.Println("Okay language - v0.2")

	filename := os.Args[2]

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error reading file:", err)
		return
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()

	fmt.Println("=== PARSED PROGRAM ===")
	for i, stmt := range program.Statements {
		fmt.Printf("%d: %#v\n", i, stmt)
	}
}
