package eval

import (
	"monkey/ast"
	_ "fmt"
)

const (
	CLASS_OBJ        = "CLASS_OBJ"
	INSTANCE_OBJ     = "INSTANCE_OBJ"
	METHODINFO_OBJ   = "METHODINFO_OBJ"
	PROPERTYINFO_OBJ = "PROPERTYINFO_OBJ"
)


type ClassComponentKind byte
const (
	ClassMemberKind ClassComponentKind = iota
	ClassMethodKind
	ClassPropertyKind
//	ClassIndexerKind    //not implemented
)


type ClassMethod interface {
	Object
	classMethod() ast.ModifierLevel
}

type Class struct {
	Name       string
	Parent     *Class
	Members    []*ast.LetStatement
	Methods    map[string]ClassMethod //BuiltinMethod or Function object
	Properties map[string]*ast.PropertyDeclStmt
	Scope      *Scope
	IsAnnotation bool //true if the class is an annotation class
}

func (c *Class) Inspect() string { return "<class:" + c.Name + ">" }

func (c *Class) Type() ObjectType { return CLASS_OBJ }

func (c *Class) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "isAnnotationPresent":
		return c.IsAnnotationPresent(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, c.Type()))
}

func (c *Class) IsAnnotationPresent(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

//	if c.HasAnnotation {
//		return TRUE
//	}
	return FALSE
}

func (c *Class) GetMethod(name string) ClassMethod {
	m, ok := c.Methods[name]
	if ok || c.Parent == nil { //check self's method
		return m //found, return self's method
	}

	//return parent's method
	return c.Parent.GetMethod(name)
}

func (c *Class) GetProperty(name string) *ast.PropertyDeclStmt {
	p, ok := c.Properties[name]
	if ok || c.Parent == nil { //check self's method
		return p //found, return self's properties
	}

	//return parent's property
	return c.Parent.GetProperty(name)
}

func (c *Class) IsStatic(val string, kind ClassComponentKind) bool {
	switch kind {
	case ClassMemberKind:
		return c.CheckMembers(val)
	case ClassMethodKind:
		return c.CheckMethods(val)
	case ClassPropertyKind:
		return c.CheckProperties(val)
	}
	return false
}

func (c *Class) CheckMembers(val string) bool {
	result := false
	for _, member := range c.Members {
		for _, name := range member.Names {
			if name.Value == val {
				if member.StaticFlag {
					result = true
				}
			}
		}
	}

	if result == false {
		if c.Parent != nil {
			result = c.Parent.CheckMembers(val)
		}
	}
	return result
}

func (c *Class) CheckMethods(val string) bool {
	v := c.GetMethod(val)
	if v == nil {
		return false
	}

	switch o := v.(type) {
	case *Function:
		if o.Literal.StaticFlag {
			return true
		}
	case *BuiltinMethod:
		return false
	}

	return false
}

func (c *Class) CheckProperties(val string) bool {
	v := c.GetProperty(val)
	if v == nil {
		return false
	}

	return v.StaticFlag
}

func (c *Class) GetModifierLevel(name string, kind ClassComponentKind) ast.ModifierLevel {
	var ret ast.ModifierLevel
	switch kind {
	case ClassMemberKind:
		for _, letStmt := range c.Members {
			for _, v := range letStmt.Names {
				if v.Value == name {
					ret = letStmt.ModifierLevel
				}
			}
		}

	case ClassMethodKind:
		clsMethod := c.GetMethod(name)
		if clsMethod != nil {
			ret = clsMethod.classMethod()
		}

	case ClassPropertyKind:
	}

	if ret == ast.ModifierPrivate || ret == ast.ModifierDefault {
		return ast.ModifierPrivate
	}
	return ret
}

type ObjectInstance struct {
	Class  *Class
	Scope *Scope
}

func (oi *ObjectInstance) Inspect() string  { return "<Instance:" + oi.Class.Name + ">" }
func (oi *ObjectInstance) Type() ObjectType { return INSTANCE_OBJ }
func (oi *ObjectInstance) GetMethod(name string) ClassMethod { return oi.Class.GetMethod(name) }
func (oi *ObjectInstance) GetProperty(name string) *ast.PropertyDeclStmt { return oi.Class.GetProperty(name) }
func (oi *ObjectInstance) IsStatic(val string, kind ClassComponentKind) bool { return oi.Class.IsStatic(val, kind) }
func (oi *ObjectInstance) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, oi.Type()))
}

func (oi *ObjectInstance) GetModifierLevel(name string, kind ClassComponentKind) ast.ModifierLevel {
	return oi.Class.GetModifierLevel(name, kind)
}

//The base class of all classes in monkey
var BASE_CLASS = &Class{
	Name:    "object",
	Parent:  nil,
	Members: nil,
	Methods: map[string]ClassMethod{},
}

//Builtin @Override annotation class
var OVERRIDE_ANNOCLASS = &Class{
	Name:    "Override",
	Parent:  BASE_CLASS,
	IsAnnotation: true,
}

//Builtin @NotNull annotation class
//The annotated element must not be {@code null}.
var NOTNULL_ANNOCLASS = &Class{
	Name:    "NotNull",
	Parent:  BASE_CLASS,
	IsAnnotation: true,
}

//Builtin @NotEmpty annotation class
//Asserts that the annotated string, tuple, map or array is not null or empty.
var NOTEMPTY_ANNOCLASS = &Class{
	Name:    "NotEmpty",
	Parent:  BASE_CLASS,
	IsAnnotation: true,
}

func initRootObject() bool {
	BASE_CLASS.Methods = map[string]ClassMethod {
		"toString": &BuiltinMethod{
			Fn:func(line string, self *ObjectInstance, scope *Scope, args ...Object) Object {
				argLen := len(args)
				if argLen != 0 {
					panic(NewError(line, ARGUMENTERROR, "0", argLen))
				}

				if self == nil {
					return NewString("")
				}

				return NewString(self.Class.Inspect())
			},
		},
		"instanceOf":
			&BuiltinMethod{
				Fn:func(line string, self *ObjectInstance, scope *Scope, args ...Object) Object {
					argLen := len(args)
					if argLen != 1 {
						panic(NewError(line, ARGUMENTERROR, "1", argLen))
					}

					if self == nil {
						return NIL
					}

					switch class := args[0].(type) {
					case *String:
						return nativeBoolToBooleanObject(InstanceOf(class.String, self))
					case *Class:
						return nativeBoolToBooleanObject(InstanceOf(class.Name, self))
					}

					panic(NewError(line, GENERICERROR, "is_a/instanceOf expected a class or string for its argument"))
				},
		},
		"classOf":
			&BuiltinMethod{
				Fn:func(line string, self *ObjectInstance, scope *Scope, args ...Object) Object {
					argLen := len(args)
					if argLen != 0 {
						panic(NewError(line, ARGUMENTERROR, "0", argLen))
					}

					if self == nil {
						return NewString("")
					}

					return NewString(self.Class.Name)
				},
		},
		"hashCode":&BuiltinMethod{
			Fn:func(line string, self *ObjectInstance, scope *Scope, args ...Object) Object {
				argLen := len(args)
				if argLen != 0 {
					panic(NewError(line, ARGUMENTERROR, "0", argLen))
				}

				if self == nil {
					return NewUInteger(0)
				}

				var v uint64 = 14695981039346656037
				data := []rune(self.Class.Name)
				for _, c := range data {
					v = (v ^ uint64(c)) * 1099511628211
				}
				return NewUInteger(v)
			},
		},
		"getMethods":&BuiltinMethod{
			Fn:func(line string, self *ObjectInstance, scope *Scope, args ...Object) Object {
				argLen := len(args)
				if argLen != 0 {
					panic(NewError(line, ARGUMENTERROR, "0", argLen))
				}

				ret := &Array{}
				if self != nil {
					for name, _ := range self.Class.Methods {
						//m := &MethodInfo{Name: name, instance: self, Method:self.GetMethod(name), Scope: scope}
						m := &MethodInfo{Name: name, Instance: self, Scope: scope}
						ret.Members = append(ret.Members, m)
					}
				}

				return ret
			},
		},
		"getProperties":&BuiltinMethod{
			Fn:func(line string, self *ObjectInstance, scope *Scope, args ...Object) Object {
				argLen := len(args)
				if argLen != 0 {
					panic(NewError(line, ARGUMENTERROR, "0", argLen))
				}

				ret := &Array{}
				if self != nil {
					for name, _ := range self.Class.Properties {
						m := &PropertyInfo{Name: name, Instance: self, Scope: scope}
						ret.Members = append(ret.Members, m)
					}
				}

				return ret
			},
		},
	}

	BASE_CLASS.Methods["is_a"] = BASE_CLASS.Methods["instanceOf"]

	return true
}

//MethodInfo object
type MethodInfo struct {
	Name     string          //Method name
	Instance *ObjectInstance //instance
	Scope    *Scope          //Method's scope
}

func (m *MethodInfo) Inspect() string { return "<method:" + m.Name + ">" }
func (m *MethodInfo) Type() ObjectType { return METHODINFO_OBJ }

func (m *MethodInfo) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "name", "getName":
		return m.GetName(line, args...)
	case "invoke":
		return m.Invoke(line, scope, args...)
	case "getAnnotations":
		return m.GetAnnotations(line, args...)
	case "getAnnotation":
		return m.GetAnnotation(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, m.Type()))
}

func (m *MethodInfo) GetName(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(m.Name)
}

func (m *MethodInfo) Invoke(line string, scope *Scope, args ...Object) Object {
	method := m.Instance.GetMethod(m.Name)
	if method != nil {
		switch meth := method.(type) {
			case *Function:
				newScope := NewScope(m.Instance.Scope)
				newScope.Set("parent", m.Instance.Class.Parent)
				return evalFunctionDirect(method, args, newScope)
			case *BuiltinMethod:
				builtinMethod :=&BuiltinMethod{Fn: meth.Fn, Instance: m.Instance}
				aScope := NewScope(m.Instance.Scope)
				return evalFunctionDirect(builtinMethod, args, aScope)
		}
	}

	return NIL
}

func (m *MethodInfo) GetAnnotations(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	method := m.Instance.GetMethod(m.Name)
	ret := &Array{}
	
	switch o := method.(type) {
	case *Function:
		for _, anno := range o.Annotations {
			ret.Members = append(ret.Members, anno)
		}
		return ret
	}
	return ret
}

func (m *MethodInfo) GetAnnotation(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	annoCls, ok := args[0].(*Class)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "getAnnotation", "*Class", args[0].Type()))
	}

	method := m.Instance.GetMethod(m.Name)
	
	switch o := method.(type) {
	case *Function:
		for _, anno := range o.Annotations {
			if anno.Class.Name == annoCls.Name {
				return anno
			}
		}
	}
	return NIL
}


//PropertyInfo object
type PropertyInfo struct {
	Name     string          //Property name
	Instance *ObjectInstance //instance
	Scope    *Scope          //Property's scope
}

func (p *PropertyInfo) Inspect() string { return "<property:" + p.Name + ">" }
func (p *PropertyInfo) Type() ObjectType { return PROPERTYINFO_OBJ }

func (p *PropertyInfo) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "name", "getName":
		return p.GetName(line, args...)
	case "getAnnotations":
		return p.GetAnnotations(line, scope, args...)
	case "value":
		return p.Value(line, args...)
//	case "getAnnotation":
//		return p.GetAnnotation(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, p.Type()))
}

func (p *PropertyInfo) GetName(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(p.Name)
}

func (p *PropertyInfo) Value(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	prop := p.Instance.GetProperty(p.Name)
	if prop != nil {
		if prop.Getter == nil { //property xxx { set; }
			panic(NewError(line, PROPERTYUSEERROR, p.Name, p.Instance.Class.Name))
		} else {
			if len(prop.Getter.Body.Statements) == 0 { //property xxx { get; }
				v, _ := p.Instance.Scope.Get("_" + p.Name)
				return v
			} else {
				results := Eval(prop.Getter.Body, p.Instance.Scope)
				if results.Type() == RETURN_VALUE_OBJ {
					return results.(*ReturnValue).Value
				}
			}
		}
	}
	panic(NewError(line, UNKNOWNIDENT, p.Name))
}

func (p *PropertyInfo) GetAnnotations(line string, scope *Scope, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret := &Array{}
	propStmt := p.Instance.GetProperty(p.Name)
	processClassAnnotation(propStmt.Annotations, scope, line, ret)

	return ret
}

//func (p *PropertyInfo) GetAnnotation(line string, args ...Object) Object {
//	if len(args) != 1 {
//		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
//	}
//
//	annoCls, ok := args[0].(*Class)
//	if !ok {
//		panic(NewError(line, PARAMTYPEERROR, "first", "getAnnotation", "*Class", args[0].Type()))
//	}
//
//	propStmt := p.Instance.GetProperty(m.Name)
//	
//	switch o := method.(type) {
//	case *Function:
//		for _, anno := range o.Annotations {
//			if anno.Class.Name == annoCls.Name {
//				return anno
//			}
//		}
//	}
//	return NIL
//}



var _ = initRootObject()

func InstanceOf(className string, oi *ObjectInstance) bool {
	if oi == nil {
		return false
	}

	cls := oi.Class
	for {
		if cls.Name == className {
			return true
		}
		if cls.Parent == nil {
			return false
		}
		cls = cls.Parent
	}
}
