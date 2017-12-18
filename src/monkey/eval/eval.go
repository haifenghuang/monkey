package eval

import (
	"bytes"
	"fmt"
	"math"
	"monkey/ast"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"unicode/utf8"
)

var (
	TRUE     = &Boolean{Bool: true, Valid: true}
	FALSE    = &Boolean{Bool: false, Valid: true}
	BREAK    = &Break{}
	CONTINUE = &Continue{}
	NIL      = &Nil{}
)

var includeScope *Scope
var importedCache map[string]Object
var mux sync.Mutex

func Eval(node ast.Node, scope *Scope) Object {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case *Error:
				//if panic is a Error Object, print its contents
				fmt.Println(r.Error())
				return
			case runtime.Error:
				const msgFmt = "%s:%d %s(): A panic occured, but was recovered; Details: %+v;\n\t*** Stack Trace ***\n\t%s*** End Stack Trace ***\n"
				funcName, filename, lineNum := "Unknown Function", "Unknown File", -1
				var programCounter uintptr
				var isKnownFunc bool
				var depth int = 4
				if programCounter, filename, lineNum, isKnownFunc = runtime.Caller(depth); isKnownFunc {
					filename = path.Base(filename)
					funcName = path.Base(runtime.FuncForPC(programCounter).Name())
				}
				prettyStack := bytes.Join(bytes.Split(debug.Stack(), []byte{'\n'})[6:], []byte{'\n', '\t'})
				fmt.Errorf(msgFmt, filename, lineNum, funcName, r, prettyStack)
				//fmt.Println(r.Error())
				//panic(r)
				return
			}
		}
	}()

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.IncludeStatement:
		return evalIncludeStatement(node, scope)
	case *ast.LetStatement:
		return evalLetStatement(node, scope)
	case *ast.ReturnStatement:
		return evalReturnStatment(node, scope)
	case *ast.DeferStmt:
		return evalDeferStatment(node, scope)
	case *ast.NamedFunction:
		return evalNamedFnStatement(node, scope)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return evalFloatLiteral(node)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.InterpolatedString:
		return evalInterpolatedString(node, scope)
	case *ast.Identifier:
		return evalIdentifier(node, scope)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, scope)
	case *ast.HashLiteral:
		return evalHashLiteral(node, scope)
	case *ast.StructLiteral:
		return evalStructLiteral(node, scope)
	case *ast.EnumLiteral:
		return evalEnumLiteral(node, scope)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, scope)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, scope)
	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		if left.Type() == ERROR_OBJ {
			return left
		}

		right := Eval(node.Right, scope)
		if right.Type() == ERROR_OBJ {
			return right
		}
		return evalInfixExpression(node, left, right)
	case *ast.PostfixExpression:
		left := Eval(node.Left, scope)
		if left.Type() == ERROR_OBJ {
			return left
		}
		return evalPostfixExpression(left, node)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, scope)
	case *ast.CallExpression:
		return evalFunctionCall(node, scope)
	case *ast.MethodCallExpression:
		return evalMethodCallExpression(node, scope)
	case *ast.IndexExpression:
		return evalIndexExpression(node, scope)
	case *ast.GrepExpr:
		return evalGrepExpression(node, scope)
	case *ast.MapExpr:
		return evalMapExpression(node, scope)
	case *ast.CaseExpr:
		return evalCaseExpression(node, scope)
	case *ast.DoLoop:
		return evalDoLoopExpression(node, scope)
	case *ast.WhileLoop:
		return evalWhileLoopExpression(node, scope)
	case *ast.ForLoop:
		return evalForLoopExpression(node, scope)
	case *ast.ForEverLoop:
		return evalForEverLoopExpression(node, scope)
	case *ast.ForEachArrayLoop:
		return evalForEachArrayExpression(node, scope)
	case *ast.ForEachDotRange:
		return evalForEachDotRangeExpression(node, scope)
	case *ast.ForEachMapLoop:
		return evalForEachMapExpression(node, scope)
	case *ast.ListComprehension:
		return evalListComprehension(node, scope)
	case *ast.ListRangeComprehension:
		return evalListRangeComprehension(node, scope)
	case *ast.MapComprehension:
		return evalMapComprehension(node, scope)
	case *ast.BreakExpression:
		return BREAK
	case *ast.ContinueExpression:
		return CONTINUE
	case *ast.ThrowStmt:
		return evalThrowStatement(node, scope)
	case *ast.AssignExpression:
		return evalAssignExpression(node, scope)
	case *ast.RegExLiteral:
		return evalRegExLiteral(node)
	case *ast.TryStmt:
		return evalTryStatement(node, scope)
	case *ast.TernaryExpression:
		return evalTernaryExpression(node, scope)
	case *ast.SpawnStmt:
		return evalSpawnStatment(node, scope)
	case *ast.YieldExpression:
		return NIL
	case *ast.NilLiteral:
		return NIL
	case *ast.Pipe:
		return evalPipeExpression(node, scope)
	}
	return nil
}

// Program Evaluation Entry Point Functions, and Helpers:
func evalProgram(program *ast.Program, scope *Scope) (results Object) {
	loadIncludes(program.Includes, scope)
	for _, statement := range program.Statements {
		results = Eval(statement, scope)
		switch s := results.(type) {
		case *ReturnValue:
			return s.Value
		case *Error:
			if s.Kind == THROWNOTHANDLED {
				panic(NewError(statement.Pos().Sline(), THROWNOTHANDLED, s.Message))
			}
			return s
			//		case *ThrowValue:
			//			//convert ThrowValue to Errors
			//			throwObj := results.(*ThrowValue).Value
			//			throwObjStr := throwObj.(*String).String
			//			panic(NewError(statement.Pos().Sline(), THROWNOTHANDLED, throwObjStr))
		}
	}
	if results == nil {
		return NIL
	}
	return results
}

func loadIncludes(includes map[string]*ast.IncludeStatement, scope *Scope) {
	if includeScope == nil {
		includeScope = NewScope(nil)
	}
	for _, p := range includes {
		Eval(p, scope)
	}
}

// Statements...
func evalIncludeStatement(i *ast.IncludeStatement, scope *Scope) Object {

	mux.Lock()
	defer mux.Unlock()

	// Check the cache
	if cache, ok := importedCache[i.IncludePath.String()]; ok {
		return cache
	}

	imported := &IncludedObject{Name: i.IncludePath.String(), Scope: NewScope(nil)}

	// capture stdout to suppress output during evaluating import
	so := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	if _, ok := includeScope.Get(i.IncludePath.String()); !ok {
		evalProgram(i.Program, imported.Scope)
		includeScope.Set(i.IncludePath.String(), imported)
	}

	// restore stdout
	w.Close()
	os.Stdout = so

	//store the evaluated result to cache
	importedCache[i.IncludePath.String()] = imported

	return imported
}

func evalLetStatement(l *ast.LetStatement, scope *Scope) (val Object) {
	valuesLen := len(l.Values)

	for idx, item := range l.Names {
		if idx >= valuesLen { //There are more Values than Names
			scope.Set(item.String(), NIL)
		} else {
			val = Eval(l.Values[idx], scope)
			if val.Type() != ERROR_OBJ {
				scope.Set(item.String(), val)
			} else {
				return
			}
		}
	}

	return
}

func evalNumAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	var leftVal float64
	var rightVal float64

	isInt := left.Type() == INTEGER_OBJ && val.Type() == INTEGER_OBJ

	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*Integer).Int64)
	} else {
		leftVal = left.(*Float).Float64
	}

	//Check `right`'s type
	if val.Type() == INTEGER_OBJ {
		rightVal = float64(val.(*Integer).Int64)
	} else {
		rightVal = val.(*Float).Float64
	}

	var ok bool
	switch a.Token.Literal {
	case "+=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal+rightVal)))
			if ok {
				return
			}
		} else {
			ret, ok = scope.Reset(name, NewFloat(leftVal+rightVal))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "-=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal-rightVal)))
			if ok {
				return
			}
		} else {
			ret, ok = scope.Reset(name, NewFloat(leftVal-rightVal))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "*=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal*rightVal)))
			if ok {
				return
			}
		} else {
			ret, ok = scope.Reset(name, NewFloat(leftVal*rightVal))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "/=":
		if rightVal == 0 {
			panic(NewError(a.Pos().Sline(), DIVIDEBYZERO))
		}

		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal/rightVal)))
			if ok {
				return
			}
		} else {
			ret, ok = scope.Reset(name, NewFloat(leftVal/rightVal))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "%=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)%int64(rightVal)))
			if ok {
				return
			}
		} else {
			ret, ok = scope.Reset(name, NewFloat(math.Mod(leftVal, rightVal)))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "^=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)^int64(rightVal)))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "|=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)|int64(rightVal)))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "&=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)&int64(rightVal)))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
	}
	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
}

func evalStrAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	leftVal := left.(*String).String
	var ok bool

	switch a.Token.Literal {
	case "=":
		switch nodeType := a.Name.(type) {
		case *ast.IndexExpression: //str[idx] = xxx
			index := Eval(nodeType.Index, scope)
			if index == NIL {
				ret = NIL
				return
			}
			if index.Type() != INTEGER_OBJ { //must be an integer
				panic(NewError(a.Pos().Sline(), GENERICERROR, "Array index value should evaluate to an integer"))
			}

			idx := index.(*Integer).Int64
			if idx < 0 || idx >= int64(len(leftVal)) {
				panic(NewError(a.Pos().Sline(), INDEXERROR, idx))
			}

			str := NewString(leftVal[:idx] + val.Inspect() + leftVal[idx+1:])
			ret, ok = scope.Reset(name, str)
			if ok {
				return
			}
		}
	}

	ret, ok = scope.Reset(name, NewString(leftVal+val.Inspect()))
	if ok {
		return
	}
	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

}

func evalArrayAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	leftVals := left.(*Array).Members

	var ok bool
	switch a.Token.Literal {
	case "+=":
		switch nodeType := a.Name.(type) {
		case *ast.Identifier:
			name = nodeType.Value
			leftVals = append(leftVals, val)
			ret, ok = scope.Reset(name, &Array{Members: leftVals})
			if ok {
				return
			}
		}
	case "=":
		switch nodeType := a.Name.(type) {
		case *ast.IndexExpression: //arr[idx] = xxx
			index := Eval(nodeType.Index, scope)
			if index == NIL {
				ret = NIL
				return
			}
			if index.Type() != INTEGER_OBJ { //must be an integer
				panic(NewError(a.Pos().Sline(), GENERICERROR, "Array index value should evaluate to an integer"))
			}

			idx := index.(*Integer).Int64
			if idx < 0 {
				panic(NewError(a.Pos().Sline(), INDEXERROR, idx))
			}

			if idx < int64(len(leftVals)) { //index is in range
				leftVals[idx] = val
				ret, ok = scope.Reset(name, &Array{Members: leftVals})
				if ok {
					return
				}
			} else { //index is out of range, we auto-expand the array
				for i := int64(len(leftVals)); i < idx; i++ {
					leftVals = append(leftVals, NIL)
				}

				leftVals = append(leftVals, val)
				ret, ok = scope.Reset(name, &Array{Members: leftVals})
				if ok {
					return
				}
			}
		}

		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
	}

	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
}

func evalHashAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	leftVals := left.(*Hash).Pairs

	var ok bool
	switch a.Token.Literal {
	case "+=":
		if _, ok := val.(*Hash); !ok { //must be hash type
			panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
		}

		rightVals := val.(*Hash).Pairs
		for k, v := range rightVals {
			leftVals[k] = v
		}
		ret, ok = scope.Reset(name, &Hash{Pairs: leftVals})
		if ok {
			return
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
	case "-=":
		hashable, ok := val.(Hashable)
		if !ok {
			panic(NewError(a.Pos().Sline(), KEYERROR, val.Type()))
		}
		if hashPair, ok := leftVals[hashable.HashKey()]; ok {
			delete(leftVals, hashable.HashKey())
			return hashPair.Value
		}
		ret, ok = scope.Reset(name, &Hash{Pairs: leftVals})
		if ok {
			return
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
	case "=":
		switch nodeType := a.Name.(type) {
		case *ast.IndexExpression: //hash[key] = xxx
			key := Eval(nodeType.Index, scope)
			hashable, ok := key.(Hashable)
			if !ok {
				panic(NewError(a.Pos().Sline(), KEYERROR, val.Type()))
			}
			leftVals[hashable.HashKey()] = HashPair{Key: key, Value: val}
			return left
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
	}

	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
}

func evalStructAssignExpression(a *ast.AssignExpression, scope *Scope, val Object) (retVal Object) {
	strArr := strings.Split(a.Name.String(), ".")
	var aObj Object
	var aVal Object
	var ok bool
	if aObj, ok = scope.Get(strArr[0]); !ok {
		panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[0]))
	}

	st, ok := aObj.(*Struct)
	if !ok {
		panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[0]))
	}

	if aVal, ok = st.Scope.Get(strArr[1]); !ok {
		panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[1]))
	}

	structScope := st.Scope

	if a.Token.Literal == "=" {
		v, ok := structScope.Reset(strArr[1], retVal)
		if ok {
			return v
		}
		panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, a.Name.String()))
	}

	switch aVal.Type() {
	case INTEGER_OBJ, FLOAT_OBJ:
		retVal = evalNumAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	case STRING_OBJ:
		retVal = evalStrAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	case ARRAY_OBJ:
		retVal = evalArrayAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	case HASH_OBJ:
		retVal = evalHashAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	}

	panic(NewError(a.Pos().Sline(), INFIXOP, aVal.Type(), a.Token.Literal, val.Type()))
}

func evalAssignExpression(a *ast.AssignExpression, scope *Scope) (val Object) {
	val = Eval(a.Value, scope)
	if val.Type() == ERROR_OBJ {
		return val
	}

	if strings.Contains(a.Name.String(), ".") {
		var aObj Object
		var ok bool

		strArr := strings.Split(a.Name.String(), ".")
		if aObj, ok = scope.Get(strArr[0]); !ok {
			panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[0]))
		}
		if aObj.Type() == ENUM_OBJ { //it's enum type
			panic(NewError(a.Pos().Sline(), GENERICERROR, "Enum value cannot be reassigned!"))
		}

		return evalStructAssignExpression(a, scope, val)
	}

	var name string
	switch nodeType := a.Name.(type) {
	case *ast.Identifier:
		name = nodeType.Value
	case *ast.IndexExpression:
		name = nodeType.Left.(*ast.Identifier).Value
	}

	if a.Token.Literal == "=" {
		switch nodeType := a.Name.(type) {
		case *ast.Identifier:
			name := nodeType.Value
			v, ok := scope.Reset(name, val)
			if ok {
				return v
			}
			panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, a.Name.String()))
		}
	}

	// Check if the variable exists or not
	var left Object
	var ok bool
	if left, ok = scope.Get(name); !ok {
		panic(NewError(a.Pos().Sline(), UNKNOWNIDENT, name))
	}

	switch left.Type() {
	case INTEGER_OBJ, FLOAT_OBJ:
		val = evalNumAssignExpression(a, name, left, scope, val)
		return
	case STRING_OBJ:
		val = evalStrAssignExpression(a, name, left, scope, val)
		return
	case ARRAY_OBJ:
		val = evalArrayAssignExpression(a, name, left, scope, val)
		return
	case HASH_OBJ:
		val = evalHashAssignExpression(a, name, left, scope, val)
		return
	}

	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
}

func evalReturnStatment(r *ast.ReturnStatement, scope *Scope) Object {
	if value := Eval(r.ReturnValue, scope); value != nil {
		return &ReturnValue{Value: value}
	}

	return NIL
}

func evalDeferStatment(d *ast.DeferStmt, scope *Scope) Object {
	frame := scope.CurrentFrame()
	if frame == nil {
		panic(NewError(d.Pos().Sline(), DEFERERROR))
	}

	if frame.CurrentCall == nil {
		panic(NewError(d.Pos().Sline(), DEFERERROR))
	}

	switch d.Call.(type) {
	case *ast.CallExpression:
		callExp := d.Call.(*ast.CallExpression)
		closure := func() {
			evalFunctionCall(callExp, scope)
		}
		frame.defers = append(frame.defers, closure)
	case *ast.MethodCallExpression:
		callExp := d.Call.(*ast.MethodCallExpression)
		closure := func() {
			evalMethodCallExpression(callExp, scope)
		}
		frame.defers = append(frame.defers, closure)
	}

	return NIL
}

func evalThrowStatement(t *ast.ThrowStmt, scope *Scope) Object {
	value := Eval(t.Expr, scope)
	if value.Type() == ERROR_OBJ {
		return value
	}

	var strObj *String
	var ok bool
	if strObj, ok = value.(*String); !ok {
		panic(NewError(t.Pos().Sline(), THROWERROR))
	}
	return &Error{Kind: THROWNOTHANDLED, Message: strObj.String}
}

// Booleans
func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Literals
func evalIntegerLiteral(i *ast.IntegerLiteral) Object {
	return NewInteger(i.Value)
}

func evalFloatLiteral(f *ast.FloatLiteral) Object {
	return NewFloat(f.Value)
}

func evalStringLiteral(s *ast.StringLiteral) Object {
	return NewString(s.Value)
}

func evalInterpolatedString(is *ast.InterpolatedString, scope *Scope) Object {
	s := &InterpolatedString{String: &String{Valid: true}, RawValue: is.Value, Expressions: is.ExprMap}
	s.Interpolate(scope)
	return s
}

func evalArrayLiteral(a *ast.ArrayLiteral, scope *Scope) Object {
	return &Array{Members: evalArgs(a.Members, scope)}
}

func evalRegExLiteral(re *ast.RegExLiteral) Object {
	regExpression, err := regexp.Compile(re.Value)
	if err != nil {
		panic(NewError(re.Pos().Sline(), INVALIDARG))
	}

	return &RegEx{RegExp: regExpression, Value: re.Value}
}

func evalIdentifier(i *ast.Identifier, scope *Scope) Object {
	//Get from global scope first
	if obj, ok := GetGlobalObj(i.String()); ok {
		return obj
	}

	val, ok := scope.Get(i.String())
	if !ok {
		//panic(NewError(i.Pos().Sline(), UNKNOWNIDENT, i.String()))

		if val, ok = includeScope.Get(i.String()); !ok {
			panic(NewError(i.Pos().Sline(), UNKNOWNIDENT, i.String()))
		}
	}
	if i, ok := val.(*InterpolatedString); ok {
		i.Interpolate(scope)
		return i
	}

	return val
}

func evalHashLiteral(hl *ast.HashLiteral, scope *Scope) Object {
	innerScope := NewScope(scope)
	hashMap := make(map[HashKey]HashPair)
	for key, value := range hl.Pairs {
		key := Eval(key, innerScope)
		if key.Type() == ERROR_OBJ {
			return key
		}

		if hashable, ok := key.(Hashable); ok {
			v := Eval(value, innerScope)
			if v.Type() == ERROR_OBJ {
				return v
			}
			hashMap[hashable.HashKey()] = HashPair{Key: key, Value: v}
		} else {
			panic(NewError(hl.Pos().Sline(), KEYERROR, key.Type()))
		}
	}
	return &Hash{Pairs: hashMap}
}

func evalStructLiteral(s *ast.StructLiteral, scope *Scope) Object {
	structScope := NewScope(nil)
	for key, value := range s.Pairs {
		if ident, ok := key.(*ast.Identifier); ok {
			aObj := Eval(value, scope)
			structScope.Set(ident.String(), aObj)
		} else {
			panic(NewError(s.Pos().Sline(), KEYERROR, "IDENT"))
		}
	}
	return &Struct{Scope: structScope, methods: make(map[string]*Function)}
}

func evalEnumLiteral(e *ast.EnumLiteral, scope *Scope) Object {
	enumScope := NewScope(nil)
	for key, value := range e.Pairs {
		if ident, ok := key.(*ast.Identifier); ok {
			aObj := Eval(value, scope)
			enumScope.Set(ident.String(), aObj)
		} else {
			panic(NewError(e.Pos().Sline(), KEYERROR, "IDENT"))
		}
	}
	return &Enum{Scope: enumScope}
}

func evalNamedFnStatement(namedFn *ast.NamedFunction, scope *Scope) Object {
	fnObj := evalFunctionLiteral(namedFn.FunctionLiteral, scope)
	scope.Set(namedFn.Ident.String(), fnObj) //save to scope

	return NIL
}

func evalFunctionLiteral(fl *ast.FunctionLiteral, scope *Scope) Object {
	fn := &Function{Literal: fl, Variadic: fl.Variadic, Scope: scope}

	if fl.Values != nil { //check for default values
		for _, item := range fl.Parameters {
			if _, ok := fl.Values[item.String()]; !ok { //if not has default value, then continue
				continue
			}
			val := Eval(fl.Values[item.String()], scope)
			if val.Type() != ERROR_OBJ {
				fn.Scope.Set(item.String(), val)
			} else {
				return val
			}
		}
	}

	return fn
}

// Prefix expressions, e.g. `!true, -5`
func evalPrefixExpression(p *ast.PrefixExpression, scope *Scope) Object {
	right := Eval(p.Right, scope)
	if right.Type() == ERROR_OBJ {
		return right
	}
	switch p.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		switch right.Type() {
		case INTEGER_OBJ:
			i := right.(*Integer)
			i.Int64 = -i.Int64
			return i
		case FLOAT_OBJ:
			f := right.(*Float)
			f.Float64 = -f.Float64
			return f
		}

	case "++":
		return evalIncrementPrefixOperatorExpression(p, right)
	case "--":
		return evalDecrementPrefixOperatorExpression(p, right)
	}
	panic(NewError(p.Pos().Sline(), PREFIXOP, p, right.Type()))
}

func evalIncrementPrefixOperatorExpression(p *ast.PrefixExpression, right Object) Object {
	switch right.Type() {
	case INTEGER_OBJ:
		rightObj := right.(*Integer)
		rightObj.Int64 = rightObj.Int64 + 1
		return NewInteger(rightObj.Int64)
	case FLOAT_OBJ:
		rightObj := right.(*Float)
		rightObj.Float64 = rightObj.Float64 + 1
		return NewFloat(rightObj.Float64)
	default:
		panic(NewError(p.Pos().Sline(), PREFIXOP, p.Operator, right.Type()))
	}
}

func evalDecrementPrefixOperatorExpression(p *ast.PrefixExpression, right Object) Object {
	switch right.Type() {
	case INTEGER_OBJ:
		rightObj := right.(*Integer)
		rightObj.Int64 = rightObj.Int64 - 1
		return NewInteger(rightObj.Int64)
	case FLOAT_OBJ:
		rightObj := right.(*Float)
		rightObj.Float64 = rightObj.Float64 - 1
		return NewFloat(rightObj.Float64)
	default:
		panic(NewError(p.Pos().Sline(), PREFIXOP, p.Operator, right.Type()))
	}
}

// Helper for evaluating Bang(!) expressions. Coerces truthyness based on object presence.
func evalBangOperatorExpression(right Object) Object {
	return nativeBoolToBooleanObject(!IsTrue(right))
}

// Evaluate infix expressions, e.g 1 + 2, a == 5, true == true, etc...
func evalInfixExpression(node *ast.InfixExpression, left, right Object) Object {
	_, leftIsNum := left.(Number)
	_, rightIsNum := right.(Number)
	//hasNumArg := leftIsNum || rightIsNum

	//Note :Here the 'switch's order is important, if you change the order, it will evaluate differently
	//e.g. 1 + [2,3] + "45" = [1,2,3,"45"](it's an array), if you change
	//`case (left.Type() == ARRAY_OBJ || right.Type() == ARRAY_OBJ)` to a lower order in the case, it will
	//return [1,2,3]"45"(that is a string)
	switch {
	case node.Operator == "and" || node.Operator == "&&":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) && objectToNativeBoolean(right))
	case node.Operator == "or" || node.Operator == "||":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) || objectToNativeBoolean(right))
	case leftIsNum && rightIsNum:
		return evalNumberInfixExpression(node, left, right)
	case (left.Type() == ARRAY_OBJ || right.Type() == ARRAY_OBJ):
		return evalArrayInfixExpression(node, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(node, left, right)
	case (left.Type() == STRING_OBJ || right.Type() == STRING_OBJ):
		return evalMixedTypeInfixExpression(node, left, right)
	case (left.Type() == HASH_OBJ && right.Type() == HASH_OBJ):
		return evalHashInfixExpression(node, left, right)
	case node.Operator == "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		//Here we need to special handling for `Boolean` object. Because most of the time `BOOLEAN` will
		//return TRUE and FALSE. But sometimes we have to returns a new `Boolean` object,
		//Here we need to compare `Boolean.Bool`，or else when we using
		//   if (aBool == true)
		//it will return false, but actually aBool is true.
		if left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ {
			l := left.(*Boolean)
			r := right.(*Boolean)
			if l.Bool == r.Bool {
				return TRUE
			}
			return FALSE
		}

		if left.Type() == NIL_OBJ && right.Type() == NIL_OBJ { //(s == nil) should return true if s is nil
			return TRUE
		}

		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ {
			l := left.(*Boolean)
			r := right.(*Boolean)
			if l.Bool != r.Bool {
				return TRUE
			}
			return FALSE
		}

		return nativeBoolToBooleanObject(left != right)
	}

	panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

func objectToNativeBoolean(o Object) bool {
	if r, ok := o.(*ReturnValue); ok {
		o = r.Value
	}
	switch obj := o.(type) {
	case *Boolean:
		return obj.Bool
	case *Nil:
		return false
	case *Integer:
		if obj.Int64 == 0 {
			return false
		}
		return true
	case *Float:
		if obj.Float64 == 0.0 {
			return false
		}
		return true
	case *Array:
		if len(obj.Members) == 0 {
			return false
		}
		return true
	case *Hash:
		if len(obj.Pairs) == 0 {
			return false
		}
		return true
	default:
		return true
	}
}

func evalNumberInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	var leftVal float64
	var rightVal float64
	isInt := left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ

	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*Integer).Int64)
	} else {
		leftVal = left.(*Float).Float64
	}

	if right.Type() == INTEGER_OBJ {
		rightVal = float64(right.(*Integer).Int64)
	} else {
		rightVal = right.(*Float).Float64
	}

	switch node.Operator {
	case "**":
		val := math.Pow(leftVal, rightVal)
		if isInt {
			return NewInteger(int64(val))
		}
		return NewFloat(val)
	case "&":
		if isInt {
			val := int64(leftVal) & int64(rightVal)
			return NewInteger(int64(val))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "|":
		if isInt {
			val := int64(leftVal) | int64(rightVal)
			return NewInteger(int64(val))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "^":
		if isInt {
			val := int64(leftVal) ^ int64(rightVal)
			return NewInteger(int64(val))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "+":
		val := leftVal + rightVal
		if isInt {
			return NewInteger(int64(val))
		}
		return NewFloat(val)
	case "-":
		val := leftVal - rightVal
		if isInt {
			return NewInteger(int64(val))
		}
		return NewFloat(val)
	case "*":
		val := leftVal * rightVal
		if isInt {
			return NewInteger(int64(val))
		}
		return NewFloat(val)
	case "/":
		if rightVal == 0 {
			panic(NewError(node.Pos().Sline(), DIVIDEBYZERO))
		}
		val := leftVal / rightVal
		//Should Always return float
		return NewFloat(val)
	case "%":
		if isInt {
			return NewInteger(int64(leftVal) % int64(rightVal))
		}
		return NewFloat(math.Mod(leftVal, rightVal))
	case ">>":
		if isInt {
			aRes := uint64(leftVal) >> uint64(rightVal)
			return NewInteger(int64(aRes)) //NOTE: CAST MAYBE NOT CORRECT
		}
	case "<<":
		if isInt {
			aRes := uint64(leftVal) << uint64(rightVal)
			return NewInteger(int64(aRes)) //NOTE: CAST MAYBE NOT CORRECT
		}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	}

	return NIL
}

func evalModuloExpression(left *Integer, right *Integer) Object {
	mod := left.Int64 % right.Int64
	if mod < 0 {
		mod += right.Int64
	}
	return NewInteger(mod)
}

func evalStringInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	l := left.(*String)
	r := right.(*String)

	switch node.Operator {
	case "=~": //match
		matched, _ := regexp.MatchString(r.String, l.String)
		if matched {
			return TRUE
		}
		return FALSE

	case "!~": //not match
		matched, _ := regexp.MatchString(r.String, l.String)
		if matched {
			return FALSE
		}
		return TRUE

	case "==":
		return nativeBoolToBooleanObject(l.String == r.String)
	case "!=":
		return nativeBoolToBooleanObject(l.String != r.String)
	case "+":
		return NewString(l.String + r.String)
	case "<":
		return nativeBoolToBooleanObject(l.String < r.String)
	case "<=":
		return nativeBoolToBooleanObject(l.String <= r.String)
	case ">":
		return nativeBoolToBooleanObject(l.String > r.String)
	case ">=":
		return nativeBoolToBooleanObject(l.String >= r.String)
	}
	panic(NewError(node.Pos().Sline(), INFIXOP, l.Type(), node.Operator, r.Type()))
}

func evalMixedTypeInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	switch node.Operator {
	case "+":
		return NewString(left.Inspect() + right.Inspect())
	case "*":
		if left.Type() == INTEGER_OBJ {
			integer := left.(*Integer).Int64
			return NewString(strings.Repeat(right.Inspect(), int(integer)))
		}
		if right.Type() == INTEGER_OBJ {
			integer := right.(*Integer).Int64
			return NewString(strings.Repeat(left.Inspect(), int(integer)))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left.Type() != STRING_OBJ || right.Type() != STRING_OBJ {
			panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
		}

		if left.(*String).String == right.(*String).String {
			return TRUE
		}
		return FALSE

	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() != STRING_OBJ || right.Type() != STRING_OBJ {
			panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
		}

		if left.(*String).String != right.(*String).String {
			return TRUE
		}
		return FALSE

	default:
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	}

	//panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

func evalArrayInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	switch node.Operator {
	case "+":
		if left.Type() == ARRAY_OBJ {
			leftVals := left.(*Array).Members

			if right.Type() == ARRAY_OBJ {
				rightVals := right.(*Array).Members
				leftVals = append(leftVals, rightVals...)
			} else {
				leftVals = append(leftVals, right)
			}
			return &Array{Members: leftVals}
		}

		//right is array
		rightVals := right.(*Array).Members
		if left.Type() == ARRAY_OBJ {
			leftVals := left.(*Array).Members
			rightVals = append(rightVals, leftVals...)
			return &Array{Members: rightVals}
		} else {
			ret := &Array{}
			ret.Members = append(ret.Members, left)
			ret.Members = append(ret.Members, rightVals...)
			return ret
		}

	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left.Type() != ARRAY_OBJ || right.Type() != ARRAY_OBJ {
			panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
		}

		leftVals := left.(*Array).Members
		rightVals := right.(*Array).Members
		if len(leftVals) != len(rightVals) {
			return FALSE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i])
			if !IsTrue(aBool) {
				return FALSE
			}
		}
		return TRUE
	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() != ARRAY_OBJ || right.Type() != ARRAY_OBJ {
			panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
		}
		leftVals := left.(*Array).Members
		rightVals := right.(*Array).Members
		if len(leftVals) != len(rightVals) {
			return TRUE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i])
			if IsTrue(aBool) {
				return TRUE
			}
		}
		return FALSE
	}
	panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

func evalHashInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	leftVals := left.(*Hash).Pairs
	rightVals := right.(*Hash).Pairs

	switch node.Operator {
	case "+":
		for k, v := range rightVals {
			leftVals[k] = v
		}

		return &Hash{Pairs: leftVals}
	case "==":
		return nativeBoolToBooleanObject(compareHashObj(leftVals, rightVals))
	case "!=":
		return nativeBoolToBooleanObject(!compareHashObj(leftVals, rightVals))
	}
	panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

func compareHashObj(left, right map[HashKey]HashPair) bool {
	if len(left) != len(right) {
		return false
	}

	found := 0
	for lk, lv := range left {
		for rk, rv := range right {
			if lk.Value == rk.Value && (lv.Key.Inspect() == rv.Key.Inspect() && lv.Value.Inspect() == rv.Value.Inspect()) {
				found += 1
				continue
			}
		}
	}

	return found == len(left)
}

// IF expressions, if (evaluates to boolean) True: { Block Statement } Optional Else: {Block Statement}
//func evalIfExpression(ie *ast.IfExpression, scope *Scope) Object {
//	condition := Eval(ie.Condition, scope)
//	if condition.Type() == ERROR_OBJ {
//		return condition
//	}
//
//	if IsTrue(condition) {
//		return evalBlockStatements(ie.Consequence.Statements, scope)
//	} else if ie.Alternative != nil {
//		return evalBlockStatements(ie.Alternative.Statements, scope)
//	}
//	return NIL
//}

func evalIfExpression(ie *ast.IfExpression, scope *Scope) Object {
	//eval "if/else-if" part
	for _, c := range ie.Conditions {
		condition := Eval(c.Cond, scope)
		if condition.Type() == ERROR_OBJ {
			return condition
		}

		if IsTrue(condition) {
			return evalBlockStatements(c.Block.Statements, scope)
		}
	}

	//eval "else" part
	if ie.Alternative != nil {
		return evalBlockStatements(ie.Alternative.Statements, scope)
	}

	return NIL
}

func evalDoLoopExpression(dl *ast.DoLoop, scope *Scope) Object {
	newScope := NewScope(scope)

	var e Object
	for {
		e = Eval(dl.Block, newScope)
		if e.Type() == ERROR_OBJ {
			return e
		}

		if _, ok := e.(*Break); ok {
			break
		}
		if _, ok := e.(*Continue); ok {
			continue
		}
		if v, ok := e.(*ReturnValue); ok {
			if v.Value != nil {
				return v.Value
			}
			break
		}
	}

	if e == nil || e.Type() == BREAK_OBJ || e.Type() == CONTINUE_OBJ {
		return NIL
	}
	return e
}

func evalWhileLoopExpression(wl *ast.WhileLoop, scope *Scope) Object {
	innerScope := NewScope(scope)

	condition := Eval(wl.Condition, innerScope)
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	var result Object
	for IsTrue(condition) {
		result = Eval(wl.Block, innerScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				return v.Value
			}
			break
		}
		condition = Eval(wl.Condition, innerScope)
		if condition.Type() == ERROR_OBJ {
			return condition
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return NIL
	}
	return result
}

func evalGrepExpression(ge *ast.GrepExpr, scope *Scope) Object {
	aValue := Eval(ge.Value, scope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	_, ok := aValue.(Listable) //must be listable
	if !ok {
		panic(NewError(ge.Pos().Sline(), NOTLISTABLE))
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	}

	result := &Array{}

	result.Members = []Object{}

	for _, item := range members {
		//Note: we must opening a new scope, because the variable is different in each iteration.
		//If not, then the next iteration will overwrite the previous assigned variable.
		newSubScope := NewScope(scope)
		newSubScope.Set(ge.Var, item)

		var cond Object

		if ge.Block != nil {
			cond = Eval(ge.Block, newSubScope)
		} else {
			cond = Eval(ge.Expr, newSubScope)
		}

		if IsTrue(cond) {
			result.Members = append(result.Members, item)
		}
	}
	return result
}

func evalMapExpression(me *ast.MapExpr, scope *Scope) Object {
	aValue := Eval(me.Value, scope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	_, ok := aValue.(Listable) //must be listable
	if !ok {
		panic(NewError(me.Pos().Sline(), NOTLISTABLE))
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	}

	result := &Array{}
	result.Members = []Object{}

	for _, item := range members {
		newSubScope := NewScope(scope)
		newSubScope.Set(me.Var, item)

		var r Object
		if me.Block != nil {
			r = Eval(me.Block, newSubScope)
		} else {
			r = Eval(me.Expr, newSubScope)
		}
		if r.Type() == ERROR_OBJ {
			return r
		}

		result.Members = append(result.Members, r)
	}
	return result
}

//[x+1 for x in arr <where cond>]
//[ str for str in strs <where cond>]
func evalListComprehension(lc *ast.ListComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)
	aValue := Eval(lc.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	_, ok := aValue.(Iterable) //must be listable
	if !ok {
		panic(NewError(lc.Pos().Sline(), NOTITERABLE))
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	}

	ret := &Array{}
	var result Object
	for idx, value := range members {
		newSubScope := NewScope(innerScope)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(lc.Var, value)
		if lc.Cond != nil {
			cond := Eval(lc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(lc.Expr, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		ret.Members = append(ret.Members, result)
	}

	return ret
}

//[x for x in a..b <where cond>]
//Almost same as evalForEachDotRangeExpression() function
func evalListRangeComprehension(lc *ast.ListRangeComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)

	startIdx := Eval(lc.StartIdx, innerScope)
	endIdx := Eval(lc.EndIdx, innerScope)

	var j int64
	arr := &Array{}

	switch startIdx.(type) {
	case *Integer:
		startVal := startIdx.(*Integer).Int64
		if endIdx.Type() != INTEGER_OBJ {
			panic(NewError(lc.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ, endIdx.Type()))
		}
		endVal := endIdx.(*Integer).Int64

		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		}
	case *String:
		startVal := startIdx.(*String).String
		if endIdx.Type() != STRING_OBJ {
			panic(NewError(lc.Pos().Sline(), RANGETYPEERROR, STRING_OBJ, endIdx.Type()))
		}
		endVal := endIdx.(*String).String

		//only support single character with lowercase
		alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"

		//convert to int for easy comparation
		leftByte := []int32(strings.ToLower(startVal))[0]
		rightByte := []int32(strings.ToLower(endVal))[0]
		if leftByte >= rightByte { // z -> a
			for i := len(alphabet) - 1; i >= 0; i-- {
				v := int32(alphabet[i])
				if v <= leftByte && v >= rightByte {
					arr.Members = append(arr.Members, NewString(string(v)))
				}
			}
		} else { // a -> z
			for _, v := range alphabet {
				if v >= leftByte && v <= rightByte {
					arr.Members = append(arr.Members, NewString(string(v)))
				}
			}
		}
	}

	ret := &Array{}
	var result Object
	for idx, value := range arr.Members {
		newSubScope := NewScope(innerScope)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(lc.Var, value)
		if lc.Cond != nil {
			cond := Eval(lc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(lc.Expr, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		ret.Members = append(ret.Members, result)
	}

	return ret
}

//[ expr for k,v in hash <where cond>]
func evalMapComprehension(mc *ast.MapComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)
	aValue := Eval(mc.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	_, ok := aValue.(Iterable) //must be listable
	if !ok {
		panic(NewError(mc.Pos().Sline(), NOTITERABLE))
	}

	//must be a *Hash, if not, panic
	hash, _ := aValue.(*Hash)

	ret := &Array{}
	var result Object
	for _, pair := range hash.Pairs {
		newSubScope := NewScope(innerScope)
		newSubScope.Set(mc.Key, pair.Key)
		newSubScope.Set(mc.Value, pair.Value)

		if mc.Cond != nil {
			cond := Eval(mc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(mc.Expr, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		ret.Members = append(ret.Members, result)
	}

	return ret
}

func evalCaseExpression(ce *ast.CaseExpr, scope *Scope) Object {
	rv := Eval(ce.Expr, scope) //case expression
	if rv.Type() == ERROR_OBJ {
		return rv
	}

	done := false
	var elseExpr *ast.CaseElseExpr
	for _, item := range ce.Matches {
		if cee, ok := item.(*ast.CaseElseExpr); ok {
			elseExpr = cee //cee: Case'Expr Else part
			continue
		}

		matchExpr := item.(*ast.CaseMatchExpr)
		matchRv := Eval(matchExpr.Expr, NewScope(scope)) //matcher expression
		if matchRv.Type() == ERROR_OBJ {
			return matchRv
		}

		//check 'rv' and 'matchRv' equality, if not equal, then continue
		if !equal(ce.IsWholeMatch, rv, matchRv) {
			continue
		}
		//Eval matcher block
		matcherScope := NewScope(scope)
		rv = Eval(matchExpr.Block, matcherScope)
		if rv.Type() == ERROR_OBJ {
			return rv
		}

		done = true
		break
	}

	if !done && elseExpr != nil {
		elseScope := NewScope(scope)
		rv = Eval(elseExpr.Block, elseScope)
		if rv.Type() == ERROR_OBJ {
			return rv
		}
	}
	return rv
}

func evalForLoopExpression(fl *ast.ForLoop, scope *Scope) Object { //fl:For Loop
	innerScope := NewScope(scope)

	init := Eval(fl.Init, innerScope)
	if init.Type() == ERROR_OBJ {
		return init
	}

	condition := Eval(fl.Cond, innerScope)
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	var result Object
	for IsTrue(condition) {
		newSubScope := NewScope(innerScope)
		result = Eval(fl.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			newVal := Eval(fl.Update, newSubScope) //Before continue, we need to call 'Update' and 'Cond'
			if newVal.Type() == ERROR_OBJ {
				return newVal
			}

			condition = Eval(fl.Cond, newSubScope)
			if condition.Type() == ERROR_OBJ {
				return condition
			}

			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				return v.Value
			}
			break
		}
		newVal := Eval(fl.Update, newSubScope)
		if newVal.Type() == ERROR_OBJ {
			return newVal
		}

		condition = Eval(fl.Cond, newSubScope)
		if condition.Type() == ERROR_OBJ {
			return condition
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return NIL
	}
	return result
}

func evalForEverLoopExpression(fel *ast.ForEverLoop, scope *Scope) Object {
	var e Object
	for {
		e = Eval(fel.Block, NewScope(scope))
		if e.Type() == ERROR_OBJ {
			return e
		}

		if _, ok := e.(*Break); ok {
			break
		}
		if _, ok := e.(*Continue); ok {
			continue
		}
		if v, ok := e.(*ReturnValue); ok {
			if v.Value != nil {
				return v.Value
			}
			break
		}
	}

	if e == nil || e.Type() == BREAK_OBJ || e.Type() == CONTINUE_OBJ {
		return NIL
	}
	return e
}

func evalForEachArrayExpression(fal *ast.ForEachArrayLoop, scope *Scope) Object { //fal:For Array Loop
	innerScope := NewScope(scope)

	aValue := Eval(fal.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	_, ok := aValue.(Iterable)
	if !ok {
		panic(NewError(fal.Pos().Sline(), NOTITERABLE))
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	}

	ret := &Array{}
	var result Object
	for idx, value := range members {
		newSubScope := NewScope(innerScope)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(fal.Var, value)
		if fal.Cond != nil {
			cond := Eval(fal.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fal.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {

			if v.Value != nil {
				ret.Members = append(ret.Members, v.Value)
				//return v.Value
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	//Here we need to check `nil`, because if the initial condition is not true, then `for`'s Body will have no
	//chance to execute, the result will be nil
	//Because is the reason why we need to check for `BREAK_OBJ` or `CONTINUE_OBJ`:
	//    for i in 5..1 where i > 2 {
	//      if (i == 3) { continue }
	//      putln('i={i}')
	//    }
	//They will output "continue", this is not we expected
	//A LONG TIME HIDDEN BUG!
	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

func evalForEachArrayWithIndex(fml *ast.ForEachMapLoop, val Object, scope *Scope) Object {
	var members []Object
	if val.Type() == STRING_OBJ {
		aStr, _ := val.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if val.Type() == ARRAY_OBJ {
		arr, _ := val.(*Array)
		members = arr.Members
	}

	ret := &Array{}
	var result Object
	for idx, value := range members {
		newSubScope := NewScope(scope)
		newSubScope.Set(fml.Key, NewInteger(int64(idx)))
		newSubScope.Set(fml.Value, value)
		if fml.Cond != nil {
			cond := Eval(fml.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fml.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {

			if v.Value != nil {
				ret.Members = append(ret.Members, v.Value)
				//return v.Value
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

func evalForEachMapExpression(fml *ast.ForEachMapLoop, scope *Scope) Object { //fml:For Map Loop
	innerScope := NewScope(scope)

	aValue := Eval(fml.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	_, ok := aValue.(Iterable)
	if !ok {
		panic(NewError(fml.Pos().Sline(), NOTITERABLE))
	}

	//for index, value in arr
	//for index, value in string
	if aValue.Type() == STRING_OBJ || aValue.Type() == ARRAY_OBJ {
		return evalForEachArrayWithIndex(fml, aValue, innerScope)
	}

	hash, _ := aValue.(*Hash)

	ret := &Array{}
	var result Object
	for _, pair := range hash.Pairs {
		newSubScope := NewScope(innerScope)
		newSubScope.Set(fml.Key, pair.Key)
		newSubScope.Set(fml.Value, pair.Value)

		if fml.Cond != nil {
			cond := Eval(fml.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fml.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				ret.Members = append(ret.Members, v.Value)
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

func evalForEachDotRangeExpression(fdr *ast.ForEachDotRange, scope *Scope) Object { //fdr:For Dot Range
	innerScope := NewScope(scope)

	startIdx := Eval(fdr.StartIdx, innerScope)
	//	if startIdx.Type() != INTEGER_OBJ {
	//		panic(NewError(fdr.Pos().Sline(), RANGETYPEERROR, startIdx.Type()))
	//	}

	endIdx := Eval(fdr.EndIdx, innerScope)
	//	if startIdx.Type() != INTEGER_OBJ {
	//		panic(NewError(fdr.Pos().Sline(), RANGETYPEERROR, endIdx.Type()))
	//	}

	var j int64
	arr := &Array{}

	switch startIdx.(type) {
	case *Integer:
		startVal := startIdx.(*Integer).Int64
		if endIdx.Type() != INTEGER_OBJ {
			panic(NewError(fdr.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ, endIdx.Type()))
		}
		endVal := endIdx.(*Integer).Int64

		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		}
	case *String:
		startVal := startIdx.(*String).String
		if endIdx.Type() != STRING_OBJ {
			panic(NewError(fdr.Pos().Sline(), RANGETYPEERROR, STRING_OBJ, endIdx.Type()))
		}
		endVal := endIdx.(*String).String

		//only support single character with lowercase
		alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"

		//convert to int for easy comparation
		leftByte := []int32(strings.ToLower(startVal))[0]
		rightByte := []int32(strings.ToLower(endVal))[0]
		if leftByte >= rightByte { // z -> a
			for i := len(alphabet) - 1; i >= 0; i-- {
				v := int32(alphabet[i])
				if v <= leftByte && v >= rightByte {
					arr.Members = append(arr.Members, NewString(string(v)))
				}
			}
		} else { // a -> z
			for _, v := range alphabet {
				if v >= leftByte && v <= rightByte {
					arr.Members = append(arr.Members, NewString(string(v)))
				}
			}
		}
	}

	ret := &Array{}
	var result Object
	for idx, value := range arr.Members {
		newSubScope := NewScope(innerScope)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(fdr.Var, value)
		if fdr.Cond != nil {
			cond := Eval(fdr.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fdr.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				arr.Members = append(arr.Members, v.Value)
			}
			break
		} else {
			arr.Members = append(arr.Members, result)
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

// Helper function IsTrue for IF evaluation - neccessity is dubious
func IsTrue(obj Object) bool {
	if b, ok := obj.(*Boolean); ok { //if it is a Boolean Object
		return b.Bool
	}

	if _, ok := obj.(*Nil); ok { //if it is a Nil Object
		return false
	}

	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		switch obj.Type() {
		case INTEGER_OBJ:
			if obj.(*Integer).Int64 == 0 {
				return false
			}
		case FLOAT_OBJ:
			if obj.(*Float).Float64 == 0.0 {
				return false
			}
		case STRING_OBJ:
			if len(obj.(*String).String) == 0 {
				return false
			}
		case ARRAY_OBJ:
			if len(obj.(*Array).Members) == 0 {
				return false
			}
		case HASH_OBJ:
			if len(obj.(*Hash).Pairs) == 0 {
				return false
			}
		}
		return true
	}
}

// Block Statement Evaluation - The innards of both IF and Function calls
// very similar to parseProgram, but because we need to leave the return
// value wrapped in it's Object, it remains, for now.
func evalBlockStatements(block []ast.Statement, scope *Scope) (results Object) {
	for _, statement := range block {
		results = Eval(statement, scope)
		if results.Type() == ERROR_OBJ {
			return results
		}

		if results != nil && results.Type() == RETURN_VALUE_OBJ {
			return
		}
		if _, ok := results.(*Break); ok {
			return
		}
		if _, ok := results.(*Continue); ok {
			return
		}
	}
	return //do not return NIL, becuase we have already set the 'results'
}

func evalTryBlockStatements(block []ast.Statement, scope *Scope) (results Object) {
	for _, statement := range block {
		results = Eval(statement, scope)
		if isTryError(results) {
			return results
		}

		if results != nil && results.Type() == RETURN_VALUE_OBJ {
			return
		}
		if _, ok := results.(*Break); ok {
			return
		}
		if _, ok := results.(*Continue); ok {
			return
		}
	}
	return //do not return NIL, becuase we have already set the 'results'
}

// Eval when a function is _called_, includes fn literal evaluation and calling builtins
func evalFunctionCall(call *ast.CallExpression, scope *Scope) Object {
	fn, ok := scope.Get(call.Function.String())
	if !ok {

		if f, ok := call.Function.(*ast.FunctionLiteral); ok {
			//let add =fn(x,y) { x+y }
			//add(2,3)
			fn = &Function{Literal: f, Scope: scope, Variadic: f.Variadic}
			scope.Set(call.Function.String(), fn)
		} else if idxExpr, ok := call.Function.(*ast.IndexExpression); ok { //index expression
			//let complex={ "add"=>fn(x,y){ x+y } }
			//complex["add"](2,3)
			aValue := Eval(idxExpr, scope)
			if aValue.Type() == ERROR_OBJ {
				return aValue
			}

			if aFn, ok := aValue.(*Function); ok { //index expression
				fn = aFn
			} else {
				panic(NewError(call.Function.Pos().Sline(), UNKNOWNIDENT, call.Function.String()))
			}
		} else if builtin, ok := builtins[call.Function.String()]; ok {
			//let complex={ "add"=>fn(x,y){ fn(z) {x+y+z} } }
			//complex["add"](2,3)(4)
			args := evalArgs(call.Arguments, scope)
			return builtin.Fn(call.Function.Pos().Sline(), args...)
		} else if callExpr, ok := call.Function.(*ast.CallExpression); ok { //call expression
			aValue := Eval(callExpr, scope)
			if aValue.Type() == ERROR_OBJ {
				return aValue
			}

			fn = aValue
		} else {
			panic(NewError(call.Function.Pos().Sline(), UNKNOWNIDENT, call.Function.String()))
		}
	}
	f := fn.(*Function)
	newScope := NewScope(f.Scope)

	//Register this function call in the call stack
	newScope.CallStack.Frames = append(newScope.CallStack.Frames, CallFrame{FuncScope: newScope, CurrentCall: call})

	//Using golang's defer mechanism, before function return, call current frame's defer method
	defer func() {
		frame := newScope.CurrentFrame()
		if len(frame.defers) != 0 {
			frame.runDefers(newScope)
		}

		//After run, must pop the frame
		stack := newScope.CallStack
		stack.Frames = stack.Frames[0 : len(stack.Frames)-1]
	}()

	variadicParam := []Object{}
	args := evalArgs(call.Arguments, scope)
	for i, _ := range call.Arguments {
		//Because of function default values, we need to check `i >= len(args)`
		if f.Variadic && i >= len(f.Literal.Parameters)-1 {
			for j := i; j < len(args); j++ {
				variadicParam = append(variadicParam, args[j])
			}
			break
		} else if i >= len(f.Literal.Parameters) {
			break
		} else {
			newScope.Set(f.Literal.Parameters[i].String(), args[i])
		}
	}

	// Variadic argument is passed as a single array
	// of parameters.
	if f.Variadic {
		newScope.Set(f.Literal.Parameters[len(f.Literal.Parameters)-1].String(), &Array{Members: variadicParam})
		if len(call.Arguments) < len(f.Literal.Parameters) {
			f.Scope.Set("@_", NewInteger(int64(len(f.Literal.Parameters)-1)))
		} else {
			f.Scope.Set("@_", NewInteger(int64(len(call.Arguments))))
		}
	} else {
		f.Scope.Set("@_", NewInteger(int64(len(f.Literal.Parameters))))
	}

	r := Eval(f.Literal.Body, newScope)
	if r.Type() == ERROR_OBJ {
		return r
	}

	if obj, ok := r.(*ReturnValue); ok {
		return obj.Value
	}
	return r
}

// Method calls for builtin Objects
func evalMethodCallExpression(call *ast.MethodCallExpression, scope *Scope) Object {
	//First check if is a stanard library object
	str := call.Object.String()
	if obj, ok := GetGlobalObj(str); ok {
		switch call.Call.(type) {
		case *ast.Identifier:
			if i, ok := GetGlobalObj(str + "." + call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			if method, ok := call.Call.(*ast.CallExpression); ok {
				args := evalArgs(method.Arguments, scope)
				return obj.CallMethod(call.Call.Pos().Sline(), scope, method.Function.String(), args...)
			}
		}
	}

	obj := Eval(call.Object, scope)
	if obj.Type() == ERROR_OBJ {
		return obj
	}
	switch m := obj.(type) {
	case *IncludedObject:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			if o.Function.String() == "Scope" {
				return obj.CallMethod(call.Call.Pos().Sline(), m.Scope, "Scope")
			}
			return evalFunctionCall(o, m.Scope)
		}
	case *Struct:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			args := evalArgs(o.Arguments, scope)
			return obj.CallMethod(call.Call.Pos().Sline(), m.Scope, o.Function.String(), args...)
		}
	case *Enum:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			args := evalArgs(o.Arguments, scope)
			return obj.CallMethod(call.Call.Pos().Sline(), m.Scope, o.Function.String(), args...)
		}
	default:
		if method, ok := call.Call.(*ast.CallExpression); ok {
			args := evalArgs(method.Arguments, scope)
			return obj.CallMethod(call.Call.Pos().Sline(), scope, method.Function.String(), args...)
		}
	}

	panic(NewError(call.Call.Pos().Sline(), NOMETHODERROR, call.String(), obj.Type()))

}

func evalArgs(args []ast.Expression, scope *Scope) []Object {
	//TODO: Refactor this to accept the params and args, go ahead and
	// update scope while looping and return the Scope object.
	e := []Object{}
	for _, v := range args {
		item := Eval(v, scope)
		e = append(e, item)
	}
	return e
}

// Index Expressions, i.e. array[0], array[2:4] or hash["mykey"]
func evalIndexExpression(ie *ast.IndexExpression, scope *Scope) Object {
	left := Eval(ie.Left, scope)

	switch iterable := left.(type) {
	case *Array:
		return evalArrayIndex(iterable, ie, scope)
	case *Hash:
		return evalHashKeyIndex(iterable, ie, scope)
	case *String:
		return evalStringIndex(iterable, ie, scope)
	}
	panic(NewError(ie.Pos().Sline(), NOINDEXERROR, left.Type()))
}

func evalStringIndex(str *String, ie *ast.IndexExpression, scope *Scope) Object {
	var idx int64
	length := int64(utf8.RuneCountInString(str.String))
	//length := int64(len(str.String))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalStringSliceExpression(str, exp, scope)
	}
	index := Eval(ie.Index, scope)
	if index.Type() == ERROR_OBJ {
		return index
	}
	idx = index.(*Integer).Int64
	if idx >= length || idx < 0 {
		panic(NewError(ie.Pos().Sline(), INDEXERROR, idx))
	}

	return NewString(string([]rune(str.String)[idx])) //support utf8,not very efficient
	//return &String{String: string(str.String[idx]), Valid:true}  //only support ASCII
}

func evalStringSliceExpression(str *String, se *ast.SliceExpression, scope *Scope) Object {
	var idx int64
	var slice int64

	length := int64(utf8.RuneCountInString(str.String))
	//length := int64(len(str.String))

	startIdx := Eval(se.StartIndex, scope)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}
	idx = startIdx.(*Integer).Int64
	if idx >= length || idx < 0 {
		panic(NewError(se.Pos().Sline(), INDEXERROR, idx))
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, scope)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}
		slice = slIndex.(*Integer).Int64
		if slice >= (length + 1) {
			panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
		}
		if slice < 0 {
			panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
		}
	}
	if idx == 0 && slice == length {
		return str
	}

	if slice < idx {
		panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
	}

	runes := []rune(str.String)
	return NewString(string(runes[idx:slice]))
}

func evalHashKeyIndex(hash *Hash, ie *ast.IndexExpression, scope *Scope) Object {
	key := Eval(ie.Index, scope)
	if key.Type() == ERROR_OBJ {
		return key
	}
	hashable, ok := key.(Hashable)
	if !ok {
		panic(NewError(ie.Pos().Sline(), KEYERROR, key.Type()))
	}
	hashPair, ok := hash.Pairs[hashable.HashKey()]
	// TODO: should we return an error here? If not, maybe arrays should return NIL as well?
	if !ok {
		return NIL
	}
	return hashPair.Value
}

func evalArraySliceExpression(array *Array, se *ast.SliceExpression, scope *Scope) Object {
	var idx int64
	var slice int64
	length := int64(len(array.Members))

	startIdx := Eval(se.StartIndex, scope)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}
	idx = startIdx.(*Integer).Int64
	if idx < 0 {
		panic(NewError(se.Pos().Sline(), INDEXERROR, idx))
	}

	if idx >= length {
		return NIL
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, scope)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}
		slice = slIndex.(*Integer).Int64
		if slice >= (length+1) || slice < 0 {
			panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
		}
	}
	if idx == 0 && slice == length {
		return array
	}

	if slice < idx {
		panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
	}

	if slice == length {
		return &Array{Members: array.Members[idx:]}
	}
	return &Array{Members: array.Members[idx:slice]}
}

func evalArrayIndex(array *Array, ie *ast.IndexExpression, scope *Scope) Object {
	var idx int64
	length := int64(len(array.Members))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalArraySliceExpression(array, exp, scope)
	}
	index := Eval(ie.Index, scope)
	if index.Type() == ERROR_OBJ {
		return index
	}
	idx = index.(*Integer).Int64
	if idx < 0 {
		panic(NewError(ie.Pos().Sline(), INDEXERROR, idx))
	}
	if idx >= length {
		return NIL
	}
	return array.Members[idx]
}

func evalPostfixExpression(left Object, node *ast.PostfixExpression) Object {
	switch node.Operator {
	case "++":
		return evalIncrementPostfixOperatorExpression(node, left)
	case "--":
		return evalDecrementPostfixOperatorExpression(node, left)
	default:
		panic(NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type()))
	}
}

func evalIncrementPostfixOperatorExpression(node *ast.PostfixExpression, left Object) Object {
	switch left.Type() {
	case INTEGER_OBJ:
		leftObj := left.(*Integer)
		returnVal := NewInteger(leftObj.Int64)
		leftObj.Int64 = leftObj.Int64 + 1
		return returnVal
	case FLOAT_OBJ:
		leftObj := left.(*Float)
		returnVal := NewFloat(leftObj.Float64)
		leftObj.Float64 = leftObj.Float64 + 1
		return returnVal
	default:
		panic(NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type()))
	}
}

func evalDecrementPostfixOperatorExpression(node *ast.PostfixExpression, left Object) Object {
	switch left.Type() {
	case INTEGER_OBJ:
		leftObj := left.(*Integer)
		returnVal := NewInteger(leftObj.Int64)
		leftObj.Int64 = leftObj.Int64 - 1
		return returnVal
	case FLOAT_OBJ:
		leftObj := left.(*Float)
		returnVal := NewFloat(leftObj.Float64)
		leftObj.Float64 = leftObj.Float64 - 1
		return returnVal
	default:
		panic(NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type()))
	}
}

func evalTryStatement(ts *ast.TryStmt, scope *Scope) Object {
	tryScope := NewScope(scope)

	rv := evalTryBlockStatements(ts.Block.Statements, tryScope) //try statement
	if isTryError(rv) {
		rvObj := rv
		var rvObjStr string

		if rv.Type() == ERROR_OBJ {
			rvObjStr = rv.(*Error).Message
		} else if rv.Type() == NIL_OBJ {
			rvObjStr = rv.(*Nil).OptionalMsg
		} else if rv.Type() == BOOLEAN_OBJ {
			rvObjStr = rv.(*Boolean).OptionalMsg
		}

		done := false
		var catchAllStmt *ast.CatchAllStmt

		catchScope := NewScope(scope)
		for _, item := range ts.Catches {
			catchSubScope := NewScope(catchScope)
			if cas, ok := item.(*ast.CatchAllStmt); ok {
				catchAllStmt = cas //cas: Catch All Statement
				continue
			}

			cs := item.(*ast.CatchStmt)
			if cs.VarType == 1 { //IDENTIFIER
				val, ok1 := catchSubScope.Get(cs.Var)
				if !ok1 {
					catchSubScope.Set(cs.Var, rvObj) //put it to scope
					cs.Var = rvObjStr
				} else {
					valStr, ok2 := val.(*String)
					if !ok2 {
						panic(NewError(cs.Pos().Sline(), THROWERROR))
					}
					cs.Var = valStr.String
				}
			}

			if cs.Var != rvObjStr {
				continue
			}

			catchSubScope.Set(cs.Var, rvObj)

			rv = evalTryBlockStatements(cs.Block.Statements, catchSubScope) //catch Block
			if rv.Type() == ERROR_OBJ {
				return rv
			}

			done = true
			break
		} //end for

		if !done && catchAllStmt != nil {
			rv = evalTryBlockStatements(catchAllStmt.Block.Statements, NewScope(scope))
			if rv.Type() == ERROR_OBJ {
				return rv
			}
		}
	} //end if

	if ts.Finally != nil { //finally
		finalScope := NewScope(scope)
		rv = evalTryBlockStatements(ts.Finally.Statements, finalScope)
		if rv.Type() == ERROR_OBJ {
			return rv
		}
	}

	return rv
}

//Evaluate ternary expression
func evalTernaryExpression(te *ast.TernaryExpression, scope *Scope) Object {
	condition := Eval(te.Condition, scope) //eval condition
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	if IsTrue(condition) {
		return Eval(te.IfTrue, scope)
	} else {
		return Eval(te.IfFalse, scope)
	}
}

func evalSpawnStatment(s *ast.SpawnStmt, scope *Scope) Object {
	newSpawnScope := NewScope(scope)

	switch callExp := s.Call.(type) {
	case *ast.CallExpression:
		go (func() {
			evalFunctionCall(callExp, newSpawnScope)
		})()
	case *ast.MethodCallExpression:
		go (func() {
			evalMethodCallExpression(callExp, newSpawnScope)
		})()
	default:
		panic(NewError(s.Pos().Sline(), SPAWNERROR))
	}

	return NIL
}

func evalPipeExpression(p *ast.Pipe, scope *Scope) Object {
	left := Eval(p.Left, scope)

	// Convert the type object back to an expression
	// so it can be passed to the FunctionCall arguments.
	argument := obj2Expression(left)
	if argument == nil {
		return NIL
	}

	// The right side operator should be a function.
	switch rightFunc := p.Right.(type) {
	case *ast.MethodCallExpression:
		// Prepend the left-hand interpreted value
		// to the function arguments.
		rightFunc.Call.(*ast.CallExpression).Arguments = append([]ast.Expression{argument}, rightFunc.Call.(*ast.CallExpression).Arguments...)
		return Eval(rightFunc, scope)
	case *ast.CallExpression:
		rightFunc.Arguments = append([]ast.Expression{argument}, rightFunc.Arguments...)
		return Eval(rightFunc, scope)

	}

	return NIL
}

// Convert a Object to an ast.Expression.
func obj2Expression(obj Object) ast.Expression {
	switch value := obj.(type) {
	case *Boolean:
		return &ast.Boolean{Value: value.Bool}
	case *Integer:
		return &ast.IntegerLiteral{Value: value.Int64}
	case *Float:
		return &ast.FloatLiteral{Value: value.Float64}
	case *String:
		return &ast.StringLiteral{Value: value.String}
	case *Nil:
		return &ast.NilLiteral{}
	case *Array:
		array := &ast.ArrayLiteral{}
		for _, v := range value.Members {
			result := obj2Expression(v)
			if result == nil {
				return nil
			}
			array.Members = append(array.Members, result)
		}
		return array
	case *Hash:
		hash := &ast.HashLiteral{}
		hash.Pairs = make(map[ast.Expression]ast.Expression)
		for _, v := range value.Pairs {
			key := &ast.StringLiteral{Value: v.Key.Inspect()}
			result := obj2Expression(v.Value)
			if result == nil {
				return nil
			}
			hash.Pairs[key] = result
		}
		return hash
	}

	return nil
}

//Returns true when lhsV and rhsV is same value.
func equal(isWholeMatch bool, lhsV, rhsV Object) bool {
	if lhsV == nil && rhsV == nil {
		return true
	}
	if (lhsV == nil && rhsV != nil) || (lhsV != nil && rhsV == nil) {
		return false
	}

	if lhsV.Type() != rhsV.Type() {
		return false
	}

	if lhsV.Type() == NIL_OBJ {
		if rhsV.Type() == NIL_OBJ {
			return true
		}
	}

	if lhsV.Type() == STRING_OBJ && rhsV.Type() == STRING_OBJ {
		leftStr := lhsV.(*String).String
		rightStr := rhsV.(*String).String

		if isWholeMatch {
			r := reflect.DeepEqual(lhsV, rhsV)
			if r {
				return true
			} else {
				return false
			}
		} else {
			matched, _ := regexp.MatchString(rightStr, leftStr)
			return matched
		}
	} else {
		r := reflect.DeepEqual(lhsV, rhsV)
		if r {
			return true
		} else {
			return false
		}
	}

	return false
}

func isTryError(o Object) bool {
	if o.Type() == ERROR_OBJ ||
		(o.Type() == NIL_OBJ && o.(*Nil).OptionalMsg != "") ||
		(o.Type() == BOOLEAN_OBJ && o.(*Boolean).Bool == false && o.(*Boolean).OptionalMsg != "") {
		return true
	}
	return false
}
