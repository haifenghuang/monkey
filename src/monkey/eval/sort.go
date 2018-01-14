package eval

import (
	"sort"
	"strings"
)

type Ordering int

const (
	Ascending Ordering = iota
	Descending
	CaseInsensitiveAscending
	CaseInsensitiveDescending
)

const (
	SORT_OBJ = "SORT_OBJ"
	sort_name = "sort"
)

func NewSortObj() Object {
	ret := &SortObj{}
	SetGlobalObj(sort_name, ret)

	SetGlobalObj(sort_name+".Ascending", NewInteger(int64(Ascending)))
	SetGlobalObj(sort_name+".Descending", NewInteger(int64(Descending)))
	SetGlobalObj(sort_name+".CaseInsensitiveAscending", NewInteger(int64(CaseInsensitiveAscending)))
	SetGlobalObj(sort_name+".CaseInsensitiveDescending", NewInteger(int64(CaseInsensitiveDescending)))

	return ret
}

type SortObj struct{}

func (s *SortObj) Inspect() string  { return "<" + sort_name + ">" }
func (s *SortObj) Type() ObjectType { return SORT_OBJ }

func (s *SortObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "sortFloats":
		return s.SortFloats(line, args...)
	case "floatsAreSorted":
		return s.FloatsAreSorted(line, args...)
	case "sortInts":
		return s.SortInts(line, args...)
	case "intsAreSorted":
		return s.IntsAreSorted(line, args...)
	case "sortUInts":
		return s.SortUInts(line, args...)
	case "uintsAreSorted":
		return s.UIntsAreSorted(line, args...)
	case "sortStrings":
		return s.SortStrings(line, args...)
	case "stringsAreSorted":
		return s.StringsAreSorted(line, args...)
	}

	panic(NewError(line, NOMETHODERROR, method, s.Type()))
}

func (s *SortObj) SortFloats(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	floatArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sortFloats", "*Array", args[0].Type()))
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "sortFloats", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(Descending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	var sortArr []float64
	for _, v := range floatArr.Members {
		_, ok := v.(Number)
		if !ok {
			panic(NewError(line, GENERICERROR, "Not all data are numbers"))
		}

		var f float64
		if v.Type() == INTEGER_OBJ {
			f = float64(v.(*Integer).Int64)
		} else if v.Type() == UINTEGER_OBJ {
			f = float64(v.(*UInteger).UInt64)
		} else {
			f = v.(*Float).Float64
		}
		sortArr = append(sortArr, f)
	}

	floatSlice := Float64Slice{FloatArr: sortArr, SortOrder: sortOrdering}
	floatSlice.Sort()

	ret := &Array{}
	for _, v := range floatSlice.FloatArr {
		ret.Members = append(ret.Members, NewFloat(v))
	}
	return ret
}

func (s *SortObj) FloatsAreSorted(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	floatArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "floatsAreSorted", "*Array", args[0].Type()))
	}

	var sortArr []float64
	for _, v := range floatArr.Members {
		_, ok := v.(Number)
		if !ok {
			panic(NewError(line, GENERICERROR, "Not all data are numbers"))
		}

		var f float64
		if v.Type() == INTEGER_OBJ {
			f = float64(v.(*Integer).Int64)
		} else if v.Type() == UINTEGER_OBJ {
			f = float64(v.(*UInteger).UInt64)
		} else {
			f = v.(*Float).Float64
		}
		sortArr = append(sortArr, f)
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "floatsAreSorted", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(Descending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	sorted := Float64sAreSorted(sortArr, sortOrdering)
	if sorted {
		return TRUE
	}
	return FALSE
}

func (s *SortObj) SortInts(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	IntsArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sortInts", "*Array", args[0].Type()))
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "sortInts", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(Descending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	var sortArr []int64
	for _, v := range IntsArr.Members {
		if v.Type() != INTEGER_OBJ {
			panic(NewError(line, GENERICERROR, "Not all data are ints"))
		}

		i := v.(*Integer).Int64
		sortArr = append(sortArr, i)
	}

	intSlice := IntSlice{IntArr: sortArr, SortOrder: sortOrdering}
	intSlice.Sort()

	ret := &Array{}
	for _, v := range intSlice.IntArr {
		ret.Members = append(ret.Members, NewInteger(int64(v)))
	}
	return ret
}

func (s *SortObj) IntsAreSorted(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	intArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "intsAreSorted", "*Array", args[0].Type()))
	}

	var sortArr []int64
	for _, v := range intArr.Members {
		if v.Type() != INTEGER_OBJ {
			panic(NewError(line, GENERICERROR, "Not all data are ints"))
		}

		i := v.(*Integer).Int64
		sortArr = append(sortArr, i)
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "intsAreSorted", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(Descending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	sorted := IntsAreSorted(sortArr, sortOrdering)
	if sorted {
		return TRUE
	}
	return FALSE
}

func (s *SortObj) SortUInts(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	UIntsArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sortUInts", "*Array", args[0].Type()))
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "sortInts", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(Descending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	var sortArr []uint64
	for _, v := range UIntsArr.Members {
		if v.Type() != UINTEGER_OBJ {
			panic(NewError(line, GENERICERROR, "Not all data are uints"))
		}

		i := v.(*UInteger).UInt64
		sortArr = append(sortArr, uint64(i))
	}

	intSlice := UIntSlice{UIntArr: sortArr, SortOrder: sortOrdering}
	intSlice.Sort()

	ret := &Array{}
	for _, v := range intSlice.UIntArr {
		ret.Members = append(ret.Members, NewUInteger(uint64(v)))
	}
	return ret
}

func (s *SortObj) UIntsAreSorted(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	uintArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "uintsAreSorted", "*Array", args[0].Type()))
	}

	var sortArr []uint64
	for _, v := range uintArr.Members {
		if v.Type() != UINTEGER_OBJ {
			panic(NewError(line, GENERICERROR, "Not all data are uints"))
		}

		i := v.(*UInteger).UInt64
		sortArr = append(sortArr, i)
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "uintsAreSorted", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(Descending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	sorted := UIntsAreSorted(sortArr, sortOrdering)
	if sorted {
		return TRUE
	}
	return FALSE
}

func (s *SortObj) SortStrings(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	stringsArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sortStrings", "*Array", args[0].Type()))
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "sortStrings", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(CaseInsensitiveDescending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}

	var sortArr []string
	for _, v := range stringsArr.Members {
		if v.Type() != STRING_OBJ {
			panic(NewError(line, GENERICERROR, "Not all data are strings"))
		}

		sortArr = append(sortArr, v.(*String).String)
	}

	strSlice := StringSlice{StrArr: sortArr, SortOrder: sortOrdering}
	strSlice.Sort()

	ret := &Array{}
	for _, v := range strSlice.StrArr {
		ret.Members = append(ret.Members, NewString(v))
	}
	return ret
}

func (s *SortObj) StringsAreSorted(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	stringArr, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "stringsAreSorted", "*Array", args[0].Type()))
	}

	var sortArr []string
	for _, v := range stringArr.Members {
		if v.Type() != STRING_OBJ {
			panic(NewError(line, GENERICERROR, "Not all data are strings"))
		}

		s := v.(*String).String
		sortArr = append(sortArr, s)
	}

	sortOrdering := Ascending
	if len(args) == 2 {
		sortOrder, ok := args[1].(*Integer)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "sortStrings", "*Integer", args[1].Type()))
		}

		if sortOrder.Int64 >= int64(Ascending) && sortOrder.Int64 <= int64(CaseInsensitiveDescending) {
			sortOrdering = Ordering(sortOrder.Int64)
		}
	}
	sorted := StringsAreSorted(sortArr, sortOrdering)
	if sorted {
		return TRUE
	}
	return FALSE

}

/////////////////////////////////////////////////////
//Below codes are mostly copied from golang source with some modifications

// IntSlice attaches the methods of Interface to []int, sorting in ascending/descending order.
type IntSlice struct {
	IntArr    []int64
	SortOrder Ordering
}

func (p IntSlice) Len() int      { return len(p.IntArr) }
func (p IntSlice) Swap(i, j int) { p.IntArr[i], p.IntArr[j] = p.IntArr[j], p.IntArr[i] }
func (p IntSlice) Less(i, j int) bool {
	switch p.SortOrder {
	case Descending:
		return p.IntArr[i] > p.IntArr[j]
	default:
		return p.IntArr[i] < p.IntArr[j]
	}
}

// Sort is a convenience method.
func (p IntSlice) Sort() { sort.Sort(p) }

// UIntSlice attaches the methods of Interface to []int, sorting in ascending/descending order.
type UIntSlice struct {
	UIntArr    []uint64
	SortOrder Ordering
}

func (p UIntSlice) Len() int      { return len(p.UIntArr) }
func (p UIntSlice) Swap(i, j int) { p.UIntArr[i], p.UIntArr[j] = p.UIntArr[j], p.UIntArr[i] }
func (p UIntSlice) Less(i, j int) bool {
	switch p.SortOrder {
	case Descending:
		return p.UIntArr[i] > p.UIntArr[j]
	default:
		return p.UIntArr[i] < p.UIntArr[j]
	}
}

// Sort is a convenience method.
func (p UIntSlice) Sort() { sort.Sort(p) }


// Float64Slice attaches the methods of Interface to []float64, sorting in ascending/descending order.
type Float64Slice struct {
	FloatArr  []float64
	SortOrder Ordering
}

func (p Float64Slice) Len() int      { return len(p.FloatArr) }
func (p Float64Slice) Swap(i, j int) { p.FloatArr[i], p.FloatArr[j] = p.FloatArr[j], p.FloatArr[i] }
func (p Float64Slice) Less(i, j int) bool {
	switch p.SortOrder {
	case Descending:
		return p.FloatArr[i] > p.FloatArr[j] || isNaN(p.FloatArr[i]) && !isNaN(p.FloatArr[j])
	default:
		return p.FloatArr[i] < p.FloatArr[j] || isNaN(p.FloatArr[i]) && !isNaN(p.FloatArr[j])
	}
}

// isNaN is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN(f float64) bool {
	return f != f
}

// Sort is a convenience method.
func (p Float64Slice) Sort() { sort.Sort(p) }

// StringSlice attaches the methods of Interface to []string, sorting in ascending/descending order.
type StringSlice struct {
	StrArr    []string
	SortOrder Ordering
}

func (p StringSlice) Len() int      { return len(p.StrArr) }
func (p StringSlice) Swap(i, j int) { p.StrArr[i], p.StrArr[j] = p.StrArr[j], p.StrArr[i] }
func (p StringSlice) Less(i, j int) bool {
	switch p.SortOrder {
	case Descending:
		return p.StrArr[i] > p.StrArr[j]
	case CaseInsensitiveDescending:
		strI := strings.ToLower(p.StrArr[i])
		strJ := strings.ToLower(p.StrArr[j])
		return strI > strJ
	case CaseInsensitiveAscending:
		strI := strings.ToLower(p.StrArr[i])
		strJ := strings.ToLower(p.StrArr[j])
		return strI < strJ
	default: //default is Ascending
		return p.StrArr[i] < p.StrArr[j]
	}
}

// Sort is a convenience method.
func (p StringSlice) Sort() { sort.Sort(p) }

// IntsAreSorted tests whether a slice of ints is sorted in ascending/descending order.
func IntsAreSorted(a []int64, o Ordering) bool { return sort.IsSorted(IntSlice{IntArr: a, SortOrder: o}) }

// UIntsAreSorted tests whether a slice of uints is sorted in ascending/descending order.
func UIntsAreSorted(a []uint64, o Ordering) bool { return sort.IsSorted(UIntSlice{UIntArr: a, SortOrder: o}) }

// Float64sAreSorted tests whether a slice of float64s is sorted in ascending/descending order.
func Float64sAreSorted(a []float64, o Ordering) bool {
	return sort.IsSorted(Float64Slice{FloatArr: a, SortOrder: o})
}

// StringsAreSorted tests whether a slice of strings is sorted in ascending/descending order.
func StringsAreSorted(a []string, o Ordering) bool {
	return sort.IsSorted(StringSlice{StrArr: a, SortOrder: o})
}
