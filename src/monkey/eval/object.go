package eval

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"monkey/ast"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

var formatMap = map[int]string{
    0: "",
    1: "%v",
}

var colorMap = map[string]string{
	"STRING": "31", //red
	"NUMBER": "32", //green
	"ARRAY":  "33", //yellow
	"HASH":   "34", //blue
	"TUPLE":  "35", //purple(magenta)
	"BOOL":   "36", //cyan
}

type ObjectType string

const (
	availFlags = "-+# 0"

	UTC   = 0
	LOCAL = 1

	INTEGER_OBJ      = "INTEGER"
	UINTEGER_OBJ     = "UINTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	BREAK_OBJ        = "BREAK"
	CONTINUE_OBJ     = "CONTINUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	TUPLE_OBJ        = "TUPLE"
	HASH_OBJ         = "HASH"
	INCLUDED_OBJ     = "INCLUDE"
	STRUCT_OBJ       = "STRUCT"
	ENUM_OBJ         = "ENUM"
	FILE_OBJ         = "FILE"
	REGEX_OBJ        = "REGEX"
	CHANNEL_OBJ = "CHANNEL"
	NIL_OBJ     = "NIL_OBJ"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	CallMethod(line string, scope *Scope, method string, args ...Object) Object
}

//Whether the Object is a number (INT, FLOAT)
type Number interface {
	number()
}

//Whether the Object is iterable (HASH, ARRAY, RANGE, STRING, TUPLE)
type Iterable interface {
	iter()
}

//Whether the Object is throwable (STRING for now)
type Throwable interface {
	throw()
}

//Whether the Object is the target of IO writer
type Writable interface {
	IOWriter() io.Writer
}

type Hashable interface {
	HashKey() HashKey
}

type Struct struct {
	Scope   *Scope
	methods map[string]*Function
}

func (s *Struct) Inspect() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for k, v := range s.Scope.store {
		out.WriteString(k)
		out.WriteString("=>")
		out.WriteString(v.Inspect())
		out.WriteString(",")
	}
	out.WriteString(" }")

	return out.String()
}

func (s *Struct) Type() ObjectType { return STRUCT_OBJ }
func (s *Struct) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	fn, ok := s.methods[method]
	if !ok {
		panic(NewError(line, NOMETHODERROR, method, s.Type()))
	}
	fn.Scope = NewScope(scope)
	fn.Scope.Set("this", s)
	for i, v := range fn.Literal.Parameters {
		fn.Scope.Set(v.String(), args[i])
	}

	r := Eval(fn.Literal.Body, fn.Scope)
	if obj, ok := r.(*ReturnValue); ok {
		return obj.Value
	}
	return r
}

type Enum struct {
	Scope *Scope
}

func (e *Enum) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	for k, v := range e.Scope.store {
		out.WriteString(k)
		out.WriteString("=")
		out.WriteString(v.Inspect())
		out.WriteString(",")
	}
	out.WriteString(" }")

	return out.String()
}

func (e *Enum) Type() ObjectType { return ENUM_OBJ }
func (e *Enum) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "getName":
		return e.GetName(line, args...)
	case "getNames":
		return e.GetNames(line, args...)
	case "getValues":
		return e.GetValues(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, e.Type()))
}

func (e *Enum) GetNames(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := &Array{}
	for k, _ := range e.Scope.store {
		ret.Members = append(ret.Members, NewString(k))
	}
	return ret
}

func (e *Enum) GetValues(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := &Array{}
	for _, v := range e.Scope.store {
		ret.Members = append(ret.Members, v)
	}
	return ret
}

func (e *Enum) GetName(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	for k, v := range e.Scope.store {
		if equal(true, args[0], v) {
			return NewString(k)
		}
	}
	return NIL
}

func (b *Builtin) Inspect() string  { return "<builtin function>" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, b.Type()))
}

type IncludedObject struct {
	Name  string
	Scope *Scope
}

func (io *IncludedObject) Inspect() string  { return fmt.Sprintf("included object: %s", io.Name) }
func (io *IncludedObject) Type() ObjectType { return INCLUDED_OBJ }
func (io *IncludedObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, io.Type()))
}

type Function struct {
	Literal  *ast.FunctionLiteral
	Variadic bool
	Scope    *Scope
	Instance *ObjectInstance //For use with class functions
	Annotations []*ObjectInstance
}

func (f *Function) Inspect() string  { return f.Literal.String() }
func (r *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) classMethod() ast.ModifierLevel { 
	return f.Literal.ModifierLevel
}

func (f *Function) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, f.Type()))
}

type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, rv.Type()))
}

//type ThrowValue struct{ Value Object }
//
//func (tv *ThrowValue) Inspect() string  { return tv.Value.Inspect() }
//func (tv *ThrowValue) Type() ObjectType { return THROW_OBJ }
//func (tv *ThrowValue) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
//	panic(NewError(line, NOMETHODERROR, method, tv.Type()))
//}

//return a Nil object with error message 's'
func NewNil(s string) *Nil {
	return &Nil{OptionalMsg: s}
}

type Nil struct {
	//sometimes when a function fails, it will return NIL. If this happens, we also need to
	//know the error reason. The error message is stored in `OptionalMsg`
	OptionalMsg string
}

func (n *Nil) Inspect() string {
	if n.OptionalMsg != "" {
		return n.OptionalMsg
	}
	return "nil"
}
func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "message":
		return n.Message(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, n.Type()))
}

func (n *Nil) Message(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(n.OptionalMsg)
}

//Json marshal handling
func (n *Nil) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}

func (n *Nil) UnmarshalJSON(b []byte) error {
	return nil
}

//Returns a valid Integer Object, that is Valid=true
func NewInteger(i int64) *Integer {
	return &Integer{Int64: i, Valid: true}
}

type Integer struct {
	Int64 int64
	Valid bool
}

func (i *Integer) Inspect() string {
	if i.Valid {
		return fmt.Sprintf("%d", i.Int64)
	}
	return "ERROR: Integer is null"
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) number()          {}
func (i *Integer) CallMethod(line string, scope *Scope, method string, args ...Object) Object {

	switch method {
	case "valid", "isValid":
		return i.IsValid(line, args...)
	case "setValid":
		return i.SetValid(line, args...)
	case "next":
		return i.Next(line, args...)
	case "prev":
		return i.Prev(line, args...)
	case "isEven":
		return i.IsEven(line, args...)
	case "isOdd":
		return i.IsOdd(line, args...)
	case "downto":
		return i.Downto(line, args...)
	case "upto":
		return i.Upto(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, i.Type()))
}

func (i *Integer) IsValid(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		return TRUE
	}
	return &Boolean{Bool: i.Valid, Valid: false}
}

func (i *Integer) SetValid(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 0 && argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", argLen))
	}

	if argLen == 0 {
		i.Int64, i.Valid = 0, true
		return i
	}

	val, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setValid", "*Integer", args[0].Type()))
	}

	i.Int64, i.Valid = val.Int64, true
	return i
}

func (i *Integer) Next(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		return NewInteger(i.Int64 + 1)
	}
	return NewFalseObj("Integer is not valid\n")
}

func (i *Integer) Prev(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		return NewInteger(i.Int64 - 1)
	}
	return NewFalseObj("Integer is not valid\n")
}

func (i *Integer) IsEven(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		if i.Int64 % 2 == 0 {
			return TRUE
		}
		return FALSE
	}
	return NewFalseObj("Integer is not valid\n")
}

func (i *Integer) IsOdd(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		if i.Int64 % 2 != 0 {
			return TRUE
		}
		return FALSE
	}
	return NewFalseObj("Integer is not valid\n")
}

func (i *Integer) Downto(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", argLen))
	}

	val, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "downto", "*Integer", args[0].Type()))
	}

	retArr := &Array{}
	for x := i.Int64; x >= val.Int64; x-- {
		retArr.Members = append(retArr.Members, NewInteger(x))
	}
	return retArr
}

func (i *Integer) Upto(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", argLen))
	}

	val, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "upto", "*Integer", args[0].Type()))
	}

	retArr := &Array{}
	for x := i.Int64; x <= val.Int64; x++ {
		retArr.Members = append(retArr.Members, NewInteger(x))
	}
	return retArr
}

//Implements sql's Scanner Interface.
//So when calling sql.Rows.Scan(xxx), or sql.Row.Scan(xxx), we could pass this object to `Scan` method
func (i *Integer) Scan(value interface{}) error {
	if value == nil {
		i.Valid = false
		return nil
	}
	i.Int64, i.Valid = value.(int64), true
	return nil
}

//Implements driver's Valuer Interface.
//So when calling sql.Exec(xx), we could pass this object to `Exec` method
func (i Integer) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

//Json marshal handling
func (i *Integer) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return []byte(fmt.Sprintf("%v", i.Int64)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (i *Integer) UnmarshalJSON(b []byte) error {
	content := string(b)
	if content == "null" {
		i.Valid = false
		i.Int64 = 0
		return nil
	}

	var err error

	if content[len(content)-1] == 'u' {
		content = content[:len(content)-1]
	}
	if strings.HasPrefix(content, "0b") {
		i.Int64, err = strconv.ParseInt(content[2:], 2, 64)
	} else if strings.HasPrefix(content, "0x") {
		i.Int64, err = strconv.ParseInt(content[2:], 16, 64)
	} else if strings.HasPrefix(content, "0c") {
		i.Int64, err = strconv.ParseInt(content[2:], 8, 64)
	} else {
		i.Int64, err = strconv.ParseInt(content, 10, 64)
	}
	if err != nil {
		i.Valid = false
		return err
	}

	i.Valid = true
	return nil
}

//Returns a valid Unsigned Integer Object, that is Valid=true
func NewUInteger(i uint64) *UInteger {
	return &UInteger{UInt64: i, Valid: true}
}

type UInteger struct {
	UInt64 uint64
	Valid bool
}

func (i *UInteger) Inspect() string {
	if i.Valid {
		return fmt.Sprintf("%d", i.UInt64)
	}
	return "ERROR: Unsigned Integer is null"
}

func (i *UInteger) Type() ObjectType { return UINTEGER_OBJ }
func (i *UInteger) number()          {}
func (i *UInteger) CallMethod(line string, scope *Scope, method string, args ...Object) Object {

	switch method {
	case "valid", "isValid":
		return i.IsValid(line, args...)
	case "setValid":
		return i.SetValid(line, args...)
	case "next":
		return i.Next(line, args...)
	case "prev":
		return i.Prev(line, args...)
	case "isEven":
		return i.IsEven(line, args...)
	case "isOdd":
		return i.IsOdd(line, args...)
	case "downto":
		return i.Downto(line, args...)
	case "upto":
		return i.Upto(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, i.Type()))
}

func (i *UInteger) IsValid(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		return TRUE
	}
	return &Boolean{Bool: i.Valid, Valid: false}
}

func (i *UInteger) SetValid(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 0 && argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", argLen))
	}

	if argLen == 0 {
		i.UInt64, i.Valid = 0, true
		return i
	}

	val, ok := args[0].(*UInteger)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setValid", "*UInteger", args[0].Type()))
	}

	i.UInt64, i.Valid = val.UInt64, true
	return i
}

func (i *UInteger) Next(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		return NewUInteger(i.UInt64 + 1)
	}
	return NewFalseObj("Unsigned Integer is not valid\n")
}

func (i *UInteger) Prev(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		return NewUInteger(i.UInt64 - 1)
	}
	return NewFalseObj("Unsigned Integer is not valid\n")
}

func (i *UInteger) IsEven(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		if i.UInt64 % 2 == 0 {
			return TRUE
		}
		return FALSE
	}
	return NewFalseObj("Unsigned Integer is not valid\n")
}

func (i *UInteger) IsOdd(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if i.Valid {
		if i.UInt64 % 2 != 0 {
			return TRUE
		}
		return FALSE
	}
	return NewFalseObj("Unsigned Integer is not valid\n")
}

func (i *UInteger) Downto(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", argLen))
	}

	val, ok := args[0].(*UInteger)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "downto", "*UInteger", args[0].Type()))
	}

	retArr := &Array{}
	for x := i.UInt64; x >= val.UInt64; x-- {
		retArr.Members = append(retArr.Members, NewUInteger(x))
	}
	return retArr
}

func (i *UInteger) Upto(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", argLen))
	}

	val, ok := args[0].(*UInteger)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "upto", "*UInteger", args[0].Type()))
	}

	retArr := &Array{}
	for x := i.UInt64; x <= val.UInt64; x++ {
		retArr.Members = append(retArr.Members, NewUInteger(x))
	}
	return retArr
}

//Implements sql's Scanner Interface.
//So when calling sql.Rows.Scan(xxx), or sql.Row.Scan(xxx), we could pass this object to `Scan` method
func (i *UInteger) Scan(value interface{}) error {
	if value == nil {
		i.Valid = false
		return nil
	}
	i.UInt64, i.Valid = value.(uint64), true
	return nil
}

//Implements driver's Valuer Interface.
//So when calling sql.Exec(xx), we could pass this object to `Exec` method
func (i UInteger) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.UInt64, nil
}

//Json marshal handling
func (i *UInteger) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return []byte(fmt.Sprintf("%v", i.UInt64)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (i *UInteger) UnmarshalJSON(b []byte) error {
	content := string(b)
	if content == "null" {
		i.Valid = false
		i.UInt64 = 0
		return nil
	}

	var err error

	if content[len(content)-1] == 'u' {
		content = content[:len(content)-1]
	}
	if strings.HasPrefix(content, "0b") {
		i.UInt64, err = strconv.ParseUint(content[2:], 2, 64)
	} else if strings.HasPrefix(content, "0x") {
		i.UInt64, err = strconv.ParseUint(content[2:], 16, 64)
	} else if strings.HasPrefix(content, "0c") {
		i.UInt64, err = strconv.ParseUint(content[2:], 8, 64)
	} else {
		i.UInt64, err = strconv.ParseUint(content, 10, 64)
	}
	if err != nil {
		i.Valid = false
		return err
	}

	i.Valid = true
	return nil
}

//Returns a valid Float Object, that is Valid=true
func NewFloat(f float64) *Float {
	return &Float{Float64: f, Valid: true}
}

type Float struct {
	Float64 float64
	Valid   bool
}

func (f *Float) Inspect() string {
	if f.Valid {
		return fmt.Sprintf("%g", f.Float64)
	}
	return "ERROR: Float is null"
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) number()          {}
func (f *Float) CallMethod(line string, scope *Scope, method string, args ...Object) Object {

	switch method {
	case "valid", "isValid":
		return f.IsValid(line, args...)
	case "setValid":
		return f.SetValid(line, args...)
	case "ceil":
		return f.Ceil(line, args...)
	case "floor":
		return f.Floor(line, args...)
	case "trunc":
		return f.Trunc(line, args...)
	case "sqrt":
		return f.Sqrt(line, args...)
	case "pow":
		return f.Pow(line, args...)
	case "round":
		return f.Round(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, f.Type()))
}

func (f *Float) IsValid(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if f.Valid {
		return TRUE
	}
	return &Boolean{Bool: f.Valid, Valid: false}
}

func (f *Float) SetValid(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 0 && argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", argLen))
	}

	if argLen == 0 {
		f.Float64, f.Valid = 0.0, true
		return f
	}

	val, ok := args[0].(*Float)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setValid", "*Float", args[0].Type()))
	}

	f.Float64, f.Valid = val.Float64, true
	return f
}

func (f *Float) Ceil(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if f.Valid {
		return NewFloat(math.Ceil(f.Float64))
	}
	return NewFalseObj("Float is not valid\n")
}

func (f *Float) Floor(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if f.Valid {
		return NewFloat(math.Floor(f.Float64))
	}

	return NewFalseObj("Float is not valid\n")
}

func (f *Float) Trunc(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if f.Valid {
		return NewFloat(math.Trunc(f.Float64))
	}

	return NewFalseObj("Float is not valid\n")
}

func (f *Float) Sqrt(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if f.Valid {
		return NewFloat(math.Sqrt(f.Float64))
	}

	return NewFalseObj("Float is not valid\n")
}

func (f *Float) Pow(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var temp float64
	switch input := args[0].(type) {
	case *Integer:
		temp = float64(input.Int64)
	case *UInteger:
		temp = float64(input.UInt64)
	case *Float:
		temp = input.Float64
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "pow", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	if f.Valid {
		return NewFloat(math.Pow(f.Float64, temp))
	}

	return NewFalseObj("Float is not valid\n")
}

func (f *Float) Round(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var precision int64
	switch o := args[0].(type) {
	case *Integer:
		precision = o.Int64
	case *UInteger:
		precision = int64(o.UInt64)
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "round", "*Integer|*UInteger", args[0].Type()))
	}

	format := fmt.Sprintf("%%.%df", precision) //'%.xf', x is the precision, e.g. %.2f
	resultStr := fmt.Sprintf(format, f.Float64)   //convert to string
	ret, err := strconv.ParseFloat(resultStr, 64) //convert string back to float
	if err != nil {
		return NewFloat(math.NaN())
	}
	return NewFloat(ret)
}

func (f *Float) Scan(value interface{}) error {
	if value == nil {
		f.Valid = false
		return nil
	}
	f.Float64, f.Valid = value.(float64), true
	return nil
}

func (f Float) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

func (f *Float) MarshalJSON() ([]byte, error) {
	if f.Valid {
		return []byte(fmt.Sprintf("%v", f.Float64)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (f *Float) UnmarshalJSON(b []byte) error {
	content := string(b)
	if content == "null" {
		f.Valid = false
		f.Float64 = 0.0
		return nil
	}

	var err error
	f.Float64, err = strconv.ParseFloat(content, 64)
	if err != nil {
		f.Valid = false
		return err
	}
	f.Valid = true
	return nil
}

func NewFalseObj(s string) *Boolean {
	return &Boolean{Bool: false, Valid: true, OptionalMsg: s}
}

type Boolean struct {
	Bool  bool
	Valid bool
	//sometimes when a function fails, it will return `false`. If this happens, we also need to
	//know the error reason. The error message is stored in `OptionalMsg`
	OptionalMsg string
}

func (b *Boolean) Inspect() string {
	if b.Valid {
		if b.Bool == false && b.OptionalMsg != "" {
			return b.OptionalMsg
		}
		return fmt.Sprintf("%v", b.Bool)
	}
	return "ERROR: Boolean is null"
}

func (b *Boolean) Scan(value interface{}) error {
	if value == nil {
		b.Valid = false
		return nil
	}
	b.Bool, b.Valid = value.(bool), true
	return nil
}

func (b Boolean) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) CallMethod(line string, scope *Scope, method string, args ...Object) Object {

	switch method {
	case "valid", "isValid":
		return b.IsValid(line, args...)
	case "setValid":
		return b.SetValid(line, args...)
	case "message":
		return b.Message(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, b.Type()))
}

func (b *Boolean) Message(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(b.OptionalMsg)
}

func (b *Boolean) IsValid(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	if b.Valid {
		return TRUE
	}
	return &Boolean{Bool: b.Valid, Valid: false}
}

func (b *Boolean) SetValid(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	val, ok := args[0].(*Boolean)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setValid", "*Boolean", args[0].Type()))
	}

	b.Bool, b.Valid = val.Bool, true

	return b
}

func (b *Boolean) MarshalJSON() ([]byte, error) {
	if b.Valid {
		if b.Bool {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	} else {
		return json.Marshal(nil)
	}
}

func (b *Boolean) UnmarshalJSON(bytes []byte) error {
	str := string(bytes)
	if str == "null" {
		b.Valid = false
		b.Bool = false
		return nil
	}

	if strings.ToLower(str) == "true" || str == "1" {
		b.Valid = true
		b.Bool = true
		return nil
	}

	if strings.ToLower(str) == "false" || str == "0" || str == "" {
		b.Valid = true
		b.Bool = false
		return nil
	}
	return errors.New(string(bytes) + " is not a valid JSON bool")
}

type Break struct{}

func (b *Break) Inspect() string  { return "break" }
func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, b.Type()))
}

type Continue struct{}

func (c *Continue) Inspect() string  { return "continue" }
func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, c.Type()))
}

func initGlobalObj() {
	//Predefine `stdin`, `stdout`, `stderr`
	SetGlobalObj("stdin", &FileObject{File: os.Stdin})
	SetGlobalObj("stdout", &FileObject{File: os.Stdout})
	SetGlobalObj("stderr", &FileObject{File: os.Stderr})

	//runtime
	SetGlobalObj("RUNTIME_ARCH", NewString(runtime.GOARCH)) //running program's architecture target: one of 386, amd64, arm, s390x, and so on
	SetGlobalObj("RUNTIME_OS", NewString(runtime.GOOS))     //running program's operating system target: one of darwin, freebsd, linux, and so on.
}

func init() {
	initGlobalObj()

	NewOsObj()
	NewNetObj()
	NewHTTPObj()
	NewTimeObj()
	NewMathObj()
	NewJsonObj()
	NewFlagObj()
	NewFilePathObj()
	NewIOUtilObj()
	NewFmtObj()
	NewLoggerObj()
	NewStringsObj()
	NewSortObj()
	NewSqlsObject()
	NewLinqObj()
	NewRegExpObj()
	NewTemplateObj()
	NewDecimalObj()
	NewUnicodeObj()
}

func marshalJsonObject(obj interface{}) (bytes.Buffer, error) {
	if obj == nil {
		return bytes.Buffer{}, errors.New("json error: maybe unsupported type or invalid data")
	}

	var out bytes.Buffer
	switch obj.(type) {
	case *Integer:
		value := obj.(*Integer)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	case *UInteger:
		value := obj.(*UInteger)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	case *Float:
		value := obj.(*Float)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	case *String:
		value := obj.(*String)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	case *Boolean:
		value := obj.(*Boolean)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	case *Array:
		value := obj.(*Array)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	case *Hash:
		value := obj.(*Hash)
		res, err := value.MarshalJSON()
		if err != nil {
			return bytes.Buffer{}, err
		}
		out.WriteString(string(res))
	default:
		return bytes.Buffer{}, errors.New("json error: maybe unsupported type or invalid data")
	}
	return out, nil
}

func unmarshalJsonObject(val interface{}) (Object, error) {
	var ret Object
	var err error = nil
	switch val.(type) {
	case []interface{}:
		ret, err = unmarshalArray(val.([]interface{}))
	case map[string]interface{}:
		ret, err = unmarshalHash(val.(map[string]interface{}))
	case float64:
		ret = NewFloat(val.(float64))
	case bool:
		b := val.(bool)
		if b {
			ret = TRUE
		} else {
			ret = FALSE
		}
	case string:
		ret = NewString(val.(string))
	case nil:
		ret = NIL
	}
	return ret, err
}

func unmarshalArray(a []interface{}) (Object, error) {
	arr := &Array{}

	for _, v := range a {
		item, err := unmarshalJsonObject(v)
		if err != nil {
			return FALSE, err
		}
		arr.Members = append(arr.Members, item)
	}

	return arr, nil
}

func unmarshalHash(m map[string]interface{}) (Object, error) {
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}

	for key, value := range m {
		keyObj, err := unmarshalJsonObject(key)
		if err != nil {
			return FALSE, err
		}

		valObj, err := unmarshalJsonObject(value)
		if err != nil {
			return FALSE, err
		}

		if hashable, ok := keyObj.(Hashable); ok {
			hash.Pairs[hashable.HashKey()] = HashPair{Key: keyObj, Value: valObj}
		} else {
			return FALSE, errors.New("key error: type is not hashable")
		}
	}
	return hash, nil
}

//this function is used in template.go for converting 'Object' to 'interface{}'
func object2RawValue(obj Object) interface{} {
	var ret interface{} = nil
	objType := obj.Type()
	switch objType {
	case HASH_OBJ:
		ret = hashObj2RawValue(obj.(*Hash))
	case ARRAY_OBJ:
		ret = arrayObj2RawValue(obj.(*Array))
	case INTEGER_OBJ:
		ret = obj.(*Integer).Int64
	case UINTEGER_OBJ:
		ret = obj.(*UInteger).UInt64
	case FLOAT_OBJ:
		ret = obj.(*Float).Float64
	case BOOLEAN_OBJ:
		ret = obj.(*Boolean).Bool
	case NIL_OBJ:
		ret = nil
	case STRING_OBJ:
		ret = obj.(*String).String
	default:
		panic("Could not convert to RawValue!")
	}
	return ret
}

func arrayObj2RawValue(arr *Array) interface{} {
	ret := make([]interface{}, len(arr.Members))
	for idx, v := range arr.Members {
		ret[idx] = object2RawValue(v)
	}
	return ret
}

func hashObj2RawValue(h *Hash) interface{} {
	ret := make(map[interface{}]interface{})
	for _, v := range h.Pairs{
		ret[object2RawValue(v.Key)] = object2RawValue(v.Value)
	}
	return ret
}


//This `Formatter` struct is mainly used to encapsulate golang
//`fmt` package's `Formatter` interface.
//When we implement this interface, our `Object` could be directed passed to fmt.Printf(xxx)
type Formatter struct {
	Obj Object
}

func (ft *Formatter) Format(s fmt.State, verb rune) {
	format := make([]byte, 0, 128)
	format = append(format, '%')
	var f byte
	for i := 0; i < len(availFlags); i++ {
		if f = availFlags[i]; s.Flag(int(f)) {
			format = append(format, f)
		}
	}
	var width, prec int
	var ok bool
	if width, ok = s.Width(); ok {
		format = strconv.AppendInt(format, int64(width), 10)
	}
	if prec, ok = s.Precision(); ok {
		format = append(format, '.')
		format = strconv.AppendInt(format, int64(prec), 10)
	}
	if verb > utf8.RuneSelf {
		format = append(format, string(verb)...)
	} else {
		//Here we use '%_' to print the object's type
		if verb == '_' {
			format = append(format, byte('T'))
		} else {
			format = append(format, byte(verb))
		}
	}

	if string(format) == "%T" {
		t := reflect.TypeOf(ft.Obj)
		strArr := strings.Split(t.String(), ".") //t.String() = "*eval.xxx"
		fmt.Fprintf(s, "%s", strArr[1])          //NEED CHECK for "index out of bounds?"
		return
	}

	formatStr := string(format)
	var reset = "\033[0m"
	switch obj := ft.Obj.(type) {
	case *Boolean:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["BOOL"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Bool)
	case *Nil:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["BOOL"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Inspect())
	case *Integer:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["NUMBER"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Int64)
	case *UInteger:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["NUMBER"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.UInt64)
	case *Float:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["NUMBER"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Float64)
	case *String:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["STRING"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.String)
	case *Array:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["ARRAY"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Inspect())
	case *Hash:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["HASH"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Inspect())
	case *Tuple:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["TUPLE"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Inspect())
	case *DecimalObj:
		if REPLColor {
			formatStr = "\033[1;" + colorMap["NUMBER"] + "m" + formatStr + reset
		}
		fmt.Fprintf(s, formatStr, obj.Inspect())
	default:
		fmt.Fprintf(s, formatStr, obj.Inspect())
	}
}

/*  when you call println/print like below:
		println("a=", 10)
		print("a=", 10)
	the golang will return "a= 10", not "a=10", this is not what we expected.
	the solution is take from:
		https://stackoverflow.com/questions/25928991/go-print-without-space-between-items
*/
func correctPrintResult(needNewLine bool, args ...Object) (string, []interface{}) {
	l := len(args)
	if s, isOk := formatMap[l]; !isOk {
		for i := 0; i < l; i++ {
			s += "%v"
		}
		formatMap[l] = s
	}

	s := formatMap[l]
	if needNewLine {
		s = s + "\n"
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	return s, wrapped
}

