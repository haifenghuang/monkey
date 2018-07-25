package eval

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strconv"
	"strings"
	"unicode"
)

type InterpolatedString struct {
	String      *String
	RawValue    string
	Expressions map[byte]ast.Expression
}

type Interpolable interface {
	Interpolate(scope *Scope)
}

func (is *InterpolatedString) Inspect() string  { return is.String.String }
func (is *InterpolatedString) Type() ObjectType { return STRING_OBJ }
func (is *InterpolatedString) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	return is.String.CallMethod(line, scope, method, args...)
}

func (is *InterpolatedString) Interpolate(scope *Scope) {
	var out bytes.Buffer

	objIndex := "0"[0]
	ol := len(is.Expressions)
	if ol == 0 {
		is.String.String = is.RawValue
		return
	}
	mStr := "{" + string(objIndex) + "}"
	sl := len(is.RawValue)
	ml := len(mStr)

	for i := 0; i < sl; i++ {
		if i+ml > sl {
			out.WriteString(is.RawValue[i:])
			break
		}
		if is.RawValue[i:i+ml] == mStr {
			v := is.evalInterpExpression(is.Expressions[objIndex], scope)
			out.WriteString(v)
			i += ml - 1
			objIndex++
			if (objIndex - 47) > byte(ol) {
				out.WriteString(is.RawValue[i+1:])
				break
			}
			mStr = "{" + string(objIndex) + "}"
			ml = len(mStr)
		} else {
			out.WriteByte(is.RawValue[i])
		}
	}
	is.String.String = out.String()
}

func (is *InterpolatedString) evalInterpExpression(exp ast.Expression, s *Scope) string {
	_, ok := exp.(*ast.Identifier)
	if ok {
		sv, ok := s.Get(exp.String())
		if ok {
			iss, ok := sv.(*InterpolatedString)
			if ok {
				if iss.RawValue == is.RawValue {
					return exp.String()
				}
			}
		}
	}
	evaluated := Eval(exp, s)
	if evaluated.Type() == ERROR_OBJ {
		return exp.String()
	}
	return evaluated.Inspect()
}

//Returns a valid String Object, that is Valid=true
func NewString(s string) *String {
	return &String{String: s, Valid: true}
}

type String struct {
	String string
	Valid  bool
}

func (s *String) iter() bool { return true }

func (s *String) throw() {}
func (s *String) Inspect() string {
	if s.Valid {
		return s.String
	}
	return ""
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.String))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

func (s *String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return []byte(fmt.Sprintf("\"%v\"", s.String)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (s *String) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		s.Valid = false
		s.String = ""
		return nil
	}

	s.String = string(b[1 : len(b)-1])
	s.Valid = true
	return nil
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) CallMethod(line string, scope *Scope, method string, args ...Object) Object {

	switch method {
	case "find", "index":
		return s.Find(line, args...)
	case "lower":
		return s.Lower(line, args...)
	case "reverse":
		return s.Reverse(line, args...)
	case "upper":
		return s.Upper(line, args...)
	case "trimLeft":
		return s.TrimLeft(line, args...)
	case "trimRight":
		return s.TrimRight(line, args...)
	case "trim":
		return s.Trim(line, args...)
	case "trimPrefix":
		return s.TrimPrefix(line, args...)
	case "trimSuffix":
		return s.TrimSuffix(line, args...)
	case "split":
		return s.Split(line, args...)
	case "replace":
		return s.Replace(line, args...)
	case "count":
		return s.Count(line, args...)
		//	case "join":
		//		return s.Join(line, args...)
	case "substr":
		return s.SubStr(line, args...)
	case "len":
		return s.Len(line, args...)
	case "compare":
		return s.Compare(line, args...)
	case "contains":
		return s.Contains(line, args...)
	case "containsAny":
		return s.ContainsAny(line, args...)
	case "fields":
		return s.Fields(line, args...)
	case "hasPrefix":
		return s.HasPrefix(line, args...)
	case "hasSuffix":
		return s.HasSuffix(line, args...)
	case "lastIndex":
		return s.LastIndex(line, args...)
	case "repeat":
		return s.Repeat(line, args...)
	case "title":
		return s.Title(line, args...)
	case "chomp":
		return s.Chomp(line, args...)
	case "parseInt":
		return s.ParseInt(line, args...)
	case "parseUInt":
		return s.ParseUInt(line, args...)
	case "parseBool":
		return s.ParseBool(line, args...)
	case "parseFloat":
		return s.ParseFloat(line, args...)
	case "atoi":
		return s.Atoi(line, args...)
	case "itoa":
		return s.Itoa(line, args...)
	case "writeLine":
		return s.WriteLine(line, args...)
	case "write":
		return s.Write(line, args...)
	case "isEmpty":
		return s.IsEmpty(line, args...)
	case "hash":
		return s.Hash(line, args...)
	case "valid", "isValid":
		return s.IsValid(line, args...)
	case "setValid":
		return s.SetValid(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, s.Type()))
}

func (s *String) Count(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	sub, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "count", "*String", args[0].Type()))
	}

	count := strings.Count(s.String, sub.String)
	return NewInteger(int64(count))
}

func (s *String) Find(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}
	sub, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "find", "*String", args[0].Type()))
	}

	idx := strings.Index(s.String, sub.String)
	return NewInteger(int64(idx))
}

//join(arr)  or  join(arr, sep)
//func (s *String) Join(line string, args ...Object) Object {
//	argLen := len(args)
//	if argLen != 1 && argLen != 2 {
//		panic(NewError(line, ARGUMENTERROR, "1|2", argLen))
//	}
//
//	sep := ""
//	a, ok := args[0].(*Array)
//	if !ok {
//		panic(NewError(line, PARAMTYPEERROR, "first", "join", "*Array", args[0].Type()))
//	}
//
//	if argLen == 2 {
//		v, ok := args[1].(*String)
//		if !ok {
//			panic(NewError(line, PARAMTYPEERROR, "second", "join", "*String", args[1].Type()))
//		}
//		sep = v.String
//	}
//
//	var tmp []string
//	for _, item := range a.Members {
//		tmp = append(tmp, item.Inspect())
//	}
//
//	ret := strings.Join(tmp, sep)
//	return NewString(ret)
//}
//

func isWhiteSpace(a rune) bool {
	return unicode.IsSpace(a)
}

func (s *String) Reverse(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	runes := []rune(s.String)
	end := len(runes) - 1
	if end < 1 {
		return s
	}

	var out bytes.Buffer
	for i := end; i >= 0; i-- {
		out.WriteRune(runes[i])
	}
	return NewString(out.String())
}

func (s *String) Replace(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	mObj, ok := args[0].(*String) //mObj: modify object
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replace", "*String", args[0].Type()))
	}

	rObj, ok := args[1].(*String) //rObj: replace object
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replace", "*String", args[1].Type()))
	}

	ret := strings.Replace(s.String, mObj.String, rObj.String, -1)
	return NewString(ret)
}

func (s *String) TrimLeft(line string, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
	}

	if len(args) == 0 {
		ret := strings.TrimLeftFunc(s.String, isWhiteSpace)
		return NewString(ret)
	}

	cutset, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimLeft", "*String", args[0].Type()))
	}

	ret := strings.TrimLeft(s.String, cutset.String)
	return NewString(ret)
}

func (s *String) TrimRight(line string, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
	}

	if len(args) == 0 {
		ret := strings.TrimRightFunc(s.String, isWhiteSpace)
		return NewString(ret)
	}

	cutset, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimRight", "*String", args[0].Type()))
	}

	ret := strings.TrimRight(s.String, cutset.String)
	return NewString(ret)
}

func (s *String) Trim(line string, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
	}

	if len(args) == 0 {
		ret := strings.TrimSpace(s.String)
		return NewString(ret)
	}

	cutset, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trim", "*String", args[0].Type()))
	}

	ret := strings.Trim(s.String, cutset.String)
	return NewString(ret)
}

func (s *String) TrimPrefix(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	prefix, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimPrefix", "*String", args[0].Type()))
	}

	ret := strings.TrimPrefix(s.String, prefix.String)
	return NewString(ret)
}

func (s *String) TrimSuffix(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	suffix, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimSuffix", "*String", args[0].Type()))
	}

	ret := strings.TrimPrefix(s.String, suffix.String)
	return NewString(ret)
}

func (s *String) Split(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sep, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "split", "*String", args[0].Type()))
	}

	a := &Array{}
	strArr := strings.Split(s.String, sep.String)
	for _, v := range strArr {
		a.Members = append(a.Members, NewString(v))
	}
	return a
}

func (s *String) Lower(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	if s.String == "" {
		return s
	}

	str := strings.ToLower(s.String)
	return NewString(str)
}

func (s *String) Upper(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	if s.String == "" {
		return s
	}

	ret := strings.ToUpper(s.String)
	return NewString(ret)
}

func (s *String) SubStr(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	if args[0].Type() != INTEGER_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "first", "subStr", "*Integer", args[0].Type()))
	}
	if args[1].Type() != INTEGER_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "second", "subStr", "*Integer", args[1].Type()))
	}

	pos := args[0].(*Integer)
	length := args[1].(*Integer)

	if pos.Int64 < 0 || length.Int64 < 0 {
		panic(NewError(line, INVALIDARG))
	}

	runes := []rune(s.String)
	l := pos.Int64 + length.Int64
	aLen := int64(len(runes))
	if l > aLen {
		l = aLen
	}

	return NewString(string(runes[pos.Int64:l]))
}

func (s *String) Len(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	rune := []rune(s.String)

	return NewInteger(int64(len(rune)))
}

func (s *String) Compare(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	b, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "compare", "*String", args[0].Type()))
	}

	ret := strings.Compare(s.String, b.String)
	return NewInteger(int64(ret))
}

func (s *String) Contains(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	substr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "contains", "*String", args[0].Type()))
	}

	ret := strings.Contains(s.String, substr.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *String) ContainsAny(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	chars, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "containsAny", "*String", args[0].Type()))
	}

	ret := strings.ContainsAny(s.String, chars.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *String) Fields(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	arr := &Array{}
	ret := strings.Fields(s.String)
	for _, v := range ret {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (s *String) HasPrefix(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	prefix, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hasPrefix", "*String", args[0].Type()))
	}

	ret := strings.HasPrefix(s.String, prefix.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *String) HasSuffix(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	suffix, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hasSuffix", "*String", args[0].Type()))
	}

	ret := strings.HasSuffix(s.String, suffix.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *String) LastIndex(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	sep, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lastIndex", "*String", args[0].Type()))
	}

	ret := strings.LastIndex(s.String, sep.String)
	return NewInteger(int64(ret))
}

func (s *String) Repeat(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	count, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "repeat", "*Integer", args[0].Type()))
	}

	ret := strings.Repeat(s.String, int(count.Int64))
	return NewString(ret)
}

func (s *String) Title(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := strings.Title(s.String)
	return NewString(ret)
}

func (s *String) Chomp(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := strings.TrimRight(s.String, "\r\n")
	return NewString(ret)
}

//If you want to check if the parse is successful, you could do this:
//    v = "abcd".parseInt(10)
//    if !v.valid() {
//        println("abcd is not an int")
//    }
// This also applies to parseFloat() and parseBool()
func (s *String) ParseInt(line string, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
	}

	var base int64 = 10

	if len(args) == 1 {
		iBaseObj, ok := args[0].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "first", "parseInt", "*Integer", args[0].Type()))
		}
		base = iBaseObj.Int64
	}

	ret, err := strconv.ParseInt(s.String, int(base), 64)
	if err != nil {
		return &Integer{Int64: 0, Valid: false}
	}
	return NewInteger(ret)
}

func (s *String) ParseUInt(line string, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
	}

	var base uint64 = 10

	if len(args) == 1 {
		iBaseObj, ok := args[0].(*UInteger)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "first", "parseUInt", "*UInteger", args[0].Type()))
		}
		base = iBaseObj.UInt64
	}

	ret, err := strconv.ParseUint(s.String, int(base), 64)
	if err != nil {
		return &UInteger{UInt64: 0, Valid: false}
	}
	return NewUInteger(ret)
}

func (s *String) ParseBool(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret, err := strconv.ParseBool(s.String)
	if err != nil {
		return &Boolean{Bool: false, Valid: false}
	}

	if ret {
		return TRUE
	}
	return FALSE
}

func (s *String) ParseFloat(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret, err := strconv.ParseFloat(s.String, 64)
	if err != nil {
		return &Float{Float64: 0, Valid: false}
	}
	return NewFloat(ret)
}

func (s *String) Atoi(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret, err := strconv.Atoi(s.String)
	if err != nil {
		return &Integer{Int64: 0, Valid: false}
	}
	return NewInteger(int64(ret))
}

func (s *String) Itoa(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var i int
	switch o := args[0].(type) {
	case *Integer:
		i = int(o.Int64)
	case *UInteger:
		i = int(o.UInt64)
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "itoa", "*Integer|*UInteger", args[0].Type()))
	}

	ret := strconv.Itoa(i)
	return NewString(ret)
}

func (s *String) WriteLine(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	fmt.Println(s.String)
	return NIL
}

func (s *String) Write(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	fmt.Print(s.String)
	return NIL
}

func (s *String) IsEmpty(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if len(s.String) == 0 {
		return TRUE
	}

	return FALSE
}

//code stolen from https://github.com/AlasdairF/Hash/blob/master/hash.go
func (s *String) Hash(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	var v uint64 = 14695981039346656037
	data := []rune(s.String)
	for _, c := range data {
		v = (v ^ uint64(c)) * 1099511628211
	}
	return NewUInteger(v)
}

func (s *String) IsValid(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if s.Valid {
		return TRUE
	}
	return &Boolean{Bool: s.Valid, Valid: false}
}

func (s *String) SetValid(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 0 && argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", argLen))
	}

	if argLen == 0 {
		s.String, s.Valid = "", true
		return s
	}

	val, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setValid", "*String", args[0].Type()))
	}

	s.String, s.Valid = val.String, true
	return s
}

func (s *String) Scan(value interface{}) error {
	if value == nil {
		s.Valid = false
		return nil
	}
	s.String, s.Valid = fmt.Sprintf("%s", value), true
	return nil
}

func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

const (
	STRINGS_OBJ = "STRINGS_OBJ"
	strings_name = "strings"
)

func NewStringsObj() Object {
	ret := &StringsObj{}
	SetGlobalObj(strings_name, ret)

	return ret
}

type StringsObj struct{}

func (s *StringsObj) Inspect() string  { return strings_name }
func (s *StringsObj) Type() ObjectType { return STRINGS_OBJ }
func (s *StringsObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {

	switch method {
	case "find", "index":
		return s.Find(line, args...)
	case "lower":
		return s.Lower(line, args...)
	case "reverse":
		return s.Reverse(line, args...)
	case "upper":
		return s.Upper(line, args...)
	case "trimLeft":
		return s.TrimLeft(line, args...)
	case "trimRight":
		return s.TrimRight(line, args...)
	case "trim":
		return s.Trim(line, args...)
	case "trimPrefix":
		return s.TrimPrefix(line, args...)
	case "trimSuffix":
		return s.TrimSuffix(line, args...)
	case "split":
		return s.Split(line, args...)
	case "replace":
		return s.Replace(line, args...)
	case "count":
		return s.Count(line, args...)
	case "join":
		return s.Join(line, args...)
	case "substr":
		return s.SubStr(line, args...)
	case "len":
		return s.Len(line, args...)
	case "compare":
		return s.Compare(line, args...)
	case "contains":
		return s.Contains(line, args...)
	case "containsAny":
		return s.ContainsAny(line, args...)
	case "fields":
		return s.Fields(line, args...)
	case "hasPrefix":
		return s.HasPrefix(line, args...)
	case "hasSuffix":
		return s.HasSuffix(line, args...)
	case "lastIndex":
		return s.LastIndex(line, args...)
	case "repeat":
		return s.Repeat(line, args...)
	case "title":
		return s.Title(line, args...)
	case "chomp":
		return s.Chomp(line, args...)
	case "parseInt":
		return s.ParseInt(line, args...)
	case "parseUInt":
		return s.ParseUInt(line, args...)
	case "parseBool":
		return s.ParseBool(line, args...)
	case "parseFloat":
		return s.ParseFloat(line, args...)
	case "atoi":
		return s.Atoi(line, args...)
	case "itoa":
		return s.Itoa(line, args...)
	case "writeLine":
		return s.WriteLine(line, args...)
	case "write":
		return s.Write(line, args...)
	case "isEmpty":
		return s.IsEmpty(line, args...)
	case "hash":
		return s.Hash(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, s.Type()))
}

func (s *StringsObj) Count(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "count", "*String", args[0].Type()))
	}

	sub, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "count", "*String", args[1].Type()))
	}

	count := strings.Count(source.String, sub.String)
	return NewInteger(int64(count))
}

func (s *StringsObj) Find(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "find", "*String", args[0].Type()))
	}

	sub, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "find", "*String", args[1].Type()))
	}

	idx := strings.Index(source.String, sub.String)
	return NewInteger(int64(idx))
}

//join(arr)  or  join(arr, sep)
func (s *StringsObj) Join(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 && argLen != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", argLen))
	}

	sep := ""
	a, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "join", "*Array", args[0].Type()))
	}

	if argLen == 2 {
		v, ok := args[1].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "join", "*String", args[1].Type()))
		}
		sep = v.String
	}

	var tmp []string
	for _, item := range a.Members {
		tmp = append(tmp, item.Inspect())
	}

	ret := strings.Join(tmp, sep)
	return NewString(ret)
}

func (s *StringsObj) Reverse(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "reverse", "*String", args[0].Type()))
	}

	runes := []rune(source.String)
	end := len(runes) - 1
	if end < 1 {
		return source
	}

	var out bytes.Buffer
	for i := end; i >= 0; i-- {
		out.WriteRune(runes[i])
	}
	return NewString(out.String())
}

func (s *StringsObj) Replace(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replace", "*String", args[0].Type()))
	}

	mObj, ok := args[1].(*String) //mObj: modify object
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replace", "*String", args[1].Type()))
	}

	rObj, ok := args[2].(*String) //rObj: replace object
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "replace", "*String", args[2].Type()))
	}

	ret := strings.Replace(source.String, mObj.String, rObj.String, -1)
	return NewString(ret)
}

func (s *StringsObj) TrimLeft(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimLeft", "*String", args[0].Type()))
	}

	if len(args) == 1 {
		ret := strings.TrimLeftFunc(source.String, isWhiteSpace)
		return NewString(ret)
	}

	cutset, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "trimLeft", "*String", args[1].Type()))
	}

	ret := strings.TrimLeft(source.String, cutset.String)
	return NewString(ret)
}

func (s *StringsObj) TrimRight(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimRight", "*String", args[0].Type()))
	}

	if len(args) == 1 {
		ret := strings.TrimRightFunc(source.String, isWhiteSpace)
		return NewString(ret)
	}

	cutset, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "trimRight", "*String", args[1].Type()))
	}

	ret := strings.TrimRight(source.String, cutset.String)
	return NewString(ret)
}

func (s *StringsObj) Trim(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trim", "*String", args[0].Type()))
	}

	if len(args) == 1 {
		ret := strings.TrimSpace(source.String)
		return NewString(ret)
	}

	cutset, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "trim", "*String", args[1].Type()))
	}

	ret := strings.Trim(source.String, cutset.String)
	return NewString(ret)
}

func (s *StringsObj) TrimPrefix(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimPrefix", "*String", args[0].Type()))
	}

	prefix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "trimPrefix", "*String", args[1].Type()))
	}

	ret := strings.TrimPrefix(source.String, prefix.String)
	return NewString(ret)
}

func (s *StringsObj) TrimSuffix(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "trimSuffix", "*String", args[0].Type()))
	}

	suffix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "trimSuffix", "*String", args[1].Type()))
	}

	ret := strings.TrimPrefix(source.String, suffix.String)
	return NewString(ret)
}

func (s *StringsObj) Split(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "split", "*String", args[0].Type()))
	}

	sep, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "split", "*String", args[1].Type()))
	}

	a := &Array{}
	strArr := strings.Split(source.String, sep.String)
	for _, v := range strArr {
		a.Members = append(a.Members, NewString(v))
	}
	return a
}

func (s *StringsObj) Lower(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lower", "*String", args[0].Type()))
	}

	if source.String == "" {
		return source
	}

	str := strings.ToLower(source.String)
	return NewString(str)
}

func (s *StringsObj) Upper(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "upper", "*String", args[0].Type()))
	}

	if source.String == "" {
		return source
	}

	ret := strings.ToUpper(source.String)
	return NewString(ret)
}

func (s *StringsObj) SubStr(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	if args[0].Type() != STRING_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "first", "subStr", "*STRING", args[0].Type()))
	}
	if args[1].Type() != INTEGER_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "second", "subStr", "*Integer", args[1].Type()))
	}
	if args[2].Type() != INTEGER_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "third", "subStr", "*Integer", args[2].Type()))
	}

	source := args[0].(*String)
	pos := args[1].(*Integer)
	length := args[2].(*Integer)

	if pos.Int64 < 0 || length.Int64 < 0 {
		panic(NewError(line, INVALIDARG))
	}

	runes := []rune(source.String)
	l := pos.Int64 + length.Int64
	aLen := int64(len(runes))
	if l > aLen {
		l = aLen
	}

	return NewString(string(runes[pos.Int64:l]))
}

func (s *StringsObj) Len(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "len", "*String", args[0].Type()))
	}

	rune := []rune(source.String)

	return NewInteger(int64(len(rune)))
}

func (s *StringsObj) Compare(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "compare", "*String", args[0].Type()))
	}

	b, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "compare", "*String", args[1].Type()))
	}

	ret := strings.Compare(source.String, b.String)
	return NewInteger(int64(ret))
}

func (s *StringsObj) Contains(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "contains", "*String", args[0].Type()))
	}

	substr, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "contains", "*String", args[1].Type()))
	}

	ret := strings.Contains(source.String, substr.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *StringsObj) ContainsAny(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "containsAny", "*String", args[0].Type()))
	}

	chars, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "containsAny", "*String", args[1].Type()))
	}

	ret := strings.ContainsAny(source.String, chars.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *StringsObj) Fields(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "fields", "*String", args[0].Type()))
	}

	arr := &Array{}
	ret := strings.Fields(source.String)
	for _, v := range ret {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (s *StringsObj) HasPrefix(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hasprefix", "*String", args[0].Type()))
	}

	prefix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "hasPrefix", "*String", args[1].Type()))
	}

	ret := strings.HasPrefix(source.String, prefix.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *StringsObj) HasSuffix(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hasSuffix", "*String", args[0].Type()))
	}

	suffix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "hasSuffix", "*String", args[1].Type()))
	}

	ret := strings.HasSuffix(source.String, suffix.String)
	if ret {
		return TRUE
	}

	return FALSE
}

func (s *StringsObj) LastIndex(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hasSuffix", "*lastIndex", args[0].Type()))
	}

	sep, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "lastIndex", "*String", args[1].Type()))
	}

	ret := strings.LastIndex(source.String, sep.String)
	return NewInteger(int64(ret))
}

func (s *StringsObj) Repeat(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "repeat", "*lastIndex", args[0].Type()))
	}

	count, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "repeat", "*Integer", args[1].Type()))
	}

	ret := strings.Repeat(source.String, int(count.Int64))
	return NewString(ret)
}

func (s *StringsObj) Title(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "title", "*String", args[0].Type()))
	}

	ret := strings.Title(source.String)
	return NewString(ret)
}

func (s *StringsObj) Chomp(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	source, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "chomp", "*String", args[0].Type()))
	}

	ret := strings.TrimRight(source.String, "\r\n")
	return NewString(ret)
}

//If you want to check if the parse is successful, you could do this:
//    v = strings.parseInt("abc", 10)
//    if !v.valid() {
//        println("abcd is not an int")
//    }
// This also applies to parseFloat() and parseBool()
func (s *StringsObj) ParseInt(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseInt", "*String", args[0].Type()))
	}

	var base int64 = 10
	if len(args) == 2 {
		iBaseObj, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "parseInt", "*Integer", args[1].Type()))
		}
		base = iBaseObj.Int64
	}

	ret, err := strconv.ParseInt(strObj.String, int(base), 64)
	if err != nil {
		return &Integer{Int64: 0, Valid: false}
	}
	return NewInteger(ret)
}

func (s *StringsObj) ParseUInt(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseUInt", "*String", args[0].Type()))
	}

	var base uint64 = 10
	if len(args) == 2 {
		iBaseObj, ok := args[1].(*UInteger)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "parseUInt", "*UInteger", args[1].Type()))
		}
		base = iBaseObj.UInt64
	}

	ret, err := strconv.ParseUint(strObj.String, int(base), 64)
	if err != nil {
		return &UInteger{UInt64: 0, Valid: false}
	}
	return NewUInteger(ret)
}

func (s *StringsObj) ParseBool(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseBool", "*String", args[0].Type()))
	}

	ret, err := strconv.ParseBool(strObj.String)
	if err != nil {
		return &Boolean{Bool: false, Valid: false}
	}

	if ret {
		return TRUE
	}
	return FALSE
}

func (s *StringsObj) ParseFloat(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseFloat", "*String", args[0].Type()))
	}

	ret, err := strconv.ParseFloat(strObj.String, 64)
	if err != nil {
		return &Float{Float64: 0, Valid: false}
	}
	return NewFloat(ret)
}

func (s *StringsObj) Atoi(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "atoi", "*String", args[0].Type()))
	}

	ret, err := strconv.Atoi(strObj.String)
	if err != nil {
		return &Integer{Int64: 0, Valid: false}
	}
	return NewInteger(int64(ret))
}

func (s *StringsObj) Itoa(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var i int
	switch o := args[0].(type) {
	case *Integer:
		i = int(o.Int64)
	case *UInteger:
		i = int(o.UInt64)
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "itoa", "*Integer|*UInteger", args[0].Type()))
	}

	ret := strconv.Itoa(i)
	return NewString(ret)
}

func (s *StringsObj) WriteLine(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeLine", "*String", args[0].Type()))
	}

	fmt.Println(strObj.String)
	return NIL
}

func (s *StringsObj) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	fmt.Print(strObj.String)
	return NIL
}

func (s *StringsObj) IsEmpty(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "isEmpty", "*String", args[0].Type()))
	}

	if len(strObj.String) == 0 {
		return TRUE
	}
	return FALSE
}

//code stolen from https://github.com/AlasdairF/Hash/blob/master/hash.go
func (s *StringsObj) Hash(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hash", "*String", args[0].Type()))
	}

	var v uint64 = 14695981039346656037
	data := []rune(strObj.String)
	for _, c := range data {
		v = (v ^ uint64(c)) * 1099511628211
	}
	return NewUInteger(v)
}
