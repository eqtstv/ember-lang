package repl

import (
	"bufio"
	"fmt"
	"io"

	"ember_lang/lexer"
	"ember_lang/parser"
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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
