package main

import (
	"fmt"
	"os"

	"github.com/tiagollopes/okay/lexer"
	"github.com/tiagollopes/okay/parser"
	"github.com/tiagollopes/okay/eval" // Importe a nova pasta
)

func main() {
	fmt.Println("Okay language - v0.2")

	if len(os.Args) < 3 {
		fmt.Println("Uso: okay build <arquivo.ok>")
		return
	}

	filename := os.Args[2]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error reading file:", err)
		return
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	fmt.Println("=== EXECUTANDO PROGRAMA ===")
	env := eval.NewEnvironment()
	eval.Eval(program, env)
}
