package eval

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/tiagollopes/okay/parser"
)

var modoSilencioso = false

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

		inicializado := false

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/favicon.ico" {
				return
			}

			// 1. CARGA INICIAL TOTALMENTE SILENCIOSA
			if !inicializado {
				modoSilencioso = true // Ativa o silêncio total
				for _, stmt := range n.Body {
					Eval(stmt, env)
				}
				modoSilencioso = false // Desativa o silêncio para a execução real
				inicializado = true
			}

			// 2. ATUALIZAÇÃO VIA URL
			query := r.URL.Query()
			if len(query) > 0 {
				fmt.Println("\n[REQUISIÇÃO]: Processando parâmetros...")
				for key, values := range query {
					if len(values) > 0 {
						env.vars[key] = values[0]
					}
				}
			}

			// 3. EXECUÇÃO REAL (Com logs liberados)
			for _, stmt := range n.Body {
				Eval(stmt, env)
			}

			fmt.Fprintf(w, "--- Okay Language API ---\nServico: %s\nStatus: OK", n.Name)
		})

		err := http.ListenAndServe(":"+n.Port, nil)
		if err != nil {
			fmt.Printf("Erro crítico: %s\n", err)
		}

	case *parser.RepeatStatement:
		countVal := n.Count
		if v, ok := env.vars[n.Count]; ok {
			countVal = v
		}

		times, err := strconv.Atoi(countVal)
		if err != nil {
			return
		}

		for i := 0; i < times; i++ {
			for _, stmt := range n.Body {
				Eval(stmt, env)
			}
		}

	case *parser.IfStatement:
		entrarNoIf := false
		leftVal := n.Condition.Left
		if v, ok := env.vars[n.Condition.Left]; ok {
			leftVal = v
		}

		if n.Condition.Operator == "" {
			entrarNoIf = (leftVal == "true")
		} else {
			rightVal := n.Condition.Right
			if v, ok := env.vars[n.Condition.Right]; ok {
				rightVal = v
			}

			leftNum, errL := strconv.Atoi(leftVal)
			rightNum, errR := strconv.Atoi(rightVal)

			if errL == nil && errR == nil {
				switch n.Condition.Operator {
				case ">": entrarNoIf = leftNum > rightNum
				case "<": entrarNoIf = leftNum < rightNum
				case "==": entrarNoIf = leftNum == rightNum
				}
			} else {
				switch n.Condition.Operator {
				case "==": entrarNoIf = leftVal == rightVal
				}
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
		if _, jaExiste := env.vars[n.Name]; jaExiste {
			if _, isExpr := n.Value.(*parser.Expression); !isExpr {
				return
			}
		}

		switch val := n.Value.(type) {
		case string:
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
		// SÓ IMPRIME SE NÃO ESTIVER NO MODO SILENCIOSO
		if !modoSilencioso {
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
}
