package eval

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const (
	TEMPLATE_OBJ = "TEMPLATE_OBJ"
	template_name = "template"
)

type TemplateObj struct {
	Template *template.Template
}

func NewTemplateObj() Object {
	ret := &TemplateObj{}
	SetGlobalObj(template_name, ret)

	return ret
}

func (t *TemplateObj) Inspect() string  { return template_name }
func (t *TemplateObj) Type() ObjectType { return TEMPLATE_OBJ }

func (t *TemplateObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "new":
		return t.New(line, args...)
	case "parse":
		return t.Parse(line, args...)
	case "parseFiles":
		return t.ParseFiles(line, args...)
	case "execute":
		return t.Execute(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, t.Type()))
}

func (t *TemplateObj) New(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 0 && argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", argLen))
	}

	var name string = "main"
	if argLen == 1 {
		strObj, ok := args[0].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "first", "new", "*String", args[0].Type()))
		}
		name = strObj.String
	}

	return &TemplateObj{Template:template.New(name)}
}

func (t *TemplateObj) Parse(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parse", "*String", args[0].Type()))
	}

	if t.Template == nil {
		return NewNil("Before calling execute, you should first call 'new|parseFiles' function")
	}

	var err error = nil
	t.Template, err = t.Template.Parse(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return t
}

func (t *TemplateObj) ParseFiles(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseFiles", "*String", args[0].Type()))
	}

	temp, err := template.ParseFiles(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{Template:temp}
}

func (t *TemplateObj) Execute(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling execute, you should first call 'new|parseFiles' function")
	}

	var data []byte
	var err error

	objType := args[0].Type()
	switch objType {
	case HASH_OBJ:
		//convert 'Hash' object to interface{}
		h := args[0].(*Hash)
		data, err = h.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case ARRAY_OBJ:
		//convert 'Array' object to interface{}
		arr := args[0].(*Array)
		data, err = arr.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case INTEGER_OBJ:
		//convert 'Integer' object to interface{}
		i := args[0].(*Integer)
		data, err = i.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case FLOAT_OBJ:
		//convert 'Float' object to interface{}
		f := args[0].(*Float)
		data, err = f.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case BOOLEAN_OBJ:
		//convert 'Boolean' object to interface{}
		b := args[0].(*Boolean)
		data, err = b.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case NIL_OBJ:
		//convert 'Nil' object to interface{}
		n := args[0].(*Nil)
		data, err = n.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case STRING_OBJ:
		//convert 'String' object to interface{}
		s := args[0].(*String)
		data, err = s.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case TIME_OBJ:
		//convert 'TimeObj' object to interface{}
		t := args[0].(*TimeObj)
		data, err = t.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "execute", "*Integer|*Float|*String|*Boolean|*Nil|*TimeObj|*Array|*Hash", objType))
	}

	var obj interface{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return NewNil(err.Error())
	}

	var out bytes.Buffer
	err = t.Template.Execute(&out, obj)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(out.String())
}

