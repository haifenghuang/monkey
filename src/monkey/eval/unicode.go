package eval

import (
	"unicode"
)

type unicodeFunc func(rune) bool

var funcMap = map[string]unicodeFunc {
	"IsControl" : unicode.IsControl,
	"IsDigit"   : unicode.IsDigit,
	"IsGraphic" : unicode.IsGraphic,
	"IsLetter"  : unicode.IsLetter,
	"IsLower"   : unicode.IsLower,
	"IsMark"    : unicode.IsMark,
	"IsNumber"  : unicode.IsNumber,
	"IsPrint"   : unicode.IsPrint,
	"IsPunct"   : unicode.IsPunct,
	"IsSpace"   : unicode.IsSpace,
	"IsSymbol"  : unicode.IsSymbol,
	"IsTitle"   : unicode.IsTitle,
	"IsUpper"   : unicode.IsUpper,
}


func NewUnicodeObj() *UnicodeObj {
	ret := &UnicodeObj{}
	SetGlobalObj(unicode_name, ret)

	return ret
}

const (
	UNICODE_OBJ = "UNICODE_OBJ"
	unicode_name = "unicode"
)

type UnicodeObj struct {
}

func (u *UnicodeObj) Inspect() string  { return "<" + unicode_name + ">" }
func (u *UnicodeObj) Type() ObjectType { return UNICODE_OBJ }

func (u *UnicodeObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "isControl":
		return u.IsControl(line, args...)
	case "isDigit":
		return u.IsDigit(line, args...)
	case "isGraphic":
		return u.IsGraphic(line, args...)
	case "isLetter":
		return u.IsLetter(line, args...)
	case "isLower":
		return u.IsLower(line, args...)
	case "isMark":
		return u.IsMark(line, args...)
	case "isNumber":
		return u.IsNumber(line, args...)
	case "isPrint":
		return u.IsPrint(line, args...)
	case "isPunct":
		return u.IsPunct(line, args...)
	case "isSpace":
		return u.IsSpace(line, args...)
	case "isSymbol":
		return u.IsSymbol(line, args...)
	case "isTitle":
		return u.IsTitle(line, args...)
	case "isUpper":
		return u.IsUpper(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, u.Type()))
}

func (u *UnicodeObj) IsControl(line string, args ...Object) Object {
	return u.checkFunc("IsControl", line, args...)
}

func (u *UnicodeObj) IsDigit(line string, args ...Object) Object {
	return u.checkFunc("IsDigit", line, args...)
}

func (u *UnicodeObj) IsGraphic(line string, args ...Object) Object {
	return u.checkFunc("IsGraphic", line, args...)
}

func (u *UnicodeObj) IsLetter(line string, args ...Object) Object {
	return u.checkFunc("IsLetter", line, args...)
}

func (u *UnicodeObj) IsLower(line string, args ...Object) Object {
	return u.checkFunc("IsLower", line, args...)
}

func (u *UnicodeObj) IsMark(line string, args ...Object) Object {
	return u.checkFunc("IsMark", line, args...)
}

func (u *UnicodeObj) IsNumber(line string, args ...Object) Object {
	return u.checkFunc("IsNumber", line, args...)
}

func (u *UnicodeObj) IsPrint(line string, args ...Object) Object {
	return u.checkFunc("IsPrint", line, args...)
}

func (u *UnicodeObj) IsPunct(line string, args ...Object) Object {
	return u.checkFunc("IsPunct", line, args...)
}

func (u *UnicodeObj) IsSpace(line string, args ...Object) Object {
	return u.checkFunc("IsSpace", line, args...)
}
func (u *UnicodeObj) IsSymbol(line string, args ...Object) Object {
	return u.checkFunc("IsSymbol", line, args...)
}

func (u *UnicodeObj) IsTitle(line string, args ...Object) Object {
	return u.checkFunc("IsTitle", line, args...)
}

func (u *UnicodeObj) IsUpper(line string, args ...Object) Object {
	return u.checkFunc("IsUpper", line, args...)
}

func (u *UnicodeObj) checkFunc(name string, line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", name, "*String", args[0].Type()))
	}

	r := []rune(strObj.String)
	if funcMap[name](r[0]) {
		return TRUE
	}

	return FALSE
}


