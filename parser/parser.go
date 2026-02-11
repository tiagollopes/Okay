package parser

import (
	"fmt"

	"github.com/tiagollopes/okay/lexer"
)

// Adicione na definição de structs
type RepeatStatement struct {
	Count       string      // Pode ser um número ou uma variável
	Body        []Statement
}
func (rs *RepeatStatement) statementNode() {}

// Representa uma operação, ex: 10 + 5
type Expression struct {
	Left     string // Valor da esquerda
	Operator string // O símbolo (ex: "+")
	Right    string // Valor da direita
}

type Parser struct {
	l      *lexer.Lexer
	curTok lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curTok = p.l.NextToken()
}

//
// ===== AST (Árvore de Sintaxe Abstrata) =====
//

type Statement interface {
	statementNode()
}

type Program struct {
	Statements []Statement
}

// Novo: Diferencia se o argumento do print é Texto, Número ou Variável
type PrintArgument struct {
	Type  string // "STRING", "IDENT" ou "NUMBER"
	Value string
}

type PrintStatement struct {
	Args []PrintArgument
}

func (ps *PrintStatement) statementNode() {}

type ServiceStatement struct {
	Name string
	Port string
	Body []Statement
}

func (ss *ServiceStatement) statementNode() {}

type VarDeclarationStatement struct {
	Name  string
	Value interface{} // Mudamos de string para interface{}
}

func (vds *VarDeclarationStatement) statementNode() {}

type IfStatement struct {
	Condition Expression // ex: preco > 100
	Consequence []Statement // IF
	Alternative []Statement // ELSE
}

func (is *IfStatement) statementNode() {}

//
// ===== PARSER (Lógica de leitura) =====
//

func (p *Parser) ParseProgram() *Program {
	program := &Program{}

	for p.curTok.Type != lexer.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program
}

func (p *Parser) ParseStatement() Statement {
	switch p.curTok.Type {
	case lexer.SERVICE:
		return p.parseService()
	case lexer.REPEAT: // <-- Porta de entrada para o Token REPEAT
		return p.parseRepeat()
	case lexer.IDENT:
		if p.curTok.Literal == "print" {
			return p.parsePrint()
		}
		if p.curTok.Literal == "if" {
			return p.parseIf()
		}
	case lexer.LET:
		return p.parseVarDeclaration()
	}

	fmt.Println("unknown statement:", p.curTok.Literal)
	p.nextToken()
	return nil
}

func (p *Parser) parseService() Statement {
	stmt := &ServiceStatement{}
	p.nextToken()

	if p.curTok.Type == lexer.IDENT {
		stmt.Name = p.curTok.Literal
		p.nextToken()
	}

	if p.curTok.Type == lexer.PORT {
		p.nextToken()
	}

	if p.curTok.Type == lexer.NUMBER {
		stmt.Port = p.curTok.Literal
		p.nextToken()
	}

	if p.curTok.Type == lexer.LBRACE {
		p.nextToken()
	} else {
		fmt.Println("Erro: esperado '{'")
		return nil
	}

	for p.curTok.Type != lexer.RBRACE && p.curTok.Type != lexer.EOF {
		innerStmt := p.ParseStatement()
		if innerStmt != nil {
			stmt.Body = append(stmt.Body, innerStmt)
		}
	}

	if p.curTok.Type == lexer.RBRACE {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePrint() Statement {
	p.nextToken() // pula 'print'

	if p.curTok.Type != lexer.LPAREN {
		fmt.Println("expected '(' after print")
		return nil
	}
	p.nextToken()

	args := []PrintArgument{}

	for p.curTok.Type != lexer.RPAREN && p.curTok.Type != lexer.EOF {
		if p.curTok.Type == lexer.STRING || p.curTok.Type == lexer.IDENT || p.curTok.Type == lexer.NUMBER {
			// Aqui salvamos o Tipo e o Valor separadamente
			args = append(args, PrintArgument{
				Type:  string(p.curTok.Type),
				Value: p.curTok.Literal,
			})
			p.nextToken()
		}

		if p.curTok.Type == lexer.COMMA {
			p.nextToken()
		}
	}

	p.nextToken() // pula ')'

	if p.curTok.Type != lexer.SEMICOLON {
		fmt.Println("expected ';' after print")
		return nil
	}
	p.nextToken()

	return &PrintStatement{Args: args}
}

func (p *Parser) parseVarDeclaration() Statement {
	stmt := &VarDeclarationStatement{}
	p.nextToken() // pula 'let'

	if p.curTok.Type != lexer.IDENT {
		fmt.Println("Erro: esperado nome da variável")
		return nil
	}
	stmt.Name = p.curTok.Literal
	p.nextToken()

	if p.curTok.Type != lexer.ASSIGN {
		fmt.Println("Erro: esperado '='")
		return nil
	}
	p.nextToken()

	// Guardamos o primeiro valor (pode ser o único ou o começo de uma conta)
	leftVal := p.curTok.Literal
	p.nextToken()

	// Aceitamos qualquer um dos 4 operadores
	if p.curTok.Type == lexer.PLUS || p.curTok.Type == lexer.MINUS ||
	   p.curTok.Type == lexer.ASTERISK || p.curTok.Type == lexer.SLASH {

		operator := p.curTok.Literal
		p.nextToken() // pula o operador (+, -, *, /)

		rightVal := p.curTok.Literal
		p.nextToken() // pula o valor da direita

		stmt.Value = &Expression{
			Left:     leftVal,
			Operator: operator,
			Right:    rightVal,
		}
	} else {
		stmt.Value = leftVal
	}

	if p.curTok.Type != lexer.SEMICOLON {
		fmt.Println("Erro: esperado ';'")
		return nil
	}
	p.nextToken()

	return stmt
}

func (p *Parser) parseIf() Statement {
	stmt := &IfStatement{}
	p.nextToken() // pula o 'if'

	// 1. Lendo a condição
	if p.curTok.Type != lexer.LPAREN {
		fmt.Println("Erro: esperado '(' após if")
		return nil
	}
	p.nextToken()

	cond := Expression{
		Left:     p.curTok.Literal, // Pode ser o nome da variável ou o valor "true"/"false"
		Operator: "",
		Right:    "",
	}
	p.nextToken()

	// MÁGICA: Se o próximo token NÃO for o fechamento de parênteses ')',
	// significa que temos uma comparação (>, <, ==)
	if p.curTok.Type != lexer.RPAREN {
		cond.Operator = p.curTok.Literal
		p.nextToken()

		cond.Right = p.curTok.Literal
		p.nextToken()
	}
	// Se for ')', a condição fica apenas com o Left preenchido, o que o Eval entende como booleano.

	stmt.Condition = cond

	if p.curTok.Type != lexer.RPAREN {
		fmt.Println("Erro: esperado ')' após condição")
		return nil
	}
	p.nextToken()

	// 2. Bloco IF { ... }
	if p.curTok.Type != lexer.LBRACE {
		fmt.Println("Erro: esperado '{' para iniciar bloco if")
		return nil
	}
	p.nextToken()

	for p.curTok.Type != lexer.RBRACE && p.curTok.Type != lexer.EOF {
		innerStmt := p.ParseStatement()
		if innerStmt != nil {
			stmt.Consequence = append(stmt.Consequence, innerStmt)
		}
	}

	if p.curTok.Type == lexer.RBRACE {
		p.nextToken() // pula '}'
	}

	// 3. BLOCO ELSE
	if p.curTok.Type == lexer.IDENT && p.curTok.Literal == "else" {
		p.nextToken() // pula 'else'

		if p.curTok.Type != lexer.LBRACE {
			fmt.Println("Erro: esperado '{' após else")
			return nil
		}
		p.nextToken() // pula '{'

		for p.curTok.Type != lexer.RBRACE && p.curTok.Type != lexer.EOF {
			innerStmt := p.ParseStatement()
			if innerStmt != nil {
				stmt.Alternative = append(stmt.Alternative, innerStmt)
			}
		}

		if p.curTok.Type == lexer.RBRACE {
			p.nextToken() // pula '}' do else
		}
	}

	return stmt
}

func (p *Parser) parseRepeat() Statement {
	stmt := &RepeatStatement{}
	p.nextToken() // pula 'repeat'

	stmt.Count = p.curTok.Literal // Pega o número ou a variável
	p.nextToken()

	if p.curTok.Type != lexer.LBRACE {
		fmt.Println("Erro: esperado '{' após quantidade no repeat")
		return nil
	}
	p.nextToken() // pula '{'

	for p.curTok.Type != lexer.RBRACE && p.curTok.Type != lexer.EOF {
		innerStmt := p.ParseStatement()
		if innerStmt != nil {
			stmt.Body = append(stmt.Body, innerStmt)
		}
	}

	if p.curTok.Type == lexer.RBRACE {
		p.nextToken() // pula '}'
	}
	return stmt
}
