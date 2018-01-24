package eval

import (
	"monkey/ast"
	_ "fmt"
)

const (
	CLASS_OBJ = "CLASS_OBJ"
	INSTANCE_OBJ = "INSTANCE_OBJ"
)


type ClassComponentKind byte
const (
	ClassMemberKind ClassComponentKind = iota
	ClassMethodKind
	ClassPropertyKind
	ClassIndexerKind    //not implemented
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
}

func (c *Class) Inspect() string { return "<class:" + c.Name + ">" }

func (c *Class) Type() ObjectType { return CLASS_OBJ }

func (c *Class) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, c.Type(), method))
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
	case ClassIndexerKind: //not implemented
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
func (oi *ObjectInstance) CallMethod(line string, scope *Scope, method string, args ...Object) Object { 
	panic(NewError(line, NOMETHODERROR, oi.Type(), method))
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
	}

	BASE_CLASS.Methods["is_a"] = BASE_CLASS.Methods["instanceOf"]

	return true
}

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
