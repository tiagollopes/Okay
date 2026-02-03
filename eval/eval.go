package eval

import (
	"fmt"
	"net/http"
	"github.com/tiagollopes/okay/parser"
)

type Environment struct {
	vars map[string]string
}

func NewEnvironment() *Environment {
	return &Environment{vars: make(map[string]string)}
}

func Eval(node interface{}, env *Environment) {
	switch n := node.(type) {

	case *parser.Program:
		for _, stmt := range n.Statements {
			Eval(stmt, env)
		}

	case *parser.ServiceStatement:
		fmt.Printf("==> [OKAY] Servidor '%s' ouvindo na porta %s...\n", n.Name, n.Port)

		// 1. Primeiro, executamos o que está dentro do serviço para carregar as variáveis
		for _, stmt := range n.Body {
			Eval(stmt, env)
		}

		// 2. Configuramos a resposta do servidor
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Vamos tentar buscar uma variável chamada 'mensagem' na memória da Okay
			msg, ok := env.vars["mensagem"]
			if !ok {
				msg = "Serviço Online (Sem variável 'mensagem' definida)"
			}

			fmt.Fprintf(w, "--- Okay Language Backend ---\n")
			fmt.Fprintf(w, "Servico: %s\n", n.Name)
			fmt.Fprintf(w, "Resposta: %s\n", msg)
		})

		// 3. Iniciamos o servidor
		err := http.ListenAndServe(":"+n.Port, nil)
		if err != nil {
			fmt.Printf("Erro crítico: %s\n", err)
		}

	case *parser.VarDeclarationStatement:
		env.vars[n.Name] = n.Value

	case *parser.PrintStatement:
		fmt.Print("[LOG]: ")
		for _, arg := range n.Args {
			if arg.Type == "IDENT" {
				val, ok := env.vars[arg.Value]
				if ok {
					fmt.Print(val, " ")
				} else {
					fmt.Printf("<%s?> ", arg.Value)
				}
			} else {
				fmt.Print(arg.Value, " ")
			}
		}
		fmt.Println()
	}
}
