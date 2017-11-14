package eval

import (
	"encoding/json"
	"monkey/ast"
	"regexp"
)

type RegEx struct {
	RegExp *regexp.Regexp
	Value         string
}

func (re *RegEx) Inspect() string  { return re.Value }
func (re *RegEx) Type() ObjectType { return REGEX_OBJ }

func (re *RegEx) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "match", "matchString":
		return re.Match(line, args...)
	case "replace", "replaceAllString":
		return re.Replace(line, args...)
	case "split":
		return re.Split(line, args...)
	case "findAllString":
		return re.FindAllString(line, args...)
	case "findAllStringIndex":
		return re.FindAllStringIndex(line, args...)
	case "findAllStringSubmatch":
		return re.FindAllStringSubmatch(line, args...)
	case "findAllStringSubmatchIndex":
		return re.FindAllStringSubmatchIndex(line, args...)
	case "findString":
		return re.FindString(line, args...)
	case "findStringIndex":
		return re.FindStringIndex(line, args...)
	case "findStringSubmatch":
		return re.FindStringSubmatch(line, args...)
	case "findStringSubmatchIndex":
		return re.FindStringSubmatchIndex(line, args...)
	case "numSubexp":
		return re.NumSubexp(line, args...)
	case "replaceAllLiteralString":
		return re.ReplaceAllLiteralString(line, args...)
	case "replaceAllStringFunc":
		return re.ReplaceAllStringFunc(line, scope, args...)
	case "string":
		return re.String(line, args...)
	case "subexpNames":
		return re.SubexpNames(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, re.Type()))
}

func (re *RegEx) Match(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if args[0].Type() != STRING_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "first", "match", "*String", args[0].Type()))
	}

	str := args[0].(*String)
	matched := re.RegExp.MatchString(str.String)
	if matched {
		return TRUE
	}
	return FALSE
}

func (re *RegEx) Replace(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	if args[0].Type() != STRING_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "first", "replace", "*String", args[0].Type()))
	}

	if args[1].Type() != STRING_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "second", "replace", "*String", args[1].Type()))
	}

	str := args[0].(*String)
	repl := args[1].(*String)
	result := re.RegExp.ReplaceAllString(str.String, repl.String)
	return NewString(result)
}

func (re *RegEx) Split(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if args[0].Type() != STRING_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "first", "split", "*String", args[0].Type()))
	}

	str := args[0].(*String)
	splitResult := re.RegExp.Split(str.String, -1)

	a := &Array{}
	for i := 0; i < len(splitResult); i++ {
		a.Members = append(a.Members, NewString(splitResult[i]))
	}
	return a
}

func (re *RegEx) FindAllString(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllString", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllString", "*Integer", args[1].Type()))
	}

	ret := &Array{}

	strArr := re.RegExp.FindAllString(strObj.String, int(intObj.Int64))
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}
	return ret
}

func (re *RegEx) FindAllStringIndex(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllStringIndex", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllStringIndex", "*Integer", args[1].Type()))
	}

	ret := &Array{}

	intArr2D := re.RegExp.FindAllStringIndex(strObj.String, int(intObj.Int64))
	for _, v1 := range intArr2D {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewInteger(int64(v2)))
		}
		ret.Members = append(ret.Members, tmpArr)
	}
	return ret
}

func (re *RegEx) FindAllStringSubmatch(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllStringSubmatch", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllStringSubmatch", "*Integer", args[1].Type()))
	}

	ret := &Array{}

	strArr2D := re.RegExp.FindAllStringSubmatch(strObj.String, int(intObj.Int64))
	for _, v1 := range strArr2D {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewString(v2))
		}
		ret.Members = append(ret.Members, tmpArr)
	}
	return ret
}

func (re *RegEx) FindAllStringSubmatchIndex(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllStringSubmatchIndex", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllStringSubmatchIndex", "*Integer", args[1].Type()))
	}

	ret := &Array{}

	intArr2D := re.RegExp.FindAllStringSubmatchIndex(strObj.String, int(intObj.Int64))
	for _, v1 := range intArr2D {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewInteger(int64(v2)))
		}
		ret.Members = append(ret.Members, tmpArr)
	}
	return ret
}

func (re *RegEx) FindString(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findString", "*String", args[0].Type()))
	}

	ret := re.RegExp.FindString(strObj.String)
	return NewString(ret)
}

func (re *RegEx) FindStringIndex(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findStringIndex", "*String", args[0].Type()))
	}

	ret := &Array{}
	intArr := re.RegExp.FindStringIndex(strObj.String)
	for _, v := range intArr {
		ret.Members = append(ret.Members, NewInteger(int64(v)))
	}

	return ret
}

func (re *RegEx) FindStringSubmatch(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findStringSubmatch", "*String", args[0].Type()))
	}

	ret := &Array{}
	strArr := re.RegExp.FindStringSubmatch(strObj.String)
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}

	return ret
}

func (re *RegEx) FindStringSubmatchIndex(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findStringSubmatchIndex", "*String", args[0].Type()))
	}

	ret := &Array{}
	intArr := re.RegExp.FindStringSubmatchIndex(strObj.String)
	for _, v := range intArr {
		ret.Members = append(ret.Members, NewInteger(int64(v)))
	}

	return ret
}

func (re *RegEx) NumSubexp(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	i := re.RegExp.NumSubexp()
	
	return NewInteger(int64(i))
}

func (re *RegEx) ReplaceAllLiteralString(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	srcObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replaceAllLiteralString", "*String", args[0].Type()))
	}

	replObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replaceAllLiteralString", "*String", args[1].Type()))
	}

	s := re.RegExp.ReplaceAllLiteralString(srcObj.String, replObj.String)
	return NewString(s)
}

func (re *RegEx) ReplaceAllStringFunc(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	srcObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replaceAllStringFunc", "*String", args[0].Type()))
	}

	block, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replaceAllStringFunc", "*Function", args[1].Type()))
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 1 {
		panic(NewError(line, FUNCCALLBACKERROR, 1, paramCount))
	}

	ret := re.RegExp.ReplaceAllStringFunc(srcObj.String, func(str string) string {
		return replFunc(scope, block, str)
	})

	return NewString(ret)
}

func (re *RegEx) String(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	s := re.RegExp.String()
	return NewString(s)
}

func (re *RegEx) SubexpNames(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := &Array{}
	strArr := re.RegExp.SubexpNames()
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}

	return ret
}

//MarshalJSON turns regex into string/json
func (re *RegEx) MarshalJSON() ([]byte, error) {
	return json.Marshal(re.RegExp.String())
}

//UnmarshalJSON turns a string into proper regex
func (re *RegEx) UnmarshalJSON(b []byte) error {
	str := new(string)
	json.Unmarshal(b, str)
	reg, err := regexp.Compile(*str)
	if err != nil {
		return err
	}
	re.RegExp = reg
	return nil
}

/* REGEXP OBJECT */
const (
	REGEXP_OBJ = "REGEXP_OBJ"
	regexp_name = "regexp"
)

type RegExpObj struct{
	RegExp *regexp.Regexp
}

func NewRegExpObj() Object {
	ret := &RegExpObj{}
	SetGlobalObj(regexp_name, ret)

	return ret
}

func (rex *RegExpObj) Inspect() string  {
	if rex.RegExp == nil {
		return "Invalid RegExpObj!"
	}
	return rex.RegExp.String()
}

func (rex *RegExpObj) Type() ObjectType { return REGEXP_OBJ }

func (rex *RegExpObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "compile":
		return rex.Compile(line, args...)
	case "compilePOSIX":
		return rex.CompilePOSIX(line, args...)
	case "mustCompile":
		return rex.MustCompile(line, args...)
	case "mustCompilePOSIX":
		return rex.MustCompilePOSIX(line, args...)
	case "findAllString":
		return rex.FindAllString(line, args...)
	case "findAllStringIndex":
		return rex.FindAllStringIndex(line, args...)
	case "findAllStringSubmatch":
		return rex.FindAllStringSubmatch(line, args...)
	case "findAllStringSubmatchIndex":
		return rex.FindAllStringSubmatchIndex(line, args...)
	case "findString":
		return rex.FindString(line, args...)
	case "findStringIndex":
		return rex.FindStringIndex(line, args...)
	case "findStringSubmatch":
		return rex.FindStringSubmatch(line, args...)
	case "findStringSubmatchIndex":
		return rex.FindStringSubmatchIndex(line, args...)
	case "matchString", "match":
		return rex.MatchString(line, args...)
	case "numSubexp":
		return rex.NumSubexp(line, args...)
	case "replaceAllLiteralString":
		return rex.ReplaceAllLiteralString(line, args...)
	case "replaceAllString", "replace":
		return rex.ReplaceAllString(line, args...)
	case "replaceAllStringFunc":
		return rex.ReplaceAllStringFunc(line, scope, args...)
	case "split":
		return rex.Split(line, args...)
	case "string":
		return rex.String(line, args...)
	case "subexpNames":
		return rex.SubexpNames(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, rex.Type()))
}

func (rex *RegExpObj) Compile(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "compile", "*String", args[0].Type()))
	}

	var err error = nil
	rex.RegExp, err = regexp.Compile(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return rex
}

func (rex *RegExpObj) CompilePOSIX(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "compile", "*String", args[0].Type()))
	}

	var err error = nil
	rex.RegExp, err = regexp.CompilePOSIX(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return rex
}

func (rex *RegExpObj) MustCompile(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "mustCompile", "*String", args[0].Type()))
	}

	//if regexp.MustCompile() panic, we capture it, and set 'reg.RegExp' to nil.
	defer func() {
		if r := recover(); r != nil {
			rex.RegExp = nil
		}
	}()

	rex.RegExp = regexp.MustCompile(strObj.String)
	return rex
}

func (rex *RegExpObj) MustCompilePOSIX(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "mustCompilePOSIX", "*String", args[0].Type()))
	}

	//if regexp.MustCompilePOSIX() panic, we capture it, and set 'reg.RegExp' to nil.
	defer func() {
		if r := recover(); r != nil {
			rex.RegExp = nil
		}
	}()

	rex.RegExp = regexp.MustCompilePOSIX(strObj.String)
	return rex
}

func (rex *RegExpObj) FindAllString(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllString", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllString", "*Integer", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findAllString, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}

	strArr := rex.RegExp.FindAllString(strObj.String, int(intObj.Int64))
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}
	return ret
}

func (rex *RegExpObj) FindAllStringIndex(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllStringIndex", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllStringIndex", "*Integer", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findAllStringIndex, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}

	intArr2D := rex.RegExp.FindAllStringIndex(strObj.String, int(intObj.Int64))
	for _, v1 := range intArr2D {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewInteger(int64(v2)))
		}
		ret.Members = append(ret.Members, tmpArr)
	}
	return ret
}

func (rex *RegExpObj) FindAllStringSubmatch(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllStringSubmatch", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllStringSubmatch", "*Integer", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findAllStringSubmatch, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}

	strArr2D := rex.RegExp.FindAllStringSubmatch(strObj.String, int(intObj.Int64))
	for _, v1 := range strArr2D {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewString(v2))
		}
		ret.Members = append(ret.Members, tmpArr)
	}
	return ret
}

func (rex *RegExpObj) FindAllStringSubmatchIndex(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findAllStringSubmatchIndex", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "findAllStringSubmatchIndex", "*Integer", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findAllStringSubmatchIndex, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}

	intArr2D := rex.RegExp.FindAllStringSubmatchIndex(strObj.String, int(intObj.Int64))
	for _, v1 := range intArr2D {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewInteger(int64(v2)))
		}
		ret.Members = append(ret.Members, tmpArr)
	}
	return ret
}

func (rex *RegExpObj) FindString(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findString", "*String", args[0].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findString, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := rex.RegExp.FindString(strObj.String)
	return NewString(ret)
}

func (rex *RegExpObj) FindStringIndex(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findStringIndex", "*String", args[0].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findStringIndex, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}
	intArr := rex.RegExp.FindStringIndex(strObj.String)
	for _, v := range intArr {
		ret.Members = append(ret.Members, NewInteger(int64(v)))
	}

	return ret
}

func (rex *RegExpObj) FindStringSubmatch(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findStringSubmatch", "*String", args[0].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findStringSubmatch, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}
	strArr := rex.RegExp.FindStringSubmatch(strObj.String)
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}

	return ret
}

func (rex *RegExpObj) FindStringSubmatchIndex(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "findStringSubmatchIndex", "*String", args[0].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling findStringSubmatchIndex, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}
	intArr := rex.RegExp.FindStringSubmatchIndex(strObj.String)
	for _, v := range intArr {
		ret.Members = append(ret.Members, NewInteger(int64(v)))
	}

	return ret
}

func (rex *RegExpObj) MatchString(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "matchString", "*String", args[0].Type()))
	}

	b := rex.RegExp.MatchString(strObj.String)
	if b {
		return TRUE
	}
	return FALSE
}

func (rex *RegExpObj) NumSubexp(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling NumSubexp, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	i := rex.RegExp.NumSubexp()
	
	return NewInteger(int64(i))
}

func (rex *RegExpObj) ReplaceAllLiteralString(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	srcObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replaceAllLiteralString", "*String", args[0].Type()))
	}

	replObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replaceAllLiteralString", "*String", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling replaceAllLiteralString, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	s := rex.RegExp.ReplaceAllLiteralString(srcObj.String, replObj.String)
	return NewString(s)
}

func (rex *RegExpObj) ReplaceAllString(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	srcObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replaceAllString", "*String", args[0].Type()))
	}

	replObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replaceAllString", "*String", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling replaceAllString, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	s := rex.RegExp.ReplaceAllString(srcObj.String, replObj.String)
	return NewString(s)
}

func (rex *RegExpObj) ReplaceAllStringFunc(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	srcObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "replaceAllStringFunc", "*String", args[0].Type()))
	}

	block, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "replaceAllStringFunc", "*Function", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling replaceAllStringFunc, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 1 {
		panic(NewError(line, FUNCCALLBACKERROR, 1, paramCount))
	}

	ret := rex.RegExp.ReplaceAllStringFunc(srcObj.String, func(str string) string {
		return replFunc(scope, block, str)
	})

	return NewString(ret)
}

func (rex *RegExpObj) Split(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "split", "*String", args[0].Type()))
	}

	intObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "split", "*Integer", args[1].Type()))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling split, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}
	strArr := rex.RegExp.Split(strObj.String, int(intObj.Int64))
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}

	return ret
}

func (rex *RegExpObj) String(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling string, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	s := rex.RegExp.String()
	return NewString(s)
}

func (rex *RegExpObj) SubexpNames(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if rex.RegExp == nil {
		return NewNil("Before calling subexpNames, you should first call 'compile|compilePOSIX|mustCompile|mustCompilePOSIX'")
	}

	ret := &Array{}
	strArr := rex.RegExp.SubexpNames()
	for _, v := range strArr {
		ret.Members = append(ret.Members, NewString(v))
	}

	return ret
}


//ReplaceAllStringFunc()'s callback function
func replFunc(scope *Scope, f *Function, str string) string {
	s := NewScope(scope)

	//Store to `Scope`，so below `Eval() could use them
	s.Set(f.Literal.Parameters[0].(*ast.Identifier).Value, NewString(str))
	r := Eval(f.Literal.Body, s)
	if obj, ok := r.(*ReturnValue); ok {
		r = obj.Value
	}

	//check for return value, must be a '*String' type
	ret, ok := r.(*String)
	if ok {
		return ret.String
	}

	return ""
}