package lexer

import "unicode"

type Lexer struct {
	input []rune
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: EOF}
	}

	ch := l.input[l.pos]

	// Symbols
	switch ch {
	case '=':
		l.pos++
		return Token{Type: ASSIGN, Literal: "="}
	case '{':
		l.pos++
		return Token{Type: LBRACE, Literal: "{"}
	case '}':
		l.pos++
		return Token{Type: RBRACE, Literal: "}"}
	case '(':
		l.pos++
		return Token{Type: LPAREN, Literal: "("}
	case ')':
		l.pos++
		return Token{Type: RPAREN, Literal: ")"}
	case ',':
		l.pos++
		return Token{Type: COMMA, Literal: ","}
	case ';':
		l.pos++
		return Token{Type: SEMICOLON, Literal: ";"}
	case '"':
		return l.readString()
	}

	// Identifier or keyword
	if unicode.IsLetter(ch) {
		start := l.pos
		for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || unicode.IsDigit(l.input[l.pos])) {
			l.pos++
		}
		lit := string(l.input[start:l.pos])

		switch lit {
			case "service":
				return Token{Type: SERVICE, Literal: lit}
			case "port":
				return Token{Type: PORT, Literal: lit}
			case "let": // <--- NOVO
				return Token{Type: LET, Literal: lit}
			default:
				return Token{Type: IDENT, Literal: lit}
			}
	}

	// Number
	if unicode.IsDigit(ch) {
		start := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
			l.pos++
		}
		return Token{Type: NUMBER, Literal: string(l.input[start:l.pos])}
	}

	// Unknown char: skip
	l.pos++
	return l.NextToken()
}

func (l *Lexer) readString() Token {
	l.pos++ // skip opening quote
	start := l.pos

	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.pos++
	}

	lit := string(l.input[start:l.pos])
	l.pos++ // skip closing quote

	return Token{Type: STRING, Literal: lit}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
}

