package repl

import (
	"bufio"
	"fmt"
	"io"

	"ember_lang/ember_lang/evaluator"
	"ember_lang/ember_lang/lexer"
	"ember_lang/ember_lang/parser"
)

const (
	CYAN   = "\033[1;96m"
	RESET  = "\033[0m"
	PROMPT = CYAN + "⟶ " + RESET
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		lexer := lexer.New(line)
		parser := parser.New(lexer)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
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
