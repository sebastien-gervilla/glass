package interpreter

import (
	"glass/language/evaluator"
	"glass/language/lexer"
	"glass/language/object"
	"glass/language/parser"
)

func Interpret(code string, environment *object.Environment) []string {
	lexer := lexer.New(code)
	parser := parser.New(lexer)

	program := parser.ParseProgram()
	if len(parser.GetErrors()) > 0 {
		return parser.GetErrors()
	}

	evaluator.Evaluate(program, environment)

	return []string{}
}
