package evaluator

import (
	"ember_lang/ember_lang/object"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(10 <= 10) == true", true},
		{"(10 >= 10) == true", true},
		{"(10 < 10) == true", false},
		{"(10 > 10) == true", false},
		{"(10 == 10) == true", true},
		{"(10 != 10) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 <= 1) { 10 } else { 20 }", 10},
		{"if (1 >= 1) { 10 } else { 20 }", 10},
		{"if (1 <= 2) { 10 } else { 20 }", 10},
		{"if (1 >= 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
			return 1;
			}
		`,
			10},
		{`
			if (10 <= 10) {
				return 10;
			}
			return 20;
		`,
			10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
		{
			`"Hello" - "World!"`,
			"Unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"Unusable as hash key: FUNCTION",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinLenFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// String
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		// Array
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func TestBuiltinPushFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{`push([], 1)`, []int64{1}},
		{`push([1, 2], 3)`, []int64{1, 2, 3}},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}
		if len(result.Elements) != len(tt.expected) {
			t.Fatalf("array has wrong num of elements. got=%d",
				len(result.Elements))
		}
		for i, expected := range tt.expected {
			testIntegerObject(t, result.Elements[i], expected)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
		{
			"[1, 2, 3][-2]",
			2,
		},
		{
			"[1, 2, 3][-8]",
			nil,
		},
		{
			`["one", "two", "three"][1]`,
			"two",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
			}
			if str.Value != expected {
				t.Errorf("String has wrong value. got=%q", str.Value)
			}
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestMapFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`map([1, 2, 3], fn(x) { x * 2; })`, []interface{}{2, 4, 6}},
		{`map([1, 2, 3, 4, 5], fn(x) { x * x; })`, []interface{}{1, 4, 9, 16, 25}},
		{`map(["one", "two", "three", "four", "five"], len)`, []interface{}{3, 3, 5, 4, 4}},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}
		if len(result.Elements) != len(tt.expected) {
			t.Fatalf("array has wrong num of elements. Got=%d, Expected=%d",
				len(result.Elements), len(tt.expected))
		}

		for i, expected := range tt.expected {
			switch expected := expected.(type) {
			case int:
				testIntegerObject(t, result.Elements[i], int64(expected))
			case string:
				str, ok := result.Elements[i].(*object.String)
				if !ok {
					t.Fatalf("object is not String. got=%T (%+v)", result.Elements[i], result.Elements[i])
				}
				if str.Value != expected {
					t.Errorf("String has wrong value. got=%q", str.Value)
				}
			}
		}
	}
}

func TestReduceFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`reduce([1, 2, 3], fn(acc, x) { acc + x; }, 0)`, 6},
		{`reduce([1, 2, 3], fn(acc, x) { acc + x; }, 1)`, 7},
		{`reduce([1, 2, 3], fn(acc, x) { acc * x; }, 1)`, 6},
		{`reduce([1, 2, 3], fn(acc, x) { acc * x; }, 0)`, 0},
		{`reduce([1, 2, 3], fn(acc, x) { acc * x; }, 1)`, 6},
		{`reduce([1, 2, 3], add, 0)`, 6},
		{`reduce([1, 2, 3], sub, 0)`, -6},
		{`reduce([1, 2, 3], mul, 1)`, 6},
		{`reduce([1, 2, 3], div, 1)`, 0},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBuiltinStackingFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`reduce(map(["one", "two", "three"], len), add, 0)`, 11},
		{`reduce(map([1, 2, 3], fn(x) { x * x; }), add, 0)`, 14},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestOrEqualOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 <= 5", true},
		{"4 <= 5", true},
		{"6 <= 5", false},
		{"10 <= 10", true},
		{"10 >= 9", true},
		{"10 >= 10", true},
		{"10 >= 11", false},
		{"10 <= 9", false},
		{"10 <= 10", true},
		{"10 <= 11", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestFibo(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
			let fib = fn(n) {
				if (n <= 1) {
					return n;
				}
				return fib(n - 1) + fib(n - 2);
			};
			fib(10);
		`, 55},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalArrayInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{`[1, 2, 3] + [4, 5, 6]`, []int64{1, 2, 3, 4, 5, 6}},
		{`[1, 2, 3] + []`, []int64{1, 2, 3}},
		{`[] + [1, 2, 3]`, []int64{1, 2, 3}},
		{`[] + []`, []int64{}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}
		if len(result.Elements) != len(tt.expected) {
			t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
		}

		for i, expected := range tt.expected {
			testIntegerObject(t, result.Elements[i], expected)
		}
	}
}

func TestEvalBuiltinConcatFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{`concat([1, 2, 3], [4, 5, 6])`, []int64{1, 2, 3, 4, 5, 6}},
		{`concat([1, 2, 3], [])`, []int64{1, 2, 3}},
		{`concat([], [1, 2, 3])`, []int64{1, 2, 3}},
		{`concat([], [])`, []int64{}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}

		for i, expected := range tt.expected {
			testIntegerObject(t, result.Elements[i], expected)
		}
	}

}
