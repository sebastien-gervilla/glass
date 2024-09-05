package object

import (
	"bytes"
	"fmt"
	"glass/language/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJECT      = "INTEGER"
	STRING_OBJECT       = "STRING"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	BUILTIN_OBJECT      = "BUILTIN"
)

type Object interface {
	GetType() ObjectType
	Inspect() string
}

// Error
type Error struct {
	Message string
}

func (e *Error) GetType() ObjectType { return ERROR_OBJECT }
func (e *Error) Inspect() string     { return "ERROR: " + e.Message }

// Environment
type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	environment := NewEnvironment()
	environment.outer = outer
	return environment
}

func (environment *Environment) Get(name string) (Object, bool) {
	obj, ok := environment.store[name]

	// Reach for outer variables
	if !ok && environment.outer != nil {
		obj, ok = environment.outer.Get(name)
	}

	return obj, ok
}

func (environment *Environment) Set(name string, val Object) Object {
	environment.store[name] = val
	return val
}

// Integer
type Integer struct {
	Value int64
}

func (integer *Integer) GetType() ObjectType { return INTEGER_OBJECT }
func (integer *Integer) Inspect() string     { return fmt.Sprintf("%d", integer.Value) }

// String
type String struct {
	Value string
}

func (str *String) GetType() ObjectType { return STRING_OBJECT }
func (str *String) Inspect() string     { return str.Value }

// Boolean
type Boolean struct {
	Value bool
}

func (boolean *Boolean) GetType() ObjectType { return BOOLEAN_OBJECT }
func (boolean *Boolean) Inspect() string     { return fmt.Sprintf("%t", boolean.Value) }

// Null
type Null struct{}

func (null *Null) GetType() ObjectType { return NULL_OBJECT }
func (null *Null) Inspect() string     { return "null" }

// Return
type ReturnValue struct {
	Value Object
}

func (value *ReturnValue) GetType() ObjectType { return RETURN_VALUE_OBJECT }
func (value *ReturnValue) Inspect() string     { return value.Value.Inspect() }

// Functions
type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (function *Function) GetType() ObjectType { return FUNCTION_OBJECT }
func (function *Function) Inspect() string {
	var buffer bytes.Buffer
	params := []string{}
	for _, p := range function.Parameters {
		params = append(params, p.String())
	}
	buffer.WriteString("function")
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(params, ", "))
	buffer.WriteString(") {\n")
	buffer.WriteString(function.Body.String())
	buffer.WriteString("\n}")
	return buffer.String()
}

type BuiltinFunction func(args ...Object) Object

// Builtins
type Builtin struct {
	Function BuiltinFunction
}

func (builtin *Builtin) GetType() ObjectType { return BUILTIN_OBJECT }
func (builtin *Builtin) Inspect() string     { return "builtin function" }
