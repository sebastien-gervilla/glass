package lexer

import (
	token "glass/language/token"
)

type Lexer struct {
	line         string
	lineNumber   int
	position     int
	readPosition int
	character    byte
	getNextLine  func() (string, bool)
}

func New(line string, getNextLine func() (string, bool)) *Lexer {
	lexer := &Lexer{
		line:        line,
		lineNumber:  1,
		getNextLine: getNextLine,
	}

	lexer.readCharacter()

	return lexer
}

func (lexer *Lexer) readCharacter() {
	if lexer.readPosition >= len(lexer.line) {
		nextLine, isFileEnd := lexer.getNextLine()
		lexer.lineNumber++

		for nextLine == "" {
			if isFileEnd {
				lexer.character = 0
				return
			}

			nextLine, isFileEnd = lexer.getNextLine()
			lexer.lineNumber++
		}

		lexer.line = nextLine
		lexer.position = 0
		lexer.readPosition = 0
		lexer.character = lexer.line[lexer.readPosition]
	} else {
		lexer.character = lexer.line[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) peekCharacter() byte {
	if lexer.readPosition >= len(lexer.line) {
		return 0
	}

	return lexer.line[lexer.readPosition]
}

func (lexer *Lexer) Next() token.Token {

	lexer.skipWhitespace()

	nextToken := token.Token{
		Literal:  string(lexer.character),
		Line:     lexer.lineNumber,
		Position: lexer.position,
	}

	switch lexer.character {

	case '=':
		if lexer.peekCharacter() == '=' {
			// Advance to peeked character
			lexer.readCharacter()
			nextToken.Type = token.EQUALS
			nextToken.Literal = "=="
		} else {
			nextToken.Type = token.ASSIGN
		}

	case '!':
		if lexer.peekCharacter() == '=' {
			// Advance to peeked character
			lexer.readCharacter()
			nextToken.Type = token.NOT_EQUALS
			nextToken.Literal = "=="
		} else {
			nextToken.Type = token.NOT
		}

	case ';':
		nextToken.Type = token.SEMICOLON

	case '(':
		nextToken.Type = token.LPAREN

	case ')':
		nextToken.Type = token.RPAREN

	case '.':
		nextToken.Type = token.DOT

	case ',':
		nextToken.Type = token.COMMA

	case '+':
		nextToken.Type = token.PLUS

	case '-':
		nextToken.Type = token.MINUS

	case '{':
		nextToken.Type = token.LBRACE

	case '}':
		nextToken.Type = token.RBRACE

	case '[':
		nextToken.Type = token.LBRACKET

	case ']':
		nextToken.Type = token.RBRACKET

	case ':':
		nextToken.Type = token.COLON

	case '/':
		nextToken.Type = token.SLASH

	case '*':
		nextToken.Type = token.ASTERISK

	case '"':
		nextToken.Type = token.STRING
		nextToken.Literal = lexer.readString()

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

		nextToken.Type = token.ILLEGAL
	}

	lexer.readCharacter()
	return nextToken
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

	return lexer.line[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.character) {
		lexer.readCharacter()
	}

	return lexer.line[position:lexer.position]
}

func (lexer *Lexer) readString() string {
	position := lexer.position + 1

	lexer.readCharacter()
	for lexer.character != '"' {
		lexer.readCharacter()
	}

	return lexer.line[position:lexer.position]
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.character == ' ' || lexer.character == '\t' || lexer.character == '\n' || lexer.character == '\r' {
		lexer.readCharacter()
	}
}
