package eval

import (
	"bytes"
	"encoding/json"
	"monkey/ast"
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
	case "parseGlob":
		return t.ParseGlob(line, args...)
	case "clone":
		return t.Clone(line, args...)
	case "definedTemplates":
		return t.DefinedTemplates(line, args...)
	case "delims":
		return t.Delims(line, args...)
	case "execute":
		return t.Execute(line, args...)
	case "executeTemplate":
		return t.ExecuteTemplate(line, args...)
	case "funcs":
		return t.Funcs(line, scope, args...)
	case "lookup":
		return t.Lookup(line, args...)
	case "name":
		return t.Name(line, args...)
	case "option":
		return t.Option(line, args...)
	case "templates":
		return t.Templates(line, args...)
	case "htmlEscape":
		return t.HTMLEscape(line, args...)
	case "htmlEscaper":
		return t.HTMLEscaper(line, args...)
	case "htmlEscapeString":
		return t.HTMLEscapeString(line, args...)
	case "jsEscapeString":
		return t.JSEscapeString(line, args...)
	case "jsEscape":
		return t.JSEscape(line, args...)
	case "jsEscaper":
		return t.JSEscaper(line, args...)
	case "urlQueryEscaper":
		return t.URLQueryEscaper(line, args...)
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
		return NewNil("Before calling parse(), you should first call 'new|parseFiles|parseGlob' function")
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

func (t *TemplateObj) ParseGlob(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseGlob", "*String", args[0].Type()))
	}

	temp, err := template.ParseGlob(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{Template:temp}
}

func (t *TemplateObj) Clone(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling clone(), you should first call 'new|parseFiles|parseGlob' function")
	}

	temp, err := t.Template.Clone()
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{Template:temp}
}

func (t *TemplateObj) DefinedTemplates(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling definedTemplates(), you should first call 'new|parseFiles|parseGlob' function")
	}

	s := t.Template.DefinedTemplates()

	return NewString(s)
}

func (t *TemplateObj) Delims(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling delims(), you should first call 'new|parseFiles|parseGlob' function")
	}

	leftStrObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "delims", "*String", args[0].Type()))
	}

	rightStrObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "delims", "*String", args[1].Type()))
	}

	t.Template.Delims(leftStrObj.String, rightStrObj.String)

	return t
}

func (t *TemplateObj) Execute(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling execute(), you should first call 'new|parseFiles|parseGlob' function")
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

//Note :ExecuteTemplate is similar to Execute function, should extract a common private function.
func (t *TemplateObj) ExecuteTemplate(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling executeTemplate(), you should first call 'new|parseFiles|parseGlob' function")
	}

	nameStrObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "executeTemplate", "*String", args[0].Type()))
	}

	var data []byte
	var err error

	objType := args[1].Type()
	switch objType {
	case HASH_OBJ:
		//convert 'Hash' object to interface{}
		h := args[1].(*Hash)
		data, err = h.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case ARRAY_OBJ:
		//convert 'Array' object to interface{}
		arr := args[1].(*Array)
		data, err = arr.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case INTEGER_OBJ:
		//convert 'Integer' object to interface{}
		i := args[1].(*Integer)
		data, err = i.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case FLOAT_OBJ:
		//convert 'Float' object to interface{}
		f := args[1].(*Float)
		data, err = f.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case BOOLEAN_OBJ:
		//convert 'Boolean' object to interface{}
		b := args[1].(*Boolean)
		data, err = b.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case NIL_OBJ:
		//convert 'Nil' object to interface{}
		n := args[1].(*Nil)
		data, err = n.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case STRING_OBJ:
		//convert 'String' object to interface{}
		s := args[1].(*String)
		data, err = s.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case TIME_OBJ:
		//convert 'TimeObj' object to interface{}
		t := args[1].(*TimeObj)
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
	err = t.Template.ExecuteTemplate(&out, nameStrObj.String, obj)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(out.String())
}

func (t *TemplateObj) Funcs(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	hashObj, ok := args[0].(*Hash)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "funcs", "*Hash", args[0].Type()))
	}

	if t.Template == nil {
		return NewNil("Before calling funcs(), you should first call 'new|parseFiles|parseGlob' function")
	}

	funcMaps := template.FuncMap{}
	for _, pair := range hashObj.Pairs {
		if pair.Key.Type() != STRING_OBJ {
			panic(NewError(line, GENERICERROR, "Hash's key type should be 'STRING', got '" + pair.Key.Type() + "'"))
		}

		if pair.Value.Type() != FUNCTION_OBJ {
			panic(NewError(line, GENERICERROR, "Hash's value type should be 'FUNCTION', got '" + pair.Value.Type() + "'"))
		}

		key := pair.Key.(*String).String

		funcMaps[key] = func(args... interface{}) interface{} {
			var innerFn *Function
			for _, innerPair := range hashObj.Pairs {
				if key == innerPair.Key.(*String).String {
					innerFn = innerPair.Value.(*Function)
					break
				}
			}

			s := NewScope(scope)
			//put all the arguments into scope for later 'Eval'
			for idx, arg := range args {
				o, _ := unmarshalJsonObject(arg) //convert go object to monkey object
				s.Set(innerFn.Literal.Parameters[idx].(*ast.Identifier).Value, o)
			}
			ret := Eval(innerFn.Literal.Body, s)
			if obj, ok := ret.(*ReturnValue); ok {
				ret = obj.Value
			}
			return object2RawValue(ret) //convert monkey object back to go object
		}
	}

	t.Template.Funcs(funcMaps)

	return t
}

func (t *TemplateObj) Lookup(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling lookup(), you should first call 'new|parseFiles|parseGlob' function")
	}

	nameStrObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lookup", "*String", args[0].Type()))
	}


	temp := t.Template.Lookup(nameStrObj.String)
	if temp == nil {
		return NIL
	}

	return &TemplateObj{Template:temp}
}

func (t *TemplateObj) Name(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling name(), you should first call 'new|parseFiles|parseGlob' function")
	}

	str := t.Template.Name()
	return NewString(str)
}

func (t *TemplateObj) Option(line string, args ...Object) Object {
	if t.Template == nil {
		return NewNil("Before calling option(), you should first call 'new|parseFiles|parseGlob' function")
	}

	strArr :=[]string{}
	for _, v := range args {
		strArr = append(strArr, v.(*String).String)
	}
	t.Template.Option(strArr...)

	return t
}

func (t *TemplateObj) Templates(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.Template == nil {
		return NewNil("Before calling templates(), you should first call 'new|parseFiles|parseGlob' function")
	}

	retArr := &Array{}
	templateArr := t.Template.Templates()
	for _, v := range templateArr {
		retArr.Members = append(retArr.Members, &TemplateObj{Template:v})
	}
	return retArr
}

func (t *TemplateObj) HTMLEscapeString(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "htmlEscapeString", "*String", args[0].Type()))
	}

	ret := template.HTMLEscapeString(strObj.String)
	return NewString(ret)
}

func (t *TemplateObj) HTMLEscape(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	writerObj, ok := args[0].(Writable)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "htmlEscape", "Writable", args[0].Type()))
	}

	strObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "htmlEscape", "*String", args[1].Type()))
	}

	b := []byte(strObj.String)
	template.HTMLEscape(writerObj.IOWriter(), b)
	return NIL
}

func (t *TemplateObj) HTMLEscaper(line string, args ...Object) Object {
	if len(args) == 0 {
		return NewString("")
	}

	arrIntf := make([]interface{}, len(args))
	for i, v := range args {
		arrIntf[i] = object2RawValue(v)
	}

	ret := template.HTMLEscaper(arrIntf...)
	return NewString(ret)
}

func (t *TemplateObj) JSEscapeString(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "jsEscapeString", "*String", args[0].Type()))
	}

	ret := template.JSEscapeString(strObj.String)
	return NewString(ret)
}

func (t *TemplateObj) JSEscape(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	writerObj, ok := args[0].(Writable)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "jsEscape", "Writable", args[0].Type()))
	}

	strObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "jsEscape", "*String", args[1].Type()))
	}

	b := []byte(strObj.String)
	template.JSEscape(writerObj.IOWriter(), b)
	return NIL
}

func (t *TemplateObj) JSEscaper(line string, args ...Object) Object {
	if len(args) == 0 {
		return NewString("")
	}

	arrIntf := make([]interface{}, len(args))
	for i, v := range args {
		arrIntf[i] = object2RawValue(v)
	}

	ret := template.JSEscaper(arrIntf...)
	return NewString(ret)
}

func (t *TemplateObj) URLQueryEscaper(line string, args ...Object) Object {
	if len(args) == 0 {
		return NewString("")
	}

	arrIntf := make([]interface{}, len(args))
	for i, v := range args {
		arrIntf[i] = object2RawValue(v)
	}

	ret := template.URLQueryEscaper(arrIntf...)
	return NewString(ret)
}

