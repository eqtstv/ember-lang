package evaluator

import (
	"ember_lang/ember_lang/object"
	"fmt"
	"math/rand"
)

var builtins map[string]*object.Builtin

func init() {
	builtins = map[string]*object.Builtin{
		"print": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL
			},
		},
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
		"concat": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				array1, ok := args[0].(*object.Array)
				if !ok {
					return newError("Invalid argument to concat. Got: %s, Expected: ARRAY", args[0].Type())
				}

				array2, ok := args[1].(*object.Array)
				if !ok {
					return newError("Invalid argument to concat. Got: %s, Expected: ARRAY", args[1].Type())
				}

				return &object.Array{Elements: append(array1.Elements, array2.Elements...)}
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

				var function object.Object

				if args[1].Type() == object.BUILTIN_OBJ {
					function, ok = args[1].(*object.Builtin)
					if !ok {
						return newError("Invalid argument to map. Got: %s, Expected: FUNCTION", args[1].Type())
					}
				} else {
					function, ok = args[1].(*object.Function)
					if !ok {
						return newError("Invalid argument to map. Got: %s, Expected: FUNCTION", args[1].Type())
					}
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

				var function object.Object

				if args[1].Type() == object.BUILTIN_OBJ {
					function, ok = args[1].(*object.Builtin)
					if !ok {
						return newError("Invalid argument to reduce. Got: %s, Expected: FUNCTION", args[1].Type())
					}
				} else {
					function, ok = args[1].(*object.Function)
					if !ok {
						return newError("Invalid argument to reduce. Got: %s, Expected: FUNCTION", args[1].Type())
					}
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
		"add": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				left, ok := args[0].(*object.Integer)
				if !ok {
					return newError("Invalid argument to add. Got: %s, Expected: INTEGER", args[0].Type())
				}

				right, ok := args[1].(*object.Integer)
				if !ok {
					return newError("Invalid argument to add. Got: %s, Expected: INTEGER", args[1].Type())
				}

				return &object.Integer{Value: left.Value + right.Value}
			},
		},
		"sub": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				left, ok := args[0].(*object.Integer)
				if !ok {
					return newError("Invalid argument to sub. Got: %s, Expected: INTEGER", args[0].Type())
				}

				right, ok := args[1].(*object.Integer)
				if !ok {
					return newError("Invalid argument to sub. Got: %s, Expected: INTEGER", args[1].Type())
				}

				return &object.Integer{Value: left.Value - right.Value}
			},
		},
		"mul": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				left, ok := args[0].(*object.Integer)
				if !ok {
					return newError("Invalid argument to mul. Got: %s, Expected: INTEGER", args[0].Type())
				}

				right, ok := args[1].(*object.Integer)
				if !ok {
					return newError("Invalid argument to mul. Got: %s, Expected: INTEGER", args[1].Type())
				}

				return &object.Integer{Value: left.Value * right.Value}
			},
		},
		"div": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Invalid number of arguments. Got: %d, Expected: 2", len(args))
				}

				left, ok := args[0].(*object.Integer)
				if !ok {
					return newError("Invalid argument to div. Got: %s, Expected: INTEGER", args[0].Type())
				}

				right, ok := args[1].(*object.Integer)
				if !ok {
					return newError("Invalid argument to div. Got: %s, Expected: INTEGER", args[1].Type())
				}

				return &object.Integer{Value: left.Value / right.Value}
			},
		},
		"type": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Invalid number of arguments. Got: %d, Expected: 1", len(args))
				}

				switch args[0].(type) {
				case *object.Integer:
					return &object.String{Value: "INTEGER"}
				case *object.String:
					return &object.String{Value: "STRING"}
				case *object.Boolean:
					return &object.String{Value: "BOOLEAN"}
				case *object.Array:
					return &object.String{Value: "ARRAY"}
				case *object.Hash:
					return &object.String{Value: "HASH"}
				case *object.Function:
					return &object.String{Value: "FUNCTION"}
				default:
					return newError("Invalid argument to type. Got: %s", args[0].Type())
				}
			},
		},
		"rand": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) > 1 {
					return newError("Invalid number of arguments. Got: %d, Expected: 0 or 1", len(args))
				}

				// If no arguments, return a random int between 0 and MaxInt32
				if len(args) == 0 {
					return &object.Integer{Value: int64(rand.Int31())}
				}

				// If one argument, it should be the max value (exclusive)
				max, ok := args[0].(*object.Integer)
				if !ok {
					return newError("Invalid argument to rand. Got: %s, Expected: INTEGER", args[0].Type())
				}

				if max.Value <= 0 {
					return newError("Argument to rand must be positive, got: %d", max.Value)
				}

				return &object.Integer{Value: int64(rand.Int63n(max.Value))}
			},
		},
	}
}
