package main

import (
	"ember_lang/ember_lang/evaluator"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/object"
	"ember_lang/ember_lang/parser"
	"ember_lang/ember_lang/repl"
	"fmt"
	"os"
	"path/filepath"
)

var debug = os.Getenv("DEBUG") == "1"

func main() {
	if len(os.Args) > 1 {
		// Execute file mode
		executeFile(os.Args[1])
	} else {
		// REPL mode
		fmt.Printf("Ember Programming Language v0.0.1 (prototype)\n")
		fmt.Printf("Type \"help\" for more information.\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func executeFile(path string) {
	// Check file extension
	if filepath.Ext(path) != ".em" {
		fmt.Printf("Error: File must have .em extension\n")
		os.Exit(1)
	}

	// Read file
	code, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	if debug {
		fmt.Printf("\n========================= Source Code =========================\n%s\n", string(code))
	}

	// Lexical analysis
	l := lexer.New(string(code))
	if debug {
		fmt.Println("\n=========================== Tokens ===========================")
		for tok := l.NextToken(); tok.Type != "EOF"; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
		l = lexer.New(string(code)) // Reset lexer for parsing
	}

	// Parsing
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	if debug {
		fmt.Printf("\n=========================== AST ===========================\n%s\n", program.String())
	}

	// Evaluation
	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if debug {
		fmt.Printf("\n=========================== Result ===========================\n")
	}

	if result != nil {
		fmt.Println(result.Inspect())
	}
}

func printParserErrors(errors []string) {
	fmt.Println("\x1b[31mParser errors:\x1b[0m")
	for _, msg := range errors {
		fmt.Printf("\t%s\n", msg)
	}
}
