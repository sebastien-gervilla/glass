package evaluator

import (
	"glass/language/ast"
	"glass/language/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evaluateProgram(node)

	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)

	case *ast.ReturnStatement:
		return &object.ReturnValue{
			Value: Evaluate(node.Expression),
		}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.Boolean:
		return newBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Evaluate(node.Expression)
		return evaluatePrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Evaluate(node.LeftExpression)
		right := Evaluate(node.RightExpression)
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

		returnValue, ok := result.(*object.ReturnValue)
		if ok {
			return returnValue.Value
		}
	}

	return result
}

func evaluateBlockStatement(blockStatement *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range blockStatement.Statements {
		result = Evaluate(statement)

		if result != nil && result.GetType() == object.RETURN_VALUE_OBJECT {
			return result
		}
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
		return NULL

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
		return NULL
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

	default:
		return NULL

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
		return NULL

	}
}

func evaluateIfExpression(expression *ast.IfExpression) object.Object {
	condition := Evaluate(expression.Condition)

	if isTruthy(condition) {
		return Evaluate(expression.Consequence)
	}

	if expression.Alternative != nil {
		return Evaluate(expression.Alternative)
	}

	return NULL
}

// Utils

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
