package evaluator

import "ember_lang/ember_lang/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Invalid number of arguments. Got: %d, Expected: 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("Invalid argument to len. Got: %s, Expected: STRING", args[0].Type())
			}
		},
	},
}
