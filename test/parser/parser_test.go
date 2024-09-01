package parser_test

import (
	"glass/language/ast"
	lexer "glass/language/lexer"
	parser "glass/language/parser"
	"testing"
)

func TestLetStatements(testing *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(testing, parser)

		if len(program.Statements) != 1 {
			testing.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		statement := program.Statements[0]
		if !testLetStatement(testing, statement, test.expectedIdentifier) {
			return
		}

		// val := statement.(*ast.LetStatement).Expression
		// if !testLiteralExpression(testing, val, test.expectedValue) {
		// 	return
		// }
	}
}

func testLetStatement(testing *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		testing.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		testing.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStatement.Identifier.Value != name {
		testing.Errorf("letStatement.Name.Value not '%s'. got=%s", name, letStatement.Identifier.Value)
		return false
	}

	if letStatement.Identifier.TokenLiteral() != name {
		testing.Errorf("s.Name not '%s'. got=%s", name, letStatement.Identifier)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, parser *parser.Parser) {
	errors := parser.GetErrors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
