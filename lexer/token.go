package lexer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Keywords
	SERVICE TokenType = "SERVICE"
	PORT    TokenType = "PORT"
	LET     TokenType = "LET"

	// Identifiers & literals
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	// Symbols
	PLUS      TokenType = "+"
	ASSIGN    TokenType = "="
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	EOF TokenType = "EOF"
)

