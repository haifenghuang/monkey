package eval

import (
	gofmt "fmt"
)

func NewFmtObj() *FmtObj {
	ret := &FmtObj{}
	SetGlobalObj(fmt_name, ret)

	return ret
}

const (
	FMT_OBJ = "FMT_OBJ"
	fmt_name = "fmt"
)

type FmtObj struct {
}

func (f *FmtObj) Inspect() string  { return fmt_name }
func (f *FmtObj) Type() ObjectType { return FMT_OBJ }

func (f *FmtObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "errorf":
		return f.Errorf(line, args...)
	case "print":
		return f.Print(line, args...)
	case "printf":
		return f.Printf(line, args...)
	case "println":
		return f.Println(line, args...)
	case "sprint":
		return f.Sprint(line, args...)
	case "sprintf":
		return f.Sprintf(line, args...)
	case "sprintln":
		return f.Sprintln(line, args...)
	case "fprint":
		return f.Fprint(line, args...)
	case "fprintf":
		return f.Fprintf(line, args...)
	case "fprintln":
		return f.Fprintln(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, f.Type()))
}

func (f *FmtObj) Errorf(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">0", len(args)))
	}

	format, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "printf", "*String", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	err := gofmt.Errorf(format.String, wrapped...)
	panic(NewError(line, GENERICERROR, err.Error()))

	return NIL
}

func (f *FmtObj) Print(line string, args ...Object) Object {
	if len(args) == 0 {
		n, err := gofmt.Print()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewInteger(int64(n))
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	n, err := gofmt.Print(wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Printf(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">0", len(args)))
	}

	format, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "printf", "*String", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	n, err := gofmt.Printf(format.String, wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Println(line string, args ...Object) Object {
	if len(args) == 0 {
		n, err := gofmt.Println()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewInteger(int64(n))
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	n, err := gofmt.Println(wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Sprint(line string, args ...Object) Object {
	if len(args) == 0 {
		return NewString("")
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	ret := gofmt.Sprint(wrapped...)
	return NewString(ret)
}

func (f *FmtObj) Sprintf(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">0", len(args)))
	}

	format, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sprintf", "*String", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	ret := gofmt.Sprintf(format.String, wrapped...)
	return NewString(ret)
}

func (f *FmtObj) Sprintln(line string, args ...Object) Object {
	if len(args) == 0 {
		ret := gofmt.Sprintln()
		return NewString(ret)
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	ret := gofmt.Sprintln(wrapped...)
	return NewString(ret)
}

func (f *FmtObj) Fprint(line string, args ...Object) Object {
	if len(args) < 2 {
		panic(NewError(line, ARGUMENTERROR, ">=2", len(args)))
	}

	w, ok := args[0].(Writable)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "Fprint", "Writable", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	n, err := gofmt.Fprint(w.(Writable).IOWriter(), wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}

func (f *FmtObj) Fprintf(line string, args ...Object) Object {
	if len(args) < 2 {
		panic(NewError(line, ARGUMENTERROR, ">=2", len(args)))
	}

	w, ok := args[0].(Writable)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "Fprintf", "Writable", args[0].Type()))
	}

	format, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "Fprintf", "*String", args[1].Type()))
	}

	var n int
	var err error
	if len(args) > 2 {
		subArgs := args[2:]
		wrapped := make([]interface{}, len(subArgs))
		for i, v := range subArgs {
			wrapped[i] = &Formatter{Obj: v}
		}
		n, err = gofmt.Fprintf(w.(Writable).IOWriter(), format.String, wrapped...)
		if err != nil {
			return NewNil(err.Error())
		}
	} else {
		n, err = gofmt.Fprintf(w.(Writable).IOWriter(), format.String)
		if err != nil {
			return NewNil(err.Error())
		}
	}
	return NewInteger(int64(n))
}

func (f *FmtObj) Fprintln(line string, args ...Object) Object {
	if len(args) < 2 {
		panic(NewError(line, ARGUMENTERROR, ">=2", len(args)))
	}

	w, ok := args[0].(Writable)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "Fprintln", "Writable", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	n, err := gofmt.Fprintln(w.(Writable).IOWriter(), wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}
