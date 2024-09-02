package ast

import (
	"bytes"
	token "glass/language/token"
)

// Interfaces
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program
type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	}

	return ""
}

func (program *Program) String() string {
	var buffer bytes.Buffer

	for _, statement := range program.Statements {
		buffer.WriteString(statement.String())
	}

	return buffer.String()
}

// Identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) String() string       { return identifier.Value }

// Let statement
type LetStatement struct {
	Token      token.Token
	Identifier *Identifier
	Expression Expression
}

func (statement *LetStatement) statementNode()       {}
func (statement *LetStatement) TokenLiteral() string { return statement.Token.Literal }
func (statement *LetStatement) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(statement.TokenLiteral() + " ")
	buffer.WriteString(statement.Identifier.String())
	buffer.WriteString(" = ")
	if statement.Expression != nil {
		buffer.WriteString(statement.Expression.String())
	}
	buffer.WriteString(";")
	return buffer.String()
}

// Return statement
type ReturnStatement struct {
	Token      token.Token
	Expression Expression
}

func (statement *ReturnStatement) statementNode()       {}
func (statement *ReturnStatement) TokenLiteral() string { return statement.Token.Literal }
func (statement *ReturnStatement) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(statement.TokenLiteral() + " ")
	if statement.Expression != nil {
		buffer.WriteString(statement.Expression.String())
	}
	buffer.WriteString(";")
	return buffer.String()
}

// Expression statement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (statement *ExpressionStatement) statementNode()       {}
func (statement *ExpressionStatement) TokenLiteral() string { return statement.Token.Literal }
func (statement *ExpressionStatement) String() string {
	if statement.Expression != nil {
		return statement.Expression.String()
	}
	return ""
}
