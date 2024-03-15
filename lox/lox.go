package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/meyegui/golox/scanner"
)

// ------------------------------------------------------------ TYPE DEFINITION

type Lox struct {
	hadError bool
}

// ------------------------------------------------------------ CODE EVALUATION

func (l *Lox) RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	l.Run(string(bytes))
	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println()
			fmt.Println("Good bye!")
			break
		} else if err != nil {
			panic(err)
		}

		line = strings.TrimSpace(line)
		if len(line) > 0 {
			l.Run(line)
			l.hadError = false
		}
	}
}

func (l *Lox) Run(source string) {
	s := scanner.NewScanner(source, l)
	tokens := s.ScanTokens()

	fmt.Println("+----------------------+----------------------+----------------------+")
	fmt.Println("|       Scanned        |        Token         |        Value         |")
	fmt.Println("+----------------------+----------------------+----------------------+")
	for _, token := range tokens {
		formatString := "| %-20s | %-20s | "
		switch token.Type {
		case scanner.STRING:
			formatString += "%-20s"

		case scanner.NUMBER:
			formatString += "% 20.5f"

		default:
			formatString += "%s"
		}
		formatString += " |\n"

		value := token.Literal
		if value == nil {
			value = strings.Repeat(" ", 20)
		}

		fmt.Printf(formatString, token.Lexeme, token.Type, value)
		// fmt.Println("+----------------------+----------------------+----------------------+")
	}
	fmt.Println("+----------------------+----------------------+----------------------+")
}

// ------------------------------------------------------------ ErrorReporter INTERFACE

func (l *Lox) Error(line int, message string) {
	l.Report(line, "", message)
}

func (l *Lox) Report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}
