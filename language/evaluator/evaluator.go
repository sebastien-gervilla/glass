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
		return evaluateProgram(node)

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

	case *ast.Boolean:
		return newBooleanObject(node.Value)

	case *ast.Identifier:
		return evaluateIdentifier(node, environment)

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
		return evaluateBlockStatement(node)

	case *ast.IfExpression:
		return evaluateIfExpression(node)

	}

	return nil
}

func evaluateProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Evaluate(statement)

		switch result := result.(type) {

		case *object.ReturnValue:
			return result.Value

		case *object.Error:
			return result

		}
	}

	return result
}

func evaluateBlockStatement(blockStatement *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range blockStatement.Statements {
		result = Evaluate(statement)

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
	if !ok {
		return newError("identifier not found: " + identifier.Value)
	}

	return value
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

func evaluateIfExpression(expression *ast.IfExpression) object.Object {
	condition := Evaluate(expression.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Evaluate(expression.Consequence)
	}

	if expression.Alternative != nil {
		return Evaluate(expression.Alternative)
	}

	return NULL
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
