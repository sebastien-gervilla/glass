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
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: glass <command> <filename>")
		return
	}

	command := os.Args[1]
	filename := os.Args[2]

	fullpath, absError := filepath.Abs(filename)
	if absError != nil {
		log.Fatal("Error getting absolute path:", absError)
		return
	}

	runDirectory := filepath.Dir(fullpath)

	if command == "run" {

		// File handling
		file, openError := os.Open(filename)
		if openError != nil {
			log.Fatal("Error reading file:", openError)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			fmt.Println("File is empty.")
		}

		firstLine := scanner.Text()

		// Interpreting
		programEnvironment := object.NewProgramEnvironment(runDirectory)
		moduleEnvironment := object.NewEnvironment(filename, programEnvironment)
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

		result := evaluator.Evaluate(program, moduleEnvironment)
		if result != nil && result.GetType() == object.ERROR_OBJECT {
			log.Fatal(result.Inspect())
		}
	} else {
		fmt.Println("Unknown command:", command)
	}
}
