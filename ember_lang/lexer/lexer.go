package lexer

import (
	"ember_lang/ember_lang/token"
)

type Lexer struct {
	input        string
	position     int  // current pos in input
	readPosition int  // current reading position
	ch           byte // current char under examination
	lineNumber   int  // current line number
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		default:
			tok = newToken(token.ASSIGN, l.ch, l.lineNumber)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch, l.lineNumber)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.lineNumber)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.lineNumber)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.lineNumber)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.lineNumber)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.lineNumber)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.lineNumber)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.lineNumber)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.lineNumber)
	case ':':
		tok = newToken(token.COLON, l.ch, l.lineNumber)
	case '-':
		tok = newToken(token.MINUS, l.ch, l.lineNumber)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.lineNumber)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.lineNumber)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch, l.lineNumber)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch, l.lineNumber)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch, l.lineNumber)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			tok.LineNumber = l.lineNumber
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			tok.LineNumber = l.lineNumber
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch, l.lineNumber)
	}

	// Move to next character
	l.readChar()

	return tok
}

func New(input string) *Lexer {
	l := &Lexer{input: input, lineNumber: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// Check if end of input
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	// Update position
	l.position = l.readPosition

	// Move to next character
	l.readPosition++
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.lineNumber++
		}
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte, lineNumber int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), LineNumber: lineNumber}
}
