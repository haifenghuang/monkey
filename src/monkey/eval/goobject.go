package eval
/*
	NOTE: MOST OF THE REFLECT CODE COME FROM `https://github.com/jxwr/doby`
	WITH MINOR MODIFICATIONS
*/
import (
	"fmt"
	"reflect"
	_ "runtime"
	"strings"
)

// Wrapper for go object
type GoObject struct {
	obj interface{}
}

func (gobj *GoObject) iter() bool {
	val := reflect.ValueOf(gobj.obj)
	kind := val.Kind()

	switch kind {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func (gobj *GoObject) Inspect() string  { return fmt.Sprint(gobj.obj) }
func (gobj *GoObject) Type() ObjectType { return GO_OBJ }

func (gobj *GoObject) Equal(another Object) bool {
	anotherObj, ok := another.(*GoObject)
	if ok {
		if gobj.obj == anotherObj.obj {
			return true
		}
		return false
	} else {
		if gobj.obj == nil {
			if another.Type() == NIL_OBJ {
				return true
			}
			return false
		}
		obj := GoValueToObject(gobj.obj)
		return reflect.DeepEqual(obj, another)
	}
}

func (gobj *GoObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	theMethod := reflect.ValueOf(gobj.obj).MethodByName(method)
	if !theMethod.IsValid() {
		panic(NewError(line, NOMETHODERROR, method, gobj.Type()))
	}

	methodType := theMethod.Type()
	theArgs := []reflect.Value{}

	if methodType.NumIn() > 0 {
		i := 0
		for ; i < methodType.NumIn()-1; i++ {
			reqTyp := methodType.In(i)
			theArgs = append(theArgs, ObjectToValue(args[i], reqTyp))
		}

		for ; i < len(args); i++ {
			var reqTyp reflect.Type
			if i < methodType.NumIn() && methodType.In(i).Kind() != reflect.Slice {
				reqTyp = methodType.In(i)
			}
			theArgs = append(theArgs, ObjectToValue(args[i], reqTyp))
		}
	}

	var results []Object
	rets := theMethod.Call(theArgs)
	for _, ret := range rets {
		results = append(results, NewGoObject(ret.Interface()))
	}
	if len(results) > 1 { // There are multiple return values
		// we need to convert results to tuples.
		return &Tuple{Members: results, IsMulti: true}
	}
	return results[0]

	panic(NewError(line, NOMETHODERROR, method, gobj.Type()))
}

func NewGoObject(obj interface{}) *GoObject {
	ret := &GoObject{obj:obj}

	val := reflect.ValueOf(obj)
	if obj != nil && reflect.Indirect(val).IsValid() && val.Kind() > reflect.Invalid && val.Kind() <= reflect.UnsafePointer {
		key := reflect.Indirect(val).Type().PkgPath() + "::" + reflect.Indirect(val).Type().String()
		SetGlobalObj(key, ret)
	}

	return ret
}

// wrapper for go functions
type GoFuncObject struct {
	name string
	typ  reflect.Type
	fn   interface{}
}

func (gfn *GoFuncObject) Inspect() string  { return gfn.name }
func (gfn *GoFuncObject) Type() ObjectType { return GFO_OBJ }

func (gfn *GoFuncObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	inNum := gfn.typ.NumIn()
	inArgs := []reflect.Value{}

	dummy := []interface{}{}
	for i := 0; i < inNum && i < len(args); i++ {
		arg := args[i]
		if i == inNum-1 && gfn.typ.In(i) == reflect.TypeOf(dummy) {
			for j := i; j < len(args); j++ {
				arg := args[j]
				switch arg := arg.(type) {
				case *Integer:
					v := reflect.ValueOf(arg.Int64)
					inArgs = append(inArgs, v)
				case *UInteger:
					v := reflect.ValueOf(arg.UInt64)
					inArgs = append(inArgs, v)
				case *Float:
					v := reflect.ValueOf(arg.Float64)
					inArgs = append(inArgs, v)
				case *String:
					v := reflect.ValueOf(arg.String)
					inArgs = append(inArgs, v)
				case *Boolean:
					v := reflect.ValueOf(arg.Bool)
					inArgs = append(inArgs, v)
				case *GoObject:
					if arg.obj == nil {
						var nilObj *Nil
						inArgs = append(inArgs, reflect.ValueOf(nilObj))
					} else {
						v := reflect.ValueOf(arg.obj)
						inArgs = append(inArgs, v)
					}
				default:
					inArgs = append(inArgs, reflect.ValueOf(arg.Inspect()))
				}
			}
		} else {
			switch arg := arg.(type) {
			case *Integer:
				v := reflect.ValueOf(arg.Int64)
				t := reflect.TypeOf(arg.Int64)
				if t.ConvertibleTo(gfn.typ.In(i)) {
					v = v.Convert(gfn.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *UInteger:
				v := reflect.ValueOf(arg.UInt64)
				t := reflect.TypeOf(arg.UInt64)
				if t.ConvertibleTo(gfn.typ.In(i)) {
					v = v.Convert(gfn.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *Float:
				v := reflect.ValueOf(arg.Float64)
				t := reflect.TypeOf(arg.Float64)
				if t.ConvertibleTo(gfn.typ.In(i)) {
					v = v.Convert(gfn.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *String:
				v := reflect.ValueOf(arg.String)
				t := reflect.TypeOf(arg.String)
				if t.ConvertibleTo(gfn.typ.In(i)) {
					v = v.Convert(gfn.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *Boolean:
				v := reflect.ValueOf(arg.Bool)
				t := reflect.TypeOf(arg.Bool)
				if t.ConvertibleTo(gfn.typ.In(i)) {
					v = v.Convert(gfn.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *GoObject:
				v := reflect.ValueOf(arg.obj)
				t := reflect.TypeOf(arg.obj)
				if t != nil {
					if t.ConvertibleTo(gfn.typ.In(i)) {
						v = v.Convert(gfn.typ.In(i))
					}
					inArgs = append(inArgs, v)
				} else {
					typ := gfn.typ.In(i)
					inArgs = append(inArgs, reflect.Zero(typ))
				}
			default:
				v := reflect.ValueOf(arg)
				t := reflect.TypeOf(arg)
				if t.ConvertibleTo(gfn.typ.In(i)) {
					v = v.Convert(gfn.typ.In(i))
				}
				inArgs = append(inArgs, v)
			}
		}
	}

	var results []Object
	//Call the function
	outVals := reflect.ValueOf(gfn.fn).Call(inArgs)
	// Convert the result back to monkey Object
	for _, val := range outVals {
		// Here we only retuns GoObject, See comment of 'RegisterVars'.
		results = append(results, NewGoObject(val.Interface()))
//		switch val.Kind() {
//		case reflect.Bool:
//			results = append(results, NewBooleanObj(val.Bool()))
//		case reflect.String:
//			results = append(results, NewString(val.String()))
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			results = append(results, NewInteger(int64(val.Int())))
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			results = append(results, NewUInteger(uint64(val.Uint())))
//		case reflect.Float64, reflect.Float32:
//			results = append(results, NewFloat(val.Float()))
//		default:
//			results = append(results, NewGoObject(val.Interface()))
//		}
	}

	if len(results) > 1 { // There are multiple return values
		// we need to convert results to tuples.
		return &Tuple{Members: results, IsMulti: true}
	}
	return results[0]
	
	// panic(NewError(line, NOMETHODERROR, method, gfn.Type()))
}

func NewGoFuncObject(fname string, fn interface{}) *GoFuncObject {
	gf := &GoFuncObject{fname, reflect.TypeOf(fn), fn}
	return gf
}

// Monkey language Object to go language Value.
func ObjectToValue(obj Object, typ reflect.Type) reflect.Value {
	var v reflect.Value
	switch obj := obj.(type) {
	case *Integer:
		if typ == nil {
			v = reflect.ValueOf(obj.Int64).Convert(reflect.TypeOf(1))
		} else {
			v = reflect.ValueOf(obj.Int64).Convert(typ)
		}
	case *UInteger:
		if typ == nil {
			v = reflect.ValueOf(obj.UInt64).Convert(reflect.TypeOf(1))
		} else {
			v = reflect.ValueOf(obj.UInt64).Convert(typ)
		}
	case *Float:
		if typ == nil {
			v = reflect.ValueOf(obj.Float64).Convert(reflect.TypeOf(0.1))
		} else {
			v = reflect.ValueOf(obj.Float64).Convert(typ)
		}
	case *String:
		v = reflect.ValueOf(obj.String)
	case *Boolean:
		v = reflect.ValueOf(obj.Bool)
	case *GoObject:
		if obj.obj == nil {
			var nilObj *Nil
			v = reflect.ValueOf(nilObj)
		} else {
			v = reflect.ValueOf(obj.obj)
		}
	default:
		v = reflect.ValueOf(obj)
	}
	return v
}

// Go language Value to monkey language Object(take care of slice object value)
func GoValueToObject(obj interface{}) Object {
	val := reflect.ValueOf(obj)
	kind := val.Kind()

	switch kind {
	case reflect.Slice, reflect.Array:
		ret := &Array{}
		for i := 0; i < val.Len(); i++ {
			ret.Members = append(ret.Members, GoValueToObject(val.Index(i).Interface()))
		}
		return ret
	case reflect.String:
		return NewString(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInteger(int64(val.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewUInteger(uint64(val.Uint()))
	case reflect.Float32, reflect.Float64:
		return NewFloat(val.Float())
	case reflect.Bool:
		if obj.(bool) == true {
			return TRUE
		} else {
			return FALSE
		}
	default:
		return NewGoObject(obj)
	}
}

func RegisterVars(name string, vars map[string]interface{}) {
	for k, v := range vars {
		//Note: Here we do not convert GoValue to Object using 'GoValueToObject()' function.
		// If we do, then we may get the result which we do not expect. e.g.
		//
		// 		let hours, _ = gtime.ParseDuration("10h")
		// 		gfmt.Println(hours) //Output: 5577006791947779410
		//
		// The expected result should be '10h0m0s', but we get '5577006791947779410'.
		// The reason is that 'ParseDuration' returns a `Duration` type, the `Duration`
		// type's kind is 'reflect.Int64', if we use 'GoValueToObject()' function to 
		// convert `Duration` type to 'Integer' type, then when we print the result, it will
		// print 5577006791947779410. For it to work as expected, we need to keep the `Duration` as it is
		// and make no conversion.
		//SetGlobalObj(name + "." + k, GoValueToObject(v))
		SetGlobalObj(name + "." + k, NewGoObject(v))
	}
}

func RegisterFunctions(name string, vars map[string]interface{}) {
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	for k, v := range vars {
		key := NewString(k)
		hash.Pairs[key.HashKey()] = HashPair{Key: key, Value: NewGoFuncObject(k, v)}
	}

	//Replace all '/' to '_'. e.g. math/rand => math_rand
	newName := strings.Replace(name, "/", "_", -1);
	SetGlobalObj(newName, hash)
}

//func RegisterFunctions(name string, vars []interface{}) {
//	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
//	for _, v := range vars {
//		// In some occasions, below code will introduce panic: 
//		// 		reflect: call of reflect.Value.Pointer on int64 Value
//		fname := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
//		xs := strings.Split(fname, ".")
//		key := NewString(xs[len(xs)-1])
//		fmt.Printf("11111111111111111, key=%v\n", key)
//		hash.Pairs[key.HashKey()] = HashPair{Key: key, Value: NewGoFuncObject(fname, v)}
//	}
//
//	//Replace all '/' to '_'. e.g. math/rand => math_rand
//	newName := strings.Replace(name, "/", "_", -1);
//	SetGlobalObj(newName, hash)
//}
