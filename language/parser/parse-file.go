package parser

import (
	"bufio"
	"fmt"
	"glass/language/ast"
	"glass/language/lexer"
	"log"
	"os"
)

func GetParsedFile(filepath string) *ast.Program {
	// File handling
	file, openError := os.Open(filepath)
	if openError != nil {
		fmt.Println("Error reading file:", openError)
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Println("File is empty.")
	}

	firstLine := scanner.Text()

	// Interpreting
	lexer := lexer.New(firstLine, func() (string, bool) {
		if !scanner.Scan() {
			return "", true
		}

		return scanner.Text(), false
	})

	parser := New(lexer)
	program := parser.ParseProgram()

	errors := parser.GetErrors()
	if len(errors) > 0 {
		for _, err := range errors {
			log.Print(err)
		}
		return nil
	}

	return program
}
