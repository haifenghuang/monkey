package eval

/*
	NOTE: MOST OF THE CODE COME FROM `https://github.com/ahmetb/go-linq`
	WITH MINOR MODIFICATIONS
*/
import (
	_ "fmt"
	"math"
	"monkey/ast"
	"reflect"
	"regexp"
	"sort"
)

type comparer func(Object, Object) int

// Iterator is an alias for function to iterate over data.
type Iterator func() (item Object, ok *Boolean)

// Query is the type returned from query functions.
type Query struct {
	Iterate func() Iterator
}

type sorter struct {
	items []Object
	less  func(i, j Object) bool
}

func (s sorter) Len() int {
	return len(s.items)
}

func (s sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s sorter) Less(i, j int) bool {
	return s.less(s.items[i], s.items[j])
}

func (q Query) sort(orders []order) (r []Object) {
	next := q.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {
		r = append(r, item)
	}

	if len(r) == 0 {
		return
	}

	for i, j := range orders {
		j.scope.Set(j.selector.Literal.Parameters[0].(*ast.Identifier).Value, r[0])
		ret := Eval(j.selector.Literal.Body, j.scope)
		if obj, ok1 := ret.(*ReturnValue); ok1 {
			ret = obj.Value
		}
		orders[i].compare = getComparer(ret)
	}

	s := sorter{
		items: r,
		less: func(i, j Object) bool {
			for _, order := range orders {
				order.scope.Set(order.selector.Literal.Parameters[0].(*ast.Identifier).Value, i)
				x := Eval(order.selector.Literal.Body, order.scope)
				if obj, ok1 := x.(*ReturnValue); ok1 {
					x = obj.Value
				}

				order.scope.Set(order.selector.Literal.Parameters[0].(*ast.Identifier).Value, j)
				y := Eval(order.selector.Literal.Body, order.scope)
				if obj, ok1 := y.(*ReturnValue); ok1 {
					y = obj.Value
				}
				switch order.compare(x, y) {
				case 0:
					continue
				case -1:
					return !order.desc
				default:
					return order.desc
				}
			}

			return false
		}}

	sort.Sort(s)
	return
}

func (q Query) lessSort(scope *Scope, less *Function) (r []Object) {
	next := q.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {
		r = append(r, item)
	}

	s := sorter{
		items: r,
		less: func(i, j Object) bool {
			scope.Set(less.Literal.Parameters[0].(*ast.Identifier).Value, i)
			scope.Set(less.Literal.Parameters[1].(*ast.Identifier).Value, j)
			cond := Eval(less.Literal.Body, scope)
			if obj, ok1 := cond.(*ReturnValue); ok1 {
				cond = obj.Value
			}
			if IsTrue(cond) {
				return true
			}
			return false
		}}

	sort.Sort(s)
	return
}

const KEYVALUE_OBJ = "KEYVALUE_OBJ"

// KeyValue is a type that is used to iterate over a map (if query is created
// from a map). This type is also used by ToMap() method to output result of a
// query into a map.
type KeyValueObj struct {
	Key   Object
	Value Object
}

func (kv *KeyValueObj) Inspect() string {
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	if hashable, ok := kv.Key.(Hashable); ok {
		hash.Pairs[hashable.HashKey()] = HashPair{Key: kv.Key, Value: kv.Value}
	} else {
		panic(NewError("", KEYERROR, kv.Key.Type()))
	}
	return hash.Inspect()
}

func (kv *KeyValueObj) Type() ObjectType { return KEYVALUE_OBJ }

func (kv *KeyValueObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "key":
		return kv.Key_(line, args...)
	case "value":
		return kv.Value_(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, kv.Type()))
}

//Note : we could not use "Key" as method name, because 'Key' is a member of KeyValueObj
func (kv *KeyValueObj) Key_(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	return kv.Key
}

//Note : we could not use "Value" as method name, because 'Value' is a member of KeyValueObj
func (kv *KeyValueObj) Value_(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return kv.Value
}

const GROUP_OBJ = "GROUP_OBJ"

// Group is a type that is used to store the result of GroupBy method.
type GroupObj struct {
	Key   Object
	Group []Object
}

func (g *GroupObj) Inspect() string {
	arr := &Array{}
	for _, v := range g.Group {
		arr.Members = append(arr.Members, v)
	}
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	if hashable, ok := g.Key.(Hashable); ok {
		hash.Pairs[hashable.HashKey()] = HashPair{Key: g.Key, Value: arr}
	} else {
		panic(NewError("", KEYERROR, g.Key.Type()))
	}
	return hash.Inspect()
}

func (g *GroupObj) Type() ObjectType { return GROUP_OBJ }

func (g *GroupObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "key":
		return g.Key_(line, args...)
	case "value":
		return g.Value(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, g.Type()))
}

//Note : we could not use "Key" as method name, because 'Key' is a member of GroupObj
func (g *GroupObj) Key_(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	return g.Key
}

func (g *GroupObj) Value(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	arr := &Array{}
	for _, v := range g.Group {
		arr.Members = append(arr.Members, v)
	}
	return arr
}

type order struct {
	scope    *Scope
	selector *Function
	compare  comparer
	desc     bool
}

// OrderedQuery is the type returned from OrderBy, OrderByDescending ThenBy and
// ThenByDescending functions.
type OrderedQuery struct {
	Query
	original Query
	orders   []order
}

const linq_name = "linq"

func NewLinqObj() *LinqObj {
	ret := &LinqObj{}
	SetGlobalObj(linq_name, ret)

	return ret
}

func toSlice(lq *LinqObj) (result *Array) {
	next := lq.Query.Iterate()

	result = &Array{}
	for item, ok := next(); ok.Bool; item, ok = next() {
		result.Members = append(result.Members, item)
	}

	return
}

const LINQ_OBJ = "LINQ_OBJ"

type LinqObj struct {
	Query        Query
	OrderedQuery OrderedQuery
}

//lq:linq
func (lq *LinqObj) Inspect() string {
	r := toSlice(lq)
	return r.Inspect()
}

func (lq *LinqObj) Type() ObjectType { return LINQ_OBJ }

func (lq *LinqObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "from":
		return lq.From(line, scope, args...)
	case "range":
		return lq.Range(line, args...)
	case "repeat":
		return lq.Repeat(line, args...)
	case "where":
		return lq.Where(line, scope, args...)
	case "select":
		return lq.Select(line, scope, args...)
	case "first":
		return lq.First(line, args...)
	case "firstWith":
		return lq.FirstWith(line, scope, args...)
	case "last":
		return lq.Last(line, args...)
	case "lastWith":
		return lq.LastWith(line, scope, args...)
	case "forEach":
		return lq.ForEach(line, scope, args...)
	case "forEachIndexed":
		return lq.ForEachIndexed(line, scope, args...)
	case "take":
		return lq.Take(line, args...)
	case "takeWhile":
		return lq.TakeWhile(line, scope, args...)
	case "takeWhileIndexed":
		return lq.TakeWhileIndexed(line, scope, args...)
	case "skip":
		return lq.Skip(line, args...)
	case "skipWhile":
		return lq.SkipWhile(line, scope, args...)
	case "skipWhileIndexed":
		return lq.SkipWhileIndexed(line, scope, args...)
	case "groupBy":
		return lq.GroupBy(line, scope, args...)
	case "join":
		return lq.Join(line, scope, args...)
	case "zip":
		return lq.Zip(line, scope, args...)
	case "union":
		return lq.Union(line, args...)
	case "selectMany":
		return lq.SelectMany(line, scope, args...)
	case "selectManyIndexed":
		return lq.SelectManyIndexed(line, scope, args...)
	case "selectManyBy":
		return lq.SelectManyBy(line, scope, args...)
	case "selectManyByIndexed":
		return lq.SelectManyByIndexed(line, scope, args...)
	case "reverse":
		return lq.Reverse(line, args...)
	case "except":
		return lq.Except(line, args...)
	case "exceptBy":
		return lq.ExceptBy(line, scope, args...)
	case "append":
		return lq.Append(line, args...)
	case "concat":
		return lq.Concat(line, args...)
	case "prepend":
		return lq.Prepend(line, args...)
	case "distinct":
		return lq.Distinct(line, args...)
	case "distinctBy":
		return lq.DistinctBy(line, scope, args...)
	case "intersect":
		return lq.Intersect(line, args...)
	case "intersectBy":
		return lq.IntersectBy(line, scope, args...)
	case "aggregate":
		return lq.Aggregate(line, scope, args...)
	case "aggregateWithSeed":
		return lq.AggregateWithSeed(line, scope, args...)
	case "aggregateWithSeedBy":
		return lq.AggregateWithSeedBy(line, scope, args...)
	case "all":
		return lq.All(line, scope, args...)
	case "any":
		return lq.Any(line, args...)
	case "anyWith":
		return lq.AnyWith(line, scope, args...)
	case "contains":
		return lq.Contains(line, args...)
	case "count":
		return lq.Count(line, args...)
	case "countWith":
		return lq.CountWith(line, scope, args...)
	case "sequenceEqual":
		return lq.SequenceEqual(line, args...)
	case "single":
		return lq.Single(line, args...)
	case "singleWith":
		return lq.SingleWith(line, scope, args...)
	case "sumInts":
		return lq.SumInts(line, args...)
	case "sumFloats":
		return lq.SumFloats(line, args...)
	case "min":
		return lq.Min(line, args...)
	case "max":
		return lq.Max(line, args...)
	case "average":
		return lq.Average(line, args...)
	case "orderBy":
		return lq.OrderBy(line, scope, args...)
	case "orderByDescending":
		return lq.OrderByDescending(line, scope, args...)
	case "thenBy":
		return lq.ThenBy(line, scope, args...)
	case "thenByDescending":
		return lq.ThenByDescending(line, scope, args...)
	case "sort":
		return lq.Sort(line, scope, args...)
	case "toSlice":
		return lq.ToSlice(line, args...)
	case "toOrderedSlice":
		return lq.ToOrderedSlice(line, args...)
	case "toMap":
		return lq.ToMap(line, scope, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, lq.Type()))
}

// From initializes a linq query with passed slice, array or map as the source.
// String, channel or struct implementing Iterable interface can be used as an
// input. In this case From delegates it to FromString, FromChannel and
// FromIterable internally.
func (lq *LinqObj) From(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 && len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "1|2|3", len(args)))
	}

	obj := args[0]
	//check object type
	if obj.Type() != STRING_OBJ && obj.Type() != ARRAY_OBJ &&
		obj.Type() != HASH_OBJ && obj.Type() != FILE_OBJ && obj.Type() != CSV_OBJ {
		panic(NewError(line, PARAMTYPEERROR, "first", "from", "*Hash|*Array|*String|*File|*CsvObj", obj.Type()))
	}

	switch obj.Type() {
	case FILE_OBJ:
		if len(args) != 1 && len(args) != 2 && len(args) != 3 {
			panic(NewError(line, GENERICERROR, "File object should have 1|2|3 parameter(s):from(file, [field-separator], [commentFn])"))
		}

		var fsStr string = "," //field separator(default is ",")
		if len(args) == 2 {
			//get the field separator
			fsObj, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "from", "*String", args[1].Type()))
			}
			fsStr = fsObj.String
		}

		var commentFn *Function = nil
		if len(args) == 3 {
			//get the comment function
			var ok bool
			commentFn, ok = args[2].(*Function)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "third", "from", "*Function", args[2].Type()))
			}
		}

		scop := NewScope(scope)
		arr := &Array{}
		var lineNo int64 = 0
		l := obj.(*FileObject).ReadLine(line)
		for l != NIL {
			//ignore comment if it has
			if commentFn != nil {
				scop.Set(commentFn.Literal.Parameters[0].(*ast.Identifier).Value, l.(*String))
				cond := Eval(commentFn.Literal.Body, scop)
				if obj, ok1 := cond.(*ReturnValue); ok1 {
					cond = obj.Value
				}
				if IsTrue(cond) {
					//read the next line
					l = obj.(*FileObject).ReadLine(line)
					lineNo++
					continue
				}
			}

			hash := &Hash{Pairs: make(map[HashKey]HashPair)}
			//-1 means the whole line
			fieldIndex := NewInteger(0)
			hash.Pairs[fieldIndex.HashKey()] = HashPair{Key: fieldIndex, Value: l}

			lineNoKey := NewString("line")
			hash.Pairs[lineNoKey.HashKey()] = HashPair{Key: lineNoKey, Value: NewInteger(lineNo)}

			//strArr := strings.Split(l.(*String).String, fsStr)
			strArr := regexp.MustCompile(fsStr).Split(l.(*String).String, -1)

			nfKey := NewString("nf") //nf : number of fields
			hash.Pairs[nfKey.HashKey()] = HashPair{Key: nfKey, Value: NewInteger(int64(len(strArr)))}

			for idx, v := range strArr {
				fieldIndex := NewInteger(int64(idx+1))
				hash.Pairs[fieldIndex.HashKey()] = HashPair{Key: fieldIndex, Value: NewString(v)}
			}
			arr.Members = append(arr.Members, hash)

			//read the next line
			l = obj.(*FileObject).ReadLine(line)
			lineNo++
		}

		//Now the 'arr' variable is like below:
		//  arr = [
		//      {"line" =>LineNo1, "nf" =>line1's number of fields, 0 => line1, 1 => field1, 2 =>field2, ...},
		//      {"line" =>LineNo2, "nf" =>line2's number of fields, 0 => line2, 1 => field1, 2 =>field2, ...}
		//  ]
		len := len(arr.Members)
		//must return a new LinqObj
		return &LinqObj{Query: Query{
			Iterate: func() Iterator {
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = arr.Members[index]
						index++
					}
					return
				}
			},
		}}
	case CSV_OBJ:
		arr := &Array{}
		str2DArr := obj.(*CsvObj).ReadAll(line)
		for _, fieldArr := range str2DArr.(*Array).Members {
			hash := &Hash{Pairs: make(map[HashKey]HashPair)}

			nfKey := NewString("nf") //nf : number of fields
			hash.Pairs[nfKey.HashKey()] = HashPair{Key: nfKey, Value: NewInteger(int64(len(fieldArr.(*Array).Members)))}

			for idx, field := range fieldArr.(*Array).Members {
				fieldIdx := NewInteger(int64(idx+1))
				hash.Pairs[fieldIdx.HashKey()] = HashPair{Key: fieldIdx, Value: field}
			}
			arr.Members = append(arr.Members, hash)
		}

		//Now the 'arr' variable is like below:
		//  arr = [
		//      {"nf" =>line1's number of fields, 1 => field1, 2 =>field2, ...},
		//      {"nf" =>line2's number of fields, 1 => field1, 2 =>field2, ...}
		//  ]
		len := len(arr.Members)
		//must return a new LinqObj
		return &LinqObj{Query: Query{
			Iterate: func() Iterator {
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = arr.Members[index]
						index++
					}
					return
				}
			},
		}}
	case STRING_OBJ:
		source := obj.(*String).String
		runes := []rune(source)
		len := len(runes)

		//must return a new LinqObj
		return &LinqObj{Query: Query{
			Iterate: func() Iterator {
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = NewString(string(runes[index]))
						index++
					}

					return
				}
			},
		}}
	case ARRAY_OBJ:
		arr := obj.(*Array)
		len := len(arr.Members)

		//must return a new LinqObj
		return &LinqObj{Query: Query{
			Iterate: func() Iterator {
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = arr.Members[index]
						index++
					}
					return
				}
			},
		}}

	case HASH_OBJ:
		hash := obj.(*Hash)
		len := len(hash.Pairs)

		return &LinqObj{Query: Query{
			Iterate: func() Iterator {
				index := 0

				keys := &Array{}
				values := &Array{}
				for _, pair := range hash.Pairs {
					keys.Members = append(keys.Members, pair.Key)
					values.Members = append(values.Members, pair.Value)
				}

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						key := keys.Members[index]
						value := values.Members[index]

						item = &KeyValueObj{Key: key, Value: value}

						index++
					}

					return
				}
			},
		}}
	default:
		return &LinqObj{Query: Query{Iterate: obj.(*LinqObj).Query.Iterate}}
	} //end switch

	return NIL
}

// Range generates a sequence of integral numbers within a specified range.
func (lq *LinqObj) Range(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	startObj, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "range", "*Integer", args[0].Type()))
	}

	countObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "range", "*Integer", args[1].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			var index int64 = 0
			current := startObj.Int64

			return func() (item Object, ok *Boolean) {
				if index >= countObj.Int64 {
					return NIL, &Boolean{Bool: false, Valid: true}
				}

				item, ok = NewInteger(current), &Boolean{Bool: true, Valid: true}

				index++
				current++
				return
			}
		},
	}}
}

// Repeat generates a sequence that contains one repeated value.
func (lq *LinqObj) Repeat(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	valueObj, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "repeat", "Object", args[0].Type()))
	}

	countObj, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "repeat", "*Integer", args[1].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			var index int64 = 0

			return func() (item Object, ok *Boolean) {
				if index >= countObj.Int64 {
					return NIL, &Boolean{Bool: false, Valid: true}
				}

				item, ok = valueObj, &Boolean{Bool: true, Valid: true}

				index++
				return
			}
		},
	}}
}

// Where filters a collection of values based on a predicate.
func (lq *LinqObj) Where(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "where", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {
					s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, item)
					cond := Eval(block.Literal.Body, s)
					if obj, ok1 := cond.(*ReturnValue); ok1 {
						cond = obj.Value
					}
					if IsTrue(cond) {
						return
					}
				}

				return
			}
		},
	}}
}

func (lq *LinqObj) Select(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "select", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()

			return func() (item Object, ok *Boolean) {
				var obj Object
				obj, ok = next()
				if ok.Bool {
					s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, obj)
					item = Eval(block.Literal.Body, s)
					if obj, ok1 := item.(*ReturnValue); ok1 {
						item = obj.Value
					}
				}
				return
			}
		},
	}}
}

// Last returns the first element of a collection.
func (lq *LinqObj) First(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	item, _ := lq.Query.Iterate()()

	if item == nil {
		return NIL
	}
	return item
}

// FirstWith returns the first element of a collection that satisfies a
// specified condition.
func (lq *LinqObj) FirstWith(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "firstWith", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	next := lq.Query.Iterate()

	for item, ok := next(); ok.Bool; item, ok = next() {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, item)
		cond := Eval(block.Literal.Body, s)
		if obj, ok1 := cond.(*ReturnValue); ok1 {
			cond = obj.Value
		}
		if IsTrue(cond) {
			return item
		}
	}

	return NIL
}

// Last returns the last element of a collection.
func (lq *LinqObj) Last(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()

	var obj Object = NIL
	for item, ok := next(); ok.Bool; item, ok = next() {
		obj = item
	}

	return obj
}

// LastWith returns the last element of a collection that satisfies a specified
// condition.
func (lq *LinqObj) LastWith(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "lastWith", "*Function", args[0].Type()))
	}

	s := NewScope(scope)
	next := lq.Query.Iterate()

	var obj Object = NIL
	for item, ok := next(); ok.Bool; item, ok = next() {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, item)
		cond := Eval(block.Literal.Body, s)
		if obj, ok1 := cond.(*ReturnValue); ok1 {
			cond = obj.Value
		}
		if IsTrue(cond) {
			obj = item
		}
	}

	return obj
}

// ForEach performs the specified action on each element of a collection.
func (lq *LinqObj) ForEach(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "forEach", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	next := lq.Query.Iterate()

	for item, ok := next(); ok.Bool; item, ok = next() {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, item)
		Eval(block.Literal.Body, s)
	}

	return NIL
}

// ForEachIndexed performs the specified action on each element of a collection.
func (lq *LinqObj) ForEachIndexed(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "forEachIndexed", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	next := lq.Query.Iterate()
	index := 0

	for item, ok := next(); ok.Bool; item, ok = next() {
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, NewInteger(int64(index)))
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, item)
		Eval(block.Literal.Body, s)
		index++
	}

	return NIL
}

// Take returns a specified number of contiguous elements from the start of a
// collection.
func (lq *LinqObj) Take(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	iObj, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "take", "*Integer", args[0].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			n := iObj.Int64

			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Bool: true, Valid: true}
				if n <= 0 {
					ok.Bool = false
					return
				}

				n--
				return next()
			}
		},
	}}
}

// TakeWhile returns elements from a collection as long as a specified condition
// is true, and then skips the remaining elements.
// FirstWith returns the first element of a collection that satisfies a
// specified condition.
func (lq *LinqObj) TakeWhile(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "takeWhile", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			done := false

			return func() (item Object, ok *Boolean) {
				if done {
					ok.Bool = false
					return
				}

				item, ok = next()
				if !ok.Bool {
					done = true
					return
				}

				s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, item)
				cond := Eval(block.Literal.Body, s)
				if obj, ok1 := cond.(*ReturnValue); ok1 {
					cond = obj.Value
				}
				if IsTrue(cond) {
					return
				}

				done = true
				return NIL, &Boolean{Bool: false, Valid: true}
			}
		},
	}}
}

// TakeWhileIndexed returns elements from a collection as long as a specified
// condition is true. The element's index is used in the logic of the predicate
// function. The first argument of predicate represents the zero-based index of
// the element within collection. The second argument represents the element to
// test.
func (lq *LinqObj) TakeWhileIndexed(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "takeWhileIndexed", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			done := false
			index := 0

			return func() (item Object, ok *Boolean) {
				if done {
					ok.Bool = false
					return
				}

				item, ok = next()
				if !ok.Bool {
					done = true
					return
				}

				s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, NewInteger(int64(index)))
				s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, item)
				cond := Eval(block.Literal.Body, s)
				if obj, ok1 := cond.(*ReturnValue); ok1 {
					cond = obj.Value
				}
				if IsTrue(cond) {
					index++
					return
				}

				done = true
				return NIL, &Boolean{Bool: false, Valid: true}
			}
		},
	}}
}

// Skip bypasses a specified number of elements in a collection and then returns
// the remaining elements..
func (lq *LinqObj) Skip(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	iObj, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "skip", "*Integer", args[0].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			n := iObj.Int64

			return func() (item Object, ok *Boolean) {
				for ; n > 0; n-- {
					item, ok = next()
					if !ok.Bool {
						return
					}
				}

				return next()
			}
		},
	}}
}

// SkipWhile bypasses elements in a collection as long as a specified condition
// is true and then returns the remaining elements.
//
// This method tests each element by using predicate and skips the element if
// the result is true. After the predicate function returns false for an
// element, that element and the remaining elements in source are returned and
// there are no more invocations of predicate.
func (lq *LinqObj) SkipWhile(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "skipWhile", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			ready := false

			return func() (item Object, ok *Boolean) {
				for !ready {
					item, ok = next()
					if !ok.Bool {
						return
					}

					s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, item)
					cond := Eval(block.Literal.Body, s)
					if obj, ok1 := cond.(*ReturnValue); ok1 {
						cond = obj.Value
					}
					ready = !IsTrue(cond)
					if ready {
						return
					}
				}

				return next()
			}
		},
	}}
}

// SkipWhileIndexed bypasses elements in a collection as long as a specified
// condition is true and then returns the remaining elements. The element's
// index is used in the logic of the predicate function.
//
// This method tests each element by using predicate and skips the element if
// the result is true. After the predicate function returns false for an
// element, that element and the remaining elements in source are returned and
// there are no more invocations of predicate.
func (lq *LinqObj) SkipWhileIndexed(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "skipWhileIndexed", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			ready := false
			index := 0

			return func() (item Object, ok *Boolean) {
				for !ready {
					item, ok = next()
					if !ok.Bool {
						return
					}

					s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, NewInteger(int64(index)))
					s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, item)
					cond := Eval(block.Literal.Body, s)
					if obj, ok1 := cond.(*ReturnValue); ok1 {
						cond = obj.Value
					}
					ready = !IsTrue(cond)
					if ready {
						return
					}
					index++
				}

				return next()
			}
		},
	}}
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
func (lq *LinqObj) GroupBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	keySelector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "groupBy", "*Function", args[0].Type()))
	}

	elementSelector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "groupBy", "*Function", args[1].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		func() Iterator {
			next := lq.Query.Iterate()
			set := make(map[Object][]Object)

			for item, ok := next(); ok.Bool; item, ok = next() {
				s.Set(keySelector.Literal.Parameters[0].(*ast.Identifier).Value, item)
				key := Eval(keySelector.Literal.Body, s)
				if obj, ok1 := key.(*ReturnValue); ok1 {
					key = obj.Value
				}

				s.Set(elementSelector.Literal.Parameters[0].(*ast.Identifier).Value, item)
				element := Eval(elementSelector.Literal.Body, s)
				if obj, ok1 := element.(*ReturnValue); ok1 {
					element = obj.Value
				}

				set[key] = append(set[key], element)
			}

			len := len(set)
			idx := 0
			groups := make([]*GroupObj, len)
			for k, v := range set {
				groups[idx] = &GroupObj{k, v}
				idx++
			}

			index := 0
			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Valid: true}
				ok.Bool = index < len
				if ok.Bool {
					item = groups[index]
					index++
				}

				return
			}
		},
	}}
}

// Join correlates the elements of two collection based on matching keys.
//
// A join refers to the operation of correlating the elements of two sources of
// information based on a common key. Join brings the two information sources
// and the keys by which they are matched together in one method callp. This
// differs from the use of SelectMany, which requires more than one method call
// to perform the same operation.
//
// Join preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
func (lq *LinqObj) Join(line string, scope *Scope, args ...Object) Object {
	if len(args) != 4 {
		panic(NewError(line, ARGUMENTERROR, "4", len(args)))
	}

	inner, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "join", "*LinqObj", args[0].Type()))
	}

	outerKeySelector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "join", "*Function", args[1].Type()))
	}

	innerKeySelector, ok := args[2].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "join", "*Function", args[2].Type()))
	}

	resultSelector, ok := args[3].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "fourth", "join", "*Function", args[3].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			outernext := lq.Query.Iterate()
			innernext := inner.Query.Iterate()

			innerLookup := make(map[Object][]Object)
			for innerItem, ok := innernext(); ok.Bool; innerItem, ok = innernext() {
				s.Set(innerKeySelector.Literal.Parameters[0].(*ast.Identifier).Value, innerItem)
				innerKey := Eval(innerKeySelector.Literal.Body, s)
				if obj, ok1 := innerKey.(*ReturnValue); ok1 {
					innerKey = obj.Value
				}
				innerLookup[innerKey] = append(innerLookup[innerKey], innerItem)
			}

			var outerItem Object
			var innerGroup []Object
			innerLen, innerIndex := 0, 0

			return func() (item Object, ok *Boolean) {
				if innerIndex >= innerLen {
					has := false
					for !has {
						outerItem, ok = outernext()
						if !ok.Bool {
							return
						}

						s.Set(outerKeySelector.Literal.Parameters[0].(*ast.Identifier).Value, outerItem)
						outerKey := Eval(outerKeySelector.Literal.Body, s)
						if obj, ok1 := outerKey.(*ReturnValue); ok1 {
							outerKey = obj.Value
						}

						for k, _ := range innerLookup {
							if reflect.DeepEqual(k, outerKey) {
								has = true
								break
							}
						}

						innerGroup, _ = innerLookup[outerKey]
						innerLen = len(innerGroup)
						innerIndex = 0
					}
				}

				s.Set(resultSelector.Literal.Parameters[0].(*ast.Identifier).Value, outerItem)
				s.Set(resultSelector.Literal.Parameters[1].(*ast.Identifier).Value, innerGroup[innerIndex])
				item = Eval(resultSelector.Literal.Body, s)
				if obj, ok1 := item.(*ReturnValue); ok1 {
					item = obj.Value
				}
				innerIndex++
				return item, &Boolean{Bool: true, Valid: true}
			}
		},
	}}
}

// Zip applies a specified function to the corresponding elements of two
// collections, producing a collection of the results.
//
// The method steps through the two input collections, applying function
// resultSelector to corresponding elements of the two collections. The method
// returns a collection of the values that are returned by resultSelector. If
// the input collections do not have the same number of elements, the method
// combines elements until it reaches the end of one of the collections. For
// example, if one collection has three elements and the other one has four, the
// result collection has only three elements.
func (lq *LinqObj) Zip(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	lq2, ok := args[0].(*LinqObj) //lq:linq
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "zip", "*LinqObj", args[0].Type()))
	}

	resultSelector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "zip", "*Function", args[1].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next1 := lq.Query.Iterate()
			next2 := lq2.Query.Iterate()

			return func() (item Object, ok *Boolean) {
				item1, ok1 := next1()
				item2, ok2 := next2()

				if ok1.Bool && ok2.Bool {
					s.Set(resultSelector.Literal.Parameters[0].(*ast.Identifier).Value, item1)
					s.Set(resultSelector.Literal.Parameters[1].(*ast.Identifier).Value, item2)
					result := Eval(resultSelector.Literal.Body, s)
					if obj, ok1 := result.(*ReturnValue); ok1 {
						result = obj.Value
					}
					return result, &Boolean{Bool: true, Valid: true}
				}

				return NIL, &Boolean{Bool: false, Valid: true}
			}
		},
	}}
}

// Union produces the set union of two collections.
//
// This method excludes duplicates from the return set. This is different
// behavior to the Concat method, which returns all the elements in the input
// collection including duplicates.
func (lq *LinqObj) Union(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	lq2, ok := args[0].(*LinqObj) //lq:linq
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "union", "*LinqObj", args[0].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			next2 := lq2.Query.Iterate()

			set := make(map[Object]bool)
			use1 := true

			return func() (item Object, ok *Boolean) {
				if use1 {
					for item, ok = next(); ok.Bool; item, ok = next() {
						for k, _ := range set {
							if reflect.DeepEqual(k, item) {
								set[item] = true
								break
							}
						}
						if _, has := set[item]; !has {
							set[item] = true
							return
						}
					}

					use1 = false
				}

				for item, ok = next2(); ok.Bool; item, ok = next2() {
					for k, _ := range set {
						if reflect.DeepEqual(k, item) {
							set[item] = true
							break
						}
					}
					if _, has := set[item]; !has {
						set[item] = true
						return
					}
				}

				return
			}
		},
	}}
}

// SelectMany projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
func (lq *LinqObj) SelectMany(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "selectMany", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			outernext := lq.Query.Iterate()

			var inner Object = NIL //initialized to 'NIL' OBJECT
			var innernext Iterator

			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Bool: false, Valid: true}
				for !ok.Bool {
					if inner == NIL {

						inner, ok = outernext()
						if !ok.Bool {
							return
						}

						s.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, inner)
						lq := Eval(selector.Literal.Body, s)
						if obj, ok1 := lq.(*ReturnValue); ok1 {
							lq = obj.Value
						}

						if lq.Type() != LINQ_OBJ {
							panic(NewError(line, GENERICERROR, "Function should return a *LinqObj"))
						}

						innernext = lq.(*LinqObj).Query.Iterate()
					}

					item, ok = innernext()
					if !ok.Bool {
						inner = NIL
					}
				}

				return
			}
		},
	}}
}

// SelectManyIndexed projects each element of a collection to a Query, iterates
// and flattens the resulting collection into one collection.
//
// The first argument to selector represents the zero-based index of that
// element in the source collection. This can be useful if the elements are in a
// known order and you want to do something with an element at a particular
// index, for example. It can also be useful if you want to retrieve the index
// of one or more elements. The second argument to selector represents the
// element to process.
func (lq *LinqObj) SelectManyIndexed(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "selectManyIndexed", "*Function", args[0].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			outernext := lq.Query.Iterate()
			index := 0
			var inner Object = NIL
			var innernext Iterator

			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Bool: false, Valid: true}
				for !ok.Bool {
					if inner == NIL {
						inner, ok = outernext()
						if !ok.Bool {
							return
						}

						s.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, NewInteger(int64(index)))
						s.Set(selector.Literal.Parameters[1].(*ast.Identifier).Value, inner)
						lq := Eval(selector.Literal.Body, s)
						if obj, ok1 := lq.(*ReturnValue); ok1 {
							lq = obj.Value
						}
						if lq.Type() != LINQ_OBJ {
							panic(NewError(line, GENERICERROR, "Function should return a *LinqObj"))
						}
						innernext = lq.(*LinqObj).Query.Iterate()
						index++
					}

					item, ok = innernext()
					if !ok.Bool {
						inner = NIL
					}
				}

				return
			}
		},
	}}
}

// SelectManyBy projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection, and invokes a result
// selector function on each element therein.
func (lq *LinqObj) SelectManyBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "selectManyBy", "*Function", args[0].Type()))
	}

	resultSelector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "selectManyBy", "*Function", args[1].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			outernext := lq.Query.Iterate()
			var outer Object = NIL
			var innernext Iterator

			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Bool: false, Valid: true}
				for !ok.Bool {
					if outer == NIL {
						outer, ok = outernext()
						if !ok.Bool {
							return
						}

						s.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, outer)
						lq := Eval(selector.Literal.Body, s)
						if obj, ok1 := lq.(*ReturnValue); ok1 {
							lq = obj.Value
						}
						if lq.Type() != LINQ_OBJ {
							panic(NewError(line, GENERICERROR, "Function should return a *LinqObj"))
						}
						innernext = lq.(*LinqObj).Query.Iterate()
					}

					item, ok = innernext()
					if !ok.Bool {
						outer = NIL
					}
				}

				s.Set(resultSelector.Literal.Parameters[0].(*ast.Identifier).Value, item)
				s.Set(resultSelector.Literal.Parameters[1].(*ast.Identifier).Value, outer)
				item = Eval(resultSelector.Literal.Body, s)
				if obj, ok1 := item.(*ReturnValue); ok1 {
					item = obj.Value
				}
				return
			}
		},
	}}
}

// SelectManyByIndexed projects each element of a collection to a Query,
// iterates and flattens the resulting collection into one collection, and
// invokes a result selector function on each element therein. The index of each
// source element is used in the intermediate projected form of that element.
func (lq *LinqObj) SelectManyByIndexed(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "selectManyByIndexed", "*Function", args[0].Type()))
	}

	resultSelector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "selectManyByIndexed", "*Function", args[1].Type()))
	}

	s := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			outernext := lq.Query.Iterate()
			index := 0
			var outer Object = NIL
			var innernext Iterator

			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Bool: false, Valid: true}
				for !ok.Bool {
					if outer == NIL {
						outer, ok = outernext()
						if !ok.Bool {
							return
						}

						s.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, NewInteger(int64(index)))
						s.Set(selector.Literal.Parameters[1].(*ast.Identifier).Value, outer)
						item := Eval(selector.Literal.Body, s)
						if obj, ok1 := item.(*ReturnValue); ok1 {
							item = obj.Value
						}
						if item.Type() != LINQ_OBJ {
							panic(NewError(line, GENERICERROR, "Function should return a *LinqObj"))
						}
						innernext = item.(*LinqObj).Query.Iterate()
						index++
					}

					item, ok = innernext()
					if !ok.Bool {
						outer = NIL
					}
				}

				s.Set(resultSelector.Literal.Parameters[0].(*ast.Identifier).Value, item)
				s.Set(resultSelector.Literal.Parameters[1].(*ast.Identifier).Value, outer)
				item = Eval(resultSelector.Literal.Body, s)
				if obj, ok1 := item.(*ReturnValue); ok1 {
					item = obj.Value
				}
				return
			}
		},
	}}
}

// Reverse inverts the order of the elements in a collection.
//
// Unlike OrderBy, this sorting method does not consider the actual values
// themselves in determining the order. Rather, it just returns the elements in
// the reverse order from which they are produced by the underlying source.
func (lq *LinqObj) Reverse(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()

			items := []Object{}
			for item, ok := next(); ok.Bool; item, ok = next() {
				items = append(items, item)
			}

			index := len(items) - 1
			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Bool: false, Valid: true}
				if index < 0 {
					return
				}

				ok.Bool = true
				item = items[index]
				index--
				return
			}
		},
	}}
}

// Except produces the set difference of two sequences. The set difference is
// the members of the first sequence that don't appear in the second sequence.
func (lq *LinqObj) Except(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	lq2, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "except", "*LinqObj", args[0].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()

			next2 := lq2.Query.Iterate()
			set := make(map[Object]bool)
			for i, ok := next2(); ok.Bool; i, ok = next2() {
				set[i] = true
			}

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {

					for k, _ := range set {
						if reflect.DeepEqual(k, item) {
							set[item] = true
							break
						}
					}
					if _, has := set[item]; !has {
						return
					}
				}

				return
			}
		},
	}}
}

// ExceptBy invokes a transform function on each element of a collection and
// produces the set difference of two sequences. The set difference is the
// members of the first sequence that don't appear in the second sequence.
func (lq *LinqObj) ExceptBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	lq2, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "exceptBy", "*LinqObj", args[0].Type()))
	}

	selector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "exceptBy", "*Function", args[1].Type()))
	}

	scop := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()

			next2 := lq2.Query.Iterate()
			set := make(map[Object]bool)
			for i, ok := next2(); ok.Bool; i, ok = next2() {
				scop.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, i)
				s := Eval(selector.Literal.Body, scop)
				if obj, ok1 := s.(*ReturnValue); ok1 {
					s = obj.Value
				}
				set[s] = true
			}

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {
					scop.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, item)
					s := Eval(selector.Literal.Body, scop)
					if obj, ok1 := s.(*ReturnValue); ok1 {
						s = obj.Value
					}

					for k, _ := range set {
						if reflect.DeepEqual(k, s) {
							set[s] = true
							break
						}
					}

					if _, has := set[s]; !has {
						return
					}
				}

				return
			}
		},
	}}
}

// Append inserts an item to the end of a collection, so it becomes the last
// item.
func (lq *LinqObj) Append(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			appended := false

			return func() (Object, *Boolean) {
				i, ok := next()
				if ok.Bool {
					return i, ok
				}

				if !appended {
					appended = true
					return args[0], &Boolean{Bool: true, Valid: true}
				}

				return NIL, &Boolean{Bool: false, Valid: true}
			}
		},
	}}
}

// Concat concatenates two collections.
//
// The Concat method differs from the Union method because the Concat method
// returns all the original elements in the input sequences. The Union method
// returns only unique elements.
func (lq *LinqObj) Concat(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	lq2, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "concat", "*LinqObj", args[0].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			next2 := lq2.Query.Iterate()
			use1 := true

			return func() (item Object, ok *Boolean) {
				if use1 {
					item, ok = next()
					if ok.Bool {
						return
					}

					use1 = false
				}

				return next2()
			}
		},
	}}
}

func (lq *LinqObj) Prepend(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			prepended := false

			return func() (Object, *Boolean) {
				if prepended {
					return next()
				}

				prepended = true
				return args[0], &Boolean{Bool: true, Valid: true}
			}
		},
	}}
}

// Distinct method returns distinct elements from a collection. The result is an
// unordered collection that contains no duplicate values.
func (lq *LinqObj) Distinct(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			set := make(map[Object]bool)

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {

					//Why we need this 'for' logic? because in monkey, every object is
					//a pointer, so if we do not have this, the below `if _, has := set[item]; !has {`
					//will always evaluated to 'true', so we need to use 'DeepEqual'
					for k, _ := range set {
						if reflect.DeepEqual(k, item) {
							set[item] = true
							break
						}
					}

					if _, has := set[item]; !has {
						set[item] = true
						return
					}
				}

				return
			}
		},
	}}
}

// DistinctBy method returns distinct elements from a collection. This method
// executes selector function for each element to determine a value to compare.
// The result is an unordered collection that contains no duplicate values.
func (lq *LinqObj) DistinctBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "distinctBy", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			set := make(map[Object]bool)

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {
					scop.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, item)
					s := Eval(selector.Literal.Body, scop)
					if obj, ok1 := s.(*ReturnValue); ok1 {
						s = obj.Value
					}

					for k, _ := range set {
						if reflect.DeepEqual(k, s) {
							set[s] = true
							break
						}
					}
					if _, has := set[s]; !has {
						set[s] = true
						return
					}
				}

				return
			}
		},
	}}
}

// Intersect produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
func (lq *LinqObj) Intersect(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	lq2, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "intersect", "*LinqObj", args[0].Type()))
	}

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			next2 := lq2.Query.Iterate()

			set := make(map[Object]bool)
			for item, ok := next2(); ok.Bool; item, ok = next2() {
				set[item] = true
			}

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {

					for k, _ := range set {
						if reflect.DeepEqual(k, item) {
							set[item] = false
							break
						}
					}
					if _, has := set[item]; has {
						delete(set, item)
						return
					}
				}

				return
			}
		},
	}}
}

// IntersectBy produces the set intersection of the source collection and the
// provided input collection. The intersection of two sets A and B is defined as
// the set that contains all the elements of A that also appear in B, but no
// other elements.
//
// IntersectBy invokes a transform function on each element of both collections.
func (lq *LinqObj) IntersectBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	lq2, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "intersectBy", "*LinqObj", args[0].Type()))
	}

	selector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "intersectBy", "*Function", args[1].Type()))
	}

	scop := NewScope(scope)

	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			next := lq.Query.Iterate()
			next2 := lq2.Query.Iterate()

			set := make(map[Object]bool)
			for item, ok := next2(); ok.Bool; item, ok = next2() {
				scop.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, item)
				s := Eval(selector.Literal.Body, scop)
				if obj, ok1 := s.(*ReturnValue); ok1 {
					s = obj.Value
				}
				set[s] = true
			}

			return func() (item Object, ok *Boolean) {
				for item, ok = next(); ok.Bool; item, ok = next() {
					scop.Set(selector.Literal.Parameters[0].(*ast.Identifier).Value, item)
					s := Eval(selector.Literal.Body, scop)
					if obj, ok1 := s.(*ReturnValue); ok1 {
						s = obj.Value
					}

					for k, _ := range set {
						if reflect.DeepEqual(k, s) {
							set[s] = false
							break
						}
					}
					if _, has := set[s]; has {
						delete(set, s)
						return
					}
				}

				return
			}
		},
	}}
}

// Aggregate applies an accumulator function over a sequence.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The first element of source is used as the initial aggregate value. The
// result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (lq *LinqObj) Aggregate(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "aggregate", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)

	next := lq.Query.Iterate()

	result, any := next()
	if !any.Bool {
		return NIL
	}

	for current, ok := next(); ok.Bool; current, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, result)
		scop.Set(fn.Literal.Parameters[1].(*ast.Identifier).Value, current)
		result = Eval(fn.Literal.Body, scop)
		if obj, ok1 := result.(*ReturnValue); ok1 {
			result = obj.Value
		}
	}

	return result
}

// AggregateWithSeed applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source
// except the first one. Each time f() is called, Aggregate passes both the
// element from the sequence and an aggregated value (as the first argument to
// f()). The value of the seed parameter is used as the initial aggregate value.
// The result of f() replaces the previous aggregated value.
//
// Aggregate returns the final result of f().
func (lq *LinqObj) AggregateWithSeed(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	//this test is not needed, but here for completeness
	seed, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "aggregateWithSeed", "Object", args[0].Type()))
	}

	fn, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "aggregateWithSeed", "*Function", args[1].Type()))
	}

	scop := NewScope(scope)

	next := lq.Query.Iterate()
	result := seed

	for current, ok := next(); ok.Bool; current, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, result)
		scop.Set(fn.Literal.Parameters[1].(*ast.Identifier).Value, current)
		result = Eval(fn.Literal.Body, scop)
		if obj, ok1 := result.(*ReturnValue); ok1 {
			result = obj.Value
		}
	}

	return result
}

// AggregateWithSeedBy applies an accumulator function over a sequence. The
// specified seed value is used as the initial accumulator value, and the
// specified function is used to select the result value.
//
// Aggregate method makes it simple to perform a calculation over a sequence of
// values. This method works by calling f() one time for each element in source.
// Each time func is called, Aggregate passes both the element from the sequence
// and an aggregated value (as the first argument to func). The value of the
// seed parameter is used as the initial aggregate value. The result of func
// replaces the previous aggregated value.
//
// The final result of func is passed to resultSelector to obtain the final
// result of Aggregate.
func (lq *LinqObj) AggregateWithSeedBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	seed, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "aggregateWithSeedBy", "Object", args[0].Type()))
	}

	fn, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "aggregateWithSeedBy", "*Function", args[1].Type()))
	}

	resultSelector, ok := args[2].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "aggregateWithSeedBy", "*Function", args[2].Type()))
	}

	scop := NewScope(scope)

	next := lq.Query.Iterate()
	result := seed

	for current, ok := next(); ok.Bool; current, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, result)
		scop.Set(fn.Literal.Parameters[1].(*ast.Identifier).Value, current)
		result = Eval(fn.Literal.Body, scop)
		if obj, ok1 := result.(*ReturnValue); ok1 {
			result = obj.Value
		}
	}

	scop.Set(resultSelector.Literal.Parameters[0].(*ast.Identifier).Value, result)
	result = Eval(resultSelector.Literal.Body, scop)
	if obj, ok1 := result.(*ReturnValue); ok1 {
		result = obj.Value
	}

	return result
}

// All determines whether all elements of a collection satisfy a condition.
func (lq *LinqObj) All(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "all", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)

	next := lq.Query.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, item)
		cond := Eval(fn.Literal.Body, scop)
		if obj, ok1 := cond.(*ReturnValue); ok1 {
			cond = obj.Value
		}
		if !IsTrue(cond) {
			return FALSE
		}
	}

	return TRUE
}

// Any determines whether any element of a collection exists.
func (lq *LinqObj) Any(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	_, ok := lq.Query.Iterate()()
	return ok
}

func (lq *LinqObj) AnyWith(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "anyWith", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)

	next := lq.Query.Iterate()

	for item, ok := next(); ok.Bool; item, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, item)
		cond := Eval(fn.Literal.Body, scop)
		if obj, ok1 := cond.(*ReturnValue); ok1 {
			cond = obj.Value
		}
		if IsTrue(cond) {
			return TRUE
		}
	}

	return FALSE
}

// Contains determines whether a collection contains a specified element.
func (lq *LinqObj) Contains(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	next := lq.Query.Iterate()

	for item, ok := next(); ok.Bool; item, ok = next() {
		//using `if item == args[0]` is not correct, we need DeepEqual
		if reflect.DeepEqual(item, args[0]) {
			return TRUE
		}
	}

	return FALSE
}

// Count returns the number of elements in a collection.
func (lq *LinqObj) Count(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()

	var r int64 = 0
	for _, ok := next(); ok.Bool; _, ok = next() {
		r++
	}

	return NewInteger(r)
}

// CountWith returns a number that represents how many elements in the specified
// collection satisfy a condition.
func (lq *LinqObj) CountWith(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "countWith", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)

	var r int64 = 0
	next := lq.Query.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, item)
		cond := Eval(fn.Literal.Body, scop)
		if obj, ok1 := cond.(*ReturnValue); ok1 {
			cond = obj.Value
		}
		if IsTrue(cond) {
			r++
		}
	}

	return NewInteger(r)
}

// SequenceEqual determines whether two collections are equalp.
func (lq *LinqObj) SequenceEqual(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	lq2, ok := args[0].(*LinqObj)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sequenceEqual", "*LinqObj", args[0].Type()))
	}

	next := lq.Query.Iterate()
	next2 := lq2.Query.Iterate()

	for item, ok := next(); ok.Bool; item, ok = next() {
		item2, ok2 := next2()
		if !ok2.Bool || !reflect.DeepEqual(item, item2) { //Note: Here we should not use `item == item2`
			return FALSE
		}
	}

	_, ok2 := next2()
	if ok2.Bool {
		return FALSE
	}
	return TRUE
}

// Single returns the only element of a collection, and nil if there is not
// exactly one element in the collection.
func (lq *LinqObj) Single(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()
	item, ok := next()
	if !ok.Bool {
		return NIL
	}

	_, ok = next()
	if ok.Bool {
		return NIL
	}

	return item
}

// SingleWith returns the only element of a collection that satisfies a
// specified condition, and nil if more than one such element exists.
func (lq *LinqObj) SingleWith(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fn, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "singleWith", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)

	var r Object = NIL
	next := lq.Query.Iterate()
	found := false
	for item, ok := next(); ok.Bool; item, ok = next() {
		scop.Set(fn.Literal.Parameters[0].(*ast.Identifier).Value, item)
		cond := Eval(fn.Literal.Body, scop)
		if obj, ok1 := cond.(*ReturnValue); ok1 {
			cond = obj.Value
		}
		if IsTrue(cond) {
			if found {
				return NIL
			}
			found = true
			r = item
		}
	}

	return r
}

// SumInts computes the sum of a collection of numeric values.
//
// Method returns zero if collection contains no elements.
func (lq *LinqObj) SumInts(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()
	item, ok := next()
	if !ok.Bool {
		return NewInteger(0)
	}

	var r int64 = item.(*Integer).Int64

	for item, ok = next(); ok.Bool; item, ok = next() {
		r += item.(*Integer).Int64
	}

	return NewInteger(r)
}

// SumFloats computes the sum of a collection of numeric values.
//
// Method returns zero if collection contains no elements.
func (lq *LinqObj) SumFloats(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()
	item, ok := next()
	if !ok.Bool {
		return NewFloat(0)
	}

	var r float64 = item.(*Float).Float64

	for item, ok = next(); ok.Bool; item, ok = next() {
		r += item.(*Float).Float64
	}

	return NewFloat(r)
}

// Min returns the minimum value in a collection of values.
func (lq *LinqObj) Min(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()
	item, ok := next()
	if !ok.Bool {
		return NIL
	}

	compare := getComparer(item)
	r := item

	for item, ok := next(); ok.Bool; item, ok = next() {
		if compare(item, r) < 0 {
			r = item
		}
	}

	return r
}

// Max returns the maximum value in a collection of values.
func (lq *LinqObj) Max(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()
	item, ok := next()
	if !ok.Bool {
		return NIL
	}

	compare := getComparer(item)
	r := item

	for item, ok := next(); ok.Bool; item, ok = next() {
		if compare(item, r) > 0 {
			r = item
		}
	}

	return r
}

// Average computes the average of a collection of numeric values.
func (lq *LinqObj) Average(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	next := lq.Query.Iterate()
	item, ok := next()
	if !ok.Bool {
		return NewFloat(math.NaN())
	}

	var r float64
	n := 1
	switch item.(type) {
	case *Integer:
		var sum int64 = item.(*Integer).Int64

		for item, ok = next(); ok.Bool; item, ok = next() {
			sum += item.(*Integer).Int64
			n++
		}

		r = float64(sum)
	case *Float:
		var r float64 = item.(*Float).Float64

		for item, ok = next(); ok.Bool; item, ok = next() {
			r += item.(*Float).Float64
			n++
		}
	}

	return NewFloat(r / float64(n))
}

// OrderBy sorts the elements of a collection in ascending order. Elements are
// sorted according to a key.
func (lq *LinqObj) OrderBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "orderBy", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)
	return &LinqObj{Query: lq.Query, OrderedQuery: OrderedQuery{
		orders:   []order{{selector: selector, scope: scop}},
		original: lq.Query,
		Query: Query{
			Iterate: func() Iterator {
				items := lq.Query.sort([]order{{selector: selector, scope: scop}})
				len := len(items)
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}}
}

// OrderByDescending sorts the elements of a collection in descending order.
// Elements are sorted according to a key.
func (lq *LinqObj) OrderByDescending(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "orderByDescending", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)
	return &LinqObj{Query: lq.Query, OrderedQuery: OrderedQuery{
		orders:   []order{{selector: selector, scope: scop, desc: true}},
		original: lq.Query,
		Query: Query{
			Iterate: func() Iterator {
				items := lq.Query.sort([]order{{selector: selector, scope: scope, desc: true}})
				len := len(items)
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}}
}

// ThenBy performs a subsequent ordering of the elements in a collection in
// ascending order. This method enables you to specify multiple sort criteria by
// applying any number of ThenBy or ThenByDescending methods.
func (lq *LinqObj) ThenBy(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "thenBy", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)
	return &LinqObj{Query: lq.Query, OrderedQuery: OrderedQuery{
		orders:   append(lq.OrderedQuery.orders, order{selector: selector, scope: scop}),
		original: lq.OrderedQuery.original,
		Query: Query{
			Iterate: func() Iterator {
				items := lq.OrderedQuery.original.sort(append(lq.OrderedQuery.orders, order{selector: selector, scope: scop}))
				len := len(items)
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}}
}

// ThenByDescending performs a subsequent ordering of the elements in a
// collection in descending order. This method enables you to specify multiple
// sort criteria by applying any number of ThenBy or ThenByDescending methods.
func (lq *LinqObj) ThenByDescending(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	selector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "thenBy", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)
	return &LinqObj{Query: lq.Query, OrderedQuery: OrderedQuery{
		orders:   append(lq.OrderedQuery.orders, order{selector: selector, scope: scop, desc: true}),
		original: lq.OrderedQuery.original,
		Query: Query{
			Iterate: func() Iterator {
				items := lq.OrderedQuery.original.sort(append(lq.OrderedQuery.orders, order{selector: selector, scope: scop, desc: true}))
				len := len(items)
				index := 0

				return func() (item Object, ok *Boolean) {
					ok = &Boolean{Valid: true}
					ok.Bool = index < len
					if ok.Bool {
						item = items[index]
						index++
					}

					return
				}
			},
		},
	}}
}

// ToSlice iterates over a collection and saves the results in the slice pointed
// by v. It overwrites the existing slice, starting from index 0.
//
// If the slice pointed by v has sufficient capacity, v will be pointed to a
// resliced slice. If it does not, a new underlying array will be allocated and
// v will point to it.
func (lq *LinqObj) ToSlice(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	res := &Array{}
	next := lq.Query.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {
		res.Members = append(res.Members, item)
	}

	return res
}

func (lq *LinqObj) ToOrderedSlice(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	res := &Array{}
	next := lq.OrderedQuery.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {
		res.Members = append(res.Members, item)
	}

	return res
}

// Sort returns a new query by sorting elements with provided less function in
// ascending order. The comparer function should return true if the parameter i
// is less than j. While this method is uglier than chaining OrderBy,
// OrderByDescending, ThenBy and ThenByDescending methods, it's performance is
// much better.
func (lq *LinqObj) Sort(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	less, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "sort", "*Function", args[0].Type()))
	}

	scop := NewScope(scope)
	return &LinqObj{Query: Query{
		Iterate: func() Iterator {
			items := lq.Query.lessSort(scop, less)
			len := len(items)
			index := 0

			return func() (item Object, ok *Boolean) {
				ok = &Boolean{Valid: true}
				ok.Bool = index < len
				if ok.Bool {
					item = items[index]
					index++
				}

				return
			}
		},
	}}
}

func (lq *LinqObj) ToMap(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	keySelector, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "toMap", "*Function", args[0].Type()))
	}

	valueSelector, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "toMap", "*Function", args[1].Type()))
	}

	scop := NewScope(scope)
	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	next := lq.Query.Iterate()
	for item, ok := next(); ok.Bool; item, ok = next() {

		scop.Set(keySelector.Literal.Parameters[0].(*ast.Identifier).Value, item)
		key := Eval(keySelector.Literal.Body, scop)
		if obj, ok1 := key.(*ReturnValue); ok1 {
			key = obj.Value
		}

		scop.Set(valueSelector.Literal.Parameters[0].(*ast.Identifier).Value, item)
		value := Eval(valueSelector.Literal.Body, scop)
		if obj, ok1 := value.(*ReturnValue); ok1 {
			value = obj.Value
		}

		hash.Push(line, key, value)
	}

	return hash
}

func getComparer(data Object) comparer {
	switch data.(type) {
	case *Integer:
		return func(x, y Object) int {
			a, b := x.(*Integer), y.(*Integer)
			switch {
			case a.Int64 > b.Int64:
				return 1
			case b.Int64 > a.Int64:
				return -1
			default:
				return 0
			}
		}
	case *Float:
		return func(x, y Object) int {
			a, b := x.(*Float), y.(*Float)
			switch {
			case a.Float64 > b.Float64:
				return 1
			case b.Float64 > a.Float64:
				return -1
			default:
				return 0
			}
		}
	case *String:
		return func(x, y Object) int {
			a, b := x.(*String), y.(*String)
			switch {
			case a.String > b.String:
				return 1
			case b.String > a.String:
				return -1
			default:
				return 0
			}
		}
	case *Boolean:
		return func(x, y Object) int {
			a, b := x.(*Boolean), y.(*Boolean)
			switch {
			case a.Bool == b.Bool:
				return 0
			case a.Bool:
				return 1
			default:
				return -1
			}
		}
	default:
		panic("Comparer not supported")
	}
}
