package repl

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"ember_lang/ember_lang/evaluator"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/object"
	"ember_lang/ember_lang/parser"
)

const (
	CYAN      = "\033[1;96m"
	RESET     = "\033[0m"
	PROMPT    = CYAN + "⟶ " + RESET
	HELP_TEXT = `
Ember Programming Language REPL
=============================

Commands:
  help              Show this help message
  exit, quit        Exit the REPL

Basic Syntax:
------------
  let x = 5;                      // Variable declaration
  let addOne = fn(x) { x + 1 };   // Function definition
  if (x > 0) { ... }              // Conditional statement
  while (x < 5) { ... }           // While loop
  for (let i = 0; i < 5; i++) {   // For loop
    ...
  }
  [1, 2, 3]                       // Array literal
  {"a": 1, "b": 2}                // Hash/map literal

Built-in Functions:
-----------------
  print(value)      Print value to console
  len(array)        Get length of array or string
  type(value)       Get type of value
  push(arr, item)   Append item to array

Type 'help' for this message, 'exit' or 'quit' to exit.
`
)

func Start(in io.Reader, out io.Writer, debug string) {
	env := object.NewEnvironment()

	readline, err := readline.NewEx(&readline.Config{
		Prompt:          PROMPT,
		HistoryFile:     os.Getenv("HOME") + "/.ember_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer readline.Close()

	fmt.Fprintf(out, "Ember Programming Language v0.0.1 (prototype)\n")
	fmt.Fprintf(out, "Type \"help\" for more information.\n")

	for {
		line, err := readline.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch line {
		case "help":
			fmt.Fprint(out, HELP_TEXT)
			continue
		case "exit", "quit":
			return
		}

		if debug == "1" {
			fmt.Printf("\n========================= Source Code =========================\n%s\n", line)
		}

		lexer := lexer.New(line)

		if debug == "1" {
			fmt.Println("\n=========================== Tokens ===========================")
			for tok := lexer.NextToken(); tok.Type != "EOF"; tok = lexer.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		}

		parser := parser.New(lexer)

		program := parser.ParseProgram()

		if debug == "1" {
			fmt.Printf("\n=========================== AST ===========================\n%s\n", program.String())
		}

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			_, _ = io.WriteString(out, evaluated.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "\033[31mFailed to parse program!\033[0m\n\n")
	_, _ = io.WriteString(out, "\033[31mPARSER ERRORS:\033[0m\n")

	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
