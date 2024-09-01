package main

import (
	lexer "glass/language/lexer"
	token "glass/language/token"
	"testing"
)

// var input = `let five = 5;
// 		let ten = 10;

// 		let add = fn(x, y) {
// 		x + y;
// 		};

// 		let result = add(five, ten);
// 		!-/*5;
// 		5 < 10 > 5;

// 		if (5 < 10) {
// 			return true;
// 		} else {
// 			return false;
// 		}

// 		10 == 10;
// 		10 != 9;
// 	`

var input = `let five = 5;`

func TestNextToken(t *testing.T) {

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "ten"},
		// {token.ASSIGN, "="},
		// {token.INT, "10"},
		// {token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "add"},
		// {token.ASSIGN, "="},
		// {token.FUNCTION, "fn"},
		// {token.LPAREN, "("},
		// {token.IDENT, "x"},
		// {token.COMMA, ","},
		// {token.IDENT, "y"},
		// {token.RPAREN, ")"},
		// {token.LBRACE, "{"},
		// {token.IDENT, "x"},
		// {token.PLUS, "+"},
		// {token.IDENT, "y"},
		// {token.SEMICOLON, ";"},
		// {token.RBRACE, "}"},
		// {token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "result"},
		// {token.ASSIGN, "="},
		// {token.IDENT, "add"},
		// {token.LPAREN, "("},
		// {token.IDENT, "five"},
		// {token.COMMA, ","},
		// {token.IDENT, "ten"},
		// {token.RPAREN, ")"},
		// {token.SEMICOLON, ";"},
	}

	lexer := lexer.New(input)

	for index, test := range tests {
		tok := lexer.Next()

		if tok.Type != test.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				index, test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				index, test.expectedLiteral, tok.Literal)
		}
	}
}
