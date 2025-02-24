package repl

import (
	"bufio"
	"fmt"
	"io"

	"ember_lang/ember_lang/evaluator"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/object"
	"ember_lang/ember_lang/parser"
)

const (
	CYAN   = "\033[1;96m"
	RESET  = "\033[0m"
	PROMPT = CYAN + "⟶ " + RESET
)

func Start(in io.Reader, out io.Writer, debug string) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

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
