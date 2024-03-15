package main

import (
	"fmt"
	"os"

	"github.com/meyegui/golox/lox"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	}

	interpreter := lox.Lox{}
	if len(args) == 2 {
		interpreter.RunFile(args[1])
	} else {
		interpreter.RunPrompt()
	}
}
