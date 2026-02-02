package parser

import (
	"fmt"

	"github.com/tiagollopes/okay/lexer"
)

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
	Value string
}

func (vds *VarDeclarationStatement) statementNode() {}

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
	case lexer.IDENT:
		if p.curTok.Literal == "print" {
			return p.parsePrint()
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
	p.nextToken()

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

	if p.curTok.Type == lexer.STRING || p.curTok.Type == lexer.NUMBER || p.curTok.Type == lexer.IDENT {
		stmt.Value = p.curTok.Literal
		p.nextToken()
	}

	if p.curTok.Type != lexer.SEMICOLON {
		fmt.Println("Erro: esperado ';'")
		return nil
	}
	p.nextToken()

	return stmt
}
