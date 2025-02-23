package main

import (
	"ember_lang/repl"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Ember Programming Language v0.0.1 (prototype)\n")
	fmt.Printf("Type \"help\" for more information.\n")

	repl.Start(os.Stdin, os.Stdout)
}
