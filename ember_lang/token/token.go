package token

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type       TokenType
	Literal    string
	LineNumber int
}

// String returns a colored string representation of the token
func (t Token) String() string {
	// Atom One Dark theme colors (VS Code TypeScript-like)
	const (
		purple = "\033[38;2;198;120;221m" // Keywords, control flow
		blue   = "\033[38;2;97;175;239m"  // Function keywords
		red    = "\033[38;2;224;108;117m" // Errors, important tokens
		cyan   = "\033[38;2;86;182;194m"  // Built-in types, numbers
		orange = "\033[38;2;209;154;102m" // Strings
		white  = "\033[38;2;171;178;191m" // Default text, identifiers
		gray   = "\033[38;2;92;99;112m"   // Punctuation, delimiters
		green  = "\033[38;2;152;195;121m" // Special keywords (true/false)
		reset  = "\033[0m"
	)

	// Get color based on token type
	var typeColor string
	switch t.Type {
	case ILLEGAL, EOF:
		typeColor = red
	case FUNCTION:
		typeColor = blue
	case LET, IF, ELSE, RETURN, WHILE, FOR:
		typeColor = purple
	case TRUE, FALSE:
		typeColor = green
	case PLUS, MINUS, BANG, ASTERISK, SLASH, LT, GT, LTE, GTE, EQ, NEQ, ASSIGN:
		typeColor = white
	case INT:
		typeColor = cyan
	case STRING:
		typeColor = orange
	case COMMA, SEMICOLON, COLON, LPAREN, RPAREN, LBRACE, RBRACE, LBRACKET, RBRACKET:
		typeColor = gray
	case IDENTIFIER:
		typeColor = white
	default:
		typeColor = white
	}

	return fmt.Sprintf("%s%q%s %s=>%s %s%s%s",
		typeColor, // Value in dynamic color
		t.Literal,
		reset,
		gray, // Arrow in gray
		reset,
		typeColor, // Type in same color as value
		t.Type,
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
	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	BANG     = "BANG"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"

	// Suffix operators
	INCREMENT = "INCREMENT"

	LT  = "LT"
	GT  = "GT"
	LTE = "LTE"
	GTE = "GTE"

	EQ  = "EQ"
	NEQ = "NEQ"

	// Delimiters
	COMMA     = "COMMA"     // ,
	SEMICOLON = "SEMICOLON" // ;
	COLON     = "COLON"     // :
	LPAREN    = "LPAREN"    // (
	RPAREN    = "RPAREN"    // )
	LBRACE    = "LBRACE"    // {
	RBRACE    = "RBRACE"    // }
	LBRACKET  = "LBRACKET"  // [
	RBRACKET  = "RBRACKET"  // ]

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// Loops
	WHILE = "WHILE"
	FOR   = "FOR"

	// Mutable
	MUT = "MUT"

	// Special
	COMMENT = "COMMENT"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
	"mut":    MUT,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENTIFIER
}
