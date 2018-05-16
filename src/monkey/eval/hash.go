package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) iter() bool { return true }

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		var key, val string
		if pair.Key.Type() == STRING_OBJ {
			key = "\"" + pair.Key.Inspect() + "\""
		} else {
			key = pair.Key.Inspect()
		}

		if pair.Value.Type() == STRING_OBJ {
			val = "\"" + pair.Value.Inspect() + "\""
		} else {
			val = pair.Value.Inspect()
		}
		pairs = append(pairs, fmt.Sprintf("%s : %s", key, val))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *Hash) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "filter":
		return h.Filter(line, scope, args...)
	case "keys":
		return h.Keys(line, args...)
	case "map":
		return h.Map(line, scope, args...)
	case "merge":
		return h.Merge(line, args...)
	case "pop":
		return h.Pop(line, args...)
	case "push":
		return h.Push(line, args...)
	case "values":
		return h.Values(line, args...)
	}

	panic(NewError(line, NOMETHODERROR, method, h.Type()))
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Bool {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Int64)}
}

func (u *UInteger) HashKey() HashKey {
	return HashKey{Type: u.Type(), Value: u.UInt64}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.String))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (t *Tuple) HashKey() HashKey {
	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash uint64 = 0
	for _, v := range t.Members {
		hashable, ok := v.(Hashable)
		if !ok {
			errStr := fmt.Sprintf("key error: type %s is not hashable", v.Type())
			panic(errStr)
		}

		h := hashable.HashKey()

		hash += h.Value
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15

	return HashKey{Type: t.Type(), Value: hash}
}

func (h *Hash) Filter(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "filter", "*Function", args[0].Type()))
	}
	if len(block.Literal.Parameters) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(block.Literal.Parameters)))
	}
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	s := NewScope(scope)
	for _, argument := range h.Pairs {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument.Key)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, argument.Value)
		cond := Eval(block.Literal.Body, s)
		if IsTrue(cond) {
			hash.Push(line, argument.Key, argument.Value)
		}
	}
	return hash

}

func (h *Hash) Keys(line string, args ...Object) Object {
	keys := &Array{}
	for _, pair := range h.Pairs {
		keys.Members = append(keys.Members, pair.Key)
	}
	return keys
}

func (h *Hash) Map(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "map", "*Function", args[0].Type()))
	}
	if len(block.Literal.Parameters) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(block.Literal.Parameters)))
	}
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	s := NewScope(scope)
	for _, argument := range h.Pairs {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument.Key)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, argument.Value)
		r := Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
		rh, ok := r.(*Hash)
		if !ok {
			NewError(line, RTERROR, HASH_OBJ)
		}
		for _, v := range rh.Pairs {
			hash.Push(line, v.Key, v.Value)
		}
	}
	return hash
}

func (h *Hash) Merge(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	m, ok := args[0].(*Hash)
	if !ok {
		panic(NewError(line, ARGUMENTERROR, args[0].Type(), "hash.merge"))
	}
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	for _, v := range h.Pairs {
		hash.Push(line, v.Key, v.Value)
	}
	for _, v := range m.Pairs {
		hash.Push(line, v.Key, v.Value)
	}
	return hash
}

func (h *Hash) Pop(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	hashable, ok := args[0].(Hashable)
	if !ok {
		panic(NewError(line, KEYERROR, args[0].Type()))
	}
	if hashPair, ok := h.Pairs[hashable.HashKey()]; ok {
		delete(h.Pairs, hashable.HashKey())
		return hashPair.Value
	}
	return NIL
}

func (h *Hash) Push(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}
	if hashable, ok := args[0].(Hashable); ok {
		h.Pairs[hashable.HashKey()] = HashPair{Key: args[0], Value: args[1]}
	} else {
		panic(NewError(line, KEYERROR, args[0].Type()))
	}
	return h
}

func (h *Hash) Values(line string, args ...Object) Object {
	values := &Array{}
	for _, pair := range h.Pairs {
		values.Members = append(values.Members, pair.Value)
	}
	return values
}

//Json marshal handling
func (h *Hash) MarshalJSON() ([]byte, error) {
	if len(h.Pairs) == 0 {
		return json.Marshal(nil)
	}

	var out bytes.Buffer

	var first = true
	out.WriteString("{")
	for _, pair := range h.Pairs {

		if first {
			first = false
		} else {
			out.WriteString(",")
		}

		var res bytes.Buffer
		var err error
		res, err = marshalJsonObject(pair.Key)
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(res.String())

		out.WriteString(":")

		res, err = marshalJsonObject(pair.Value)
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(res.String())

	} //end for
	out.WriteString("}")

	return out.Bytes(), nil
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		h = &Hash{Pairs: make(map[HashKey]HashPair)}
		return nil
	}

	var obj interface{}
	err := json.Unmarshal(b, &obj)
	if err != nil {
		h = &Hash{Pairs: make(map[HashKey]HashPair)}
		return err
	}

	if _, ok := obj.(map[string]interface{}); !ok {
		h = &Hash{Pairs: make(map[HashKey]HashPair)}
		return errors.New("object is not a hash")
	}

	ret, err := unmarshalHash(obj.(map[string]interface{}))
	h = ret.(*Hash)
	return nil
}
