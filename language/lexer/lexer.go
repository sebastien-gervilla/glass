package lexer

import token "glass/language/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	character    byte
}

func New(input string) *Lexer {
	lexer := &Lexer{
		input: input,
	}

	lexer.readCharacter()

	return lexer
}

func (lexer *Lexer) readCharacter() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.character = 0
	} else {
		lexer.character = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) Next() token.Token {
	var nextToken token.Token

	lexer.skipWhitespace()

	switch lexer.character {

	case '=':
		nextToken = newToken(token.ASSIGN, lexer.character)

	case ';':
		nextToken = newToken(token.SEMICOLON, lexer.character)

	case '(':
		nextToken = newToken(token.LPAREN, lexer.character)

	case ')':
		nextToken = newToken(token.RPAREN, lexer.character)

	case ',':
		nextToken = newToken(token.COMMA, lexer.character)

	case '+':
		nextToken = newToken(token.PLUS, lexer.character)

	case '{':
		nextToken = newToken(token.LBRACE, lexer.character)

	case '}':
		nextToken = newToken(token.RBRACE, lexer.character)

	case 0:
		nextToken.Literal = ""
		nextToken.Type = token.EOF

	default:
		if isValidCharacter(lexer.character) {
			nextToken.Literal = lexer.readIdentifier()
			nextToken.Type = token.LookupIdentifier(nextToken.Literal)
			return nextToken
		}

		if isDigit(lexer.character) {
			nextToken.Type = token.INT
			nextToken.Literal = lexer.readNumber()
			return nextToken
		}

		nextToken = newToken(token.ILLEGAL, lexer.character)
	}

	lexer.readCharacter()
	return nextToken
}

func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(character),
	}
}

func isValidCharacter(character byte) bool {
	return ('a' <= character && character <= 'z') ||
		('A' <= character && character <= 'Z') ||
		character == '_'
}

func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isValidCharacter(lexer.character) {
		lexer.readCharacter()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.character) {
		lexer.readCharacter()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.character == ' ' || lexer.character == '\t' || lexer.character == '\n' || lexer.character == '\r' {
		lexer.readCharacter()
	}
}
