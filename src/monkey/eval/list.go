package eval

import (
	"container/list"
)

const (
	LIST_OBJ     = "LIST_OBJ"
	LISTELEM_OBJ = "LISTELEM_OBJ"
)

//List object
type ListObject struct {
	List *list.List
}

func (l *ListObject) Inspect() string  { return "<list>" }
func (l *ListObject) Type() ObjectType { return LIST_OBJ }
func (l *ListObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "back":
		return l.Back(line, args...)
	case "front":
		return l.Front(line, args...)
	case "init":
		return l.Init(line, args...)
	case "insertAfter":
		return l.InsertAfter(line, args...)
	case "insertBefore":
		return l.InsertBefore(line, args...)
	case "len":
		return l.Len(line, args...)
	case "moveToBack":
		return l.MoveToBack(line, args...)
	case "moveToFront":
		return l.MoveToFront(line, args...)
	case "pushBack":
		return l.PushBack(line, args...)
	case "pushBackList":
		return l.PushBackList(line, args...)
	case "pushFront":
		return l.PushFront(line, args...)
	case "pushFrontList":
		return l.PushFrontList(line, args...)
	case "remove":
		return l.Remove(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, l.Type()))
	}
}

func (l *ListObject) Back(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &ListElemObject{Elem: l.List.Back()}
}

func (l *ListObject) Front(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &ListElemObject{Elem: l.List.Front()}
}

func (l *ListObject) Init(line string, args ...Object) Object {
	l.List = l.List.Init()
	return l
}

func (l *ListObject) InsertAfter(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	value, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "insertAfter", "Object", args[0].Type()))
	}

	mark, ok := args[1].(*ListElemObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "insertAfter", "*ListElemObject", args[1].Type()))
	}

	return &ListElemObject{Elem: l.List.InsertAfter(value, mark.Elem)}
}

func (l *ListObject) InsertBefore(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	value, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "insertBefore", "Object", args[0].Type()))
	}

	mark, ok := args[1].(*ListElemObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "insertBefore", "*ListElemObject", args[1].Type()))
	}

	return &ListElemObject{Elem: l.List.InsertBefore(value, mark.Elem)}
}

func (l *ListObject) Len(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewInteger(int64(l.List.Len()))
}

func (l *ListObject) MoveToBack(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	e, ok := args[0].(*ListElemObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "moveToBack", "*ListElemObject", args[0].Type()))
	}

	l.List.MoveToBack(e.Elem)

	return NIL
}

func (l *ListObject) MoveToFront(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	e, ok := args[0].(*ListElemObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "moveToFront", "*ListElemObject", args[0].Type()))
	}

	l.List.MoveToFront(e.Elem)

	return NIL
}

func (l *ListObject) PushBack(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	value, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "pushBack", "Object", args[0].Type()))
	}

	return &ListElemObject{Elem: l.List.PushBack(value)}
}

func (l *ListObject) PushBackList(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	value, ok := args[0].(*ListObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "pushBackList", "*ListObject", args[0].Type()))
	}

	l.List.PushBackList(value.List)

	return NIL
}

func (l *ListObject) PushFront(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	value, ok := args[0].(Object)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "pushFront", "Object", args[0].Type()))
	}

	return &ListElemObject{Elem: l.List.PushFront(value)}
}

func (l *ListObject) PushFrontList(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	value, ok := args[0].(*ListObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "pushFrontList", "*ListObject", args[0].Type()))
	}

	l.List.PushFrontList(value.List)

	return NIL
}

func (l *ListObject) Remove(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	e, ok := args[0].(*ListElemObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "remove", "*ListElemObject", args[0].Type()))
	}

	r := l.List.Remove(e.Elem)
	return r.(Object)
}

//List element object
type ListElemObject struct {
	Elem *list.Element
}

func (e *ListElemObject) Inspect() string {
	obj, ok := e.Elem.Value.(Object)
	if ok {
		return obj.Inspect()
		//fmt.Sprintf("%v", e.Elem.Value)
	}
	return "Unknown Value"
}
func (e *ListElemObject) Type() ObjectType { return LISTELEM_OBJ }
func (e *ListElemObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "next":
		return e.Next(line, args...)
	case "prev":
		return e.Prev(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, e.Type()))
	}
}

func (e *ListElemObject) Next(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	elem := e.Elem.Next()
	if elem == nil {
		return NIL
	}
	return &ListElemObject{Elem: elem}
}

func (e *ListElemObject) Prev(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	elem := e.Elem.Prev()
	if elem == nil {
		return NIL
	}
	return &ListElemObject{Elem: elem}
}
