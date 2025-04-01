package evaluator

import (
	"ember_lang/ember_lang/ast"
	"ember_lang/ember_lang/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val, node.Name.Mutable)
		return val

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.IncrementExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalIncrementExpression(left)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)
	case *ast.PointerReferenceExpression:
		return evalPointerReferenceExpression(node, env)
	case *ast.PointerDereferenceExpression:
		return evalPointerDereferenceExpression(node, env)
	}

	return nil
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		if result != nil {
			resultType := result.Type()

			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("Not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], param.Mutable)
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}

	integer, ok := right.(*object.Integer)
	if !ok {
		return newError("Unknown operator: -%s", right.Type())
	}

	value := integer.Value
	return &object.Integer{Value: -value}
}

func evalPlusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator: +%s", right.Type())
	}

	integer, ok := right.(*object.Integer)
	if !ok {
		return newError("Unknown operator: +%s", right.Type())
	}

	value := integer.Value
	return &object.Integer{Value: value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("Index operator not supported: %s %s", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array object.Object, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	indexValue := index.(*object.Integer).Value

	if indexValue < 0 {
		indexValue = int64(len(arrayObject.Elements)) + indexValue
	}

	if indexValue < 0 || indexValue >= int64(len(arrayObject.Elements)) {
		return NULL
	}

	return arrayObject.Elements[indexValue]
}

func evalHashIndexExpression(hash object.Object, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("Unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func evalArrayInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftArray := left.(*object.Array)
	rightArray := right.(*object.Array)

	return &object.Array{Elements: append(leftArray.Elements, rightArray.Elements...)}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("Identifier not found: %s", node.Value)
}

func evalIncrementExpression(left object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	return &object.Integer{Value: leftVal + 1}
}

func evalWhileExpression(node *ast.WhileExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	for isTruthy(condition) {
		result := Eval(node.Body, env)
		if isError(result) {
			return result
		}

		condition = Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}
	}

	return NULL
}

func evalForExpression(node *ast.ForExpression, env *object.Environment) object.Object {
	letStatement := node.LetStatement

	// Set the initial value of the loop variable
	env.Set(letStatement.Name.Value, Eval(letStatement.Value, env), true)

	for {
		condition := Eval(node.Condition, env)
		if !isTruthy(condition) {
			break
		}

		result := Eval(node.Body, env)
		if isError(result) {
			return result
		}

		increment := Eval(node.Increment, env)
		if isError(increment) {
			return increment
		}

		env.Set(letStatement.Name.Value, increment, true)
	}

	return NULL
}

func evalAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) object.Object {
	// Assignment to a variable
	if identifier, ok := node.Left.(*ast.Identifier); ok {
		// Check if the variable is mutable
		if !env.IsMutable(identifier.Value) {
			return newError("(line %d) Cannot assign to immutable variable: %s", identifier.Token.LineNumber, identifier.Value)
		}

		// Evaluate the right-hand side of the assignment
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		// Set the value of the variable
		env.Set(identifier.Value, right, true)

		// Return the value of the assignment
		return right
	}

	// Assignment to a dereferenced pointer
	if pointerDeref, ok := node.Left.(*ast.PointerDereferenceExpression); ok {
		// Evaluate the pointer
		pointerObj := Eval(pointerDeref.Right, env)
		if isError(pointerObj) {
			return pointerObj
		}

		// Check if it's a pointer
		pointer, ok := pointerObj.(*object.Pointer)
		if !ok {
			return newError("(line %d) Cannot dereference non-pointer value: %s", node.Token.LineNumber, pointerObj.Type())
		}

		// Check if the variable is mutable
		if !env.IsMutable(pointer.Name) {
			return newError("(line %d) Cannot assign to immutable variable: %s", node.Token.LineNumber, pointer.Name)
		}

		// Evaluate the right-hand side of the assignment
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		// Set the value of the variable
		env.Set(pointer.Name, right, true)

		// Update the pointer's value
		pointer.Value = right

		return right
	}

	// Assignment to an Index Expression
	// eg. arr[0] = 1, map[key] = value
	if indexExpression, ok := node.Left.(*ast.IndexExpression); ok {
		// Evaluate the left-hand side of the assignment
		left := Eval(indexExpression.Left, env)
		if isError(left) {
			return left
		}

		// Check if the variable is mutable
		if identifier, ok := indexExpression.Left.(*ast.Identifier); ok {
			if !env.IsMutable(identifier.Value) {
				return newError("(line %d) Cannot assign to immutable variable: %s", identifier.Token.LineNumber, identifier.Value)
			}
		} else {
			// Handle nested index expressions or other complex left sides
			// This allows for cases like: arr[i][j] = value
			return newError("(line %d) Complex index expressions not yet supported for assignment", node.Token.LineNumber)
		}

		// Evaluate the index of the assignment
		index := Eval(indexExpression.Index, env)
		if isError(index) {
			return index
		}

		// Evaluate the right-hand side of the assignment
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		// Assign the value to the index of the array or map
		switch left := left.(type) {
		// Array Assignment
		case *object.Array:
			indexValue, ok := index.(*object.Integer)
			if !ok {
				return newError("(line %d) Array index must be an integer", node.Token.LineNumber)
			}

			idx := indexValue.Value
			// Support negative indices (like Python)
			if idx < 0 {
				idx = int64(len(left.Elements)) + idx
			}

			// Check bounds
			if idx < 0 || idx >= int64(len(left.Elements)) {
				return newError("(line %d) Array index out of bounds: %d", node.Token.LineNumber, idx)
			}

			left.Elements[idx] = right
			return right

		// Map Assignment
		case *object.Hash:
			key, ok := index.(object.Hashable)
			if !ok {
				return newError("(line %d) Unusable as hash key: %s", node.Token.LineNumber, index.Type())
			}
			left.Pairs[key.HashKey()] = object.HashPair{Key: index, Value: right}
			return right

		default:
			return newError("(line %d) Cannot index into type: %s", node.Token.LineNumber, left.Type())
		}
	}

	return newError("(line %d) invalid assignment target", node.Token.LineNumber)
}

func evalPointerReferenceExpression(node *ast.PointerReferenceExpression, env *object.Environment) object.Object {
	// We can only take the address of identifiers
	if identifier, ok := node.Right.(*ast.Identifier); ok {
		// Check if the variable exists
		val, exists := env.Get(identifier.Value)
		if !exists {
			return newError("Cannot take address of undefined variable: %s", identifier.Value)
		}

		// Create a pointer to the variable
		return &object.Pointer{
			Name:  identifier.Value,
			Value: val,
		}
	}

	return newError("Cannot take address of non-identifier expression")
}

func evalPointerDereferenceExpression(node *ast.PointerDereferenceExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	if pointer, ok := right.(*object.Pointer); ok {
		// Get the current value from the environment
		val, exists := env.Get(pointer.Name)
		if !exists {
			return newError("Pointer references undefined variable: %s", pointer.Name)
		}
		return val
	}

	return newError("Cannot dereference non-pointer value: %s", right.Type())
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
