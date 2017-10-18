package eval

import (
	"encoding/json"
	"regexp"
)

type RegEx struct {
	RegExpression *regexp.Regexp
	Value         string
}

func (re *RegEx) Inspect() string  { return re.Value }
func (re *RegEx) Type() ObjectType { return REGEX_OBJ }

func (re *RegEx) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "match":
		return re.Match(line, args...)
	case "replace":
		return re.Replace(line, args...)
	case "split":
		return re.Split(line, args...)
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
	matched := re.RegExpression.MatchString(str.String)
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
	result := re.RegExpression.ReplaceAllString(str.String, repl.String)
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
	splitResult := re.RegExpression.Split(str.String, -1)

	a := &Array{}
	for i := 0; i < len(splitResult); i++ {
		a.Members = append(a.Members, NewString(splitResult[i]))
	}
	return a
}

//MarshalJSON turns regex into string/json
func (re *RegEx) MarshalJSON() ([]byte, error) {
	return json.Marshal(re.RegExpression.String())
}

//UnmarshalJSON turns a string into proper regex
func (re *RegEx) UnmarshalJSON(b []byte) error {
	str := new(string)
	json.Unmarshal(b, str)
	reg, err := regexp.Compile(*str)
	if err != nil {
		return err
	}
	re.RegExpression = reg
	return nil
}
