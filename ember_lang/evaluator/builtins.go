package evaluator

import (
	"ember_lang/ember_lang/object"
)

var builtins map[string]*object.Builtin

func init() {
	builtins = map[string]*object.Builtin{
		"len": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Invalid number of arguments. Got: %d, Expected: 1", len(args))
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("Invalid argument to len. Got: %s, Expected: STRING or ARRAY", args[0].Type())
				}
			},
		},
		"push": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				array, ok := args[0].(*object.Array)
				if !ok {
					return newError("Invalid argument to push. Got: %s, Expected: ARRAY", args[0].Type())
				}

				newArray := &object.Array{Elements: append(array.Elements, args[1])}
				return newArray
			},
		},
		"map": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				array, ok := args[0].(*object.Array)
				if !ok {
					return newError("Invalid argument to map. Got: %s, Expected: ARRAY", args[0].Type())
				}

				function, ok := args[1].(*object.Function)
				if !ok {
					return newError("Invalid argument to map. Got: %s, Expected: FUNCTION", args[1].Type())
				}

				newArray := &object.Array{Elements: make([]object.Object, 0, len(array.Elements))}

				for _, elem := range array.Elements {
					result := applyFunction(function, []object.Object{elem})
					if isError(result) {
						return result
					}
					newArray.Elements = append(newArray.Elements, result)
				}

				return newArray
			},
		},
		"reduce": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("Invalid number of arguments. Got: %d, Expected: 3", len(args))
				}

				array, ok := args[0].(*object.Array)
				if !ok {
					return newError("Invalid argument to reduce. Got: %s, Expected: ARRAY", args[0].Type())
				}

				function, ok := args[1].(*object.Function)
				if !ok {
					return newError("Invalid argument to reduce. Got: %s, Expected: FUNCTION", args[1].Type())
				}

				initialValue, ok := args[2].(*object.Integer)
				if !ok {
					return newError("Invalid argument to reduce. Got: %s, Expected: INTEGER", args[1].Type())
				}

				for index, elem := range array.Elements {
					if elem.Type() != object.INTEGER_OBJ {
						return newError("Invalid argument to reduce at index %d. Got: %s, Expected: INTEGER", index, elem.Type())
					}

					result := applyFunction(function, []object.Object{initialValue, elem})
					if isError(result) {
						return result
					}

					resultInt, ok := result.(*object.Integer)
					if !ok {
						return newError("Reduce function must return INTEGER, got: %s", result.Type())
					}
					initialValue = resultInt
				}

				return initialValue
			},
		},
	}
}
