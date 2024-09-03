package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"let":      LET,
	"function": FUNCTION,
	"return":   RETURN,
}

const (
	// Misc
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	// Operators
	ASSIGN       = "="
	EQUALS       = "=="
	NOT_EQUALS   = "!="
	GREATER_THAN = ">"
	LESS_THAN    = "<"
	PLUS         = "+"
	MINUS        = "-"
	SLASH        = "/"
	ASTERISK     = "*"
	NOT          = "!"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

func LookupIdentifier(identifier string) TokenType {
	keyword, found := keywords[identifier]
	if found {
		return keyword
	}

	return IDENTIFIER
}
