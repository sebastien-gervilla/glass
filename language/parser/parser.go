package parser

import (
	"fmt"
	ast "glass/language/ast"
	lexer "glass/language/lexer"
	token "glass/language/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -expression or !expression
	CALL        // myFunction(expression, expression)
)

var precedences = map[token.TokenType]int{
	token.EQUALS:       EQUALS,
	token.NOT_EQUALS:   EQUALS,
	token.LESS_THAN:    LESSGREATER,
	token.GREATER_THAN: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.LPAREN:       CALL,
}

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

func (parser *Parser) registerPrefix(tokenType token.TokenType, function prefixParsingFunction) {
	parser.prefixParsingFunctions[tokenType] = function
}
func (parser *Parser) registerInfix(tokenType token.TokenType, function infixParsingFunction) {
	parser.infixParsingFunctions[tokenType] = function
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := ast.Program{
		Statements: []ast.Statement{},
	}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.nextToken()
	}

	return nil
}

// Statements

func (parser *Parser) parseStatement() ast.Statement {
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
func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{
		Token:      parser.currentToken,
		Statements: []ast.Statement{},
	}

	parser.nextToken()
	for !parser.isCurrentToken(token.RBRACE) && !parser.isCurrentToken(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}

		parser.nextToken()
	}

	return blockStatement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{}

	statement.Expression = parser.parseExpression(LOWEST)

	if parser.isPeekToken(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

// Expressions

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParsingFunctions[parser.currentToken.Type]

	if prefix == nil {
		message := fmt.Sprintf("no prefix parse function found for %q token", parser.currentToken.Type)
		parser.errors = append(parser.errors, message)
		return nil
	}

	leftExpression := prefix()

	for !parser.isPeekToken(token.SEMICOLON) && precedence < parser.getPeekPrecedence() {
		infix := parser.infixParsingFunctions[parser.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		parser.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
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

func (parser *Parser) getCurrentPrecedence() int {
	precedence, ok := precedences[parser.currentToken.Type]
	if !ok {
		return LOWEST
	}

	return precedence
}

func (parser *Parser) getPeekPrecedence() int {
	precedence, ok := precedences[parser.peekToken.Type]
	if !ok {
		return LOWEST
	}

	return precedence
}

func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) addParseError(token token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", token, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}
