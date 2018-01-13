package eval

import (
	"bytes"
	"encoding/json"
	"monkey/ast"
	text "text/template"
	html "html/template"
)

const (
	TEMPLATE_OBJ = "TEMPLATE_OBJ"
	template_name = "template"
	T_TEXT = 0
	T_HTML = 1
)

type TemplateObj struct {
	TmplType         int64  //template type
	TextTemplate *text.Template
	HTMLTemplate *html.Template
}

func NewTemplateObj() Object {
	ret := &TemplateObj{}
	SetGlobalObj(template_name, ret)

	SetGlobalObj(template_name+".TEXT", NewInteger(0))
	SetGlobalObj(template_name+".HTML", NewInteger(1))
	
	return ret
}

func (t *TemplateObj) Inspect() string  { return "<" + template_name +">" }
func (t *TemplateObj) Type() ObjectType { return TEMPLATE_OBJ }

func (t *TemplateObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "newText","text":
		return t.NewText(line, args...)
	case "newHtml", "html":
		return t.NewHtml(line, args...)
	case "new":
		return t.New(line, args...)
	case "parse":
		return t.Parse(line, args...)
	case "parseTextFiles":
		return t.ParseTextFiles(line, args...)
	case "parseHtmlFiles":
		return t.ParseHtmlFiles(line, args...)
	case "parseFiles":
		return t.ParseFiles(line, args...)
	case "parseTextGlob":
		return t.ParseTextGlob(line, args...)
	case "parseHtmlGlob":
		return t.ParseHtmlGlob(line, args...)
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

func (t *TemplateObj) NewText(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", argLen))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "newText", "*String", args[0].Type()))
	}

	t.TmplType = T_TEXT
	t.TextTemplate = text.New(strObj.String)
	return t
}

func (t *TemplateObj) NewHtml(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", argLen))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "newHtml", "*String", args[0].Type()))
	}

	t.TmplType = T_HTML
	t.HTMLTemplate = html.New(strObj.String)
	return t
}

func (t *TemplateObj) New(line string, args ...Object) Object {
	argLen := len(args)
	if argLen != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", argLen))
	}

	intObj, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "new", "*Integer", args[0].Type()))
	}

	tmplType := intObj.Int64
	if tmplType != T_TEXT && tmplType != T_HTML {
		panic(NewError(line, GENERICERROR, "First parameter of new() should be 0(text)|1(html)."))
	}

	var name string = ""
	strObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "new", "*String", args[1].Type()))
	}
	name = strObj.String

	if tmplType == T_TEXT {
		t.TmplType = tmplType
		t.TextTemplate = text.New(name)
	} else if tmplType == T_HTML {
		t.TmplType = tmplType
		t.HTMLTemplate = html.New(name)
	}
	return t
}

func (t *TemplateObj) Parse(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parse", "*String", args[0].Type()))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewNil("Before calling parse(), you should first call 'new|parseFiles|parseGlob' function")
	}

	if t.TmplType == T_TEXT {
		temp, err := t.TextTemplate.Parse(strObj.String)
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType: t.TmplType, TextTemplate :temp}

	} else if t.TmplType == T_HTML {
		temp, err := t.HTMLTemplate.Parse(strObj.String)
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType: t.TmplType, HTMLTemplate :temp}
	}

	return NIL
}

func (t *TemplateObj) ParseTextFiles(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseTextFiles", "*String", args[0].Type()))
	}

	temp, err := text.ParseFiles(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{TmplType: T_TEXT, TextTemplate: temp}
}

func (t *TemplateObj) ParseHtmlFiles(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseHtmlFiles", "*String", args[0].Type()))
	}

	temp, err := html.ParseFiles(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{TmplType: T_HTML, HTMLTemplate: temp}
}

func (t *TemplateObj) ParseFiles(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	intObj, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseFiles", "*Integer", args[0].Type()))
	}

	tmplType := intObj.Int64
	if tmplType != T_TEXT && tmplType != T_HTML {
		panic(NewError(line, GENERICERROR, "First parameter of parseFiles() should be 0(text)|1(html)."))
	}

	strObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "parseFiles", "*String", args[1].Type()))
	}

	if tmplType == T_TEXT {
		temp, err := text.ParseFiles(strObj.String)
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType : tmplType, TextTemplate: temp}
	} else if tmplType == T_HTML {
		temp, err := html.ParseFiles(strObj.String)
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType : tmplType, HTMLTemplate: temp}
	}

	return NIL
}

func (t *TemplateObj) ParseTextGlob(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseTextGlob", "*String", args[0].Type()))
	}

	temp, err := text.ParseGlob(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{TmplType:t.TmplType, TextTemplate:temp}
}

func (t *TemplateObj) ParseHtmlGlob(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseHtmlGlob", "*String", args[0].Type()))
	}

	temp, err := html.ParseGlob(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:temp}
}

func (t *TemplateObj) ParseGlob(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "parseGlob", "*String", args[0].Type()))
	}

	if t.TmplType == T_TEXT {
		temp, err := text.ParseGlob(strObj.String)
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType:t.TmplType, TextTemplate:temp}
	} else if t.TmplType == T_HTML {
		temp, err := html.ParseGlob(strObj.String)
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:temp}
	}

	return NIL
}

func (t *TemplateObj) Clone(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewNil("Before calling clone(), you should first call 'new|parseFiles|parseGlob' function")
	}

	if t.TmplType == T_TEXT {
		temp, err := t.TextTemplate.Clone()
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType:t.TmplType, TextTemplate:temp}
	} else if t.TmplType == T_HTML {
		temp, err := t.HTMLTemplate.Clone()
		if err != nil {
			return NewNil(err.Error())
		}
		return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:temp}
	}

	return NIL
}

func (t *TemplateObj) DefinedTemplates(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	var s string = ""
	if t.TextTemplate == nil && t.HTMLTemplate == nil {
			return NewNil("Before calling definedTemplates(), you should first call 'new|parseFiles|parseGlob' function")
	}

	if t.TmplType == T_TEXT {
		s = t.TextTemplate.DefinedTemplates()
	} else if t.TmplType == T_HTML {
		s = t.HTMLTemplate.DefinedTemplates()
	}

	return NewString(s)
}

func (t *TemplateObj) Delims(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
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

	if t.TmplType == T_TEXT {
		return &TemplateObj{TmplType:t.TmplType, TextTemplate:t.TextTemplate.Delims(leftStrObj.String, rightStrObj.String)}
	} else if t.TmplType == T_HTML {
		return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:t.HTMLTemplate.Delims(leftStrObj.String, rightStrObj.String)}
	}

	return NIL
}

func (t *TemplateObj) Execute(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewFalseObj("Before calling execute(), you should first call 'new|parseFiles|parseGlob' function")
	}

	var isWriter bool = false
	var strObj *String
	writerObj, ok := args[0].(Writable)
	if !ok { //is not 'Writable', check if it's a '*String'
		var ok bool
		strObj, ok = args[0].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "first", "execute", "Writable|*String", args[0].Type()))
		}
	} else {
		isWriter = true
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
			return NewFalseObj(err.Error())
		}
	case ARRAY_OBJ:
		//convert 'Array' object to interface{}
		arr := args[1].(*Array)
		data, err = arr.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case INTEGER_OBJ:
		//convert 'Integer' object to interface{}
		i := args[1].(*Integer)
		data, err = i.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case FLOAT_OBJ:
		//convert 'Float' object to interface{}
		f := args[1].(*Float)
		data, err = f.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case BOOLEAN_OBJ:
		//convert 'Boolean' object to interface{}
		b := args[1].(*Boolean)
		data, err = b.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case NIL_OBJ:
		//convert 'Nil' object to interface{}
		n := args[1].(*Nil)
		data, err = n.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case STRING_OBJ:
		//convert 'String' object to interface{}
		s := args[1].(*String)
		data, err = s.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case TIME_OBJ:
		//convert 'TimeObj' object to interface{}
		t := args[1].(*TimeObj)
		data, err = t.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "execute", "*Integer|*Float|*String|*Boolean|*Nil|*TimeObj|*Array|*Hash", objType))
	}

	var obj interface{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	if isWriter {
		if t.TmplType == T_TEXT {
			err = t.TextTemplate.Execute(writerObj.IOWriter(), obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		} else if t.TmplType == T_HTML {
			err = t.HTMLTemplate.Execute(writerObj.IOWriter(), obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		}

	} else {
		var out bytes.Buffer
		if t.TmplType == T_TEXT {
			err = t.TextTemplate.Execute(&out, obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		} else if t.TmplType == T_HTML {
			err = t.HTMLTemplate.Execute(&out, obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		}
		//set result to the first parameter
		strObj.String = out.String()
	}

	return TRUE
}

//Note :ExecuteTemplate is similar to Execute function, should extract a common private function.
func (t *TemplateObj) ExecuteTemplate(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewFalseObj("Before calling executeTemplate(), you should first call 'new|parseFiles|parseGlob' function")
	}

	var isWriter bool = false
	var strObj *String
	writerObj, ok := args[0].(Writable)
	if !ok { //is not 'Writable', check if it's a '*String'
		var ok bool
		strObj, ok = args[0].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "first", "executeTemplate", "Writable|*String", args[0].Type()))
		}
	} else {
		isWriter = true
	}

	nameStrObj, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "executeTemplate", "*String", args[1].Type()))
	}

	var data []byte
	var err error

	objType := args[2].Type()
	switch objType {
	case HASH_OBJ:
		//convert 'Hash' object to interface{}
		h := args[2].(*Hash)
		data, err = h.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case ARRAY_OBJ:
		//convert 'Array' object to interface{}
		arr := args[2].(*Array)
		data, err = arr.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
	case INTEGER_OBJ:
		//convert 'Integer' object to interface{}
		i := args[2].(*Integer)
		data, err = i.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case FLOAT_OBJ:
		//convert 'Float' object to interface{}
		f := args[2].(*Float)
		data, err = f.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case BOOLEAN_OBJ:
		//convert 'Boolean' object to interface{}
		b := args[2].(*Boolean)
		data, err = b.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case NIL_OBJ:
		//convert 'Nil' object to interface{}
		n := args[2].(*Nil)
		data, err = n.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case STRING_OBJ:
		//convert 'String' object to interface{}
		s := args[2].(*String)
		data, err = s.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	case TIME_OBJ:
		//convert 'TimeObj' object to interface{}
		t := args[2].(*TimeObj)
		data, err = t.MarshalJSON()
		if err != nil {
			return NewFalseObj(err.Error())
		}
	default:
		panic(NewError(line, PARAMTYPEERROR, "third", "executeTemplate", "*Integer|*Float|*String|*Boolean|*Nil|*TimeObj|*Array|*Hash", objType))
	}

	var obj interface{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	if isWriter {
		if t.TmplType == T_TEXT {
			err = t.TextTemplate.ExecuteTemplate(writerObj.IOWriter(), nameStrObj.String, obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		} else if t.TmplType == T_HTML {
			err = t.HTMLTemplate.ExecuteTemplate(writerObj.IOWriter(), nameStrObj.String, obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		}
	} else {
		var out bytes.Buffer
		if t.TmplType == T_TEXT {
			err = t.TextTemplate.ExecuteTemplate(&out, nameStrObj.String, obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		} else if t.TmplType == T_HTML {
			err = t.HTMLTemplate.ExecuteTemplate(&out, nameStrObj.String, obj)
			if err != nil {
				return NewFalseObj(err.Error())
			}
		}
		//set result to the first parameter
		strObj.String = out.String()
	}

	return TRUE
}

func (t *TemplateObj) Funcs(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	hashObj, ok := args[0].(*Hash)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "funcs", "*Hash", args[0].Type()))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewNil("Before calling funcs(), you should first call 'new|parseFiles|parseGlob' function")
	}

	funcMaps := make(map[string]interface{})
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

	if t.TmplType == T_TEXT {
		return &TemplateObj{TmplType:t.TmplType, TextTemplate: t.TextTemplate.Funcs(funcMaps)}
	} else if t.TmplType == T_HTML {
		return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:t.HTMLTemplate.Funcs(funcMaps)}
	}

	return NIL
}

func (t *TemplateObj) Lookup(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewNil("Before calling lookup(), you should first call 'new|parseFiles|parseGlob' function")
	}

	nameStrObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lookup", "*String", args[0].Type()))
	}

	if t.TmplType == T_TEXT {
		temp := t.TextTemplate.Lookup(nameStrObj.String)
		if temp == nil {
			return NIL
		}
		return &TemplateObj{TmplType:t.TmplType, TextTemplate:temp}
	} else if t.TmplType == T_HTML {
		temp := t.HTMLTemplate.Lookup(nameStrObj.String)
		if temp == nil {
			return NIL
		}
		return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:temp}
	}

	return NIL
}

func (t *TemplateObj) Name(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewNil("Before calling name(), you should first call 'new|parseFiles|parseGlob' function")
	}

	var str string = ""
	if t.TmplType == T_TEXT {
		str = t.TextTemplate.Name()
	} else if t.TmplType == T_HTML {
		str = t.HTMLTemplate.Name()
	}
	return NewString(str)
}

func (t *TemplateObj) Option(line string, args ...Object) Object {
	if t.TextTemplate == nil  && t.HTMLTemplate == nil {
		return NewNil("Before calling option(), you should first call 'new|parseFiles|parseGlob' function")
	}

	strArr :=[]string{}
	for _, v := range args {
		strArr = append(strArr, v.(*String).String)
	}

	if t.TmplType == T_TEXT {
		return &TemplateObj{TmplType:t.TmplType, TextTemplate:t.TextTemplate.Option(strArr...)}
	} else if t.TmplType == T_HTML {
		return &TemplateObj{TmplType:t.TmplType, HTMLTemplate:t.HTMLTemplate.Option(strArr...)}
	}

	return NIL
}

func (t *TemplateObj) Templates(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if t.TextTemplate == nil && t.HTMLTemplate == nil {
		return NewNil("Before calling templates(), you should first call 'new|parseFiles|parseGlob' function")
	}

	retArr := &Array{}

	if t.TmplType == T_TEXT {
		templateArr := t.TextTemplate.Templates()
		for _, v := range templateArr {
			retArr.Members = append(retArr.Members, &TemplateObj{TmplType:t.TmplType, TextTemplate:v})
		}
	} else if t.TmplType == T_HTML {
		templateArr := t.HTMLTemplate.Templates()
		for _, v := range templateArr {
			retArr.Members = append(retArr.Members, &TemplateObj{TmplType:t.TmplType, HTMLTemplate:v})
		}
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

	ret := text.HTMLEscapeString(strObj.String)
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
	text.HTMLEscape(writerObj.IOWriter(), b)
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

	//use text.HTMLEscaper() is ok, because html.HTMLEscaper() just call text.HTMLEscaper().
	ret := text.HTMLEscaper(arrIntf...)
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

	ret := text.JSEscapeString(strObj.String)
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
	text.JSEscape(writerObj.IOWriter(), b)
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

	ret := text.JSEscaper(arrIntf...)
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

	ret := text.URLQueryEscaper(arrIntf...)
	return NewString(ret)
}

