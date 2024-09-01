package ast

import (
	token "glass/language/token"
)

// Interfaces
type Node interface {
	TokenLiteral() string
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

// Identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }

// Let statement
type LetStatement struct {
	Token      token.Token
	Identifier *Identifier
	Expression *Expression
}

func (statement *LetStatement) statementNode()       {}
func (statement *LetStatement) TokenLiteral() string { return statement.Token.Literal }

// Return statement
type ReturnStatement struct {
	Token      token.Token
	Expression *Expression
}

func (statement *ReturnStatement) statementNode()       {}
func (statement *ReturnStatement) TokenLiteral() string { return statement.Token.Literal }
