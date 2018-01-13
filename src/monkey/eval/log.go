package eval

import (
	"log"
)

const (
	LOGGER_OBJ = "LOGGER_OBJ"
	logger_name = "logger"
)

func NewLoggerObj() Object {
	ret := &LoggerObj{}
	SetGlobalObj(logger_name, ret)

	SetGlobalObj(logger_name+".LDATE", NewInteger(int64(log.Ldate)))
	SetGlobalObj(logger_name+".LTIME", NewInteger(int64(log.Ltime)))
	SetGlobalObj(logger_name+".LMICROSECONDS", NewInteger(int64(log.Lmicroseconds)))
	SetGlobalObj(logger_name+".LLONGFILE", NewInteger(int64(log.Llongfile)))
	SetGlobalObj(logger_name+".LSHORTFILE", NewInteger(int64(log.Lshortfile)))
	SetGlobalObj(logger_name+".LUTC", NewInteger(int64(log.LUTC)))
	SetGlobalObj(logger_name+".LSTDFLAGS", NewInteger(int64(log.LstdFlags)))

	return ret
}

type LoggerObj struct {
	Logger *log.Logger
}

func (l *LoggerObj) Inspect() string  { return "<" + logger_name + ">" }
func (l *LoggerObj) Type() ObjectType { return LOGGER_OBJ }

func (l *LoggerObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "print":
		return l.Print(line, args...)
	case "printf":
		return l.Printf(line, args...)
	case "println":
		return l.Println(line, args...)
	case "fatal":
		return l.Fatal(line, args...)
	case "fatalf":
		return l.Fatalf(line, args...)
	case "fatalln":
		return l.Fatalln(line, args...)
	case "panic":
		return l.Panic(line, args...)
	case "panicf":
		return l.Panicf(line, args...)
	case "panicln":
		return l.Panicln(line, args...)
	case "flags":
		return l.Flags(line, args...)
	case "output":
		return l.Output(line, args...)
	case "prefix":
		return l.Prefix(line, args...)
	case "setFlags":
		return l.SetFlags(line, args...)
	case "setOutput":
		return l.SetOutput(line, args...)
	case "setPrefix":
		return l.SetPrefix(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, l.Type()))
}

func (l *LoggerObj) Print(line string, args ...Object) Object {
	if len(args) == 0 {
		l.Logger.Print()
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Print(wrapped...)
	return NIL
}

func (l *LoggerObj) Printf(line string, args ...Object) Object {
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

	l.Logger.Printf(format.String, wrapped...)
	return NIL
}

func (l *LoggerObj) Println(line string, args ...Object) Object {
	if len(args) == 0 {
		l.Logger.Println()
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Println(wrapped...)
	return NIL
}

func (l *LoggerObj) Fatal(line string, args ...Object) Object {
	if len(args) == 0 {
		l.Logger.Fatal()
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Fatal(wrapped...)
	return NIL
}

func (l *LoggerObj) Fatalf(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">0", len(args)))
	}

	format, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "fatalf", "*String", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Fatalf(format.String, wrapped...)
	return NIL
}

func (l *LoggerObj) Fatalln(line string, args ...Object) Object {
	if len(args) == 0 {
		l.Logger.Fatalln()
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Fatalln(wrapped...)
	return NIL
}

func (l *LoggerObj) Panic(line string, args ...Object) Object {
	if len(args) == 0 {
		l.Logger.Panic()
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Panic(wrapped...)
	return NIL
}

func (l *LoggerObj) Panicf(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">0", len(args)))
	}

	format, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "panicf", "*String", args[0].Type()))
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Panicf(format.String, wrapped...)
	return NIL
}

func (l *LoggerObj) Panicln(line string, args ...Object) Object {
	if len(args) == 0 {
		l.Logger.Panicln()
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	l.Logger.Panicln(wrapped...)
	return NIL
}

func (l *LoggerObj) Flags(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	f := l.Logger.Flags()
	return NewInteger(int64(f))
}

func (l *LoggerObj) Output(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	calldepth, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "output", "*Integer", args[0].Type()))
	}

	s, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "output", "*String", args[1].Type()))
	}

	err := l.Logger.Output(int(calldepth.Int64), s.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (l *LoggerObj) Prefix(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	p := l.Logger.Prefix()
	return NewString(p)
}

func (l *LoggerObj) SetFlags(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	flag, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setFlags", "*Integer", args[0].Type()))
	}

	l.Logger.SetFlags(int(flag.Int64))
	return NIL
}

func (l *LoggerObj) SetOutput(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	writer, ok := args[0].(Writable)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setOutput", "Writable", args[0].Type()))
	}

	l.Logger.SetOutput(writer.IOWriter())
	return NIL

}

func (l *LoggerObj) SetPrefix(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	prefix, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setPrefix", "*String", args[0].Type()))
	}

	l.Logger.SetPrefix(prefix.String)
	return NIL
}
