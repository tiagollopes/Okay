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
	case '>':
		l.pos++
		return Token{Type: GT, Literal: ">"}
	case '<':
		l.pos++
		return Token{Type: LT, Literal: "<"}
	case '+':
		l.pos++
		return Token{Type: PLUS, Literal: "+"}
	case '-':
		l.pos++
		return Token{Type: MINUS, Literal: "-"}
	case '*':
		l.pos++
		return Token{Type: ASTERISK, Literal: "*"}
	case '/':
		// Espiamos o próximo caractere sem mover o cursor ainda
		if l.peekChar() == '/' {
			// É um comentário! Vamos pular até o fim da linha
			l.skipComment()
			return l.NextToken() // Reinicia a busca pelo próximo token real
		}
		// Se for apenas um '/', é divisão normal
		l.pos++
		return Token{Type: SLASH, Literal: "/"}
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

	// Identifier or keyword (Suporta letras e Underline _)
	if unicode.IsLetter(ch) || ch == '_' {
		start := l.pos
		for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
			l.pos++
		}
		lit := string(l.input[start:l.pos])

		switch lit {
		case "service":
			return Token{Type: SERVICE, Literal: lit}
		case "port":
			return Token{Type: PORT, Literal: lit}
		case "let":
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

// peekChar olha o próximo caractere sem avançar a posição
func (l *Lexer) peekChar() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

// skipComment pula todos os caracteres até encontrar uma quebra de linha
func (l *Lexer) skipComment() {
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.pos++
	}
}
