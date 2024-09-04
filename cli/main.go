package main

import (
	"fmt"
	"glass/language/interpreter"
	"glass/language/object"
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
		// Your logic to interpret the file
		code, err := os.ReadFile(filename)

		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		environment := object.NewEnvironment()
		interpreter.Interpret(string(code), environment) // Replace with your interpreter's function
	} else {
		fmt.Println("Unknown command:", command)
	}
}
