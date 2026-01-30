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
// ===== AST =====
//

type Statement interface {
	statementNode()
}

type Program struct {
	Statements []Statement
}

type PrintStatement struct {
	Args []string
}

func (ps *PrintStatement) statementNode() {}

//
// ===== PARSER =====
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
	if p.curTok.Type == lexer.IDENT && p.curTok.Literal == "print" {
		return p.parsePrint()
	}

	fmt.Println("unknown statement:", p.curTok.Literal)
	p.nextToken()
	return nil
}

func (p *Parser) parsePrint() Statement {
	// consume 'print'
	p.nextToken()

	if p.curTok.Type != lexer.LPAREN {
		fmt.Println("expected '(' after print")
		return nil
	}
	p.nextToken()

	args := []string{}

	for p.curTok.Type != lexer.RPAREN {
		if p.curTok.Type == lexer.STRING || p.curTok.Type == lexer.IDENT {
			args = append(args, p.curTok.Literal)
			p.nextToken()
		}

		if p.curTok.Type == lexer.COMMA {
			p.nextToken()
		}
	}

	// consume ')'
	p.nextToken()

	if p.curTok.Type != lexer.SEMICOLON {
		fmt.Println("expected ';' after print")
		return nil
	}
	p.nextToken()

	return &PrintStatement{Args: args}
}
