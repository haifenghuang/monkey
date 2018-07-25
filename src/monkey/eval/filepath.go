package eval

import (
	_ "fmt"
	"monkey/ast"
	"os"
	"path/filepath"
)

func NewFilePathObj() *FilePathObj {
	ret := &FilePathObj{}
	SetGlobalObj(filepath_name, ret)

	SetGlobalObj(filepath_name+".SKIP_DIR", NewInteger(SKIP_DIR))
	SetGlobalObj(filepath_name+".SEPARATOR", NewString(string(filepath.Separator)))
	SetGlobalObj(filepath_name+".LISTSEPARATOR", NewString(string(filepath.ListSeparator)))

	return ret
}

const (
	SKIP_DIR = 1
	FILEPATH_OBJ = "FILEPATH_OBJ"
	filepath_name = "filepath"
)

type FilePathObj struct{}

func (f *FilePathObj) Inspect() string  { return "<" + filepath_name + ">" }
func (f *FilePathObj) Type() ObjectType { return FILEPATH_OBJ }

func (f *FilePathObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "abs":
		return f.Abs(line, args...)
	case "base":
		return f.Base(line, args...)
	case "clean":
		return f.Clean(line, args...)
	case "dir":
		return f.Dir(line, args...)
	case "evalSymlinks":
		return f.EvalSymlinks(line, args...)
	case "ext":
		return f.Ext(line, args...)
	case "fromSlash":
		return f.FromSlash(line, args...)
	case "glob":
		return f.Glob(line, args...)
	case "hasPrefix":
		return f.HasPrefix(line, args...)
	case "isAbs":
		return f.IsAbs(line, args...)
	case "join":
		return f.Join(line, args...)
	case "match":
		return f.Match(line, args...)
	case "rel":
		return f.Rel(line, args...)
	case "split":
		return f.Split(line, args...)
	case "splitList":
		return f.SplitList(line, args...)
	case "toSlash":
		return f.ToSlash(line, args...)
	case "volumeName":
		return f.VolumeName(line, args...)
	case "walk":
		return f.Walk(scope, line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, f.Type()))
}

func (f *FilePathObj) Abs(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "abs", "*String", args[0].Type()))
	}

	ret, err := filepath.Abs(path.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(ret)
}

func (f *FilePathObj) Base(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "base", "*String", args[0].Type()))
	}

	ret := filepath.Base(path.String)
	return NewString(ret)
}

func (f *FilePathObj) Clean(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "clean", "*String", args[0].Type()))
	}

	ret := filepath.Clean(path.String)
	return NewString(ret)
}

func (f *FilePathObj) Dir(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "dir", "*String", args[0].Type()))
	}

	ret := filepath.Dir(path.String)
	return NewString(ret)
}

func (f *FilePathObj) EvalSymlinks(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "evalSymlinks", "*String", args[0].Type()))
	}

	ret, err := filepath.EvalSymlinks(path.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(ret)
}

func (f *FilePathObj) Ext(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "ext", "*String", args[0].Type()))
	}

	ret := filepath.Ext(path.String)
	return NewString(ret)
}

func (f *FilePathObj) FromSlash(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "fromSlash", "*String", args[0].Type()))
	}

	ret := filepath.FromSlash(path.String)
	return NewString(ret)
}

func (f *FilePathObj) Glob(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	pattern, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "glob", "*String", args[0].Type()))
	}

	ret, err := filepath.Glob(pattern.String)
	if err != nil {
		return NewNil(err.Error())
	}

	arr := &Array{}
	for _, v := range ret {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (f *FilePathObj) HasPrefix(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	p, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "hasPrefix", "*String", args[0].Type()))
	}

	prefix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "hasPrefix", "*String", args[1].Type()))
	}

	ret := filepath.HasPrefix(p.String, prefix.String)
	if ret {
		return TRUE
	}
	return FALSE
}

func (f *FilePathObj) IsAbs(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "isAbs", "*String", args[0].Type()))
	}

	ret := filepath.IsAbs(path.String)
	if ret {
		return TRUE
	}
	return FALSE
}

func (f *FilePathObj) Join(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">=1", len(args)))
	}

	joinList := []string{}
	for idx, v := range args {
		path, ok := v.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "all", "join", "*String", args[idx].Type()))
		}

		joinList = append(joinList, path.String)
	}

	ret := filepath.Join(joinList...)
	return NewString(ret)
}

func (f *FilePathObj) Match(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	pattern, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "match", "*String", args[0].Type()))
	}

	name, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "match", "*String", args[1].Type()))
	}

	matched, err := filepath.Match(pattern.String, name.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	if matched {
		return TRUE
	}
	return FALSE
}

func (f *FilePathObj) Rel(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	basepath, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "rel", "*String", args[0].Type()))
	}

	targpath, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "rel", "*String", args[1].Type()))
	}

	relative, err := filepath.Rel(basepath.String, targpath.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(relative)
}

func (f *FilePathObj) Split(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "split", "*String", args[0].Type()))
	}

	dir, file := filepath.Split(path.String)

	dirVal := NewString(dir)
	fileVal := NewString(file)

	dirKey := NewString("dir")
	fileKey := NewString("file")

	hash := NewHash()
	hash.Push(line, dirKey, dirVal)
	hash.Push(line, fileKey, fileVal)
	//hash.Pairs[dirKey.HashKey()] = HashPair{Key: dirKey, Value: dirVal}
	//hash.Pairs[fileKey.HashKey()] = HashPair{Key: fileKey, Value: fileVal}

	return hash
}

func (f *FilePathObj) SplitList(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "splitList", "*String", args[0].Type()))
	}

	ret := filepath.SplitList(path.String)

	arr := &Array{}
	for _, v := range ret {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (f *FilePathObj) ToSlash(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "toSlash", "*String", args[0].Type()))
	}

	ret := filepath.ToSlash(path.String)
	return NewString(ret)
}

func (f *FilePathObj) VolumeName(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "volumeName", "*String", args[0].Type()))
	}

	ret := filepath.VolumeName(path.String)
	return NewString(ret)
}

func (f *FilePathObj) Walk(scope *Scope, line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	root, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "walk", "*String", args[0].Type()))
	}

	block, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "walk", "*Function", args[1].Type()))
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 2 {
		panic(NewError(line, FUNCCALLBACKERROR, 2, paramCount))
	}

	err := filepath.Walk(root.String, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return walkFunc(scope, block, path, info)
	})
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func walkFunc(scope *Scope, f *Function, path string, info os.FileInfo) error {
	s := NewScope(scope)

	//Store to `Scope`，so below `Eval() could use them
	s.Set(f.Literal.Parameters[0].(*ast.Identifier).Value, NewString(path))
	s.Set(f.Literal.Parameters[1].(*ast.Identifier).Value, &FileInfoObj{Info: info})
	r := Eval(f.Literal.Body, s)
	if obj, ok := r.(*ReturnValue); ok {
		r = obj.Value
	}

	//check for return value, is it a skipDir?
	iObj, ok := r.(*Integer)
	if ok {
		if iObj.Int64 == SKIP_DIR {
			return filepath.SkipDir
		}
	}
	return nil
}
