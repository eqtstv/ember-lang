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

func TestIncrementExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let mut i = 0; i = i++; return i;`, 1},
		{`let mut i = 0; i = i++; i = i++; return i;`, 2},
		{`let mut i = 0; i = i++; i = i++; i = i++; return i;`, 3},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestWhileExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
			let mut i = 0;
			while (i < 10) {
				i = i++;
			}
			return i;
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

}

func TestForExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
			for (let i = 0; i < 10; i++) {
				let x = 0;
			}
			return x;
		`, 0},
		{`
			let mut x = 0;
			for (let i = 0; i < 10; i++) {
				let x = i;
			}
			return x;
		`, 9},
		{`
			for (let i = 0; i <= 10; i++) {
				let x = i;
			}
			return x;
		`, 10},
		// Nested loops
		{`
			let sum = 0;
			for (let i = 0; i < 3; i++) {
				for (let j = 0; j < 2; j++) {
					let sum = sum + 1;
				}
			}
			return sum;
		`, 6},
		// Empty loop (no iterations)
		{`
			let mut x = 5;
			for (let i = 0; i < 0; i++) {
				let x = 10;
			}
			return x;
		`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBuiltinTypeFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`type(1)`, "INTEGER"},
		{`type("hello")`, "STRING"},
		{`type(true)`, "BOOLEAN"},
		{`type([1, 2, 3])`, "ARRAY"},
		{`type({"a": 1, "b": 2})`, "HASH"},
		{`type(fn(x) { x + 1; })`, "FUNCTION"},
		{`type(1 + 2)`, "INTEGER"},
		{`type("hello" + " world")`, "STRING"},
		{`type([1, 2, 3] + [4, 5, 6])`, "ARRAY"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}

}

func TestAssignmentExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let mut x = 10; x = 5; return x;`, 5},
		{`let mut x = 10; x = 5; let mut x = 6; return x;`, 6},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// Add this test function to test mutability in assignments
func TestMutabilityInAssignments(t *testing.T) {
	tests := []struct {
		input           string
		expectedResult  any
		expectedError   bool
		expectedMessage string
	}{
		// Mutable variables can be reassigned
		{"let mut x = 5; x = 10; return x;", 10, false, ""},
		// // Immutable variables cannot be reassigned
		{"let x = 5; x = 10; return x;", nil, true, "(line 1) Cannot assign to immutable variable: x"},
		// // Nested scopes respect mutability
		{"let mut x = 5; if (true) { x = 20; }; return x;", 20, false, ""},
		{"let x = 5; if (true) { x = 20; }; return x;", nil, true, "(line 1) Cannot assign to immutable variable: x"},
		// // Function parameters are immutable by default
		{"let f = fn(x) { x = 20; return x; }; f(5);", nil, true, "(line 1) Cannot assign to immutable variable: x"},
		// // Loop variables can be mutable
		{"let mut sum = 0; for (let i = 0; i < 5; i++) { sum = sum + i; }; return sum;", 10, false, ""},
		// // Complex example with multiple variables
		{`
			let mut counter = 0;
			let immutable = 100;
			let mut result = 0;

			while (counter < 5) {
				result = result + counter;
				counter = counter + 1;
			}

			return result;
		`, 10, false, ""},
		// // Reassigning to a variable multiple times
		{"let mut x = 5; x = 10; x = 15; return x;", 15, false, ""},
		// // Assigning the result of an expression
		{"let mut x = 5; x = 2 + 3 * 4; return x;", 14, false, ""},
		// // Assigning the value of another variable
		{"let y = 10; let mut x = 5; x = y; return x;", 10, false, ""},
		// // Assigning in a nested block
		{"let mut x = 5; if (true) { if (true) { x = 20; } }; return x;", 20, false, ""},
		// // Shadowing variables
		// {"let mut x = 5; let y = 10; if (true) { let x = 20; }; return x;", 5, false, ""},
		// // Trying to assign to a shadowed variable
		// {"let x = 5; if (true) { let mut x = 10; x = 20; }; return x;", 5, false, ""},
		// Mutability in arrays
		{`
			let numbers = [1, 2, 3];
			numbers[0] = 10;
		`, nil, true, "(line 3) Cannot assign to immutable variable: numbers"},
		{`
			let mut numbers = [1, 2, 3];
			numbers[0] = 10;
			return numbers[0];
		`, 10, false, ""},
		{`
			let mapping = {"a": 1, "b": 2, "c": 3};
			mapping["a"] = 10;
			return mapping["a"];
		`, nil, true, "(line 3) Cannot assign to immutable variable: mapping"},
		{`
			let mut mapping = {"a": 1, "b": 2, "c": 3};
			mapping["a"] = 10;
			return mapping["a"];
		`, 10, false, ""},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if tt.expectedError {
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Expected error for input: %s, got %T (%+v)", tt.input, evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expectedMessage {
				t.Errorf("Wrong error message. expected=%q, got=%q",
					tt.expectedMessage, errObj.Message)
			}
		} else {
			if intObj, ok := tt.expectedResult.(int); ok {
				testIntegerObject(t, evaluated, int64(intObj))
			} else {
				t.Errorf("Unexpected result type: %T", tt.expectedResult)
			}
		}
	}
}

// Test environment mutability tracking
func TestEnvironmentMutabilityTracking(t *testing.T) {
	env := object.NewEnvironment()

	// Set immutable variable
	env.Set("x", &object.Integer{Value: 5}, false)

	// Set mutable variable
	env.Set("y", &object.Integer{Value: 10}, true)

	// Check mutability
	if env.IsMutable("x") {
		t.Errorf("Expected x to be immutable")
	}

	if !env.IsMutable("y") {
		t.Errorf("Expected y to be mutable")
	}

	// Test with enclosed environment
	enclosed := object.NewEnclosedEnvironment(env)

	// Check if outer variables are accessible and maintain mutability
	x, ok := enclosed.Get("x")
	if !ok {
		t.Errorf("Expected to find x in enclosed environment")
	}
	if x.(*object.Integer).Value != 5 {
		t.Errorf("Wrong value for x. got=%d, want=%d", x.(*object.Integer).Value, 5)
	}
	if enclosed.IsMutable("x") {
		t.Errorf("Expected x to be immutable in enclosed environment")
	}

	y, ok := enclosed.Get("y")
	if !ok {
		t.Errorf("Expected to find y in enclosed environment")
	}
	if y.(*object.Integer).Value != 10 {
		t.Errorf("Wrong value for y. got=%d, want=%d", y.(*object.Integer).Value, 10)
	}
	if !enclosed.IsMutable("y") {
		t.Errorf("Expected y to be mutable in enclosed environment")
	}

	// Set new variables in enclosed environment
	enclosed.Set("z", &object.Integer{Value: 15}, false)
	enclosed.Set("w", &object.Integer{Value: 20}, true)

	// Check new variables
	if enclosed.IsMutable("z") {
		t.Errorf("Expected z to be immutable")
	}
	if !enclosed.IsMutable("w") {
		t.Errorf("Expected w to be mutable")
	}

	// Outer environment shouldn't see enclosed variables
	_, ok = env.Get("z")
	if ok {
		t.Errorf("Outer environment shouldn't see z")
	}
}

// Test assignment in complex expressions
func TestAssignmentInComplexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Assignment in if condition
		{`
			let mut x = 5;
			if ((x = 10) > 5) {
				return x;
			}
			return 0;
		`, 10},

		// Assignment in loop condition
		{`
			let mut x = 0;
			while (x < 5) {
				x = x + 1;
			}
			return x;
		`, 5},

		// Assignment with function calls
		{`
			let mut x = 0;
			let f = fn() { return 10; };
			x = f();
			return x;
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalQuickSortAlgorithm(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{`
		let partition = fn(arr, low, high) {
			let mut arr = arr;
			let pivot = arr[high];
			
			let mut i = low - 1;
			
			for (let j = low; j < high; j++) {
				if (arr[j] < pivot) {
					i = i + 1;
					
					let temp = arr[i];
					arr[i] = arr[j];
					arr[j] = temp;
				}
			}
			
			let temp = arr[i + 1];
			arr[i + 1] = arr[high];
			arr[high] = temp;
			
			return i + 1;
		};

		let quicksort = fn(arr, low, high) {
			if (low < high) {
				let pi = partition(arr, low, high);
				
				quicksort(arr, low, pi - 1);
				quicksort(arr, pi + 1, high);
			}
			
			return arr;
		};

		let array = [10, 7, 8, 9, 1, 5, 3, 2, 6, 4];
		let sorted = quicksort(array, 0, len(array) - 1);

		return sorted;
		`, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated == nil {
			t.Errorf("Expected evaluated to be non-nil for input: %s", tt.input)
			continue
		}

		if array, ok := evaluated.(*object.Array); ok {
			if len(array.Elements) != len(tt.expected) {
				t.Errorf("Expected array length %d, got %d", len(tt.expected), len(array.Elements))
				continue
			}

			for i, element := range array.Elements {
				if element == nil {
					t.Errorf("Expected element %d to be non-nil", i)
					continue
				}

				if intObj, ok := element.(*object.Integer); ok {
					if intObj.Value != tt.expected[i] {
						t.Errorf("Expected element %d to be %d, got %d", i, tt.expected[i], intObj.Value)
					}
				} else {
					t.Errorf("Expected element %d to be an integer, got %T", i, element)
				}
			}
		} else {
			t.Errorf("Expected evaluated to be an array, got %T", evaluated)
		}
	}
}
