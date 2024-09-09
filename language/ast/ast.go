package ast

import (
	"bytes"
	token "glass/language/token"
	"strings"
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

// Block statement
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (statement *BlockStatement) statementNode()       {}
func (statement *BlockStatement) TokenLiteral() string { return statement.Token.Literal }
func (statement *BlockStatement) String() string {
	var buffer bytes.Buffer
	for _, statement := range statement.Statements {
		buffer.WriteString(statement.String())
	}
	return buffer.String()
}

// Integer literal
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (integer *IntegerLiteral) expressionNode()      {}
func (integer *IntegerLiteral) TokenLiteral() string { return integer.Token.Literal }
func (integer *IntegerLiteral) String() string       { return integer.Token.Literal }

// String literal
type StringLiteral struct {
	Token token.Token
	Value string
}

func (literal *StringLiteral) expressionNode()      {}
func (literal *StringLiteral) TokenLiteral() string { return literal.Token.Literal }
func (literal *StringLiteral) String() string       { return literal.Token.Literal }

// Prefix operator
type PrefixExpression struct {
	Token      token.Token
	Operator   string
	Expression Expression
}

func (expression *PrefixExpression) expressionNode()      {}
func (expression *PrefixExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *PrefixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(expression.Operator)
	buffer.WriteString(expression.Expression.String())
	buffer.WriteString(")")

	return buffer.String()
}

// Infix expression
type InfixExpression struct {
	Token           token.Token
	LeftExpression  Expression
	Operator        string
	RightExpression Expression
}

func (expression *InfixExpression) expressionNode()      {}
func (expression *InfixExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *InfixExpression) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("(")
	buffer.WriteString(expression.LeftExpression.String())
	buffer.WriteString(" " + expression.Operator + " ")
	buffer.WriteString(expression.RightExpression.String())
	buffer.WriteString(")")

	return buffer.String()
}

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (boolean *Boolean) expressionNode()      {}
func (boolean *Boolean) TokenLiteral() string { return boolean.Token.Literal }
func (boolean *Boolean) String() string       { return boolean.Token.Literal }

// If expression
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (statement *IfExpression) expressionNode()      {}
func (statement *IfExpression) TokenLiteral() string { return statement.Token.Literal }
func (statement *IfExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("if")
	buffer.WriteString(statement.Condition.String())
	buffer.WriteString(" ")
	buffer.WriteString(statement.Consequence.String())
	if statement.Alternative != nil {
		buffer.WriteString("else ")
		buffer.WriteString(statement.Alternative.String())
	}
	return buffer.String()
}

// Function
type Function struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (function *Function) expressionNode()      {}
func (function *Function) TokenLiteral() string { return function.Token.Literal }
func (function *Function) String() string {
	var buffer bytes.Buffer
	parameters := []string{}
	for _, p := range function.Parameters {
		parameters = append(parameters, p.String())
	}

	buffer.WriteString(function.TokenLiteral())
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(parameters, ", "))
	buffer.WriteString(") ")
	buffer.WriteString(function.Body.String())
	return buffer.String()
}

// Call expression
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (expression *CallExpression) expressionNode()      {}
func (expression *CallExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *CallExpression) String() string {
	var buffer bytes.Buffer
	args := []string{}
	for _, a := range expression.Arguments {
		args = append(args, a.String())
	}

	buffer.WriteString(expression.Function.String())
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(args, ", "))
	buffer.WriteString(")")
	return buffer.String()
}

// Array
type ArrayLiteral struct {
	Token    token.Token // LBRACKET '[' token
	Elements []Expression
}

func (array *ArrayLiteral) expressionNode()      {}
func (array *ArrayLiteral) TokenLiteral() string { return array.Token.Literal }
func (array *ArrayLiteral) String() string {
	var buffer bytes.Buffer
	elements := []string{}
	for _, element := range array.Elements {
		elements = append(elements, element.String())
	}

	buffer.WriteString("[")
	buffer.WriteString(strings.Join(elements, ", "))
	buffer.WriteString("]")
	return buffer.String()
}

// Index
type IndexExpression struct {
	Token token.Token // LBRACKET '[' token
	Left  Expression
	Index Expression
}

func (expression *IndexExpression) expressionNode()      {}
func (expression *IndexExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *IndexExpression) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	buffer.WriteString(expression.Left.String())
	buffer.WriteString("[")
	buffer.WriteString(expression.Index.String())
	buffer.WriteString("])")
	return buffer.String()
}
