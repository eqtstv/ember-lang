package repl

import (
	"bufio"
	"fmt"
	"io"

	"ember_lang/lexer"
	"ember_lang/token"
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
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
