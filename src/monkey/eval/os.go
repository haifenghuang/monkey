package eval

import (
	_ "fmt"
	"io"
	"monkey/ast"
	"os"
	"os/exec"
	"strings"
)

const (
	OS_OBJ = "OS_OBJ"
	os_name = "os"
)

type Os struct{}

func NewOsObj() Object {
	ret := &Os{}
	SetGlobalObj(os_name, ret)

	SetGlobalObj(os_name+".O_APPEND", NewInteger(int64(os.O_APPEND)))
	SetGlobalObj(os_name+".O_CREATE", NewInteger(int64(os.O_CREATE)))
	SetGlobalObj(os_name+".O_EXCL", NewInteger(int64(os.O_EXCL)))
	SetGlobalObj(os_name+".O_RDONLY", NewInteger(int64(os.O_RDONLY)))
	SetGlobalObj(os_name+".O_RDWR", NewInteger(int64(os.O_RDWR)))
	SetGlobalObj(os_name+".O_SYNC", NewInteger(int64(os.O_SYNC)))
	SetGlobalObj(os_name+".O_TRUNC", NewInteger(int64(os.O_TRUNC)))
	SetGlobalObj(os_name+".O_WRONLY", NewInteger(int64(os.O_WRONLY)))

	SetGlobalObj(os_name+".SEEK_CUR", NewInteger(int64(os.SEEK_CUR)))
	SetGlobalObj(os_name+".SEEK_END", NewInteger(int64(os.SEEK_END)))
	SetGlobalObj(os_name+".SEEK_SET", NewInteger(int64(os.SEEK_SET)))

	SetGlobalObj(os_name+".PATH_SEPARATOR", NewString(string(os.PathSeparator)))
	SetGlobalObj(os_name+".PATH_LISTSEPARATOR", NewString(string(os.PathListSeparator)))

	SetGlobalObj(os_name+".MODE_DIR", NewInteger(int64(os.ModeDir)))
	SetGlobalObj(os_name+".MODE_APPEND", NewInteger(int64(os.ModeAppend)))
	SetGlobalObj(os_name+".MODE_EXCLUSIVE", NewInteger(int64(os.ModeExclusive)))
	SetGlobalObj(os_name+".MODE_TEMPORARY", NewInteger(int64(os.ModeTemporary)))
	SetGlobalObj(os_name+".MODE_SYMLINK", NewInteger(int64(os.ModeSymlink)))
	SetGlobalObj(os_name+".MODE_DEVICE", NewInteger(int64(os.ModeDevice)))
	SetGlobalObj(os_name+".MODE_NAMEDPIPE", NewInteger(int64(os.ModeNamedPipe)))
	SetGlobalObj(os_name+".MODE_SOCKET", NewInteger(int64(os.ModeSocket)))
	SetGlobalObj(os_name+".MODE_SETUID", NewInteger(int64(os.ModeSetuid)))
	SetGlobalObj(os_name+".MODE_SETGID", NewInteger(int64(os.ModeSetgid)))
	SetGlobalObj(os_name+".MODE_CHARDEVICE", NewInteger(int64(os.ModeCharDevice)))
	SetGlobalObj(os_name+".MODE_STICKY", NewInteger(int64(os.ModeSticky)))
	SetGlobalObj(os_name+".MODE_TYPE", NewInteger(int64(os.ModeType)))
	SetGlobalObj(os_name+".MODE_PERM", NewInteger(int64(os.ModePerm)))

	return ret
}

func (o *Os) Inspect() string  { return "<" + os_name + ">" }
func (o *Os) Type() ObjectType { return OS_OBJ }

func (o *Os) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "args":
		return o.Args(line, args...)
	case "getenv":
		return o.Getenv(line, args...)
	case "setenv":
		return o.Setenv(line, args...)
	case "clearenv":
		return o.Clearenv(line, args...)
	case "unsetenv":
		return o.Unsetenv(line, args...)
	case "getwd":
		return o.Getwd(line, args...)
	case "chdir":
		return o.Chdir(line, args...)
	case "chmod":
		return o.Chmod(line, args...)
	case "chown":
		return o.Chown(line, args...)
	case "mkdir":
		return o.Mkdir(line, args...)
	case "mkdirAll":
		return o.MkdirAll(line, args...)
	case "truncate":
		return o.Truncate(line, args...)
	case "remove":
		return o.Remove(line, args...)
	case "removeAll":
		return o.RemoveAll(line, args...)
	case "stat":
		return o.Stat(line, args...)
	case "runCmd":
		return o.RunCmd(line, args...)
	case "environ":
		return o.Environ(line, args...)
	case "exit":
		return o.Exit(line, args...)
	case "expand":
		return o.Expand(scope, line, args...)
	case "expandEnv":
		return o.ExpandEnv(line, args...)
	case "hostname":
		return o.Hostname(line, args...)
	case "link":
		return o.Link(line, args...)
	case "readlink":
		return o.Readlink(line, args...)
	case "rename":
		return o.Rename(line, args...)
	case "tempDir":
		return o.TempDir(line, args...)
	case "copyFile":
		return o.CopyFile(line, args...)
	case "isExist":
		return o.IsExist(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, o.Type()))
}

func (o *Os) Args(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	arr := &Array{}
	for _, v := range os.Args[1:] {
		arr.Members = append(arr.Members, NewString(v))
	}

	return arr
}

func (o *Os) Getenv(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "getenv", "*String", args[0].Type()))
	}

	ret := os.Getenv(key.String)
	return NewString(ret)
}

func (o *Os) Setenv(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setenv", "*String", args[0].Type()))
	}

	value, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "setenv", "*String", args[1].Type()))
	}

	err := os.Setenv(key.String, value.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) Clearenv(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	os.Clearenv()
	return NIL
}

func (o *Os) Getwd(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret, err := os.Getwd()
	if err != nil {
		return NewNil(err.Error())
	}
	return NewString(ret)
}

func (o *Os) Chdir(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	newDir, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "chdir", "*String", args[0].Type()))
	}

	err := os.Chdir(newDir.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) Chmod(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "chmod", "*String", args[0].Type()))
	}

	mode, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "chmod", "*Integer", args[1].Type()))
	}

	err := os.Chmod(name.String, os.FileMode(mode.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (o *Os) Chown(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "chown", "*String", args[0].Type()))
	}

	uid, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "chown", "*Integer", args[1].Type()))
	}

	gid, ok := args[2].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "chown", "*Integer", args[2].Type()))
	}

	err := os.Chown(name.String, int(uid.Int64), int(gid.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) Mkdir(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "mkdir", "*String", args[0].Type()))
	}

	perm, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "mkdir", "*Integer", args[1].Type()))
	}

	err := os.Mkdir(name.String, os.FileMode(perm.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) MkdirAll(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "mkdirAll", "*String", args[0].Type()))
	}

	perm, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "mkdirAll", "*Integer", args[1].Type()))
	}

	err := os.MkdirAll(path.String, os.FileMode(perm.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) Truncate(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "truncate", "*String", args[0].Type()))
	}

	size, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "truncate", "*Integer", args[1].Type()))
	}

	err := os.Truncate(name.String, size.Int64)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) Remove(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "remove", "*String", args[0].Type()))
	}

	err := os.Remove(name.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) RemoveAll(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	path, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "removeAll", "*String", args[0].Type()))
	}

	err := os.RemoveAll(path.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//The result is a hash object:
//{name: bb.my, size: 240, mode: 438, isDir: false}
func (o *Os) Stat(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "stat", "*String", args[0].Type()))
	}

	fi, err := os.Stat(name.String)
	if err != nil {
		return NewNil(err.Error())
	}

	nameVal := NewString(fi.Name())
	sizeVal := NewInteger(fi.Size())
	modeVal := NewInteger(int64(fi.Mode()))

	var isDirVal *Boolean
	if fi.IsDir() {
		isDirVal = TRUE
	} else {
		isDirVal = FALSE
	}

	nameKey := NewString("name")
	sizeKey := NewString("size")
	modeKey := NewString("mode")
	isDirKey := NewString("isDir")

	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	hash.Pairs[nameKey.HashKey()] = HashPair{Key: nameKey, Value: nameVal}
	hash.Pairs[sizeKey.HashKey()] = HashPair{Key: sizeKey, Value: sizeVal}
	hash.Pairs[modeKey.HashKey()] = HashPair{Key: modeKey, Value: modeVal}
	hash.Pairs[isDirKey.HashKey()] = HashPair{Key: isDirKey, Value: isDirVal}

	return hash
}

func (o *Os) RunCmd(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	cmd, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "runCmd", "*String", args[0].Type()))
	}

	var params []string
	arr := args[1:]
	for idx, arg := range arr {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "The remaining", "runCmd", "*String", arr[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	var cmdOut []byte
	var cmdErr error
	if cmdOut, cmdErr = exec.Command(cmd.String, params...).Output(); cmdErr != nil {
		return NewNil(cmdErr.Error())
	}

	return NewString(string(cmdOut))
}

func (o *Os) Unsetenv(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	key, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "unsetenv", "*String", args[0].Type()))
	}

	err := os.Unsetenv(key.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (o *Os) Environ(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	envs := os.Environ()

	hash := &Hash{Pairs: make(map[HashKey]HashPair)}
	for _, env := range envs {
		kvPair := strings.Split(env, "=")
		key := kvPair[0]
		val := kvPair[1]

		keyKey := NewString(key)
		valKey := NewString(val)
		hash.Pairs[keyKey.HashKey()] = HashPair{Key: keyKey, Value: valKey}
	}

	return hash
}

func (o *Os) Exit(line string, args ...Object) Object {
	if len(args) != 0 && len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
	}

	if len(args) == 0 {
		os.Exit(0)
		return NIL
	}

	code, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "exit", "*Integer", args[0].Type()))
	}

	os.Exit(int(code.Int64))

	return NIL
}

func (o *Os) Expand(scope *Scope, line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	s, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "expand", "*String", args[0].Type()))
	}

	block, ok := args[1].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "expand", "*Function", args[1].Type()))
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 1 {
		panic(NewError(line, FUNCCALLBACKERROR, 1, paramCount))
	}

	ret := os.Expand(s.String, func(s string) string {
		return evalExpandMappingFunc(scope, block, s)
	})

	return NewString(ret)
}

func evalExpandMappingFunc(scope *Scope, f *Function, s string) string {
	sc := NewScope(scope)
	sc.Set(f.Literal.Parameters[0].(*ast.Identifier).Value, NewString(s))
	r := Eval(f.Literal.Body, sc)
	if obj, ok := r.(*ReturnValue); ok {
		r = obj.Value
	}

	ret, ok := r.(*String)
	if !ok {
		return ""
	}
	return ret.String
}

func (o *Os) ExpandEnv(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	s, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "expandEnv", "*String", args[0].Type()))
	}

	ret := os.ExpandEnv(s.String)

	return NewString(ret)
}

func (o *Os) Hostname(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	ret, err := os.Hostname()
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(ret)
}

func (o *Os) Link(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	oldName, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "link", "*String", args[0].Type()))
	}

	newName, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "link", "*String", args[1].Type()))
	}

	err := os.Link(oldName.String, newName.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (o *Os) Readlink(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	name, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "readlink", "*String", args[0].Type()))
	}

	ret, err := os.Readlink(name.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(ret)
}

func (o *Os) Rename(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	oldPath, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "rename", "*String", args[0].Type()))
	}

	newPath, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "rename", "*String", args[1].Type()))
	}

	err := os.Rename(oldPath.String, newPath.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (o *Os) TempDir(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(os.TempDir())
}

func (o *Os) CopyFile(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	dst, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "copyFile", "*String", args[0].Type()))
	}

	src, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "copyFile", "*String", args[1].Type()))
	}

	in, err := os.Open(src.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	defer in.Close()

	out, err := os.Create(dst.String)
	if err != nil {
		return NewFalseObj(err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (o *Os) IsExist(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fp, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "isExist", "*String", args[0].Type()))
	}

	_, err := os.Stat(fp.String)
	if err == nil || os.IsExist(err) {
		return TRUE
	}
	return NewFalseObj(err.Error())
}

const FILEINFO_OBJ = "FILEINFO_OBJ"
//FileInfo Object
type FileInfoObj struct {
	Info os.FileInfo
}

func (fi *FileInfoObj) Inspect() string  { return FILEINFO_OBJ }
func (fi *FileInfoObj) Type() ObjectType { return FILEINFO_OBJ }

func (fi *FileInfoObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "name":
		return fi.Name(line, args...)
	case "size":
		return fi.Size(line, args...)
	case "mode":
		return fi.Mode(line, args...)
	case "modTime":
		return fi.ModTime(line, args...)
	case "isDir":
		return fi.IsDir(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, fi.Type()))
}

func (fi *FileInfoObj) Name(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(fi.Info.Name())
}

func (fi *FileInfoObj) Size(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewInteger(fi.Info.Size())
}

func (fi *FileInfoObj) Mode(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewInteger(int64(fi.Info.Mode()))
}

func (fi *FileInfoObj) ModTime(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return &TimeObj{Tm: fi.Info.ModTime(), Valid: true}
}

func (fi *FileInfoObj) IsDir(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if fi.Info.IsDir() {
		return TRUE
	}
	return FALSE
}

const PIPE_OBJ = "PIPE_OBJ"
//Pipe Object
type PipeObj struct {
	Reader *io.PipeReader
	Writer *io.PipeWriter
}

func (p *PipeObj) Inspect() string  { return PIPE_OBJ }
func (p *PipeObj) Type() ObjectType { return PIPE_OBJ }

func (p *PipeObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "read":
		return p.Read(line, args...)
	case "readClose":
		return p.ReadClose(line, args...)
	case "write":
		return p.Write(line, args...)
	case "writeClose":
		return p.WriteClose(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, p.Type()))
}

func (p *PipeObj) Read(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	var readlen int
	switch o := args[0].(type) {
	case *Integer:
		readlen = int(o.Int64)
	case *UInteger:
		readlen = int(o.UInt64)
	default:
		panic(NewError(line, PARAMTYPEERROR, "first", "read", "*Integer|*UInteger", args[0].Type()))
	}

	buffer := make([]byte, readlen)
	_, err := p.Reader.Read(buffer)
	if err != io.EOF && err != nil {
		return NewNil(err.Error())
	}

	if err == io.EOF {
		return NIL
	}

	return NewString(string(buffer))
}

func (p *PipeObj) ReadClose(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := p.Reader.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (p *PipeObj) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	n, err := p.Writer.Write([]byte(content.String))
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (p *PipeObj) WriteClose(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := p.Writer.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}
