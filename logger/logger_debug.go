package logger

import (
	"ember_lang/ember_lang/ast"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/object"
	"fmt"
)

type LoggerDebug struct {
	debug bool
}

func (l *LoggerDebug) Debug(message string) {
	if l.debug {
		fmt.Println(message)
	}
}

const separator = "========================================"

func LogSourceCode(sourceCode []byte) {
	fmt.Printf("\n%s Source Code %s\n\n", separator, separator)
	fmt.Println(string(sourceCode))
}

func LogTokens(lexer *lexer.Lexer) {
	fmt.Printf("\n%s Tokens %s\n\n", separator, separator)
	for tok := lexer.NextToken(); tok.Type != "EOF"; tok = lexer.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}

func LogAST(ast *ast.Program) {
	fmt.Printf("\n%s AST %s\n\n", separator, separator)
	for _, statement := range ast.Statements {
		fmt.Printf("%s\n", statement.String())
	}
}

func LogResult(result object.Object) {
	fmt.Printf("\n%s Result %s\n", separator, separator)
	fmt.Println(result.Inspect())
}
