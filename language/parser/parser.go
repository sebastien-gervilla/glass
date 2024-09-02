package parser

import (
	"fmt"
	ast "glass/language/ast"
	lexer "glass/language/lexer"
	token "glass/language/token"
)

type (
	prefixParsingFunction func() ast.Expression
	infixParsingFunction  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string

	prefixParsingFunctions map[token.TokenType]prefixParsingFunction
	infixParsingFunctions  map[token.TokenType]infixParsingFunction
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens, so currentToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.Next()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := ast.Program{
		Statements: []ast.Statement{},
	}

	for parser.currentToken.Type != token.EOF {
		statement := parseStatement(parser)

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.nextToken()
	}

	return nil
}

func parseStatement(parser *Parser) ast.Statement {
	switch parser.currentToken.Type {

	case token.LET:
		return parser.parseLetStatement()

	case token.RETURN:
		return parser.parseReturnStatement()

	default:
		return nil

	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: parser.currentToken,
	}

	if parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Identifier = &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}

	if parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Expresssions
	for parser.isCurrentToken(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: parser.currentToken,
	}

	parser.nextToken()
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !parser.isCurrentToken(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

// Utils

func (parser *Parser) isCurrentToken(token token.TokenType) bool {
	return parser.currentToken.Type == token
}

func (parser *Parser) isPeekToken(token token.TokenType) bool {
	return parser.peekToken.Type == token
}

func (parser *Parser) expectPeek(token token.TokenType) bool {
	if parser.isPeekToken(token) {
		parser.nextToken()
		return true
	}

	parser.addParseError(token)
	return false
}

func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) addParseError(token token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", token, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}
