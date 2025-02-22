package main

import (
	"ember_lang/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Ember programming language!\n", user.Username)

	repl.Start(os.Stdin, os.Stdout)
}
