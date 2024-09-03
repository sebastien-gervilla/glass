package repl

import (
	"bufio"
	"fmt"
	"glass/language/evaluator"
	"glass/language/lexer"
	"glass/language/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lexer := lexer.New(line)
		parser := parser.New(lexer)

		program := parser.ParseProgram()
		if len(parser.GetErrors()) != 0 {
			printParserErrors(out, parser.GetErrors())
			continue
		}

		evaluated := evaluator.Evaluate(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, parserErrors []string) {
	for _, message := range parserErrors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
