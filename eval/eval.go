package eval

import (
	"fmt"
	"net/http"
	"strconv" // <--- converte strings em números
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
		// Verificamos o que tem dentro do valor
		switch val := n.Value.(type) {

		case string:
			// Se for texto puro, salva direto
			env.vars[n.Name] = val

		case *parser.Expression:
			// 1. Resolver o lado ESQUERDO
			leftVal := val.Left
			// Se for um nome de variável, buscamos o valor dela
			if v, ok := env.vars[val.Left]; ok {
				leftVal = v
			}

			// 2. Resolver o lado DIREITO
			rightVal := val.Right
			// Se for um nome de variável, buscamos o valor dela
			if v, ok := env.vars[val.Right]; ok {
				rightVal = v
			}

			// 3. Agora sim convertemos e somamos
			leftNum, _ := strconv.Atoi(leftVal)
			rightNum, _ := strconv.Atoi(rightVal)

			var resultado int
			if val.Operator == "+" {
				resultado = leftNum + rightNum
			}

			env.vars[n.Name] = strconv.Itoa(resultado)
		}
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
