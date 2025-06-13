package language_test

import (
	"bufio"
	"glass/language/evaluator"
	"glass/language/lexer"
	"glass/language/object"
	"glass/language/parser"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLanguage(testing *testing.T) {

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		testing.Fatalf("unable to get caller info")
	}

	currentDirectory := filepath.Dir(currentFile)

	mainDirectory := filepath.Join(currentDirectory, "../../glass")
	mainFile := filepath.Join(mainDirectory, "./main.glass")

	runDirectory := filepath.Dir(mainDirectory)

	// File handling
	file, openError := os.Open(mainFile)
	if openError != nil {
		testing.Fatalf("could not find glass file %s: %s", mainFile, openError)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		testing.Fatal("File is empty.")
	}

	firstLine := scanner.Text()

	// Interpreting
	programEnvironment := object.NewProgramEnvironment(runDirectory)
	moduleEnvironment := object.NewEnvironment(mainFile, programEnvironment)
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
		testing.Fatal(result.Inspect())
	}
}
