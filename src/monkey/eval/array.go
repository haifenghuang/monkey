package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	_ "fmt"
	"monkey/ast"
	"strings"
)

type Array struct {
	Members []Object
}

func (a *Array) iter() {}
func (a *Array) list() {}

func (a *Array) Inspect() string {
	var out bytes.Buffer
	members := []string{}
	for _, m := range a.Members {
		if m.Type() == STRING_OBJ {
			members = append(members, "\""+m.Inspect()+"\"")
		} else {
			members = append(members, m.Inspect())
		}
	}
	out.WriteString("[")
	out.WriteString(strings.Join(members, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "count":
		return a.Count(line, args...)
	case "filter", "grep":
		return a.Filter(line, scope, args...)
	case "index":
		return a.Index(line, args...)
	case "map":
		return a.Map(line, scope, args...)
	case "merge":
		return a.Merge(line, args...)
	case "push":
		return a.Push(line, args...)
	case "pop":
		return a.Pop(line, args...)
	case "reduce":
		return a.Reduce(line, scope, args...)
	case "empty":
		return a.Empty(line, args...)
	case "len":
		return a.Len(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, a.Type(), method))
}

func (a *Array) Len(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	return NewInteger(int64(len(a.Members)))
}

func (a *Array) Count(line string, args ...Object) Object {
	if len(args) < 1 || len(args) > 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	count := 0
	for _, v := range a.Members {
		switch c := args[0].(type) {
		case *Integer:
			if c.Int64 == v.(*Integer).Int64 {
				count++
			}
		case *String:
			if c.String == v.(*String).String {
				count++
			}
		default:
			if c == v {
				count++
			}
		}
	}
	return NewInteger(int64(count))
}

func (a *Array) Filter(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "filter", "*Function", args[0].Type()))
	}
	arr := &Array{}
	arr.Members = []Object{}
	s := NewScope(scope)
	for _, argument := range a.Members {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
		cond := Eval(block.Literal.Body, s)
		if IsTrue(cond) {
			arr.Members = append(arr.Members, argument)
		}
	}
	return arr
}

func (a *Array) Index(line string, args ...Object) Object {
	if len(args) < 1 || len(args) > 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	for i, v := range a.Members {
		switch c := args[0].(type) {
		case *Integer:
			if c.Int64 == v.(*Integer).Int64 {
				return NewInteger(int64(i))
			}
		case *String:
			if c.String == v.(*String).String {
				return NewInteger(int64(i))
			}
		default:
			if c == v {
				return NewInteger(int64(i))
			}
		}
	}
	return NIL
}

func (a *Array) Map(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "map", "*Function", args[0].Type()))
	}
	arr := &Array{}
	s := NewScope(scope)
	for _, argument := range a.Members {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
		r := Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
		arr.Members = append(arr.Members, r)
	}
	return arr
}

func (a *Array) Merge(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	m, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "merge", "*Array", args[0].Type()))
	}
	arr := &Array{}
	for _, v := range a.Members {
		arr.Members = append(arr.Members, v)
	}
	for _, v := range m.Members {
		arr.Members = append(arr.Members, v)
	}
	return arr
}

func (a *Array) Pop(line string, args ...Object) Object {
	last := len(a.Members) - 1
	if len(args) == 0 {
		if last < 0 {
			panic(NewError(line, INDEXERROR, last))
		}
		popped := a.Members[last]
		a.Members = a.Members[:last]
		return popped
	}
	idx := args[0].(*Integer).Int64
	if idx < 0 {
		idx = idx + int64(last+1)
	}
	if idx < 0 || idx > int64(last) {
		panic(NewError(line, INDEXERROR, idx))
	}
	popped := a.Members[idx]
	a.Members = append(a.Members[:idx], a.Members[idx+1:]...)
	return popped
}

func (a *Array) Push(line string, args ...Object) Object {
	l := len(args)
	if l != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", l))
	}
	a.Members = append(a.Members, args[0])
	return a
}

func (a *Array) Reduce(line string, scope *Scope, args ...Object) Object {
	l := len(args)
	if 1 != 2 && l != 1 {
		panic(NewError(line, ARGUMENTERROR, "1|2", l))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "reduce", "*Function", args[0].Type()))
	}
	s := NewScope(scope)
	start := 1
	if l == 1 {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, a.Members[0])
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, a.Members[1])
		start += 1
	} else {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, args[1])
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, a.Members[0])
	}
	r := Eval(block.Literal.Body, s)
	for i := start; i < len(a.Members); i++ {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, r)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, a.Members[i])
		r = Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
	}
	return r

}

func (a *Array) Empty(line string, args ...Object) Object {
	l := len(args)
	if l != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", l))
	}

	if len(a.Members) == 0 {
		return TRUE
	}
	return FALSE
}

//Json marshal handling
func (a *Array) MarshalJSON() ([]byte, error) {
	if len(a.Members) == 0 {
		return json.Marshal(nil)
	}

	var out bytes.Buffer

	out.WriteString("[")
	for idx, v := range a.Members {
		if idx != 0 {
			out.WriteString(",")
		}

		res, err := marshalJsonObject(v)
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(res.String())
	} //end for
	out.WriteString("]")

	return out.Bytes(), nil
}

func (a *Array) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		a = &Array{}
		return nil
	}

	var obj interface{}
	err := json.Unmarshal(b, &obj)
	if err != nil {
		a = &Array{}
		return err
	}

	if _, ok := obj.([]interface{}); !ok {
		a = &Array{}
		return errors.New("object is not a array")
	}

	ret, err := unmarshalArray(obj.([]interface{}))
	a = ret.(*Array)
	return nil
}
