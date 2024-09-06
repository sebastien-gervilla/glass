package main

import (
	"bufio"
	"fmt"
	"glass/language/evaluator"
	"glass/language/lexer"
	"glass/language/object"
	"glass/language/parser"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: glass run <filename>")
		return
	}

	command := os.Args[1]
	filename := os.Args[2]

	if command == "run" {

		// File handling
		file, openError := os.Open(filename)
		if openError != nil {
			fmt.Println("Error reading file:", openError)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			fmt.Println("File is empty.")
		}

		firstLine := scanner.Text()

		// Interpreting
		environment := object.NewEnvironment()
		lexer := lexer.New(firstLine, func() (string, bool) {
			if !scanner.Scan() {
				return "", true
			}

			return scanner.Text(), false
		})

		parser := parser.New(lexer)
		program := parser.ParseProgram()

		errors := parser.GetErrors()
		if len(errors) > 0 {
			for _, err := range errors {
				log.Print(err)
			}
			return
		}

		evaluator.Evaluate(program, environment)
	} else {
		fmt.Println("Unknown command:", command)
	}
}
