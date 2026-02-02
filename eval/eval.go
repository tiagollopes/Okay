package eval

import (
	"fmt"
	"github.com/tiagollopes/okay/parser"
)

// Environment guarda a memória das variáveis
type Environment struct {
	vars map[string]string
}

func NewEnvironment() *Environment {
	return &Environment{vars: make(map[string]string)}
}

// Eval executa o programa
func Eval(node interface{}, env *Environment) {
	switch n := node.(type) {

	case *parser.Program:
		for _, stmt := range n.Statements {
			Eval(stmt, env)
		}

	case *parser.ServiceStatement:
		fmt.Printf("==> Iniciando Serviço: %s na porta %s\n", n.Name, n.Port)
		for _, stmt := range n.Body {
			Eval(stmt, env)
		}

	case *parser.VarDeclarationStatement:
		// Guarda o valor na memória (mapa)
		env.vars[n.Name] = n.Value

	case *parser.PrintStatement:
		for _, arg := range n.Args {
			if arg.Type == "IDENT" {
				// Se for variável, busca na memória
				val, ok := env.vars[arg.Value]
				if ok {
					fmt.Print(val, " ")
				} else {
					fmt.Printf("<erro: %s indefinida> ", arg.Value)
				}
			} else {
				// Se for texto puro ou número, imprime direto
				fmt.Print(arg.Value, " ")
			}
		}
		fmt.Println() // Quebra de linha no fim do print
	}
}
