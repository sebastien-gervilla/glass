package evaluator

import (
	"fmt"
	"glass/language/ast"
	"glass/language/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Evaluate(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evaluateProgram(node, environment)

	case *ast.LetStatement:
		value := Evaluate(node.Expression, environment)
		if isError(value) {
			return value
		}

		environment.Set(node.Identifier.Value, value)

	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, environment)

	case *ast.ReturnStatement:
		value := Evaluate(node.Expression, environment)
		if isError(value) {
			return value
		}

		return &object.ReturnValue{
			Value: Evaluate(node.Expression, environment),
		}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.StringLiteral:
		return &object.String{
			Value: node.Value,
		}

	case *ast.Boolean:
		return newBooleanObject(node.Value)

	case *ast.Identifier:
		return evaluateIdentifier(node, environment)

	case *ast.ArrayLiteral:
		elements := evaluateExpressions(node.Elements, environment)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.PrefixExpression:
		right := Evaluate(node.Expression, environment)
		if isError(right) {
			return right
		}
		return evaluatePrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Evaluate(node.LeftExpression, environment)
		if isError(left) {
			return left
		}

		right := Evaluate(node.RightExpression, environment)
		if isError(right) {
			return right
		}

		return evaluateInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evaluateBlockStatement(node, environment)

	case *ast.IfExpression:
		return evaluateIfExpression(node, environment)

	case *ast.IndexExpression:
		left := Evaluate(node.Left, environment)
		if isError(left) {
			return left
		}

		index := Evaluate(node.Index, environment)
		if isError(index) {
			return index
		}

		return evaluateIndexExpression(left, index)

	case *ast.Function:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters:  params,
			Environment: environment,
			Body:        body,
		}

	case *ast.CallExpression:
		function := Evaluate(node.Function, environment)
		if isError(function) {
			return function
		}

		arguments := evaluateExpressions(node.Arguments, environment)
		if len(arguments) == 1 && isError(arguments[0]) {
			return arguments[0]
		}

		return applyFunction(function, arguments)

	}

	return nil
}

func evaluateProgram(program *ast.Program, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Evaluate(statement, environment)

		switch result := result.(type) {

		case *object.ReturnValue:
			return result.Value

		case *object.Error:
			return result

		}
	}

	return result
}

func evaluateBlockStatement(blockStatement *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range blockStatement.Statements {
		result = Evaluate(statement, environment)

		if result != nil {
			resultType := result.GetType()
			if resultType == object.RETURN_VALUE_OBJECT || resultType == object.ERROR_OBJECT {
				return result
			}
		}
	}

	return result
}

func evaluateIdentifier(identifier *ast.Identifier, environment *object.Environment) object.Object {
	value, ok := environment.Get(identifier.Value)
	if ok {
		return value
	}

	if builtin, ok := builtins[identifier.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + identifier.Value)
}

func evaluateExpressions(expressions []ast.Expression, environment *object.Environment) []object.Object {
	var result []object.Object
	for _, expression := range expressions {
		evaluated := Evaluate(expression, environment)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {

	case "!":
		return evaluateNotOperatorExpression(right)

	case "-":
		return evaluateMinusPrefixExpression(right)

	default:
		return newError("unknown operator: %s%s", operator, right.GetType())

	}
}

func evaluateNotOperatorExpression(right object.Object) object.Object {
	switch right {

	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE

	}
}

func evaluateMinusPrefixExpression(right object.Object) object.Object {
	if right.GetType() != object.INTEGER_OBJECT {
		return newError("unknown operator: -%s", right.GetType())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evaluateInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {

	case left.GetType() == object.INTEGER_OBJECT && right.GetType() == object.INTEGER_OBJECT:
		return evaluateIntegerInfixExpression(operator, left, right)

	case left.GetType() == object.STRING_OBJECT && right.GetType() == object.STRING_OBJECT:
		return evaluateStringInfixExpression(operator, left, right)

	case operator == "==":
		return newBooleanObject(left == right)

	case operator == "!=":
		return newBooleanObject(left != right)

	case left.GetType() != right.GetType():
		return newError("type mismatch: %s %s %s", left.GetType(), operator, right.GetType())

	default:
		return newError("unknown operator: %s %s %s", left.GetType(), operator, right.GetType())

	}
}

func evaluateIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {

	case "+":
		return &object.Integer{Value: leftValue + rightValue}

	case "-":
		return &object.Integer{Value: leftValue - rightValue}

	case "*":
		return &object.Integer{Value: leftValue * rightValue}

	case "/":
		return &object.Integer{Value: leftValue / rightValue}

	case "<":
		return newBooleanObject(leftValue < rightValue)

	case ">":
		return newBooleanObject(leftValue > rightValue)

	case "==":
		return newBooleanObject(leftValue == rightValue)

	case "!=":
		return newBooleanObject(leftValue != rightValue)

	default:
		return newError("unknown operator: %s %s %s", left.GetType(), operator, right.GetType())

	}
}

func evaluateStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.GetType(), operator, right.GetType())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evaluateIfExpression(expression *ast.IfExpression, environment *object.Environment) object.Object {
	condition := Evaluate(expression.Condition, environment)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Evaluate(expression.Consequence, environment)
	}

	if expression.Alternative != nil {
		return Evaluate(expression.Alternative, environment)
	}

	return NULL
}

func applyFunction(fn object.Object, arguments []object.Object) object.Object {
	switch function := fn.(type) {

	case *object.Function:
		extendedEnvironment := extendFunctionEnvironment(function, arguments)
		evaluated := Evaluate(function.Body, extendedEnvironment)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return function.Function(arguments...)

	default:
		return newError("not a function: %s", fn.GetType())

	}
}

func extendFunctionEnvironment(function *object.Function, arguments []object.Object) *object.Environment {
	environment := object.NewEnclosedEnvironment(function.Environment)
	for index, param := range function.Parameters {
		environment.Set(param.Value, arguments[index])
	}
	return environment
}
func unwrapReturnValue(obj object.Object) object.Object {
	returnValue, ok := obj.(*object.ReturnValue)
	if ok {
		return returnValue.Value
	}

	return obj
}

func evaluateIndexExpression(left object.Object, index object.Object) object.Object {
	switch {

	case left.GetType() == object.ARRAY_OBJECT && index.GetType() == object.INTEGER_OBJECT:
		return evaluateArrayIndexExpression(left, index)

	default:
		return newError("index operator not supported: %s", left.GetType())

	}
}

func evaluateArrayIndexExpression(array object.Object, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	arrayIndex := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if arrayIndex < 0 || arrayIndex > max {
		return NULL
	}

	return arrayObject.Elements[arrayIndex]
}

// Utils

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func newBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}

	return FALSE
}

func isTruthy(object object.Object) bool {
	switch object {

	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default:
		return true

	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.GetType() == object.ERROR_OBJECT
	}

	return false
}
