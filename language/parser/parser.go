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
	INDEX       // array[indexs]
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
	token.LBRACKET:     INDEX,
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

	// Registering prefixes
	parser.prefixParsingFunctions = make(map[token.TokenType]prefixParsingFunction)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.NOT, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunction)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
	parser.registerPrefix(token.LBRACE, parser.parseHashLiteral)

	// Registering infixes
	parser.infixParsingFunctions = make(map[token.TokenType]infixParsingFunction)
	parser.registerInfix(token.EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.GREATER_THAN, parser.parseInfixExpression)
	parser.registerInfix(token.LESS_THAN, parser.parseInfixExpression)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

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
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.nextToken()
	}

	return program
}

// Statements

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {

	case token.LET:
		return parser.parseLetStatement()

	case token.RETURN:
		return parser.parseReturnStatement()

	default:
		return parser.parseExpressionStatement()

	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: parser.currentToken,
	}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Identifier = &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()

	statement.Expression = parser.parseExpression(LOWEST)

	for parser.isPeekToken(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: parser.currentToken,
	}

	parser.nextToken()

	statement.Expression = parser.parseExpression(LOWEST)

	for parser.isPeekToken(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

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
		if parser.currentToken.Type == token.ILLEGAL {
			parser.addIllegalTokenError(parser.currentToken)
			return nil
		}

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

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: parser.currentToken,
	}

	value, conversionError := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if conversionError != nil {
		message := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, message)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: parser.currentToken,
		Value: parser.isCurrentToken(token.TRUE),
	}
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.nextToken()

	expression.Expression = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:          parser.currentToken,
		LeftExpression: leftExpression,
		Operator:       parser.currentToken.Literal,
	}

	precedence := parser.getCurrentPrecedence()
	parser.nextToken()
	expression.RightExpression = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	expression := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: parser.currentToken,
	}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if !parser.isPeekToken(token.ELSE) {
		return expression
	}

	// Parsing else expression
	parser.nextToken()
	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Alternative = parser.parseBlockStatement()
	return expression
}

func (parser *Parser) parseFunction() ast.Expression {
	function := &ast.Function{
		Token: parser.currentToken,
	}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	function.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	function.Body = parser.parseBlockStatement()
	return function
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}

	if parser.isPeekToken(token.RPAREN) {
		parser.nextToken()
		return parameters
	}

	parser.nextToken()

	identifier := &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
	parameters = append(parameters, identifier)

	for parser.isPeekToken(token.COMMA) {
		parser.nextToken()
		parser.nextToken()

		identifier := &ast.Identifier{
			Token: parser.currentToken,
			Value: parser.currentToken.Literal,
		}
		parameters = append(parameters, identifier)
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:     parser.currentToken,
		Function:  function,
		Arguments: parser.parseExpressions(token.RPAREN),
	}

	return expression
}

func (parser *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{
		Token: parser.currentToken,
	}

	array.Elements = parser.parseExpressions(token.RBRACKET)

	return array
}

func (parser *Parser) parseExpressions(pairType token.TokenType) []ast.Expression {
	expressions := []ast.Expression{}

	if parser.isPeekToken(pairType) {
		parser.nextToken()
		return expressions
	}

	parser.nextToken()
	expressions = append(expressions, parser.parseExpression(LOWEST))

	for parser.isPeekToken(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		expressions = append(expressions, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(pairType) {
		return nil
	}

	return expressions
}

func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{
		Token: parser.currentToken,
		Left:  left,
	}

	parser.nextToken()
	expression.Index = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RBRACKET) {
		return nil
	}

	return expression
}

func (parser *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token: parser.currentToken,
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	for !parser.isPeekToken(token.RBRACE) {
		parser.nextToken()
		key := parser.parseExpression(LOWEST)

		if !parser.isPeekToken(token.COLON) {
			hash.Pairs[key] = key
		} else {
			parser.nextToken()
			parser.nextToken()

			value := parser.parseExpression(LOWEST)
			hash.Pairs[key] = value
		}

		if !parser.isPeekToken(token.RBRACE) && !parser.isPeekToken(token.COMMA) {
			parser.addUnexepectedTokenError(token.COMMA, parser.currentToken)
			return nil
		}
	}

	return hash
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

	parser.addUnexepectedTokenError(parser.peekToken.Type, parser.currentToken)
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

// Errors

func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) addUnexepectedTokenError(expectedType token.TokenType, unexpected token.Token) {
	message := fmt.Sprintf(
		"Expected token %s, got %s instead (l.%d:p.%d)",
		expectedType,
		unexpected.Type,
		unexpected.Line,
		unexpected.Position,
	)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) addIllegalTokenError(token token.Token) {
	message := fmt.Sprintf(
		"Illegal token %q found (l.%d:p.%d)",
		token.Literal,
		token.Line,
		token.Position,
	)
	parser.errors = append(parser.errors, message)
}
