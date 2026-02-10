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

		// Variável de controle para saber se já carregamos os valores padrão
		inicializado := false

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/favicon.ico" {
				return
			}

			// 1. CARGA INICIAL (Apenas no primeiro acesso e sem imprimir nada)
			if !inicializado {
				for _, stmt := range n.Body {
					// Durante a carga inicial, não executamos os Prints
					if _, isPrint := stmt.(*parser.PrintStatement); !isPrint {
						Eval(stmt, env)
					}
				}
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

			// 3. EXECUÇÃO REAL (Aqui sim os Logs aparecem no terminal)
			for _, stmt := range n.Body {
				Eval(stmt, env)
			}

			fmt.Fprintf(w, "--- Okay Language API ---\nServico: %s\nStatus: OK", n.Name)
		})

		// Removido qualquer Eval(n.Body) daqui de fora para garantir silêncio total no boot
		err := http.ListenAndServe(":"+n.Port, nil)
		if err != nil {
			fmt.Printf("Erro crítico: %s\n", err)
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
		// Se a variável já existe (veio da URL), não sobrescrevemos com o valor default
		if _, jaExiste := env.vars[n.Name]; jaExiste {
			// Exceto se for uma expressão matemática (precisamos recalcular com os novos dados)
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
