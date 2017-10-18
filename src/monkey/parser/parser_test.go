package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"os"
	"strings"
	"testing"
)

var path, _ = os.Getwd()

func TestParsingDoLoopExpression(t *testing.T) {
	input := `do {}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	_, ok := stmt.Expression.(*ast.DoLoop)
	if !ok {
		t.Fatalf("exp is not ast.DoLoop. got=%T", stmt)
	}
}

func TestParsingWhileLoopExpression(t *testing.T) {
	input := `while (5 < 10 ){}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	_, ok := stmt.Expression.(*ast.WhileLoop)
	if !ok {
		t.Fatalf("exp is not ast.WhileLoop. got=%T", stmt)
	}
}

func TestParsingForLoopExpression(t *testing.T) {
	//input := `for (i = 0; i< 10; i = i+1) {}`
	input := `for (i; i<10; i=i+1) {}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	_, ok := stmt.Expression.(*ast.ForLoop)
	if !ok {
		t.Fatalf("exp is not ast.ForLoop. got=%T", stmt)
	}
}

func TestParsingForEachArrayLoopExpression(t *testing.T) {
	input := `for x in array where x > 5 {}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	a, ok := stmt.Expression.(*ast.ForEachArrayLoop)
	if !ok {
		t.Fatalf("exp is not ast.ForEachArrayLoop. got=%T", stmt)
	}
	if a.Var != "x" {
		t.Fatalf("a.Var is not 'x'. got=%T", a.Var)
	}
}

func TestParsingForEachMapLoopExpression(t *testing.T) {
	input := `for key, value in hash {}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	a, ok := stmt.Expression.(*ast.ForEachMapLoop)
	if !ok {
		t.Fatalf("exp is not ast.ForEachMapLoop. got=%T", stmt)
	}
	if a.Key != "key" {
		t.Fatalf("a.Key is not 'key'. got=%T", a.Key)
	}
	if a.Value != "value" {
		t.Fatalf("a.Value is not 'value'. got=%T", a.Value)
	}
}

func TestParsingGrepExpression(t *testing.T) {
	input := `grep { $_ > 5 } [2,4,6,8,10]`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	a, ok := stmt.Expression.(*ast.GrepExpr)
	if !ok {
		t.Fatalf("exp is not ast.GrepExpr. got=%T", stmt)
	}
	if a.Var != "$_" {
		t.Fatalf("a.Var is not '$_'. got=%T", a.Var)
	}
	t.Logf(a.Block.String())
	t.Logf(a.Value.String())
}
func TestParsingAssignmentExpressions(t *testing.T) {
	input := `x = 5`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	a, ok := stmt.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("exp is not ast.AssignExpression. got=%T", stmt)
	}
	testIdentifier(t, a.Name, "x")

	testIntegerLiteral(t, a.Value, int64(5))
}

func TestParsingFloatAssignmentExpressions(t *testing.T) {
	input := `x = 5.234`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	a, ok := stmt.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("exp is not ast.AssignExpression. got=%T", stmt)
	}
	testIdentifier(t, a.Name, "x")

	testFloatLiteral(t, a.Value, float64(5.234))
}

func TestParsingEmptyHashLiteralExpressions(t *testing.T) {
	input := `{}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt)
	}
	if len(hash.Pairs) != 0 {
		t.Fatalf("wrong number of hash pairs. expected=0, got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralExpressions(t *testing.T) {
	input := `{"one" => 1, "two" => 2, "three"=> 3}`
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("wrong number of hash pairs. expected=3, got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		if literal, ok := key.(*ast.StringLiteral); ok {
			expectedValue := expected[literal.String()]
			testIntegerLiteral(t, value, expectedValue)
		} else {
			t.Fatalf("key not *ast.StringLiteral. got=%T", key)
		}
	}

}

func TestParsingMethodExpressions(t *testing.T) {
	input := "array.len(1, 2)"
	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	methodCall, ok := stmt.Expression.(*ast.MethodCallExpression)
	if !ok {
		t.Fatalf("exp not *ast.MethodCallExpression. got=%T (%+v)", stmt.Expression, stmt.Expression)
	}
	if len(methodCall.Call.(*ast.CallExpression).Arguments) != 2 {
		t.Fatalf("wrong number of arguments. expected=2, got=%d", len(methodCall.Call.(*ast.CallExpression).Arguments))
	}
}

func TestParsingSliceExpressions(t *testing.T) {
	tests := []struct {
		input string
		start interface{}
		end   interface{}
	}{
		{"myArray[-1:]", -1, nil},
		{"myArray[1:3];", 1, 3},
		{"myArray[:3];", 0, 3},
		{"myArray[1:];", 1, nil},
		{"myArray[fn(){5}():5]", "fn () { 5 }", 5},
		{"myArray[a:3];", "a", 3},
		{"myArray[:-1]", 0, -1},
		{"myArray[5:fn(){5}()]", 5, "fn () { 5 }"},
		{"myArray[3:a];", 3, "a"},
		{"myArray[1 + 1:0];", nil, 0},
		{"myArray[0:1 + 1];", 0, nil},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		indexExp, ok := stmt.Expression.(*ast.IndexExpression)
		if !ok {
			t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
		}
		testIdentifier(t, indexExp.Left, "myArray")
		sliceExp, ok := indexExp.Index.(*ast.SliceExpression)
		if !ok {
			t.Fatalf("exp not *ast.SliceExpression. got=%T", indexExp.Index)
		}
		switch start := sliceExp.StartIndex.(type) {
		case *ast.PrefixExpression:
			s := tt.start.(int)
			testIntegerLiteral(t, start.Right, int64(-s))
		case *ast.IntegerLiteral:
			s := tt.start.(int)
			testIntegerLiteral(t, start, int64(s))
		case *ast.CallExpression:
			s := tt.start.(string)
			if s != start.Function.String() {
				t.Fatalf("callExp index not %s. got=%s", s, start)
			}
		case *ast.Identifier:
			s := tt.start.(string)
			if s != start.String() {
				t.Fatalf("callExp index not %s. got=%s", s, start)
			}
		case *ast.InfixExpression:
			testInfixExpression(t, start, 1, "+", 1)

		}
		switch end := sliceExp.EndIndex.(type) {
		case *ast.PrefixExpression:
			e := tt.end.(int)
			testIntegerLiteral(t, end.Right, int64(-e))
		case *ast.IntegerLiteral:
			e := tt.end.(int)
			testIntegerLiteral(t, end, int64(e))
		case *ast.CallExpression:
			e := tt.end.(string)
			if e != end.Function.String() {
				t.Fatalf("callExp index not %s. got=%s", e, end)
			}
		case *ast.Identifier:
			e := tt.end.(string)
			if e != end.String() {
				t.Fatalf("callExp index not %s. got=%s", e, end)
			}
		case *ast.InfixExpression:
			testInfixExpression(t, end, 1, "+", 1)
		}
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestArrayExpression(t *testing.T) {
	input := "[1, 2 * 3, 2 + 2]"
	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", program.Statements[0])
	}
	if len(array.Members) != 3 {
		t.Fatalf("array length not 3. got=%d", len(array.Members))
	}
	testIntegerLiteral(t, array.Members[0], 1)
	testInfixExpression(t, array.Members[1], 2, "*", 3)
	testInfixExpression(t, array.Members[2], 2, "+", 2)
}

func TestRegExLiteralExpression(t *testing.T) {
	input := `~\d+(\w)+.*$~;`

	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.RegExLiteral)
	if !ok {
		t.Fatalf("exp not *ast.RegExLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != `\d+(\w)+.*$` {
		t.Fatalf("literal.Value not correct, got=%s", literal.Value)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello, world";`

	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello, world" {
		t.Fatalf("literal.Value not 'hello, world', got=%s", literal.Value)
	}

}

func TestInterpolatedString(t *testing.T) {
	tests := []struct {
		input     string
		expected  string
		argscount int
	}{
		{`'{hello}{x}{world}'`, "{0}{1}{2}", 3},
		{"'abc{x}'", "abc{0}", 1},
		{"'{x}{x}{x}'", "{0}{1}{2}", 3},
		{"'aa{x(x)}{x}{x+1}aa'", "aa{0}{1}{2}aa", 3},
		{"'aa{[1,2,3]}{x()}{x+1}aa'", "aa{0}{1}{2}aa", 3},
		{"'aa{[1,2,3]}b{x()}c{x+1}aa'", "aa{0}b{1}c{2}aa", 3},
		{"'aa{x+1}abc'", "aa{0}abc", 1},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		is, ok := stmt.Expression.(*ast.InterpolatedString)
		if !ok {
			t.Fatalf("exp not *ast.InterpolatedString. got=%T", stmt.Expression)
		}
		if is.Value != tt.expected {
			t.Fatalf("is.Value not %s, got=%s", tt.expected, is.Value)
		}
		if len(is.ExprMap) != tt.argscount {
			t.Fatalf("is.ExprList has wrong number of expressions. expected=%d, got=%d", tt.argscount, len(is.ExprMap))
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Values[0]
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestIncludeStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue string
		ismodule      bool
	}{
		{"include test_files", "test_files", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Includes) != 1 {
			t.Fatalf("program.Includes does not contain 1 statements. got=%d", len(program.Includes))
		}
		for _, v := range program.Includes {
			if !testLiteralExpression(t, v.IncludePath, tt.expectedValue) {
				return
			}
			if v.IsModule != tt.ismodule {
				t.Fatalf("Included value not %v. got=%v", v.IsModule, tt.ismodule)
			}
			if len(v.Program.Includes) != 3 {
				t.Fatalf("Included Program had wrong number of modules. expected=3, got=%d", len(v.Program.Includes))
			}
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.Tokenliteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	if letStmt, ok := s.(*ast.LetStatement); !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	} else if letStmt.Names[0].Value != name {
		t.Errorf("ltStmt.Names[0].Value not '%s'. got=%s", name, letStmt.Names[0].Value)
		return false
	} else if letStmt.Names[0].TokenLiteral() != name {
		t.Errorf("s.Name not %s. got=%s", name, letStmt.Names[0])
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
		{"return;", nil},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if returnStmt, ok := stmt.(*ast.ReturnStatement); ok {
			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLIteral not 'return', got %q", returnStmt.TokenLiteral())
			}
			if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
				return
			}
		} else {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got %s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got %s", "foobar", ident.TokenLiteral())
	}
}

func testBooleanExpression(t *testing.T, exp ast.Expression, value bool) bool {
	be, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if be.Value != value {
		t.Errorf("boolean.Value not '%t'. got %t", value, be.Value)
		return false
	}
	if be.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got %s", value, be.TokenLiteral())
		return false
	}
	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("ident.Value not %d. got %d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got %s", "5", literal.TokenLiteral())
	}

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l, path)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statments does not contain %d statements. got=%d", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d, got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

var EPSILON = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	flt, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}

	if !floatEquals(flt.Value, value) {
		t.Errorf("flt.Value not %g. got=%g", value, flt.Value)
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 % 5;", 5, "%", 5},
		{"true == true", true, "==", true},
		{"true != true", true, "!=", true},
		{"false == false", false, "==", false},
		{"true or true", true, "or", true},
		{"true and false", true, "and", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l, path)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not equal %d. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a * b % c",
			"((a * b) % c)",
		},
		{
			"a % b / c",
			"((a % b) / c)",
		},
		{
			"a + b % c",
			"(a + (b % c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false;",
			"((3 > 5) == false)",
		},
		{
			"3 > 5 == true;",
			"((3 > 5) == true)",
		},
		{
			"1 - (2 + 3) + 4",
			"((1 - (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"add(a) or b",
			"(add(a) or b)",
		},
		{
			"x == y or x == z",
			"((x == y) or (x == z))",
		},
		{
			"x == y and x == z",
			"((x == y) and (x == z))",
		},
		{
			"x or y and (x and z)",
			"(x or (y and (x and z)))",
		},
		{
			"(x or y) and (x or z)",
			"((x or y) and (x or z))",
		},
		{
			"(x and y) or (x and z)",
			"((x and y) or (x and z))",
		},
		{
			"(x or y) ==  (x and z)",
			"((x or y) == (x and z))",
		},
		{
			"a[0] and x",
			"((a[0]) and x)",
		},
		{
			`str(x) or i.find("abc")`,
			"(str(x) or i.find(abc))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Logf("%v, %T", value, exp)
	if ident, ok := exp.(*ast.Identifier); !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	} else if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	} else if ident.TokenLiteral() != value {
		t.Errorf("identTokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanExpression(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	if opExp, ok := exp.(*ast.InfixExpression); !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	} else if !testLiteralExpression(t, opExp.Left, left) {
		return false
	} else if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
	} else if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program body does not include %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("consequence does not include %d statements. got=%d", 1, len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testIdentifier(t, consequence.Expression, "x")
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative != nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program body does not include %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("consequence does not include %d statements. got=%d", 1, len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testIdentifier(t, consequence.Expression, "x")
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("alternative does not include %d statements. got=%d", 1, len(exp.Alternative.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		testIdentifier(t, alternative.Expression, "x")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l, path)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("function is not FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements wrong. want 1, got=%d", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l, path)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of parameters wrong. want=%d, got=%d", len(function.Parameters), len(tt.expectedParams))
		}
		t.Logf("want=%s, got=%s", function.Parameters, tt.expectedParams)
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("function is not ast.CallExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("function literal parameters wrong. want 3, got=%d", len(exp.Arguments))
	}
	t.Logf("want=%s, got=%s", strings.Split(input[4:len(input)-1], ", "), exp.Arguments)
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestTryExpressionParsing(t *testing.T) {
	input := `
                  try {
                      let th = 1 + 2
                      if (th == 3) { throw "SUMERROR" }
                  }
                  catch "OTHERERROR" {
                      putln("Catched OTHERERROR")
                  }
                  catch "SUMERROR" {
                      putln("Catched SUMERROR")
                  }
                  catch {
                      putln("Catched ALL")
                  }
                  finally {
                      putln("Finally running")
                  }

`

	l := lexer.New(input)
	p := New(l, path)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	_, ok = stmt.Expression.(*ast.TryStmt)
	if !ok {
		t.Fatalf("function is not ast.TryStmt. got=%T", stmt.Expression)
	}
}
