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

		_, _ = io.WriteString(out, program.String())
		_, _ = io.WriteString(out, "\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "\x1b[31mFailed to parse program!\x1b[0m\n")
	_, _ = io.WriteString(out, "\x1b[31mParser errors:\x1b[0m\n")

	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
