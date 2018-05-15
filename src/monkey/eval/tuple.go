package eval

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"monkey/ast"
	"reflect"
	"strings"
)

/*
	Most of the code are copied from array.go with minor modifications.
*/

type Tuple struct {
	 // Used in function return values.
	 // if a function returns multiple values, they will wrap the results into a tuple,
	 // the flag will be set to true
	IsMulti bool
	Members []Object
}

func (t *Tuple) iter() bool { return true }

func (t *Tuple) Inspect() string {
	var out bytes.Buffer
	members := []string{}
	for _, m := range t.Members {
		if m.Type() == STRING_OBJ {
			members = append(members, "\""+m.Inspect()+"\"")
		} else {
			members = append(members, m.Inspect())
		}
	}
	out.WriteString("(")
	out.WriteString(strings.Join(members, ", "))
	out.WriteString(")")

	return out.String()
}

func (t *Tuple) Type() ObjectType { return TUPLE_OBJ }

func (t *Tuple) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "count":
		return t.Count(line, args...)
	case "filter", "grep":
		return t.Filter(line, scope, args...)
	case "index":
		return t.Index(line, args...)
	case "map":
		return t.Map(line, scope, args...)
	case "merge":
		return t.Merge(line, args...)
	case "reduce":
		return t.Reduce(line, scope, args...)
	case "empty":
		return t.Empty(line, args...)
	case "len":
		return t.Len(line, args...)
	case "first", "head":
		return t.First(line, args...)
	case "last":
		return t.Last(line, args...)
	case "tail","rest":
		return t.Tail(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, t.Type()))
}

func (t *Tuple) Len(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	return NewInteger(int64(len(t.Members)))
}

func (t *Tuple) Count(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	count := 0
	for _, v := range t.Members {
		switch c := args[0].(type) {
		case *Integer:
			if c.Int64 == v.(*Integer).Int64 {
				count++
			}
		case *UInteger:
			if c.UInt64 == v.(*UInteger).UInt64 {
				count++
			}
		case *Float:
			if c.Float64 == v.(*Float).Float64 {
				count++
			}
		case *String:
			if c.String == v.(*String).String {
				count++
			}
		case *Boolean:
			if c.Bool == v.(*Boolean).Bool {
				count++
			}
		default:
			r := reflect.DeepEqual(c, v)
			if r {
				count++
			}
		}
	}
	return NewInteger(int64(count))
}

func (t *Tuple) Filter(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "filter", "*Function", args[0].Type()))
	}

	tuple := &Tuple{}
	tuple.Members = []Object{}
	s := NewScope(scope)
	for _, argument := range t.Members {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
		cond := Eval(block.Literal.Body, s)
		if IsTrue(cond) {
			tuple.Members = append(tuple.Members, argument)
		}
	}
	return tuple
}

func (t *Tuple) Index(line string, args ...Object) Object {
	if len(args) < 1 || len(args) > 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	for i, v := range t.Members {
		switch c := args[0].(type) {
		case *Integer:
			if c.Int64 == v.(*Integer).Int64 {
				return NewInteger(int64(i))
			}
		case *UInteger:
			if c.UInt64 == v.(*UInteger).UInt64 {
				return NewInteger(int64(i))
			}
		case *String:
			if c.String == v.(*String).String {
				return NewInteger(int64(i))
			}
		default:
			r := reflect.DeepEqual(c, v)
			if r {
				return NewInteger(int64(i))
			}
		}
	}
	return NIL
}

func (t *Tuple) Map(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "map", "*Function", args[0].Type()))
	}
	
	tuple := &Tuple{}
	s := NewScope(scope)
	for _, argument := range t.Members {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument)
		r := Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
		tuple.Members = append(tuple.Members, r)
	}
	return tuple
}

func (t *Tuple) Merge(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	m, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "merge", "*Array", args[0].Type()))
	}

	tuple := &Tuple{}
	for _, v := range t.Members {
		tuple.Members = append(tuple.Members, v)
	}
	for _, v := range m.Members {
		tuple.Members = append(tuple.Members, v)
	}
	return tuple
}

func (t *Tuple) Reduce(line string, scope *Scope, args ...Object) Object {
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
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, t.Members[0])
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, t.Members[1])
		start += 1
	} else {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, args[1])
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, t.Members[0])
	}
	r := Eval(block.Literal.Body, s)
	for i := start; i < len(t.Members); i++ {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, r)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, t.Members[i])
		r = Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
	}
	return r

}

func (t *Tuple) Empty(line string, args ...Object) Object {
	l := len(args)
	if l != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", l))
	}

	if len(t.Members) == 0 {
		return TRUE
	}
	return FALSE
}

func (t *Tuple) First(line string, args ...Object) Object {
	l := len(args)
	if l != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", l))
	}

	if len(t.Members) == 0 {
		return NIL
	}
	return t.Members[0]
}

func (t *Tuple) Last(line string, args ...Object) Object {
	l := len(args)
	if l != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", l))
	}

	length := len(t.Members)
	if length == 0 {
		return NIL
	}
	return t.Members[length - 1]
}

func (t *Tuple) Tail(line string, args ...Object) Object {
	l := len(args)
	if l != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", l))
	}

	length := len(t.Members)
	if length == 0 {
		return NIL
	}

	newMembers := make([]Object, length+1, length+1)
	copy(newMembers, t.Members)
	return &Tuple{Members: newMembers}
}

//Json marshal handling
func (t *Tuple) MarshalJSON() ([]byte, error) {
	if len(t.Members) == 0 {
		return json.Marshal(nil)
	}

	var out bytes.Buffer

	out.WriteString("[")  //NOTE HERE: we do not use "(", because json doesn't support tuple
	for idx, v := range t.Members {
		if idx != 0 {
			out.WriteString(",")
		}

		res, err := marshalJsonObject(v)
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(res.String())
	} //end for
	out.WriteString("]") //NOTE HERE: we do not use "(", because json doesn't support tuple
	
	return out.Bytes(), nil
}

//NO UnmarshalJSON method for Tuple, the tuple will be treated as an array.
//func (t *Tuple) UnmarshalJSON(b []byte) error {
//
//}

