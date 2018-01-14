package eval

import (
	"flag"
	_ "fmt"
)

func NewFlagObj() *FlagObj {
	ret := &FlagObj{arguments: make(map[Object]interface{})}
	SetGlobalObj(flag_name, ret)

	return ret
}

const (
	FLAG_OBJ = "FLAG_OBJ"
	flag_name = "flag"
)

type FlagObj struct {
	arguments map[Object]interface{}
}

func (f *FlagObj) Inspect() string  { return "<" + flag_name + ">" }
func (f *FlagObj) Type() ObjectType { return FLAG_OBJ }

func (f *FlagObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "arg":
		return f.Arg(line, args...)
	case "args":
		return f.Args(line, args...)
	case "nArg":
		return f.NArg(line, args...)
	case "nFlag":
		return f.NFlag(line, args...)
	case "bool":
		return f.Bool(line, args...)
	case "int":
		return f.Int(line, args...)
	case "uint":
		return f.UInt(line, args...)
	case "float":
		return f.Float(line, args...)
	case "string":
		return f.String(line, args...)
	case "set":
		return f.Set(line, args...)
	case "parse":
		return f.Parse(line, args...)
	case "parsed":
		return f.Parsed(line, args...)
	case "printDefaults":
		return f.PrintDefaults(line, args...)
	case "isSet":
		return f.IsSet(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, f.Type()))
}

func (f *FlagObj) Arg(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	i, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "arg", "*Integer", args[0].Type()))
	}

	return NewString(flag.Arg(int(i.Int64)))
}

func (f *FlagObj) Args(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := &Array{}
	strlist := flag.Args()
	for _, v := range strlist {
		ret.Members = append(ret.Members, NewString(v))
	}
	return ret
}

func (f *FlagObj) NArg(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	n := flag.NArg()
	return NewInteger(int64(n))
}

func (f *FlagObj) NFlag(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	n := flag.NFlag()
	return NewInteger(int64(n))
}

func (f *FlagObj) Bool(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "bool", "*String", args[0].Type()))
	}

	value, ok := args[1].(*Boolean)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "bool", "*Boolean", args[1].Type()))
	}

	usage, ok := args[2].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "bool", "*String", args[2].Type()))
	}

	//Because golang's flag module's `Bool` method returns pointer, only after calling `flag.Parse()`,
	//the returned pointer will be assigned, so we need store this pointer.
	b := flag.Bool(name.String, value.Bool, usage.String)
	ret := &Boolean{Bool: value.Bool, Valid: true}
	f.arguments[ret] = b

	return ret
}

func (f *FlagObj) Int(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "int", "*String", args[0].Type()))
	}

	value, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "int", "*Integer", args[1].Type()))
	}

	usage, ok := args[2].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "int", "*String", args[2].Type()))
	}

	i := flag.Int64(name.String, value.Int64, usage.String)
	ret := NewInteger(value.Int64)
	f.arguments[ret] = i

	return ret
}

func (f *FlagObj) UInt(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "uint", "*String", args[0].Type()))
	}

	value, ok := args[1].(*UInteger)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "uint", "*UInteger", args[1].Type()))
	}

	usage, ok := args[2].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "uint", "*String", args[2].Type()))
	}

	i := flag.Uint64(name.String, value.UInt64, usage.String)
	ret := NewUInteger(value.UInt64)
	f.arguments[ret] = i

	return ret
}

func (f *FlagObj) Float(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "float", "*String", args[0].Type()))
	}

	value, ok := args[1].(*Float)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "float", "*Float", args[1].Type()))
	}

	usage, ok := args[2].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "float", "*String", args[2].Type()))
	}

	flt := flag.Float64(name.String, value.Float64, usage.String)
	ret := NewFloat(value.Float64)
	f.arguments[ret] = flt

	return ret
}

func (f *FlagObj) String(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "string", "*String", args[0].Type()))
	}

	value, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "string", "*String", args[1].Type()))
	}

	usage, ok := args[2].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "string", "*String", args[2].Type()))
	}

	s := flag.String(name.String, value.String, usage.String)
	ret := NewString(value.String)
	f.arguments[ret] = s

	return ret
}

func (f *FlagObj) Set(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "set", "*String", args[0].Type()))
	}

	value, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "set", "*String", args[1].Type()))
	}

	err := flag.Set(name.String, value.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (f *FlagObj) Parse(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	flag.Parse()

	//Now `flag.Parse()` is called, we need to assined these pointers back to our objects.
	for key, value := range f.arguments {
		switch key.(type) {
		case *Boolean:
			b := key.(*Boolean)
			v := value.(*bool)
			b.Bool = *v

		case *Integer:
			i := key.(*Integer)
			v := value.(*int64)
			i.Int64 = *v

		case *UInteger:
			i := key.(*UInteger)
			v := value.(*uint64)
			i.UInt64 = *v

		case *Float:
			f := key.(*Float)
			v := value.(*float64)
			f.Float64 = *v

		case *String:
			s := key.(*String)
			v := value.(*string)
			s.String = *v
		}
	}
	return NIL
}

func (f *FlagObj) Parsed(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	b := flag.Parsed()
	if b {
		return TRUE
	}
	return FALSE
}

func (f *FlagObj) PrintDefaults(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	flag.PrintDefaults()
	return NIL
}

func (f *FlagObj) IsSet(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "isSet", "*String", args[0].Type()))
	}

	alreadySet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})

	if _, ok := alreadySet[name.String]; ok {
		return TRUE
	}
	return FALSE
}
