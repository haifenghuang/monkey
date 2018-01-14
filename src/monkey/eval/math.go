package eval

import (
	_ "fmt"
	"math"
	"math/rand"
	"time"
)

const (
	MATH_OBJ = "MATH_OBJ"
	math_name = "math"
)

type Math struct{}

func NewMathObj() Object {
	ret := &Math{}
	SetGlobalObj(math_name, ret)

	SetGlobalObj(math_name+".Pi", NewFloat(float64(math.Pi)))

	SetGlobalObj(math_name+".E", NewFloat(math.E))
	SetGlobalObj(math_name+".Phi", NewFloat(math.Phi))
	SetGlobalObj(math_name+".Sqrt2", NewFloat(math.Sqrt2))
	SetGlobalObj(math_name+".SqrtE", NewFloat(math.SqrtE))
	SetGlobalObj(math_name+".SqrtPi", NewFloat(math.SqrtPi))
	SetGlobalObj(math_name+".SqrtPhi", NewFloat(math.SqrtPhi))
	SetGlobalObj(math_name+".Ln2", NewFloat(math.Ln2))
	SetGlobalObj(math_name+".Log2E", NewFloat(math.Log2E))
	SetGlobalObj(math_name+".Ln10", NewFloat(math.Ln10))
	SetGlobalObj(math_name+".Log10E", NewFloat(math.Log10E))
	SetGlobalObj(math_name+".NaN", NewFloat(math.NaN()))

	SetGlobalObj(math_name+".MaxFloat64", NewFloat(math.MaxFloat64))
	SetGlobalObj(math_name+".MaxInt64", NewInteger(math.MaxInt64))
	SetGlobalObj(math_name+".MinInt64", NewInteger(math.MinInt64))
	SetGlobalObj(math_name+".MaxUint64", NewUInteger(math.MaxUint64))

	return ret
}

func (m *Math) Inspect() string  { return "<" + math_name + ">" }
func (m *Math) Type() ObjectType { return MATH_OBJ }

func (m *Math) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "abs":
		return m.Abs(line, args...)
	case "acos":
		return m.Acos(line, args...)
	case "acosh":
		return m.Acosh(line, args...)
	case "asin":
		return m.Asin(line, args...)
	case "asinh":
		return m.Asinh(line, args...)
	case "atan":
		return m.Atan(line, args...)
	case "atan2":
		return m.Atan2(line, args...)
	case "atanh":
		return m.Atanh(line, args...)
	case "ceil":
		return m.Ceil(line, args...)
	case "cos":
		return m.Cos(line, args...)
	case "cosh":
		return m.Cosh(line, args...)
	case "exp":
		return m.Exp(line, args...)
	case "floor":
		return m.Floor(line, args...)
	case "inf":
		return m.Inf(line, args...)
	case "isInf":
		return m.IsInf(line, args...)
	case "isNaN":
		return m.IsNaN(line, args...)
	case "max":
		return m.Max(line, args...)
	case "min":
		return m.Min(line, args...)
	case "NaN":
		return m.NaN(line, args...)
	case "pow":
		return m.Pow(line, args...)
	case "sin":
		return m.Sin(line, args...)
	case "sinh":
		return m.Sinh(line, args...)
	case "sqrt":
		return m.Sqrt(line, args...)
	case "tan":
		return m.Tan(line, args...)
	case "tanh":
		return m.Tanh(line, args...)
	case "randSeed":
		return m.RandSeed(line, args...)
	case "rand":
		return m.Rand(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, m.Type()))
}

func (m *Math) Abs(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "abs", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Abs(val.Float64))
}

func (m *Math) Acos(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "acos", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Acos(val.Float64))
}

func (m *Math) Acosh(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "acosh", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Acosh(val.Float64))
}

func (m *Math) Asin(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "asin", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Asin(val.Float64))
}

func (m *Math) Asinh(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "asinh", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Asinh(val.Float64))
}

func (m *Math) Atan(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "atan", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Atan(val.Float64))
}

func (m *Math) Atan2(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	var val1 *Float
	switch input := args[0].(type) {
	case *Integer:
		val1 = NewFloat(float64(input.Int64))
	case *UInteger:
		val1 = NewFloat(float64(input.UInt64))
	case *Float:
		val1 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "atan2", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	var val2 *Float
	switch input := args[1].(type) {
	case *Integer:
		val2 = NewFloat(float64(input.Int64))
	case *UInteger:
		val2 = NewFloat(float64(input.UInt64))
	case *Float:
		val2 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "atan2", "*Integer|*UInteger|*Float", args[1].Type()))
	}

	return NewFloat(math.Atan2(val1.Float64, val2.Float64))
}

func (m *Math) Atanh(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "atanh", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Atanh(val.Float64))
}

func (m *Math) Ceil(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "ceil", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Ceil(val.Float64))
}

func (m *Math) Cos(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "cos", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Cos(val.Float64))
}

func (m *Math) Cosh(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "cosh", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Cosh(val.Float64))
}

func (m *Math) Exp(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "exp", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Exp(val.Float64))
}

func (m *Math) Floor(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "floor", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Floor(val.Float64))
}

func (m *Math) Inf(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Integer
	switch input := args[0].(type) {
	case *Integer:
		val = input
	case *UInteger:
		val = NewInteger(int64(input.UInt64))
	case *Float:
		val = NewInteger(int64(input.Float64))
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "inf", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Inf(int(val.Int64)))
}

func (m *Math) IsInf(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	var val1 *Float
	switch input := args[0].(type) {
	case *Integer:
		val1 = NewFloat(float64(input.Int64))
	case *UInteger:
		val1 = NewFloat(float64(input.UInt64))
	case *Float:
		val1 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "isInf", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	var val2 *Integer
	switch input := args[1].(type) {
	case *Integer:
		val2 = input
	case *UInteger:
		val2 = NewInteger(int64(input.UInt64))
	case *Float:
		val2 = NewInteger(int64(input.Float64))
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "isInf", "*Integer|*UInteger|*Float", args[1].Type()))
	}

	ret := math.IsInf(val1.Float64, int(val2.Int64))
	if ret {
		return TRUE
	}
	return FALSE
}

func (m *Math) IsNaN(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "isNaN", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	ret := math.IsNaN(val.Float64)
	if ret {
		return TRUE
	}
	return FALSE
}

func (m *Math) Max(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	var val1 *Float
	switch input := args[0].(type) {
	case *Integer:
		val1 = NewFloat(float64(input.Int64))
	case *UInteger:
		val1 = NewFloat(float64(input.UInt64))
	case *Float:
		val1 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "max", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	var val2 *Float
	switch input := args[1].(type) {
	case *Integer:
		val2 = NewFloat(float64(input.Int64))
	case *UInteger:
		val2 = NewFloat(float64(input.UInt64))
	case *Float:
		val2 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "max", "*Integer|*UInteger|*Float", args[1].Type()))
	}

	return NewFloat(math.Max(val1.Float64, val2.Float64))
}

func (m *Math) Min(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	var val1 *Float
	switch input := args[0].(type) {
	case *Integer:
		val1 = NewFloat(float64(input.Int64))
	case *UInteger:
		val1 = NewFloat(float64(input.UInt64))
	case *Float:
		val1 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "min", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	var val2 *Float
	switch input := args[1].(type) {
	case *Integer:
		val2 = NewFloat(float64(input.Int64))
	case *UInteger:
		val2 = NewFloat(float64(input.UInt64))
	case *Float:
		val2 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "min", "*Integer|*UInteger|*Float", args[1].Type()))
	}

	return NewFloat(math.Min(val1.Float64, val2.Float64))
}

func (m *Math) NaN(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewFloat(math.NaN())
}

func (m *Math) Pow(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	var val1 *Float
	switch input := args[0].(type) {
	case *Integer:
		val1 = NewFloat(float64(input.Int64))
	case *UInteger:
		val1 = NewFloat(float64(input.UInt64))
	case *Float:
		val1 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "pow", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	var val2 *Float
	switch input := args[1].(type) {
	case *Integer:
		val2 = NewFloat(float64(input.Int64))
	case *UInteger:
		val2 = NewFloat(float64(input.UInt64))
	case *Float:
		val2 = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "second", "pow", "*Integer|*UInteger|*Float", args[1].Type()))
	}

	return NewFloat(math.Pow(val1.Float64, val2.Float64))
}

func (m *Math) Sin(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "sin", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Sin(val.Float64))
}

func (m *Math) Sinh(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "sinh", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Sinh(val.Float64))
}

func (m *Math) Sqrt(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "sqrt", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Sqrt(val.Float64))
}

func (m *Math) Tan(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "tan", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Tan(val.Float64))
}

func (m *Math) Tanh(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Float
	switch input := args[0].(type) {
	case *Integer:
		val = NewFloat(float64(input.Int64))
	case *UInteger:
		val = NewFloat(float64(input.UInt64))
	case *Float:
		val = input
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "tanh", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	return NewFloat(math.Tanh(val.Float64))
}

func (m *Math) RandSeed(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var val *Integer
	switch input := args[0].(type) {
	case *Integer:
		val = input
	case *UInteger:
		val = NewInteger(int64(input.UInt64))
	case *Float:
		val = NewInteger(int64(input.Float64))
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "randseed", "*Integer|*UInteger|*Float", args[0].Type()))
	}

	rand.Seed(val.Int64)
	return NIL
}

func (m *Math) Rand(line string, args ...Object) Object {
	switch len(args) {
	case 0:
		return NewInteger(int64(rand.Int()))
	case 1:
		var val *Integer
		switch input := args[0].(type) {
		case *Integer:
			val = input
		case *UInteger:
			val = NewInteger(int64(input.UInt64))
		case *Float:
			val = NewInteger(int64(input.Float64))
		default:
			panic(NewError(line, PARAMTYPEERROR, "first", "rand", "*Integer|*UInteger|*Float", args[0].Type()))
		}
		return NewInteger(int64(rand.Intn(int(val.Int64))))
	default:
		var low *Integer
		switch input := args[0].(type) {
		case *Integer:
			low = input
		case *UInteger:
			low = NewInteger(int64(input.UInt64))
		case *Float:
			low = NewInteger(int64(input.Float64))
		default:
			panic(NewError(line, PARAMTYPEERROR, "first", "rand", "*Integer|*UInteger|*Float", args[0].Type()))
		}

		var high *Integer
		switch input := args[1].(type) {
		case *Integer:
			high = input
		case *UInteger:
			high = NewInteger(int64(input.UInt64))
		case *Float:
			high = NewInteger(int64(input.Float64))
		default:
			panic(NewError(line, PARAMTYPEERROR, "second", "rand", "*Integer|*UInteger|*Float", args[1].Type()))
		}

		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(int(high.Int64 - low.Int64))
		return NewInteger(int64(n) + low.Int64)
	}
}
