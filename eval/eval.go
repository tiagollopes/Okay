package eval

import (
	"fmt"
	"net/http"
	"strconv"
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
		for _, stmt := range n.Body {
			Eval(stmt, env)
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			msg, ok := env.vars["mensagem"]
			if !ok {
				msg = "Serviço Online"
			}
			fmt.Fprintf(w, "--- Okay Language Backend ---\n")
			fmt.Fprintf(w, "Servico: %s\n", n.Name)
			fmt.Fprintf(w, "Resposta: %s\n", msg)
		})
		err := http.ListenAndServe(":"+n.Port, nil)
		if err != nil {
			fmt.Printf("Erro crítico: %s\n", err)
		}

	// --- LÓGICA DO IF / ELSE COM SUPORTE A BOOLEANO ---
	case *parser.IfStatement:
		entrarNoIf := false

		// Caso 1: Condição simples, sem operador (ex: if (logado))
		if n.Condition.Operator == "" {
			val := n.Condition.Left
			if v, ok := env.vars[n.Condition.Left]; ok {
				val = v
			}
			// Se o valor na variável for a string "true", a condição é verdadeira
			entrarNoIf = (val == "true")
		} else {
			// Caso 2: Comparação matemática (ex: if (idade > 18))
			leftVal := n.Condition.Left
			if v, ok := env.vars[n.Condition.Left]; ok {
				leftVal = v
			}
			rightVal := n.Condition.Right
			if v, ok := env.vars[n.Condition.Right]; ok {
				rightVal = v
			}

			leftNum, _ := strconv.Atoi(leftVal)
			rightNum, _ := strconv.Atoi(rightVal)

			switch n.Condition.Operator {
			case ">":
				entrarNoIf = leftNum > rightNum
			case "<":
				entrarNoIf = leftNum < rightNum
			case "==":
				entrarNoIf = leftVal == rightVal // Comparação de strings/valores
			}
		}

		if entrarNoIf {
			for _, stmt := range n.Consequence {
				Eval(stmt, env)
			}
		} else if len(n.Alternative) > 0 {
			for _, stmt := range n.Alternative {
				Eval(stmt, env)
			}
		}

	case *parser.VarDeclarationStatement:
		switch val := n.Value.(type) {
		case string:
			// Salva o valor direto (seja número, texto ou "true"/"false")
			env.vars[n.Name] = val

		case *parser.Expression:
			leftVal := val.Left
			if v, ok := env.vars[val.Left]; ok {
				leftVal = v
			}
			rightVal := val.Right
			if v, ok := env.vars[val.Right]; ok {
				rightVal = v
			}

			leftNum, _ := strconv.Atoi(leftVal)
			rightNum, _ := strconv.Atoi(rightVal)

			var resultado int
			switch val.Operator {
			case "+": resultado = leftNum + rightNum
			case "-": resultado = leftNum - rightNum
			case "*": resultado = leftNum * rightNum
			case "/":
				if rightNum != 0 {
					resultado = leftNum / rightNum
				}
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
