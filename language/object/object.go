package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
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

// Integers
type Integer struct {
	Value int64
}

func (integer *Integer) GetType() ObjectType { return INTEGER_OBJECT }
func (integer *Integer) Inspect() string     { return fmt.Sprintf("%d", integer.Value) }

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
