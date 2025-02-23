package token

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// String returns a colored string representation of the token
func (t Token) String() string {
	// ANSI escape codes for colors
	const (
		blue    = "\033[34m" // Blue instead of cyan
		magenta = "\033[35m" // Magenta instead of yellow
		cyan    = "\033[36m" // Cyan instead of green
		reset   = "\033[0m"
	)
	return fmt.Sprintf("{%sType: %s%s%s, %sLiteral: %s%q%s}",
		magenta, // Type: in magenta
		cyan,    // Type value in cyan
		t.Type,
		reset,
		magenta, // Literal: in magenta
		cyan,    // Literal value in cyan
		t.Literal,
		reset, // Reset color at the end
	)
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER" // add, foobar, x, y
	INT        = "INT"
	STRING     = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ  = "=="
	NEQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENTIFIER
}
