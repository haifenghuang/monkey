package eval

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
	"os"
	"testing"
)

func TestDoLoop(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"let a = 0; do { if(a == 10) { break } a = a + 1 }; a", 10},
		{"let a = 0; let b = 0; do { if(a == 10) { break } a = a + 1 do { if(b == 3) { break } b = b + 1 } }; a + b", 13},
	}

	for _, tt := range test {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestWhileLoop(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"let a = 0; while (a < 10) { a = a + 1 }; a", 10},
		{"let a = 0; while (a < 10) { a = a + 1 if (a == 5) { break } }; a", 5},
	}

	for _, tt := range test {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestReassignment(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a = 4;a", 4},
		{"let a = 5 * 5; a = 5;a", 5},
		{"let a = 5; let b = a * 5; a = b;a", 25},
		{"let a = 5; let b = a; let c = a + b + 5; c; b = c;b", 15},
	}

	for _, tt := range test {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFileObjects(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let f = open("../parser/test_files/module.my");str(f)`, "<file object: ../parser/test_files/module.my>"},
		{`let f = open("../parser/test_files/module.my");f.read()`, `include eval
include test
include sub_package`},
		{`let f = open("../parser/test_files/module.my");f.readline()`, "include eval"},
		{`let f = open("../parser/test_files/module.my");f.readline();f.readline()`, "include test"},
		{`let f = open("../parser/test_files/module.my");f.readline();f.readline();f.readline()`, "include sub_package"},
	}
	d, _ := os.Getwd()
	fmt.Println(d)
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestChainedCalled(t *testing.T) {
	input := `[1,2,3].map(fn(x) { x + 1 }).map(fn(x) { x * 5 }).filter(fn(x) { x > 10 }).pop()`
	testEval(input)
}

func TestStructObjects(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`struct (a=>15).a`, 15},
		{`let st = struct {a=>15}; type(addm(st, "get", fn() { self.a })) == "NIL"`, true},
		{`let st = struct {a=>15}; addm(st, "get", fn() { self.a }); st.get()`, 15},
		{`let st = struct {a=>15}; addm(st, "get", fn() { a }); type(st.get()) == "ERROR"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		default:
			t.Errorf("evaluted not %T. got=%T", evaluated, expected)
		}
	}
}

//func TestIncludeObjects(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected int64
//	}{
//		{"include test_files;eval.a;", 5},
//		{"include test_files;eval.b;", 10},
//		{"include test_files;eval.c(5);", 10},
//		{"include test_files;eval.d(4,4);", 8},
//		{"include test_files;test.d;", 25},
//		{"include test_files;pkg.testfn(10,25);", 35},
//	}
//
//	for _, tt := range tests {
//		l := lexer.New(tt.input)
//		path, _ := os.Getwd()
//		path = path + "/../parser"
//		p := parser.New(l, path)
//		s := NewScope(nil)
//		program := p.ParseProgram()
//		if len(program.Includes) == 0 {
//			t.Errorf("Parsed program has no statements or included objects.\n")
//			os.Exit(1)
//		}
//		results := Eval(program, s)
//		if len(includeScope.store) != 3 {
//			t.Fatalf("program doesn't inlcude 3 modules. got=%d", len(includeScope.store))
//		}
//		testIntegerObject(t, results, tt.expected)
//	}
//}

func TestStringMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"string".find("s")`, 0},
		{`"string".find("string")`, 0},
		{`"string".find("g")`, 5},
		{`"string".find("tr")`, 1},
		{`"string".find("ng")`, 4},
		{`"string".find("x")`, NIL},
		{`"".find("stringstring")`, NIL},
		{`"string".find("")`, NIL},
		{`"string".find(1)`, NewError("0", INPUTERROR, INTEGER_OBJ, "find")},
		{`"string".find([])`, NewError("0", INPUTERROR, ARRAY_OBJ, "find")},
		{`"string".reverse()`, "gnirts"},
		{`"".reverse()`, ""},
		{`"ab".reverse()`, "ba"},
		{`"".reverse(1)`, NewError("0", ARGUMENTERROR, "0", 1)},
		{`"".upper()`, ""},
		{`"abc".upper()`, "ABC"},
		{`"a b c".upper()`, "A B C"},
		{`"a%b!c".upper()`, "A%B!C"},
		{`"ABC".lower()`, "abc"},
		{`"A B C".lower()`, "a b c"},
		{`"A%B!C".lower()`, "a%b!c"},
		{`" string".lstrip()`, "string"},
		{`"strsing".lstrip("s")`, "trsing"},
		{`" 	".lstrip()`, ""},
		{`"
			string".lstrip()`, "string"},
		{`"` + string('\r') + `string".lstrip()`, "string"},
		{`"string".lstrip("s")`, "tring"},
		{`"string".lstrip("st")`, "ring"},
		{`"ststring".lstrip("st")`, "ring"},
		{`"string ".rstrip()`, "string"},
		{`"` + string('\r') + string('\n') + "	 " + `".rstrip()`, ""},
		{`"string".rstrip()`, "string"},
		{`"string".rstrip("g")`, "strin"},
		{`"strging".rstrip("g")`, "strgin"},
		{`"string".rstrip("ng")`, "stri"},
		{`"string
			".rstrip()`, "string"},
		// strip just calls lstrip and rstrip consecutively, we can
		// have fewer tests here since the above is pretty comprehensive
		// just make sure it calls both
		{`" string ".strip()`, "string"},
		{`"ssstringss".strip("s")`, "tring"},
		{`let s = "1 2 3".split(); s[0] + s[1] + s[2]`, "123"},
		{`let s = "1,2,3".split(","); s[0] + s[1] + s[2]`, "123"},
		{`let s = "1&_2&_3&_".split("&_"); s[0] + s[1] + s[2] + s[3]`, "123"},
		{`"abc".replace("a", "A")`, "Abc"},
		{`"this is a story and this story tells the story of this".replace("this","that")`, "that is a story and that story tells the story of that"},
		{`" A B C ".replace(" ", "!")`, "!A!B!C!"},
		{`"eee".count("e")`, 3},
		{`"These are the days of summer".count("e")`, 5},
		{`"These are the days of summer".count(" ")`, 5},
		{`" ".join(["a", "b", "c"])`, "a b c"},
		{`"!".join(["a", "b", "c"])`, "a!b!c"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case *Null:
			testNullObject(t, expected)
		case string:
			testStringObject(t, evaluated, expected)
		case *Error:
			if evaluated.(*Error).Message != expected.Message {
				t.Fatalf("wrong error message. expected=%s, got=%s", expected.Message, evaluated.(*Error).Message)
			}
		}
	}
}
func TestStringIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"string"[0]`, "s"},
		{`"string"[-1]`, "g"},
		{`"string"[2]`, "r"},
		{`"string"[0:]`, "string"},
		{`"string"[1:]`, "tring"},
		{`"string"[2:5]`, "rin"},
		{`"string"[-5:-1]`, "trin"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo"=>5}["foo"]`, 5},
		{`{"foo"=>5}["bar"]`, nil},
		{`let key = "foo";{"foo"=>5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5=>5}[5]`, 5},
		{`{true=>5}[true]`, 5},
		{`{false=>5}[false]`, 5},
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

func TestHashLiterals(t *testing.T) {
	input := `
	let two = "two";
	{
		"one"        => 10 - 9,
		two          => 1 + 1,
		"thr" + "ee" => 6 /2,
		4            => 4,
		true         => 5,
		false        => 6
	}`

	evaluated := testEval(input)
	hash, ok := evaluated.(*Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T, (%+v)", evaluated, evaluated)
	}
	expected := map[HashKey]int64{
		(&String{String: "one", Valid: true}).HashKey():   1,
		(&String{String: "two", Valid: true}).HashKey():   2,
		(&String{String: "three", Valid: true}).HashKey(): 3,
		(&Integer{Int64: 4}).HashKey():                    4,
		TRUE.HashKey():                                    5,
		FALSE.HashKey():                                   6,
	}
	if len(hash.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. expected=%d, got=%d", len(expected), len(hash.Pairs))
	}
}

func TestStringHashKey(t *testing.T) {
	hello1 := &String{String: "Hello World", Valid: true}
	hello2 := &String{String: "Hello World", Valid: true}
	diff1 := &String{String: "My name is johnny", Valid: true}
	diff2 := &String{String: "My name is johnny", Valid: true}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() == hello1.HashKey() {
		t.Errorf("strings with different content have same hash key")
	}
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
			"let myArray = [1, 2, 3, 4, 5]; let i = myArray[0:]; let mySlice = myArray[1:]; mySlice[0]",
			2,
		},
		{
			"let myArray = [1, 2, 3, 4, 5]; let i = myArray[0]; let mySlice = myArray[:1]; mySlice[0]",
			1,
		},
		{
			"let myArray = [1, 2, 3, 4, 5]; let mySlice = myArray[:]; mySlice[0]",
			1,
		},
		{
			"let myArray = [1, 2, 3, 4, 5];let mySlice = myArray[:]; mySlice[-1]",
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			if err, ok := evaluated.(*Error); ok {
				if err.Message != "index error: '3' out of range" {
					t.Errorf("wrong error message. got=%s", err.Message)
				}
			} else {
				t.Errorf("evaluated not array or error. got=%T", evaluated)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	results, ok := evaluated.(*Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T", evaluated)
	}
	if len(results.Members) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(results.Members))
	}

	testIntegerObject(t, results.Members[0], 1)
	testIntegerObject(t, results.Members[1], 4)
	testIntegerObject(t, results.Members[2], 6)
}

func TestFunctionalMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`let a = [1,2].map(fn(x) {x + 1}); fn(x) { if (x[0] == 2) { if (x[1] == 3) { return true; }} else { return false }}(a)`, true},
		{`let a = [1,2].filter(fn(x) {x == 1}); fn(x) { if (x.len() == 1) { if (x[0] == 1) { return true; }} else { return false }}(a)`, true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestHashMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{1=>"a", 2=>"b"}.pop(1)`, "a"},
		{`let a = {1=>"a", 2=>"b"}; a.pop(1); str(a)`, `{2=> b}`},
		{`let a = {1=>"a", 2=>"b"}.push(3, "c"); a[3]`, `c`},
		{`let a = {1=>"a", 2=>"b"}; let b = {3=>"c"} let c = a.merge(b); c[3]`, `c`},
		{`let a = {1=>"a", 2=>"b"}; let b = {3=>"c"} let c = a.merge(b); str(a[3])`, `null`},
		{`let a = {1=>"a", 2=>"b"}; let b = {3=>"c"} let c = a.merge(b); str(b[1])`, `null`},
		{`let a = {"a"=>1}.map(fn(k, v){ {k.upper()=>v+1} } ); str(a)`, `{A=> 2}`},
		{`let a = {"a"=>1, "b"=>2}.filter(fn(k, v){ v > 1 } ); str(a)`, `{b=> 2}`},
		{`str({"a"=>1}.keys())`, `[a]`},
		{`str({"a"=>1}.values())`, `[1]`},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestArrayMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1,2,3].pop()`, 3},
		{`[1,2,3].pop(0)`, 1},
		{`[1,2,3].pop(2)`, 3},
		{`let a = [1,2,3].push(4);a.pop()`, 4},
		{`let a = [1,2,3].pop(1)`, 2},
		{`let a = [1,2,3]; a.pop(1); len(a)`, 2},
		{`let a = [1,2,3].filter(fn(x) { x > 1}); str(a)`, `[2, 3]`},
		{`let a = [1,2,3].map(fn(x) { x + 1}); str(a)`, `[2, 3, 4]`},
		{`let a = [1,2,3].merge([4]); str(a)`, `[1, 2, 3, 4]`},
		{`let a = ["a","b","c","d"].map(fn(x){ x.upper() }); str(a)`, `[A, B, C, D]`},
		{`["a","b","c","d"].index("d")`, 3},
		{`[1,1,1,2,3].count(1)`, 3},
		{`[1,2,3,4,5].reduce(fn(x, y) { x + y})`, 15},
		{`str([[0,1],[2,3],[4,5]].reduce(fn(acc, val) { acc.merge(val) }))`, "[0, 1, 2, 3, 4, 5]"},
		{`str([[0,1],[2,3],[4,5]].reduce(fn(acc, val) { acc.merge(val) },[]))`, "[0, 1, 2, 3, 4, 5]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len([1, 3, 5])`, 3},
		{`len([1,2,3])`, 3},
		{`"string".plus()`, "undefined method 'plus' for object STRING"},
		{`"string".plus`, "undefined method 'string.plus' for object STRING"},
		{`len("one", "two")`, "wrong number of arguments. expected=1, got=2"},
		{`len(1)`, "undefined method 'len' for object INTEGER"},
		{`int("1")`, 1},
		{`int("100")`, 100},
		{`int(1)`, 1},
		{`int("one")`, `unsupported input type 'STRING: one' for function or method: int`},
		{`int([])`, `unsupported input type 'ARRAY' for function or method: int`},
		{`int({})`, `unsupported input type 'HASH' for function or method: int`},
		{`str(1)`, "1"},
		{`str(true)`, `true`},
		{`str(false)`, `false`},
		{`str("string")`, `string`},
		{`str("one")`, `one`},
		{`str([])`, `[]`},
		{`str({})`, `{}`},
		{`type([])`, ARRAY_OBJ},
		{`type({})`, HASH_OBJ},
		{`type("")`, STRING_OBJ},
		{`type(true)`, BOOLEAN_OBJ},
		{`type(fn(x){x})`, FUNCTION_OBJ},
		{`type(1)`, INTEGER_OBJ},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch s := evaluated.(type) {
			case *Boolean:

			case *String:
				testStringObject(t, evaluated, expected)
			case *Error:
				if s.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, s.Message)
				}
			default:
				t.Errorf("object is not error. got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}
func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello, World!"`, "Hello, World!"},
		{`"Hello," + " " + "World!"`, "Hello, World!"},
	}

	for _, tt := range tests {
		testStringObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
let first = 10;
let second = 10;
let third = 10;

let ourFunction = fn(first) {
  let second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
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
		{"let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{"let fact = fn(n) { if(n==1) { return n } else { return n * fact(n-1) } }; fact(5);", 120},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Literal.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Literal.Parameters)
	}
	if fn.Literal.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Literal.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Literal.Body.String() != expectedBody {
		t.Fatalf("body is not '(x + 2). got=%q", fn.Literal.Body)
	}
}

func TestLetStatements(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5", 25},
		{"let a = 5; let b = a;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range test {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"unsupported operator for infix expression: '+' and types INTEGER and BOOLEAN",
		},
		{
			"5 + true; 5;",
			"unsupported operator for infix expression: '+' and types INTEGER and BOOLEAN",
		},
		{
			"-true",
			"unsupported operator for prefix expression:'-' and type: BOOLEAN",
		},
		{
			"true + false;",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
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
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{"foobar", "unknown identifier: 'foobar' is not defined"},
		//{`"abc" + 2`, "unsupported operator for infix expression: '+' and types STRING and INTEGER"},
		{`"abc" - "abc"`, "unsupported operator for infix expression: '-' and types STRING and STRING"},
		{`"abc" * "abc"`, "unsupported operator for infix expression: '*' and types STRING and STRING"},
		{`"abc" / "abc"`, "unsupported operator for infix expression: '/' and types STRING and STRING"},
		{`"abc" > "abc"`, "unsupported operator for infix expression: '>' and types STRING and STRING"},
		{`"abc" < "abc"`, "unsupported operator for infix expression: '<' and types STRING and STRING"},
		{`{"name"=>"Monkey"}[fn(x) {x}];`, "key error: type FUNCTION is not hashable"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
		{"let x = 5; return x;", 5},
		{"return;", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if integer, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
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
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 > 2) {10}", nil},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 < 2) {10} else {20}", 10},
		{"let x = 5;if(x == 5) { return;}", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if integer, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj Object) bool {
	if obj != NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestEvalBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(2 < 1) == true", false},
		{"let x = 5;x == 5", true},
		{"let x = 5; x != 5", false},
		{"let x = 5; x > 5", false},
		{"let x = 4; x < 5", true},
		{"let x = 4; (x + 5) > 5", true},
		{`"abc" == "abc"`, true},
		{`"abc" == "bc"`, false},
		{`"abc" != "abc"`, false},
		{`"abc" != "bc"`, true},
		{`let x = "abc"; x == "abc"`, true},
		{`let x = fn(){ "abc" }; x() == "abc"`, true},
		{"true and true", true},
		{"true and false", false},
		{"true or true", true},
		{"true or false", true},
		{`"string" and false`, false},
		{`[] or false`, true},
		{`len([1,2,3]) > 2 and false`, false},
		{`type([]) == "ARRAY" and len([1234]) == 4`, false},
		{`type([]) == "ARRAY" and len("1234") == 4`, true},
		{"(true and true) or (true or false)", true},
		{"(true and true) and (true and false)", false},
		{`!!"abc".find("d")`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Bool != expected {
		t.Errorf("object has wrong value. got=%v, want=%v", result.Bool, expected)
		return false
	}
	return true
}

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
		{"20 % 4", 0},
		{"20 % 3", 2},
		{"5 * 4 % 3", 2},
		{"20 % 3 * 1", 2},
		{"5 * (3 % 1)", 0},
		{"3 % 20", 3},
		{"-4 % 20", -4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) Object {
	l := lexer.New(input)
	path, _ := os.Getwd()
	p := parser.New(l, path)
	s := NewScope(nil)
	program := p.ParseProgram()
	if len(program.Statements) == 0 {
		if len(program.Includes) == 0 {
			fmt.Printf("Parsed program has no statements or included objects.\n")
			os.Exit(1)
		}
	}
	return Eval(program, s)
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Int64 != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Int64, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj Object, expected string) bool {
	result, ok := obj.(*String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.String != expected {
		t.Errorf("object has wrong value. got='%s', want='%s'", result.String, expected)
		return false
	}
	return true
}

func TestInterpolation(t *testing.T) {
	input := []struct {
		input    string
		expected string
	}{
		{`let x = 5; 'abc{x}'`, "abc5"},
		{`'abc{x}'`, "abcx"},
		{`'abc{5 + 5}abc'`, "abc10abc"},
		{`let x = fn(x) { x * 5 };'{x(1)}{x(5)}{x(10)}'`, "52550"},
		{`let x = fn(x) { x * 5 };'abcdef{x(10)}'`, "abcdef50"},
		{`'abcdef{(10 * 5)}'`, "abcdef50"},
		{`'{10 + 10}abcdef{(10 * 5)}'`, "20abcdef50"},
		{`let x = 5; let y = '{x}';'{y}abcdef{(10 * x)}'`, "5abcdef50"},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)
		testInterpolatedStringObject(t, evaluated, tt.expected)
	}
}

func testInterpolatedStringObject(t *testing.T, obj Object, expected string) bool {
	result, ok := obj.(*InterpolatedString)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Inspect() != expected {
		t.Errorf("object has wrong value. got='%s', want='%s'", result.Inspect(), expected)
		return false
	}
	return true
}
