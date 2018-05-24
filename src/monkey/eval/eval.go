package eval

import (
	"bytes"
	"fmt"
	"math"
	"monkey/ast"
	"monkey/token"
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
var currentInstance *ObjectInstance

var mux sync.Mutex
var classMux sync.Mutex

//REPL with color support
var REPLColor bool

func Eval(node ast.Node, scope *Scope) (val Object) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case *Error:
				//if panic is a Error Object, print its contents
				fmt.Fprintf(os.Stderr, "\x1b[31m%s\x1b[0m\n", r.Error())
				//debug.PrintStack() //debug only

				//WHY return NIL? if we do not return 'NIL', we may get something like below:
				//    PANIC=runtime error: invalid memory address or nil pointer
				val = NIL
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
				val = NIL
				return
			}
		}
	}()

	//fmt.Printf("node.Type=%T, node=<%s>\n", node, node.String()) //debugging
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
		return evalReturnStatement(node, scope)
	case *ast.DeferStmt:
		return evalDeferStatement(node, scope)
	case *ast.FunctionStatement:
		return evalFunctionStatement(node, scope)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.UIntegerLiteral:
		return evalUIntegerLiteral(node)
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
	case *ast.TupleLiteral:
		return evalTupleLiteral(node, scope)
	case *ast.HashLiteral:
		return evalHashLiteral(node, scope)
	case *ast.StructLiteral:
		return evalStructLiteral(node, scope)
	case *ast.EnumLiteral:
		return evalEnumLiteral(node, scope)
	case *ast.EnumStatement:
		return evalEnumStatement(node, scope)
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
		return evalInfixExpression(node, left, right, scope)
	case *ast.PostfixExpression:
		left := Eval(node.Left, scope)
		if left.Type() == ERROR_OBJ {
			return left
		}
		return evalPostfixExpression(left, node)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.UnlessExpression:
		return evalUnlessExpression(node, scope)
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
	case *ast.ListMapComprehension:
		return evalListMapComprehension(node, scope)
	case *ast.HashComprehension:
		return evalHashComprehension(node, scope)
	case *ast.HashRangeComprehension:
		return evalHashRangeComprehension(node, scope)
	case *ast.HashMapComprehension:
		return evalHashMapComprehension(node, scope)
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
		return evalSpawnStatement(node, scope)
	case *ast.NilLiteral:
		return NIL
	case *ast.Pipe:
		return evalPipeExpression(node, scope)

	//Class related
	case *ast.ClassStatement:
		return evalClassStatement(node, scope)
	case *ast.ClassLiteral:
		return evalClassLiteral(node, scope)
	case *ast.NewExpression:
		return evalNewExpression(node, scope)
	//using
	case *ast.UsingStmt:
		return evalUsingStatement(node, scope)
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
	if l.DestructingFlag {
		v := Eval(l.Values[0], scope)
		valType := v.Type()
		switch valType {
		case HASH_OBJ:
			h := v.(*Hash)
			for _, item := range l.Names {
				if item.Token.Type == token.UNDERSCORE {
					continue
				}
				found := false
				for _, pair := range h.Pairs{
					if item.String() == pair.Key.Inspect() {
						val = pair.Value
						scope.Set(item.String(), pair.Value)
						found = true
					}
				}
				if !found {
					val = NIL
					scope.Set(item.String(), val)
				}
			}

		case ARRAY_OBJ:
			arr := v.(*Array)
			valuesLen := len(arr.Members)
			for idx, item := range l.Names {
				if idx >= valuesLen { //There are more Names than Values
					if item.Token.Type != token.UNDERSCORE {
						val = NIL
						scope.Set(item.String(), val)
					}
				} else {
					if item.Token.Type == token.UNDERSCORE {
						continue
					}
					val = arr.Members[idx]
					if val.Type() != ERROR_OBJ {
						scope.Set(item.String(), val)
					} else {
						return
					}
				}
			}

		case TUPLE_OBJ:
			tup := v.(*Tuple)
			valuesLen := len(tup.Members)
			for idx, item := range l.Names {
				if idx >= valuesLen { //There are more Names than Values
					if item.Token.Type != token.UNDERSCORE {
						val = NIL
						scope.Set(item.String(), val)
					}
					val = NIL
					scope.Set(item.String(), val)
				} else {
					if item.Token.Type == token.UNDERSCORE {
						continue
					}
					val = tup.Members[idx]
					if val.Type() != ERROR_OBJ {
						scope.Set(item.String(), val)
					} else {
						return
					}
				}
			}

		default:
			panic(NewError(l.Pos().Sline(), GENERICERROR, "Only Array|Tuple|Hash is allowed!"))
		}
	
		return
	}

	values := []Object{}
	valuesLen := 0
	for _, value := range l.Values {
		val := Eval(value, scope)
		if val.Type() == TUPLE_OBJ {
			tupleObj := val.(*Tuple)
			if tupleObj.IsMulti {
				valuesLen += len(tupleObj.Members)
				for _, tupleItem := range tupleObj.Members {
					values = append(values, tupleItem)
				}
			} else {
				valuesLen += 1
				values = append(values, tupleObj)
			}
			
		} else {
			valuesLen += 1
			values = append(values, val)
		}
	}

	for idx, item := range l.Names {
		if idx >= valuesLen { //There are more Names than Values
			if item.Token.Type != token.UNDERSCORE {
				val = NIL
				scope.Set(item.String(), val)
			}
		} else {
			if item.Token.Type == token.UNDERSCORE {
				continue
			}
			val = values[idx];
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
	isUInt := left.Type() == UINTEGER_OBJ && val.Type() == UINTEGER_OBJ

	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*Integer).Int64)
	} else if left.Type() == UINTEGER_OBJ {
		leftVal = float64(left.(*UInteger).UInt64)
	} else {
		leftVal = left.(*Float).Float64
	}

	//Check `right`'s type
	if val.Type() == INTEGER_OBJ {
		rightVal = float64(val.(*Integer).Int64)
	} else if val.Type() == UINTEGER_OBJ {
		rightVal = float64(val.(*UInteger).UInt64)
	} else {
		rightVal = val.(*Float).Float64
	}

	var ok bool
	switch a.Token.Literal {
	case "+=":
		result := leftVal + rightVal
		if isInt { //only 'INTEGER + INTEGER'
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
			if ok {
				return
			}
		} else if isUInt { //only 'UINTEGER + UINTEGER'
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
			if ok {
				return
			}
		} else {
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "-=":
		result := leftVal-rightVal
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
			if ok {
				return
			}
		} else {
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "*=":
		result := leftVal*rightVal
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
			if ok {
				return
			}
		} else {
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "/=":
		if rightVal == 0 {
			panic(NewError(a.Pos().Sline(), DIVIDEBYZERO))
		}

		result := leftVal/rightVal
		//Always return Float
		ret, ok = scope.Reset(name, NewFloat(result))
		if ok {
			return
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))

	case "%=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)%int64(rightVal)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)%uint64(rightVal)))
			if ok {
				return
			}
		} else {
			result := math.Mod(leftVal, rightVal)
			ret, ok = checkNumAssign(scope, name, left, val, result)
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
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)^uint64(rightVal)))
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
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)|uint64(rightVal)))
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
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)&uint64(rightVal)))
			if ok {
				return
			}
		}
		panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
	}
	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
}

func checkNumAssign(scope *Scope, name string, left Object, right Object, result float64) (ret Object, ok bool) {
	if (left.Type() == FLOAT_OBJ || right.Type() == FLOAT_OBJ) {
		ret, ok = scope.Reset(name, NewFloat(result))
		return
	}

	if (left.Type() == INTEGER_OBJ && right.Type() == UINTEGER_OBJ) ||
		(left.Type() == UINTEGER_OBJ && right.Type() == INTEGER_OBJ) {
		if result > math.MaxInt64 {
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
		} else {
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
		}
	}
	return
}

//str[idx] = item
//str += item
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

			var idx int64
			switch o := index.(type) {
			case *Integer:
				idx = o.Int64
			case *UInteger:
				idx = int64(o.UInt64)
			}
			
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

//array[idx] = item
//array += item
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

			var idx int64
			switch o := index.(type) {
			case *Integer:
				idx = o.Int64
			case *UInteger:
				idx = int64(o.UInt64)
			}
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


func evalTupleAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	//Tuple is an immutable sequence of values
	if a.Token.Literal == "=" { //tuple[idx] = item
		str := fmt.Sprintf("%s[IDX]", TUPLE_OBJ)
		panic(NewError(a.Pos().Sline(), INFIXOP, str, a.Token.Literal, val.Type()))
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
		case *ast.IndexExpression: //hashObj[key] = val
			key := Eval(nodeType.Index, scope)
			hashable, ok := key.(Hashable)
			if !ok {
				panic(NewError(a.Pos().Sline(), KEYERROR, val.Type()))
			}
			leftVals[hashable.HashKey()] = HashPair{Key: key, Value: val}
			return left
		case *ast.Identifier: //hashObj.key = val
			key := strings.Split(a.Name.String(), ".")[1]
			keyObj := NewString(key)
			leftVals[keyObj.HashKey()] = HashPair{Key: keyObj, Value: val}
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
	case INTEGER_OBJ, UINTEGER_OBJ, FLOAT_OBJ:
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

//instanceObj[x] = xxx
//instanceObj[x,y] = xxx
func evalClassIndexerAssignExpression(a *ast.AssignExpression, obj Object, indexExpr *ast.IndexExpression, val Object, scope *Scope) Object {
	instanceObj := obj.(*ObjectInstance)

	var num int
	switch o := indexExpr.Index.(type) {
	case *ast.ClassIndexerExpression:
		num = len(o.Parameters)
	default:
		num = 1
	}

	propName := "this" + fmt.Sprintf("%d", num)

	//check if the Indexer is static
	if instanceObj.IsStatic(propName, ClassPropertyKind) {
		panic(NewError(a.Pos().Sline(), INDEXERSTATICERROR, instanceObj.Class.Name))
	}

	p := instanceObj.GetProperty(propName)
	if p != nil {
		//no setter or setter block is empty, e.g. 'property xxx { set; }'
		if p.Setter == nil || len(p.Setter.Body.Statements) == 0 {
			panic(NewError(a.Pos().Sline(), INDEXERUSEERROR, instanceObj.Class.Name))
		} else {
			newScope := NewScope(instanceObj.Scope)
			newScope.Set("value", val)

			switch o := indexExpr.Index.(type) {
			case *ast.ClassIndexerExpression:
				for i, v := range o.Parameters {
					index := Eval(v, scope)
					newScope.Set(p.Indexes[i].Value, index)
				}
			default:
				index := Eval(indexExpr.Index, scope)
				newScope.Set(p.Indexes[0].Value, index)
			}

			results := Eval(p.Setter.Body, newScope)
			if results.Type() == RETURN_VALUE_OBJ {
				return results.(*ReturnValue).Value
			}
			return results
		}
	} else {
		panic(NewError(a.Pos().Sline(), INDEXNOTFOUNDERROR, instanceObj.Class.Name))
	}

	return val
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
		} else if aObj.Type() == HASH_OBJ { //e.g. hash.key = value
			return evalHashAssignExpression(a, strArr[0], aObj, scope, val)
		} else if aObj.Type() == INSTANCE_OBJ { //e.g. this.var = xxxx
			instanceObj := aObj.(*ObjectInstance)

//			//get variable's modifier level
//			ml := instanceObj.GetModifierLevel(strArr[1], ClassMemberKind) //ml:modifier level
//			if ml == ast.ModifierPrivate {
//				panic(NewError(a.Pos().Sline(), CLSMEMBERPRIVATE, strArr[1], instanceObj.Class.Name))
//			}

			//check if it's a property
			p := instanceObj.GetProperty(strArr[1])
			if p == nil { //not property, return value from scope
				instanceObj.Scope.Set(strArr[1], val)
			} else {
					if p.Setter == nil { //property xxx { get; }
						_, ok := instanceObj.Scope.Get(strArr[1])
						if !ok { //it's the first time assignment
							if currentInstance != nil { //inside class body assignment
								instanceObj.Scope.Set(strArr[1], val)
							} else {  //outside class body assignment
								panic(NewError(a.Pos().Sline(), PROPERTYUSEERROR, strArr[1], instanceObj.Class.Name))
							}
						} else {
							panic(NewError(a.Pos().Sline(), PROPERTYUSEERROR, strArr[1], instanceObj.Class.Name))
						}
				} else {
					if len(p.Setter.Body.Statements) == 0 { // property xxx { set; }
						instanceObj.Scope.Set("_" + strArr[1], val)
					} else {
						newScope := NewScope(instanceObj.Scope)
						newScope.Set("value", val)
						results := Eval(p.Setter.Body, newScope)
						if results.Type() == RETURN_VALUE_OBJ {
							return results.(*ReturnValue).Value
						}
					}
				}
			}
			return
		} else if aObj.Type() == CLASS_OBJ { //e.g. parent.var = xxxx
			clsObj := aObj.(*Class)
			if currentInstance != nil { //inside class body
			currentInstance.Scope.Set(strArr[1], val)
			//return evalFunctionCall(o, currentInstance.Scope) //should be Function's scope
			} else { //outside class body
				clsObj.Scope.Set(strArr[1], val)
			}
			return
		}

		return evalStructAssignExpression(a, scope, val)
	}

	var name string
	switch nodeType := a.Name.(type) {
	case *ast.Identifier:
		name = nodeType.Value
	case *ast.IndexExpression:
		name = nodeType.Left.(*ast.Identifier).Value

		//check if it's a class indexer assignment, e.g. 'clsObj[index] = xxx'
		if aObj, ok := scope.Get(name); ok {
			if aObj.Type() == INSTANCE_OBJ {
				return evalClassIndexerAssignExpression(a, aObj, nodeType, val, scope)
			}
		}
//	case *ast.TupleLiteral:
//		if val.Type() != TUPLE_OBJ {
//			panic(NewError(a.Pos().Sline(), GENERICERROR, "The right part must be a tuple"))
//		}
//		
//		tupleObj := val.(*Tuple)
//		valuesLen := len(tupleObj.Members)
//
//		for idx, item := range nodeType.Members {
//			if idx >= valuesLen { //There are more Names than Values
//				val = NIL
//				scope.Set(item.String(), val)
//			} else {
//				val = tupleObj.Members[idx];
//				if val.Type() != ERROR_OBJ {
//					scope.Set(item.String(), val)
//				} else {
//					return
//				}
//			}
//		}
//		return
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
	case INTEGER_OBJ, UINTEGER_OBJ, FLOAT_OBJ:
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
	case TUPLE_OBJ:
		val = evalTupleAssignExpression(a, name, left, scope, val)
		return
	}

	panic(NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type()))
}

func evalReturnStatement(r *ast.ReturnStatement, scope *Scope) Object {
	ret := &ReturnValue{Values: []Object{}}
	for _, value := range r.ReturnValues {
		ret.Values = append(ret.Values, Eval(value, scope))
	}

	// for old campatibility
	if len(ret.Values) > 0 {
		ret.Value = ret.Values[0]
		return ret
	}

	return NIL
}

func evalDeferStatement(d *ast.DeferStmt, scope *Scope) Object {
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

func evalUIntegerLiteral(i *ast.UIntegerLiteral) Object {
	return NewUInteger(i.Value)
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
	if a.CreationCount == nil {
		return &Array{Members: evalArgs(a.Members, scope)}
	}

	var i int64
	ret := &Array{}
	for i = 0; i < a.CreationCount.Value; i++ {
		ret.Members = append(ret.Members, NIL)
	}

	return ret
}

func evalTupleLiteral(t *ast.TupleLiteral, scope *Scope) Object {
	return &Tuple{Members: evalArgs(t.Members, scope)}
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

func evalEnumStatement(enumStmt *ast.EnumStatement, scope *Scope) Object {
	enumLiteral := evalEnumLiteral(enumStmt.EnumLiteral, scope)
	scope.Set(enumStmt.Name.String(), enumLiteral) //save to scope
	return enumLiteral
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

func evalFunctionStatement(FnStmt *ast.FunctionStatement, scope *Scope) Object {
	fnObj := evalFunctionLiteral(FnStmt.FunctionLiteral, scope)
	fn := fnObj.(*Function)

	processClassAnnotation(FnStmt.Annotations, scope, FnStmt.Pos().Sline(), fn)
	scope.Set(FnStmt.Name.String(), fnObj) //save to scope

	return fnObj
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

// Prefix expression for User Defined Operator
func evalPrefixExpressionUDO(p *ast.PrefixExpression, right Object, scope *Scope) Object {
	if fn, ok := scope.Get(p.Operator); ok {
		f := fn.(*Function)
		// set functions's parameters
		scope.Set(f.Literal.Parameters[0].String(), right)
		r := Eval(f.Literal.Body, scope)
		if r.Type() == ERROR_OBJ {
			return r
		}

		if obj, ok := r.(*ReturnValue); ok {
			// if function returns multiple-values
			// returns a tuple instead.
			if len(obj.Values) > 1 {
				return &Tuple{Members: obj.Values, IsMulti: true}
			}
			return obj.Value
		}
		return r
	}
	panic(NewError(p.Pos().Sline(), PREFIXOP, p, right.Type()))
}

// Prefix expression for Meta-Operators
func evalMetaOperatorPrefixExpression(p *ast.PrefixExpression, right Object, scope *Scope) Object {
	if right.Type() != ARRAY_OBJ {
		panic(NewError(p.Pos().Sline(), PREFIXOP, p, right.Type()))
	}

	//convert prefix operator to infix operator,
	//Because 'evalNumberInfixExpression' function need a InfixExpression
	infixExp := &ast.InfixExpression{Token: p.Token, Operator: p.Operator, Right: p.Right}

	members := right.(*Array).Members
	if len(members) == 0 {
		return NewInteger(0)
	}

	result := members[0]
	var leftIsNum bool
	var leftIsStr bool

	_, leftIsNum = result.(Number)
	if !leftIsNum {
		_, leftIsStr = result.(*String)
		if !leftIsStr {
			panic(NewError(p.Pos().Sline(), METAOPERATORERROR))
		}
	}

	for i := 1; i < len(members); i++ {
		var rightIsNum bool
		var rightIsStr bool
		_, rightIsNum = members[i].(Number)
		if !rightIsNum {
			_, rightIsStr = members[i].(*String)
			if !rightIsStr {
				panic(NewError(p.Pos().Sline(), METAOPERATORERROR))
			}
		}

		if leftIsStr || rightIsStr {
			result = evalMixedTypeInfixExpression(infixExp, result, members[i])
		} else {
			result = evalNumberInfixExpression(infixExp, result, members[i])
		}

		_, leftIsNum = result.(Number)
		if !leftIsNum {
		_, leftIsStr = result.(*String)
		if !leftIsStr {
			panic(NewError(p.Pos().Sline(), METAOPERATORERROR))
		}
	}

	} // end for

	return result
}

// Prefix expressions, e.g. `!true, -5`
func evalPrefixExpression(p *ast.PrefixExpression, scope *Scope) Object {
	right := Eval(p.Right, scope)
	if right.Type() == ERROR_OBJ {
		return right
	}

	//User Defined Operator
	if p.Token.Type == token.UDO {
		return evalPrefixExpressionUDO(p, right, scope)
	}

	if isMetaOperators(p.Token.Type) {
		return evalMetaOperatorPrefixExpression(p, right, scope)
	}

	if right.Type() == INSTANCE_OBJ {
		/* e.g. p.operator = '-':
			class vector {
				let x;
				let y;
				fn init (a, b) {
					this.x = a
					this.y = b
				}
				fn -() {
					return new Vector(-x,-y)
				}
			}
			v1 = new vector(3,4)
			v2 = -v1
		*/
		instanceObj := right.(*ObjectInstance)
		method := instanceObj.GetMethod(p.Operator)
		if method != nil {
			switch method.(type) {
				case *Function:
					newScope := NewScope(instanceObj.Scope)
					newScope.Set("parent", instanceObj.Class.Parent)
					args := []Object{right}
					return evalFunctionDirect(method, args, newScope)
				case *BuiltinMethod:
					//do nothing for now
			}
		}
		panic(NewError(p.Pos().Sline(), PREFIXOP, p, right.Type()))
	}

	switch p.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "+":
		return right
	case "-":
		switch right.Type() {
		case INTEGER_OBJ:
			i := right.(*Integer)
			return NewInteger(-i.Int64)
			//bug : we need to return a new 'Integer' object, we should not change the original 'Integer' object.
			//i.Int64 = -i.Int64
			//return i
		case UINTEGER_OBJ:
			i := right.(*UInteger)
			if i.UInt64 == 0 {
				return i
			} else {
				panic(NewError(p.Pos().Sline(), PREFIXOP, p, right.Type()))
			}
		case FLOAT_OBJ:
			f := right.(*Float)
			return NewFloat(-f.Float64)
			//bug : we need to return a new 'Float' object, we should not change the original 'Float' object.
			//f.Float64 = -f.Float64
			//return f
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
	case UINTEGER_OBJ:
		rightObj := right.(*UInteger)
		rightObj.UInt64 = rightObj.UInt64 + 1
		return NewUInteger(rightObj.UInt64)
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
	case UINTEGER_OBJ:
		rightObj := right.(*UInteger)
		rightObj.UInt64 = rightObj.UInt64 - 1
		return NewUInteger(rightObj.UInt64)
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

// Infix expression for User Defined Operator
func evalInfixExpressionUDO(p *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	if fn, ok := scope.Get(p.Operator); ok {
		f := fn.(*Function)
		// set functions two parameters
		scope.Set(f.Literal.Parameters[0].String(), left)
		scope.Set(f.Literal.Parameters[1].String(), right)
		r := Eval(f.Literal.Body, scope)
		if r.Type() == ERROR_OBJ {
			return r
		}

		if obj, ok := r.(*ReturnValue); ok {
			// if function returns multiple-values
			// returns a tuple instead.
			if len(obj.Values) > 1 {
				return &Tuple{Members: obj.Values, IsMulti: true}
			}
			return obj.Value
		}
		return r
	}

	panic(NewError(p.Pos().Sline(), INFIXOP, left.Type(), p.Operator, right.Type()))
}

// Infix expression for Meta-Operators
func evalMetaOperatorInfixExpression(p *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	//1. [1,2,3] ~+ [4,5,6] = [1+4, 2+5, 3+6]
	//2. [1,2,3] ~+ 4 = [1+4, 2+4, 3+4]
	//left must be an array
	if left.Type() != ARRAY_OBJ {
		panic(NewError(p.Pos().Sline(), INFIXOP, left.Type(), p.Operator, right.Type()))
	}

	leftMembers := left.(*Array).Members
	leftNumLen := len(leftMembers)

	//right could be an array or a number
	var rightMembers []Object
	_, rightIsNum := right.(Number)
	if rightIsNum {
		for i := 0; i < leftNumLen; i++ {
			rightMembers = append(rightMembers, right)
		}
	} else {
		if right.Type() == ARRAY_OBJ {
			rightMembers = right.(*Array).Members
		} else {
			panic(NewError(p.Pos().Sline(), INFIXOP, left.Type(), p.Operator, right.Type()))
		}
	}
	rightNumLen := len(rightMembers)

	if leftNumLen != rightNumLen {
		panic(NewError(p.Pos().Sline(), GENERICERROR, "Number of items not equal for Meta-Operators!"))
	}

	resultArr := &Array{}
	if leftNumLen == 0 {
		return resultArr
	}

	for idx, item := range leftMembers {
		var leftIsNum, rightIsNum bool
		var leftIsStr, rightIsStr bool
		_, leftIsNum = item.(Number)
		if !leftIsNum {
			_, leftIsStr = item.(*String)
			if !leftIsStr {
				panic(NewError(p.Pos().Sline(), METAOPERATORERROR))
			}
		}

		_, rightIsNum = rightMembers[idx].(Number)
		if !rightIsNum {
			_, rightIsStr = rightMembers[idx].(*String)
			if !rightIsStr {
				panic(NewError(p.Pos().Sline(), METAOPERATORERROR))
			}
		}

		var result Object
		if leftIsNum && rightIsNum {
			result = evalNumberInfixExpression(p, item, rightMembers[idx])
		} else if leftIsStr && rightIsStr{
			result = evalStringInfixExpression(p, item, rightMembers[idx])
		} else if leftIsStr || rightIsStr {
			result = evalMixedTypeInfixExpression(p, item, rightMembers[idx])
		} else {
			panic(NewError(p.Pos().Sline(), METAOPERATORERROR))
		}
		resultArr.Members = append(resultArr.Members, result)
	} // end for

	return resultArr

}

// Evaluate infix expressions, e.g 1 + 2, a == 5, true == true, etc...
func evalInfixExpression(node *ast.InfixExpression, left, right Object, scope *Scope) Object {
	//User Defined Operator
	if node.Token.Type == token.UDO {
		return evalInfixExpressionUDO(node, left, right, scope)
	}

	if isMetaOperators(node.Token.Type) {
		return evalMetaOperatorInfixExpression(node, left, right, scope)
	}

	// Check if left is 'Writable'
	if _, ok := left.(Writable); ok { //There are two Writeables in monkey: FileObject, HttpResponseWriter.
		if node.Operator == ">>" { // '>>' is refered as 'extraction operator'. e.g.
			// Left is a file object
			if left.Type() == FILE_OBJ { // FileObject is also readable
				//    let a;
				//    stdin >> a
			}

			//right should be an identifier
			var rightVar *ast.Identifier
			var ok bool
			if rightVar, ok = node.Right.(*ast.Identifier); !ok { //not an identifier
				panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
			}
			f := left.(*FileObject)
			if f.File == os.Stdin {
				// 8192 is enough?
				ret := f.Read(node.Pos().Sline(), NewInteger(8192))
				scope.Set(rightVar.String(), ret)
				return ret
			} else {
				ret := f.ReadLine(node.Pos().Sline())
				scope.Set(rightVar.String(), ret)
				return ret
			}
		}

		if node.Operator == "<<" { // '<<' is refered as 'insertion operator'
			if f, ok := left.(*FileObject); ok { // It's a FileOject
				f.Write(node.Pos().Sline(), NewString(right.Inspect()))
				//Here we return left, so we can chain multiple '<<'.
				// e.g.
				//     stdout << "hello " << "world!"
				return left;
			}
			if httpResp, ok := left.(*HttpResponseWriter); ok { // It's a HttpResponseWriter
				httpResp.Write(node.Pos().Sline(), NewString(right.Inspect()))
				return left;
			}
		}
	}

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
		return evalArrayInfixExpression(node, left, right, scope)
	case (left.Type() == TUPLE_OBJ || right.Type() == TUPLE_OBJ):
		return evalTupleInfixExpression(node, left, right, scope)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(node, left, right)
	case (left.Type() == STRING_OBJ || right.Type() == STRING_OBJ):
		return evalMixedTypeInfixExpression(node, left, right)
	case (left.Type() == HASH_OBJ && right.Type() == HASH_OBJ):
		return evalHashInfixExpression(node, left, right)
	case left.Type() == INSTANCE_OBJ:
		return evalInstanceInfixExpression(node, left, right)
	case node.Operator == "==":
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return TRUE
			}
			return FALSE
		}

		if left.Type() != right.Type() {
			return FALSE
		}

		//Here we need to special handling for `Boolean` object. Because most of the time `BOOLEAN` will
		//return TRUE and FALSE. But sometimes we have to returns a new `Boolean` object,
		//Here we need to compare `Boolean.Bool` or else when we using
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
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return FALSE
			}
			return TRUE
		}

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

func isMetaOperators(tokenType token.TokenType) bool {
	return tokenType == token.TILDEPLUS ||               // ~+
		   tokenType == token.TILDEMINUS ||     // ~-
		   tokenType == token.TILDEASTERISK ||  // ~*
		   tokenType == token.TILDESLASH ||     // ~/
		   tokenType == token.TILDEMOD ||       // ~%
		   tokenType == token.TILDECARET        // ~^
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
	case *UInteger:
		if obj.UInt64 == 0 {
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
	case *Tuple:
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
	isUInt := left.Type() == UINTEGER_OBJ && right.Type() == UINTEGER_OBJ


	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*Integer).Int64)
	} else if left.Type() == UINTEGER_OBJ {
		leftVal = float64(left.(*UInteger).UInt64)
	} else {
		leftVal = left.(*Float).Float64
	}

	if right.Type() == INTEGER_OBJ {
		rightVal = float64(right.(*Integer).Int64)
	} else if right.Type() == UINTEGER_OBJ {
		rightVal = float64(right.(*UInteger).UInt64)
	} else {
		rightVal = right.(*Float).Float64
	}

	switch node.Operator {
	case "**", "~^":
		val := math.Pow(leftVal, rightVal)
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "&":
		if isInt {
			val := int64(leftVal) & int64(rightVal)
			return NewInteger(int64(val))
		} else if isUInt {
			val := uint64(leftVal) & uint64(rightVal)
			return NewUInteger(uint64(val))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "|":
		if isInt {
			val := int64(leftVal) | int64(rightVal)
			return NewInteger(int64(val))
		} else if isUInt {
			val := uint64(leftVal) | uint64(rightVal)
			return NewUInteger(uint64(val))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "^":
		if isInt {
			val := int64(leftVal) ^ int64(rightVal)
			return NewInteger(int64(val))
		} else if isUInt {
			val := uint64(leftVal) ^ uint64(rightVal)
			return NewUInteger(uint64(val))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "+", "~+":
		val := leftVal + rightVal
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "-", "~-":
		val := leftVal - rightVal
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "*", "~*":
		val := leftVal * rightVal
		if isInt {
			return NewInteger(int64(val))
		} else 	if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "/", "~/":
		if rightVal == 0 {
			panic(NewError(node.Pos().Sline(), DIVIDEBYZERO))
		}
		val := leftVal / rightVal
		//Should Always return float
		return NewFloat(val)
	case "%", "~%":
		if isInt {
			return NewInteger(int64(leftVal) % int64(rightVal))
		} else if isUInt {
			return NewUInteger(uint64(leftVal) % uint64(rightVal))
		}
		return NewFloat(math.Mod(leftVal, rightVal))
	case ">>":
		if isInt {
			aRes := uint64(leftVal) >> uint64(rightVal)
			return NewInteger(int64(aRes)) //NOTE: CAST MAYBE NOT CORRECT
		} else if isUInt {
			aRes := uint64(leftVal) >> uint64(rightVal)
			return NewUInteger(uint64(aRes))
		}
	case "<<":
		if isInt {
			aRes := uint64(leftVal) << uint64(rightVal)
			return NewInteger(int64(aRes)) //NOTE: CAST MAYBE NOT CORRECT
		} else if isUInt {
			aRes := uint64(leftVal) << uint64(rightVal)
			return NewUInteger(uint64(aRes))
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

func checkNumInfix(left Object, right Object, val float64) Object {
	if (left.Type() == INTEGER_OBJ && right.Type() == UINTEGER_OBJ) ||
		(left.Type() == UINTEGER_OBJ && right.Type() == INTEGER_OBJ) {
		if val > math.MaxInt64 {
			return NewUInteger(uint64(val))
		} else {
			return NewInteger(int64(val))
		}
	}

	return NewFloat(val)
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
	case "+", "~+":
		return NewString(left.Inspect() + right.Inspect())
	case "*", "~*":
		if left.Type() == INTEGER_OBJ {
			integer := left.(*Integer).Int64
			return NewString(strings.Repeat(right.Inspect(), int(integer)))
		} else if left.Type() == UINTEGER_OBJ {
			uinteger := left.(*UInteger).UInt64
			return NewString(strings.Repeat(right.Inspect(), int(uinteger)))
		}
		if right.Type() == INTEGER_OBJ {
			integer := right.(*Integer).Int64
			return NewString(strings.Repeat(left.Inspect(), int(integer)))
		} else if right.Type() == UINTEGER_OBJ {
			uinteger := right.(*UInteger).UInt64
			return NewString(strings.Repeat(left.Inspect(), int(uinteger)))
		}
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	case "==":
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return TRUE
			}
			return FALSE
		}

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
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return FALSE
			}
			return TRUE
		}

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

	case "=~": //match
		var str string
		if left.Type() == INTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*Integer).Int64)
		} else if left.Type() == UINTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*UInteger).UInt64)
		} else if left.Type() == FLOAT_OBJ {
			str = fmt.Sprintf("%g", left.(*Float).Float64)
		}
		matched, _ := regexp.MatchString(right.(*String).String, str)
		if matched {
			return TRUE
		}
		return FALSE

	case "!~": //not match
		var str string
		if left.Type() == INTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*Integer).Int64)
		} else if left.Type() == UINTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*UInteger).UInt64)
		} else if left.Type() == FLOAT_OBJ {
			str = fmt.Sprintf("%g", left.(*Float).Float64)
		}
		matched, _ := regexp.MatchString(right.(*String).String, str)
		if matched {
			return FALSE
		}
		return TRUE

	default:
		panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
	}

	//panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

//array + item
//item + array
//array + array
//array == array
//array != array
//array << item (<< item)
func evalArrayInfixExpression(node *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
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
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
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
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if IsTrue(aBool) {
				return TRUE
			}
		}
		return FALSE
	case "<<":
		if left.Type() == ARRAY_OBJ {
			leftVals := left.(*Array).Members
			leftVals = append(leftVals, right)
			left.(*Array).Members = leftVals //Change the array itself
			return left //return the original array, so it could be chained by another '<<'
		}
	}
	panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

//Almost same as evalArrayInfixExpression
//tuple + item
//item + tuple
//tuple + tuple
//tuple == tuple
//tuple != tuple
func evalTupleInfixExpression(node *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	switch node.Operator {
	case "+":
		if left.Type() == TUPLE_OBJ {
			leftVals := left.(*Tuple).Members

			if right.Type() == TUPLE_OBJ {
				rightVals := right.(*Tuple).Members
				leftVals = append(leftVals, rightVals...)
			} else {
				leftVals = append(leftVals, right)
			}
			return &Tuple{Members: leftVals}
		}

		//right is array
		rightVals := right.(*Tuple).Members
		if left.Type() == TUPLE_OBJ {
			leftVals := left.(*Tuple).Members
			rightVals = append(rightVals, leftVals...)
			return &Tuple{Members: rightVals}
		} else {
			ret := &Tuple{}
			ret.Members = append(ret.Members, left)
			ret.Members = append(ret.Members, rightVals...)
			return ret
		}

	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left.Type() != TUPLE_OBJ || right.Type() != TUPLE_OBJ {
			panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
		}

		leftVals := left.(*Tuple).Members
		rightVals := right.(*Tuple).Members
		if len(leftVals) != len(rightVals) {
			return FALSE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if !IsTrue(aBool) {
				return FALSE
			}
		}
		return TRUE
	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() != TUPLE_OBJ || right.Type() != TUPLE_OBJ {
			panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
		}
		leftVals := left.(*Tuple).Members
		rightVals := right.(*Tuple).Members
		if len(leftVals) != len(rightVals) {
			return TRUE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if IsTrue(aBool) {
				return TRUE
			}
		}
		return FALSE
	}
	panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
}

//hash + item
//hast == hash
//hash != hash
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

// for operaotor overloading, e.g.
//    class Vector {
//        fn +(v) { xxxxx }
//    }
//
//    v1 = new Vector()
//    v2 = new Vector()
//    v3 = v1 + v2   //here is the operator overloading, same as 'v3 = v1.+(v2)
func evalInstanceInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	instanceObj := left.(*ObjectInstance)

	switch node.Operator {
	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left == right {
			return TRUE
		}
		return FALSE
	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left != right {
			return TRUE
		}
		return FALSE
	}
	//get methods's modifier level
//	ml := instanceObj.GetModifierLevel(node.Operator, ClassMethodKind) //ml:modifier level
//	if ml == ast.ModifierPrivate {
//		panic(NewError(node.Pos().Sline(), CLSCALLPRIVATE, node.Operator, instanceObj.Class.Name))
//	}

	method := instanceObj.GetMethod(node.Operator)
	if method != nil {
		switch m := method.(type) {
			case *Function:
				fn := &Function{
					Literal:    m.Literal,
					Scope:       NewScope(instanceObj.Scope),
					Instance:   instanceObj,
				}

				fn.Scope.Set("parent", instanceObj.Class.Parent)

				args := []Object{right}
				return evalFunctionDirect(method, args, fn.Scope)
			case *BuiltinMethod:
				args := []Object{right}
				builtinMethod :=&BuiltinMethod{Fn: m.Fn, Instance: instanceObj}
				aScope := NewScope(instanceObj.Scope)
				return evalFunctionDirect(builtinMethod, args, aScope)
		}
	}
	panic(NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type()))
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
			switch o := c.Body.(type) {
			case *ast.BlockStatement:
				return evalBlockStatements(o.Statements, scope)
			}
			return Eval(c.Body, scope)
		}
	}

	//eval "else" part
	if ie.Alternative != nil {
		switch o := ie.Alternative.(type) {
		case *ast.BlockStatement:
			return evalBlockStatements(o.Statements, scope)
		}
		return Eval(ie.Alternative, scope)
	}

	return NIL
}


func evalUnlessExpression(ie *ast.UnlessExpression, scope *Scope) Object {
	condition := Eval(ie.Condition, scope)
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	if !IsTrue(condition) {
		return evalBlockStatements(ie.Consequence.Statements, scope)
	} else if ie.Alternative != nil {
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
				//return v.Value
				return v
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
				/*
					BUG: DO NOT RETURN 'v.Value', instead we should return 'v'.

					If we return 'v.Value' then below code will print '5', not '6'(which is expected)
					let add = fn(x,y){
					    let i = 0
					    while (i++ < 10) {
					        return x * y
					    }
					    return x + y 
					}
					println(add(2,3))
				*/

				//return v.Value
				return v
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

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		panic(NewError(ge.Pos().Sline(), GREPMAPNOTITERABLE))
	}
	if !iterObj.iter() {
		panic(NewError(ge.Pos().Sline(), GREPMAPNOTITERABLE))
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
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
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

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		panic(NewError(me.Pos().Sline(), GREPMAPNOTITERABLE))
	}
	if !iterObj.iter() {
		panic(NewError(me.Pos().Sline(), GREPMAPNOTITERABLE))
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
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
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

//[ x+1 for x in arr <where cond> ]
//[ str for str in strs <where cond> ]
//[ x for x in tuple <where cond> ]
func evalListComprehension(lc *ast.ListComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)
	aValue := Eval(lc.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		panic(NewError(lc.Pos().Sline(), NOTITERABLE))
	}
	if !iterObj.iter() {
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
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
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

//[ x for x in a..b <where cond> ]
//Almost same as evalForEachDotRangeExpression() function
func evalListRangeComprehension(lc *ast.ListRangeComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)

	startIdx := Eval(lc.StartIdx, innerScope)
	endIdx := Eval(lc.EndIdx, innerScope)

	arr := &Array{}

	switch startIdx.(type) {
	case *Integer:
		startVal := startIdx.(*Integer).Int64

		var endVal int64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = o.Int64
		case *UInteger:
			endVal = int64(o.UInt64)
		default:
			panic(NewError(lc.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ + "|" + UINTEGER_OBJ, endIdx.Type()))
		}

		var j int64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		}
	case *UInteger:
		startVal := startIdx.(*UInteger).UInt64

		var endVal uint64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = uint64(o.Int64)
		case *UInteger:
			endVal = o.UInt64
		default:
			panic(NewError(lc.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ + "|" + UINTEGER_OBJ, endIdx.Type()))
		}

		var j uint64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
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

//[ expr for k,v in hash <where cond> ]
func evalListMapComprehension(mc *ast.ListMapComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)
	aValue := Eval(mc.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		panic(NewError(mc.Pos().Sline(), NOTITERABLE))
	}
	if !iterObj.iter() {
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

//{ k:v for x in arr <where cond> }
//{ k:v for str in strs <where cond> }
//{ k:v for x in tuple <where cond> }
//Almost same as evalListComprehension
func evalHashComprehension(hc *ast.HashComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)
	aValue := Eval(hc.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		panic(NewError(hc.Pos().Sline(), NOTITERABLE))
	}
	if !iterObj.iter() {
		panic(NewError(hc.Pos().Sline(), NOTITERABLE))
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
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	}

	ret := &Hash{Pairs: make(map[HashKey]HashPair)}

	for idx, value := range members {
		newSubScope := NewScope(innerScope)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(hc.Var, value)
		if hc.Cond != nil {
			cond := Eval(hc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		keyResult := Eval(hc.KeyExpr, newSubScope)
		if keyResult.Type() == ERROR_OBJ {
			return keyResult
		}

		valueResult := Eval(hc.ValExpr, newSubScope)
		if valueResult.Type() == ERROR_OBJ {
			return valueResult
		}

		if hashable, ok := keyResult.(Hashable); ok {
			ret.Pairs[hashable.HashKey()] = HashPair{Key: keyResult, Value: valueResult}
		} else {
			panic(NewError("", KEYERROR, keyResult.Type()))
		}
	}

	return ret
}

//{ k:v for x in a..b <where cond> }
//Almost same as evalListRangeComprehension() function
func evalHashRangeComprehension(hc *ast.HashRangeComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)

	startIdx := Eval(hc.StartIdx, innerScope)
	endIdx := Eval(hc.EndIdx, innerScope)

	arr := &Array{}

	switch startIdx.(type) {
	case *Integer:
		startVal := startIdx.(*Integer).Int64

		var endVal int64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = o.Int64
		case *UInteger:
			endVal = int64(o.UInt64)
		default:
			panic(NewError(hc.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ + "|" + UINTEGER_OBJ, endIdx.Type()))
		}

		var j int64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		}
	case *UInteger:
		startVal := startIdx.(*UInteger).UInt64

		var endVal uint64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = uint64(o.Int64)
		case *UInteger:
			endVal = o.UInt64
		default:
			panic(NewError(hc.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ + "|" + UINTEGER_OBJ, endIdx.Type()))
		}

		var j uint64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
			}
		}
	case *String:
		startVal := startIdx.(*String).String
		if endIdx.Type() != STRING_OBJ {
			panic(NewError(hc.Pos().Sline(), RANGETYPEERROR, STRING_OBJ, endIdx.Type()))
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

	ret := &Hash{Pairs: make(map[HashKey]HashPair)}

	for idx, value := range arr.Members {
		newSubScope := NewScope(innerScope)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(hc.Var, value)
		if hc.Cond != nil {
			cond := Eval(hc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		keyResult := Eval(hc.KeyExpr, newSubScope)
		if keyResult.Type() == ERROR_OBJ {
			return keyResult
		}

		valueResult := Eval(hc.ValExpr, newSubScope)
		if valueResult.Type() == ERROR_OBJ {
			return valueResult
		}

		if hashable, ok := keyResult.(Hashable); ok {
			ret.Pairs[hashable.HashKey()] = HashPair{Key: keyResult, Value: valueResult}
		} else {
			panic(NewError("", KEYERROR, keyResult.Type()))
		}
	}

	return ret
}

//{ k:v for k,v in hash <where cond> }
//Almost same as evalListMapComprehension
func evalHashMapComprehension(mc *ast.HashMapComprehension, scope *Scope) Object {
	innerScope := NewScope(scope)
	aValue := Eval(mc.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		panic(NewError(mc.Pos().Sline(), NOTITERABLE))
	}
	if !iterObj.iter() {
		panic(NewError(mc.Pos().Sline(), NOTITERABLE))
	}

	//must be a *Hash, if not, panic
	hash, _ := aValue.(*Hash)

	ret := &Hash{Pairs: make(map[HashKey]HashPair)}

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

		keyResult := Eval(mc.KeyExpr, newSubScope)
		if keyResult.Type() == ERROR_OBJ {
			return keyResult
		}

		valueResult := Eval(mc.ValExpr, newSubScope)
		if valueResult.Type() == ERROR_OBJ {
			return valueResult
		}

		if hashable, ok := keyResult.(Hashable); ok {
			ret.Pairs[hashable.HashKey()] = HashPair{Key: keyResult, Value: valueResult}
		} else {
			panic(NewError("", KEYERROR, keyResult.Type()))
		}
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

	if fl.Init != nil {
		init := Eval(fl.Init, innerScope)
		if init.Type() == ERROR_OBJ {
			return init
		}
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
				//return v.Value
				return v
			}
			break
		}

		if fl.Update != nil {
			newVal := Eval(fl.Update, newSubScope)
			if newVal.Type() == ERROR_OBJ {
				return newVal
			}
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
	newScope := NewScope(scope)
	for {
		e = Eval(fel.Block, newScope)
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
				//return v.Value
				return v
			}
			break
		}
	}

	if e == nil || e.Type() == BREAK_OBJ || e.Type() == CONTINUE_OBJ {
		return NIL
	}
	return e
}

//for item in array
//for item in string
//for item in tuple
//for item in channel
func evalForEachArrayExpression(fal *ast.ForEachArrayLoop, scope *Scope) Object { //fal:For Array Loop
	innerScope := NewScope(scope)

	aValue := Eval(fal.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable)
	if !ok {
		panic(NewError(fal.Pos().Sline(), NOTITERABLE))
	}
	if !iterObj.iter() {
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
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	} else if aValue.Type() == GO_OBJ { // GoObject
		goObj := aValue.(*GoObject)
		arr := GoValueToObject(goObj.obj).(*Array)
		members = arr.Members
	} else if aValue.Type() == CHANNEL_OBJ {
		chanObj := aValue.(*ChanObject)
		ret := &Array{}
		var result Object

		idx := 0
		for value := range chanObj.ch {
			scope.Set("$_", NewInteger(int64(idx)))
			idx++
			scope.Set(fal.Var, value)
			result = Eval(fal.Block, scope)
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
					return v
				}
				break
			} else {
				ret.Members = append(ret.Members, result)
			}

		} //end for
		if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
			return ret
		}
		return ret
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
				//ret.Members = append(ret.Members, v.Value)
				return v
				//return v.Value
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	//Here we need to check `nil`, because if the initial condition is not true, then `for`'s Body will have no
	//chance to execute, the result will be nil
	//this is the reason why we need to check for `BREAK_OBJ` or `CONTINUE_OBJ`:
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

//for index, value in string
//for index, value in array
//for index, value in tuple
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
	} else if val.Type() == TUPLE_OBJ {
		tuple, _ := val.(*Tuple)
		members = tuple.Members
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
				return v
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

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members:[]Object{}}
	}

	iterObj, ok := aValue.(Iterable)
	if !ok {
		panic(NewError(fml.Pos().Sline(), NOTITERABLE))
	}
	if !iterObj.iter() {
		panic(NewError(fml.Pos().Sline(), NOTITERABLE))
	}

	//for index, value in arr
	//for index, value in string
	if aValue.Type() == STRING_OBJ || aValue.Type() == ARRAY_OBJ || aValue.Type() == TUPLE_OBJ {
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
				return v
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

	arr := &Array{}

	switch startIdx.(type) {
	case *Integer:
		startVal := startIdx.(*Integer).Int64

		var endVal int64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = o.Int64
		case *UInteger:
			endVal = int64(o.UInt64)
		default:
			panic(NewError(fdr.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ + "|" + UINTEGER_OBJ, endIdx.Type()))
		}

		var j int64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		}
	case *UInteger:
		startVal := startIdx.(*UInteger).UInt64

		var endVal uint64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = uint64(o.Int64)
		case *UInteger:
			endVal = o.UInt64
		default:
			panic(NewError(fdr.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ + "|" + UINTEGER_OBJ, endIdx.Type()))
		}

		var j uint64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
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
				//ret.Members = append(ret.Members, v.Value)
				return v
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
		case UINTEGER_OBJ:
			if obj.(*UInteger).UInt64 == 0 {
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
		case TUPLE_OBJ:
			if len(obj.(*Tuple).Members) == 0 {
				return false
			}
		case GO_OBJ:
			goObj := obj.(*GoObject)
			return goObj.obj != nil
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
			//let complex={ "add" : fn(x,y){ x+y } }
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
			args := evalArgs(call.Arguments, scope)
			return builtin.Fn(call.Function.Pos().Sline(), args...)
		} else if callExpr, ok := call.Function.(*ast.CallExpression); ok { //call expression
			//let complex={ "add" : fn(x,y){ fn(z) {x+y+z} } }
			//complex["add"](2,3)(4)
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

	//check if it's static function
	aVal, ok := scope.Get("this")
	if ok && aVal.Type() == CLASS_OBJ {
		if !f.Literal.StaticFlag { //not static
			panic(NewError(call.Function.Pos().Sline(), CALLNONSTATICERROR))
		}
	}

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

	oldInstance := currentInstance
	currentInstance = f.Instance

	r := Eval(f.Literal.Body, newScope)
	currentInstance = oldInstance
	if r.Type() == ERROR_OBJ {
		return r
	}

	if obj, ok := r.(*ReturnValue); ok {
		// if function returns multiple-values
		// returns a tuple instead.
		if len(obj.Values) > 1 {
			return &Tuple{Members: obj.Values, IsMulti: true}
		}
		return obj.Value
	}
	return r
}

// Method calls for builtin Objects
func evalMethodCallExpression(call *ast.MethodCallExpression, scope *Scope) Object {
	//First check if is a stanard library object
	str := call.Object.String()
	if obj, ok := GetGlobalObj(str); ok {
		switch o := call.Call.(type) {
		case *ast.IndexExpression: // e.g. 'if gos.Args[0] == "hello" {'
			if arr, ok := GetGlobalObj(str + "." + o.Left.String()); ok {
				return evalArrayIndex(arr.(*Array), o, scope)
			}
		case *ast.Identifier: //e.g. os.O_APPEND
			if i, ok := GetGlobalObj(str + "." + o.String()); ok {
				return i
			} else { //e.g. method call like 'os.environ'
				if obj.Type() == HASH_OBJ { // It's a GoFuncObject
					hashPairs := obj.(*Hash).Pairs
					for _, pair := range hashPairs {
						funcName := pair.Key.(*String).String
						if funcName == o.String() {
							goFuncObj := pair.Value.(*GoFuncObject)
							return goFuncObj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
						}
					}
				} else {
					return obj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
				}
			}
		case *ast.CallExpression: //e.g. method call like 'os.environ()'
			if method, ok := call.Call.(*ast.CallExpression); ok {
				args := evalArgs(method.Arguments, scope)
				if obj.Type() == HASH_OBJ { // It's a GoFuncObject
					hashPairs := obj.(*Hash).Pairs
					for _, pair := range hashPairs {
						funcName := pair.Key.(*String).String
						if funcName == o.Function.String() {
							goFuncObj := pair.Value.(*GoFuncObject)
							return goFuncObj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
						}
					}
				} else {
					return obj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
				}
			}
		}
	} else {
		// if 'GetGlobalObj(str)' returns nil, then try below
		// e.g.
		//     eval.RegisterVars("runtime", map[string]interface{}{
		//          "GOOS": runtime.GOOS,
		//       })
		// The eval.RegisterVars will call SetGlobalObj("runtime.GOOS"), so the
		// global scope's name is 'runtime.GOOS', not 'runtime', therefore, the above 
		// GetGlobalObj('runtime') will returns false.
		if obj, ok := GetGlobalObj(str + "." + call.Call.String()); ok {
			return obj
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
	case *Hash:
		switch o := call.Call.(type) {
		//e.g.:
		//hashObj.key1=10
		//println(hashObj.key1)
		case *ast.Identifier:
			keyObj := NewString(call.Call.String())
			hashPair, ok := m.Pairs[keyObj.HashKey()]
			// TODO: should we return an error here?
			if !ok {
				return NIL
			}
			return hashPair.Value

		case *ast.CallExpression:
			//we need to get the hash key
			keyStr := strings.Split(call.Call.String(), "(")[0]
			keyObj := NewString(keyStr)
			hashPair, ok := m.Pairs[keyObj.HashKey()]
			if !ok {
				//Check if it's a hash object's builtin method(e.g. hashObj.keys(), hashObj.values())
				if method, ok := call.Call.(*ast.CallExpression); ok {
					args := evalArgs(method.Arguments, scope)
					return obj.CallMethod(call.Call.Pos().Sline(), scope, method.Function.String(), args...)
				}
			}

			//e.g.:
			//hashObj = {}
			//hashObj.str = fn() { return 10 }
			//hashObj.str()

			// we need 'FunctionLiteral' here, so we need to change 'o.Function',
			// because the o.Function's Type is '*ast.Identifier' which is the Hash's key
			o.Function = hashPair.Value.(*Function).Literal
			//return evalFunctionCall(o, scope)   This is a bug: not 'scope'
			return evalFunctionCall(o, hashPair.Value.(*Function).Scope) //should be Function's scope
		}
	case *ObjectInstance:
		instanceObj := m
		switch o := call.Call.(type) {
		//e.g.: instanceObj.key1
		case *ast.Identifier:
//			//get methods's modifier level
//			ml := currentInstance.GetModifierLevel(o.Value, ClassMemberKind) //ml:modifier level
//			if currentInstance.Class.Name != instanceObj.Class.Name && ml == ast.ModifierPrivate {
//				panic(NewError(call.Call.Pos().Sline(), CLSMEMBERPRIVATE, o.Value, instanceObj.Class.Name))
//			}

			val, ok := instanceObj.Scope.Get(o.Value)
			if ok {
				switch val.(type) {
				case *Function: //Function without parameter. e.g. obj.getMonth(), could be called using 'obj.getMonth'
					return evalFunctionDirect(val, []Object{}, instanceObj.Scope)
				default:
					return val
				}
			}

			//See if it's a property
			p := instanceObj.GetProperty(o.Value)
			if p != nil {
				if p.Getter == nil { //property xxx { set; }
					panic(NewError(call.Call.Pos().Sline(), PROPERTYUSEERROR, o.Value, instanceObj.Class.Name))
				} else {
					if len(p.Getter.Body.Statements) == 0 { //property xxx { get; }
						v, _ := instanceObj.Scope.Get("_" + o.Value)
						//instanceObj.Scope.Set("_" + o.Value, v)
						return v
					} else {
						results := Eval(p.Getter.Body, instanceObj.Scope)
						if results.Type() == RETURN_VALUE_OBJ {
							return results.(*ReturnValue).Value
						}
					}
				}
			}
			panic(NewError(call.Call.Pos().Sline(), UNKNOWNIDENT, o.Value))

		case *ast.CallExpression:
			//e.g. instanceObj.method()
			fname := o.Function.String() // get function name

			//get methods's modifier level
//			ml := instanceObj.GetModifierLevel(fname, ClassMethodKind) //ml:modifier level
//			if ml == ast.ModifierPrivate {
//				panic(NewError(call.Call.Pos().Sline(), CLSCALLPRIVATE, fname, instanceObj.Class.Name))
//			}

			method := instanceObj.GetMethod(fname)
			if method != nil {
				switch m := method.(type) {
					case *Function:
						fn := &Function{
							Literal:    m.Literal,
							Scope:       NewScope(instanceObj.Scope),
							Instance:   instanceObj,
						}
						fn.Scope.Set("parent", instanceObj.Class.Parent)

						args := evalArgs(o.Arguments, fn.Scope)
						return evalFunctionDirect(method, args, fn.Scope)
					case *BuiltinMethod:
						builtinMethod :=&BuiltinMethod{Fn: m.Fn, Instance: instanceObj}
						aScope := NewScope(instanceObj.Scope)
						args := evalArgs(o.Arguments, aScope)
						return evalFunctionDirect(builtinMethod, args, aScope)
				}
			}
			panic(NewError(call.Call.Pos().Sline(), NOMETHODERROR, call.String(), obj.Type()))
			return NIL
		}
	case *Class:
		clsObj := m
		switch o := call.Call.(type) {
		case *ast.Identifier: //e.g.: classObj.key1
//			//get methods's modifier level
//			ml := currentInstance.GetModifierLevel(o.Value, ClassMemberKind) //ml:modifier level
//			if currentInstance.Class.Name != clsObj.Name && ml == ast.ModifierPrivate {
//				panic(NewError(call.Call.Pos().Sline(), CLSMEMBERPRIVATE, o.Value, clsObj.Name))
//			}

			var val Object
			var ok bool
			if currentInstance != nil {
				val, ok = currentInstance.Scope.Get(o.Value)
			} else {
				//check if it is static
				if clsObj.IsStatic(o.Value, ClassMemberKind) {
					val, ok = clsObj.Scope.Get(o.Value)
				} else if clsObj.IsStatic(o.Value, ClassPropertyKind) {
					val, ok = clsObj.Scope.Get(o.Value)
				} else {
					//val, ok = NIL, false
					panic(NewError(call.Call.Pos().Sline(), CALLNONSTATICERROR))
				}
			}
			if ok {
				return val
			}
			return NIL

		case *ast.CallExpression: //e.g. classObj.method()
			fname := o.Function.String() // get function name
			isStatic := clsObj.IsStatic(fname, ClassMethodKind)
			if currentInstance == nil { //outside class body
				if !isStatic {
					panic(NewError(call.Call.Pos().Sline(), CALLNONSTATICERROR))
				}
			} else {
				if isStatic {
					panic(NewError(call.Call.Pos().Sline(), CALLNONSTATICERROR))
				}
			}

			method := clsObj.GetMethod(fname)
			if method != nil {
				switch m := method.(type) {
					case *Function:
						args := evalArgs(o.Arguments, scope)
						if currentInstance != nil { //inside class body call
							currentInstance.Scope.Set("parent", clsObj.Parent) //needed ?
							return evalFunctionDirect(m, args, currentInstance.Scope)
						} else { //outside class body call
							return evalFunctionDirect(m, args, clsObj.Scope)
						}
					case *BuiltinMethod:
						builtinMethod :=&BuiltinMethod{Fn: m.Fn, Instance: currentInstance}
						aScope := NewScope(currentInstance.Scope)
						args := evalArgs(o.Arguments, aScope)
						return evalFunctionDirect(builtinMethod, args, aScope)
					
				}
			} else {
				args := evalArgs(o.Arguments, scope)
				return clsObj.CallMethod(call.Call.Pos().Sline(), scope, fname, args...)
			}
		}

	default:
		switch o := call.Call.(type) {
		case *ast.Identifier:      //e.g. method call like '[1,2,3].first'
			return obj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
		case *ast.CallExpression:  //e.g. method call like '[1,2,3].first()'
			args := evalArgs(o.Arguments, scope)
			return obj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
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

// Index Expressions, i.e. array[0], array[2:4], tuple[3] or hash["mykey"]
func evalIndexExpression(ie *ast.IndexExpression, scope *Scope) Object {
	left := Eval(ie.Left, scope)
	switch iterable := left.(type) {
	case *Array:
		return evalArrayIndex(iterable, ie, scope)
	case *Hash:
		return evalHashKeyIndex(iterable, ie, scope)
	case *String:
		return evalStringIndex(iterable, ie, scope)
	case *Tuple:
		return evalTupleIndex(iterable, ie, scope)
	case *ObjectInstance: //class indexer's getter
		return evalClassInstanceIndex(iterable, ie, scope)
	}
	panic(NewError(ie.Pos().Sline(), NOINDEXERROR, left.Type()))
}

func evalClassInstanceIndex(instanceObj *ObjectInstance, ie *ast.IndexExpression, scope *Scope) Object {
	var num int
	switch o := ie.Index.(type) {
	case *ast.ClassIndexerExpression:
		num = len(o.Parameters)
	default:
		num = 1
	}

	propName := "this" + fmt.Sprintf("%d", num)
	p := instanceObj.GetProperty(propName)
	if p != nil {
		//no getter or getter block is empty, e.g. 'property xxx { get; }'
		if p.Getter == nil || len(p.Getter.Body.Statements) == 0 {
			panic(NewError(ie.Pos().Sline(), INDEXERUSEERROR, instanceObj.Class.Name))
		} else {
			newScope := NewScope(instanceObj.Scope)

			switch o := ie.Index.(type) {
			case *ast.ClassIndexerExpression:
				for i, v := range o.Parameters {
					index := Eval(v, scope)
					newScope.Set(p.Indexes[i].Value, index)
				}
			default:
				index := Eval(ie.Index, scope)
				newScope.Set(p.Indexes[0].Value, index)
			}

			results := Eval(p.Getter.Body, newScope)
			if results.Type() == RETURN_VALUE_OBJ {
				return results.(*ReturnValue).Value
			}
			return results
		}
	}

	panic(NewError(ie.Pos().Sline(), INDEXNOTFOUNDERROR, instanceObj.Class.Name))
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

	switch o := index.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(index) {
			idx = 1
		}
	}

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

	switch o := startIdx.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	}
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

		switch o := slIndex.(type) {
		case *Integer:
			slice = o.Int64
		case *UInteger:
			slice = int64(o.UInt64)
		}
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

	switch o := startIdx.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(o) {
			idx = 1
		}
	}

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

		switch o := slIndex.(type) {
		case *Integer:
			slice = o.Int64
		case *UInteger:
			slice = int64(o.UInt64)
		default:
			slice = 0
			if IsTrue(o) {
				slice = 1
			}
		}
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

	switch o := index.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(index) {
			idx = 1
		}
	}
	if idx < 0 {
		panic(NewError(ie.Pos().Sline(), INDEXERROR, idx))
	}
	if idx >= length {
		return NIL
	}
	return array.Members[idx]
}

//Almost same as evalArraySliceExpression
func evalTupleSliceExpression(tuple *Tuple, se *ast.SliceExpression, scope *Scope) Object {
	var idx int64
	var slice int64
	length := int64(len(tuple.Members))

	startIdx := Eval(se.StartIndex, scope)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}

	switch o := startIdx.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	}
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

		switch o := slIndex.(type) {
		case *Integer:
			slice = o.Int64
		case *UInteger:
			slice = int64(o.UInt64)
		}
		if slice >= (length+1) || slice < 0 {
			panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
		}
	}
	if idx == 0 && slice == length {
		return tuple
	}

	if slice < idx {
		panic(NewError(se.Pos().Sline(), SLICEERROR, idx, slice))
	}

	if slice == length {
		return &Tuple{Members: tuple.Members[idx:]}
	}
	return &Tuple{Members: tuple.Members[idx:slice]}
}

//Almost same as evalArrayIndex
func evalTupleIndex(tuple *Tuple, ie *ast.IndexExpression, scope *Scope) Object {
	var idx int64
	length := int64(len(tuple.Members))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalTupleSliceExpression(tuple, exp, scope)
	}
	index := Eval(ie.Index, scope)
	if index.Type() == ERROR_OBJ {
		return index
	}

	switch o := index.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(index) {
			idx = 1
		}
	}

	if idx < 0 {
		panic(NewError(ie.Pos().Sline(), INDEXERROR, idx))
	}
	if idx >= length {
		return NIL
	}
	return tuple.Members[idx]
}

func evalPostfixExpression(left Object, node *ast.PostfixExpression) Object {
	if left.Type() == INSTANCE_OBJ { //operator overloading
		instanceObj := left.(*ObjectInstance)
		method := instanceObj.GetMethod(node.Operator)
		if method != nil {
			switch method.(type) {
				case *Function:
					newScope := NewScope(instanceObj.Scope)
					newScope.Set("parent", instanceObj.Class.Parent)
					args := []Object{left}
					return evalFunctionDirect(method, args, newScope)
				case *BuiltinMethod:
					//do nothing for now
			}
		}
		panic(NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type()))
	}

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
	case UINTEGER_OBJ:
		leftObj := left.(*UInteger)
		returnVal := NewUInteger(leftObj.UInt64)
		leftObj.UInt64 = leftObj.UInt64 + 1
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
	case UINTEGER_OBJ:
		leftObj := left.(*UInteger)
		returnVal := NewUInteger(leftObj.UInt64)
		leftObj.UInt64 = leftObj.UInt64 - 1
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

func evalSpawnStatement(s *ast.SpawnStmt, scope *Scope) Object {
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

//class name : parent { block }
//class name (categoryname) { block }
func evalClassStatement(c *ast.ClassStatement, scope *Scope) Object {
	if c.CategoryName != nil { //it's a class literal
		clsObj, ok := scope.Get(c.Name.Value)
		if !ok {
			panic(NewError(c.Pos().Sline(), CLASSCATEGORYERROR, c.Name, c.CategoryName))
		}

		//category only support methods and properties
		cls := clsObj.(*Class)
		for k, f := range c.ClassLiteral.Methods { //f :function
			cls.Methods[k] = Eval(f, scope).(ClassMethod)
		}
		for k, p := range c.ClassLiteral.Properties { //p :property
			cls.Properties[k] = p
		}

		return NIL
	}

	var clsObj Object
	if c.IsAnnotation {
		clsObj = evalClassLiterlForAnno(c.ClassLiteral, scope)
	} else {
		clsObj = evalClassLiteral(c.ClassLiteral, scope)
	}

	scope.Set(c.Name.Value, clsObj) //save to scope

	return NIL
}

//let name = class : parent { block }
func evalClassLiteral(c *ast.ClassLiteral, scope *Scope) Object {
	var parentClass = BASE_CLASS //base class is the root of all classes in monkey
	if c.Parent != "" {

		parent, ok := scope.Get(c.Parent)
		if !ok {
			panic(NewError(c.Pos().Sline(), PARENTNOTDECL, c.Parent))
		}

		parentClass, ok = parent.(*Class)
		if !ok {
			panic(NewError(c.Pos().Sline(), NOTCLASSERROR, c.Parent))
		}
	}

	clsObj := &Class{
		Name:       c.Name,
		Parent:     parentClass,
		Members:    c.Members,
		Properties: c.Properties,
		Methods:    make(map[string]ClassMethod, len(c.Methods)),
	}

	tmpClass := clsObj
	classChain := make([]*Class, 0, 3)
	classChain = append(classChain, clsObj)
	for tmpClass.Parent != nil {
		classChain = append(classChain, tmpClass.Parent)
		tmpClass = tmpClass.Parent
	}

	//create a new Class scope
	newScope := NewScope(scope)
	//evaluate the 'Members' fields of class with proper scope.
	for idx := len(classChain) - 1; idx >= 0; idx-- {
		for _, member := range classChain[idx].Members {
			Eval(member, newScope) //evaluate the 'Members' fields of class
		}
		newScope = NewScope(newScope)
	}
	clsObj.Scope = newScope.parentScope
	clsObj.Scope.Set("this", clsObj) //make 'this' refer to class object itself

	for k, f := range c.Methods {
		clsObj.Methods[k] = Eval(f, scope).(ClassMethod)
	}

	//check if the method has @Override annotation, if so, search
	//the method in parent hierarchical, if not found, then panic.
	for methodName, fnStmt := range c.Methods {
		for _, anno := range fnStmt.Annotations {
			if anno.Name.Value == OVERRIDE_ANNOCLASS.Name {
				if clsObj.Parent.GetMethod(methodName) == nil {
					panic(NewError(fnStmt.FunctionLiteral.Pos().Sline(), OVERRIDEERROR, methodName, clsObj.Name))
				}
			}
		}
	}

	return clsObj
}

func evalClassLiterlForAnno(c *ast.ClassLiteral, scope *Scope) Object {
	var parentClass = BASE_CLASS //base class is the root of all classes in monkey
	if c.Parent != "" {
		parent, ok := scope.Get(c.Parent)
		if !ok {
			panic(NewError(c.Pos().Sline(), PARENTNOTDECL, c.Parent))
		}

		parentClass, ok = parent.(*Class)
		if !ok {
			panic(NewError(c.Pos().Sline(), NOTCLASSERROR, c.Parent))
		}
	}

	if parentClass != BASE_CLASS && !parentClass.IsAnnotation { //parent not annotation
		panic(NewError(c.Pos().Sline(), PARENTNOTANNOTATION, c.Name, parentClass.Name))
	}

	clsObj := &Class{
		Name:         c.Name,
		Parent:       parentClass,
		Properties:   c.Properties,
		IsAnnotation: true,
	}

	//create a new Class scope
	clsObj.Scope = NewScope(scope)
	clsObj.Scope.Set("this", clsObj) //make 'this' refer to class object itself

	return clsObj
}

//new classname(parameters)
func evalNewExpression(n *ast.NewExpression, scope *Scope) Object {
	class := Eval(n.Class, scope)
	if class == NIL || class == nil {
		panic(NewError(n.Pos().Sline(), CLSNOTDEFINE, n.Class))
	}

	clsObj, ok := class.(*Class)
	if !ok {
		panic(NewError(n.Pos().Sline(), NOTCLASSERROR, n.Class))
	}

	tmpClass := clsObj
	classChain := make([]*Class, 0, 3)
	classChain = append(classChain, clsObj)
	for tmpClass.Parent != nil {
		classChain = append(classChain, tmpClass.Parent)
		tmpClass = tmpClass.Parent
	}

	//create a new Class scope
	newScope := NewScope(scope)
	//evaluate the 'Members' fields of class with proper scope.
	for idx := len(classChain) - 1; idx >= 0; idx-- {
		for _, member := range classChain[idx].Members {
			Eval(member, newScope) //evaluate the 'Members' fields of class
		}
		newScope = NewScope(newScope)
	}

	instance := &ObjectInstance{Class: clsObj, Scope: newScope.parentScope}
	instance.Scope.Set("this", instance) //make 'this' refer to instance
	instance.Scope.Set("parent", classChain[1]) //make 'parent' refer to instance's parent

	//Is it has a constructor ?
	init := clsObj.GetMethod("init")
	if init == nil {
		return instance
	}

	classMux.Lock()
	defer classMux.Unlock()

	args := evalArgs(n.Arguments, scope)
	if len(args) == 1 && args[0].Type() == ERROR_OBJ {
		return args[0]
	}

	oldInstance := currentInstance
	currentInstance = instance
	ret := evalFunctionDirect(init, args, instance.Scope); 
	if ret.Type() == ERROR_OBJ {
		currentInstance = oldInstance
		return ret //return the error object
	}
	currentInstance = oldInstance
	return instance
}

func processClassAnnotation(Annotations []*ast.AnnotationStmt, scope *Scope, line string, obj Object) {
	for _, anno := range Annotations {  //for each annotation
		annoClass, ok := scope.Get(anno.Name.Value)
		if !ok {
			panic(NewError(line, CLSNOTDEFINE, anno.Name.Value))
		}

		annoClsObj := annoClass.(*Class)

		//create the annotation instance
		newScope := NewScope(scope)
		annoInstanceObj := &ObjectInstance{Class: annoClsObj, Scope: newScope}
		annoInstanceObj.Scope.Set("this", annoInstanceObj) //make 'this' refer to annoObj

		switch o := obj.(type) {
			case *Function:
				o.Annotations = append(o.Annotations, annoInstanceObj)
			case *Array:
				o.Members = append(o.Members, annoInstanceObj)
		}

		defaultPropMap := make(map[string]ast.Expression)
		//get all propertis which have default value in the annotation class
		tmpCls := annoClsObj
		for tmpCls != nil {
			for name, item := range tmpCls.Properties {
				if item.Default != nil {
					defaultPropMap[name] = item.Default
				}
			}

			tmpCls = tmpCls.Parent
		}

		//check if the property(which has default value) exists in anno.Attribues
		for name, item := range defaultPropMap {
			if _, ok := anno.Attributes[name]; !ok { //not exists
				anno.Attributes[name] = item
			}
		}

		for k, v := range anno.Attributes { //for each annotation attribute
			val := Eval(v, annoInstanceObj.Scope)
			p := annoClsObj.GetProperty(k)
			if p == nil {
				annoInstanceObj.Scope.Set(k, val)
			} else {
				annoInstanceObj.Scope.Set("_" + k, val)
			}
		}
	}
}

func evalFunctionDirect(fn Object, args []Object, scope *Scope) Object {
	switch fn := fn.(type) {
	case *Function:
//		if len(args) < len(fn.Literal.Parameters) {
//			panic(NewError("", GENERICERROR, "Not enough parameters to call function"))
//		}

		newScope := NewScope(scope)
		variadicParam := []Object{}
		for i, _ := range args {
			//Because of function default values, we need to check `i >= len(args)`
			if fn.Variadic && i >= len(fn.Literal.Parameters)-1 {
				for j := i; j < len(args); j++ {
					variadicParam = append(variadicParam, args[j])
				}
				break
			} else if i >= len(fn.Literal.Parameters) {
				break
			} else {
				newScope.Set(fn.Literal.Parameters[i].String(), args[i])
			}
		}

		// Variadic argument is passed as a single array
		// of parameters.
		if fn.Variadic {
			newScope.Set(fn.Literal.Parameters[len(fn.Literal.Parameters)-1].String(), &Array{Members: variadicParam})
			if len(args) < len(fn.Literal.Parameters) {
				newScope.Set("@_", NewInteger(int64(len(fn.Literal.Parameters)-1)))
			} else {
				newScope.Set("@_", NewInteger(int64(len(args))))
			}
		} else {
			newScope.Set("@_", NewInteger(int64(len(fn.Literal.Parameters))))
		}

		//newScope.DebugPrint("    ") //debug
		results := Eval(fn.Literal.Body, newScope)
		if results.Type() == RETURN_VALUE_OBJ {
			return results.(*ReturnValue).Value
		}
		return results
	case *Builtin:
		return fn.Fn("", args...)
	case *BuiltinMethod:
		return fn.Fn("", fn.Instance, scope, args...)
	}

	panic(NewError("", GENERICERROR, fn.Type() + " is not a function"))
}

//evaluate 'using' statement
func evalUsingStatement(u *ast.UsingStmt, scope *Scope) Object {
	//evaluate the assignment expression
	obj := evalAssignExpression(u.Expr, scope)
	fn := func() {
		if obj.Type() != NIL_OBJ {
			// Check if val is 'Closeable'
			if c, ok := obj.(Closeable); ok {
				//call the 'Close' method of the object
				c.close(u.Pos().Sline())
			}
		}
	}
	defer func() {
		if r := recover(); r != nil { // if there is panic, we need to call fn()
			fn()
		} else { //no panic, we also need to call fn()
			fn()
		}
	}()

	//evaluate the 'using' block statement
	Eval(u.Block, scope)

	return NIL
}

// Convert a Object to an ast.Expression.
func obj2Expression(obj Object) ast.Expression {
	switch value := obj.(type) {
	case *Boolean:
		return &ast.Boolean{Value: value.Bool}
	case *Integer:
		return &ast.IntegerLiteral{Value: value.Int64}
	case *UInteger:
		return &ast.UIntegerLiteral{Value: value.UInt64}
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
	case *Tuple:
		tuple := &ast.TupleLiteral{}
		for _, v := range value.Members {
			result := obj2Expression(v)
			if result == nil {
				return nil
			}
			tuple.Members = append(tuple.Members, result)
		}
		return tuple
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

	if isGoObj(lhsV) || isGoObj(rhsV) { // if it's GoObject
		return compareGoObj(lhsV, rhsV)
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

func isGoObj(o Object) bool {
	return o.Type() == GO_OBJ
}

func compareGoObj(left, right Object) bool {
	if left.Type() == GO_OBJ || right.Type() == GO_OBJ {
		var goObj *GoObject
		var another Object
		if left.Type() == GO_OBJ {
			goObj = left.(*GoObject)
			another = right
		} else {
			goObj = right.(*GoObject)
			another = left
		}

		return goObj.Equal(another)
	}

	//left and right both are GoObject
	return left.(*GoObject).Equal(right)
}
