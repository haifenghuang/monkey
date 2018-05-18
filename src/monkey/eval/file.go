package eval

import (
	"bufio"
	_ "fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	IOUTIL_OBJ = "IOUTIL_OBJ"
	ioutil_name = "ioutil"
)

type IOUtilObj struct{}

func NewIOUtilObj() Object {
	ret := &IOUtilObj{}
	SetGlobalObj(ioutil_name, ret)

	return ret
}

func (i *IOUtilObj) Inspect() string  { return "<" + ioutil_name + ">" }
func (i *IOUtilObj) Type() ObjectType { return IOUTIL_OBJ }
func (i *IOUtilObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "readAll":
		return i.ReadAll(line, args...)
	case "readDir":
		return i.ReadDir(line, args...)
	case "readFile":
		return i.ReadFile(line, args...)
	case "tempDir":
		return i.TempDir(line, args...)
	case "tempFile":
		return i.TempFile(line, args...)
	case "writeFile":
		return i.WriteFile(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, i.Type()))
	}
}

func (i *IOUtilObj) ReadAll(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	fObj, ok := args[0].(*FileObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "readAll", "*FileObject", args[0].Type()))
	}

	reader := bufio.NewReader(fObj.File)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return NewNil(err.Error())
	}
	return NewString(string(b))
}

func (i *IOUtilObj) ReadDir(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	dirname, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "readDir", "*String", args[0].Type()))
	}

	files, err := ioutil.ReadDir(dirname.String)
	if err != nil {
		return NewNil(err.Error())
	}

	arr := &Array{}
	for _, file := range files {
		arr.Members = append(arr.Members, &FileInfoObj{Info: file})
	}

	return arr
}

func (i *IOUtilObj) ReadFile(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	filename, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "readFile", "*String", args[0].Type()))
	}

	b, err := ioutil.ReadFile(filename.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(string(b))
}

func (i *IOUtilObj) TempDir(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	dir, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "tempDir", "*String", args[0].Type()))
	}

	prefix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "tempDir", "*String", args[1].Type()))
	}

	name, err := ioutil.TempDir(dir.String, prefix.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewString(name)
}

func (i *IOUtilObj) TempFile(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	dir, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "tempFile", "*String", args[0].Type()))
	}

	prefix, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "tempFile", "*String", args[1].Type()))
	}

	f, err := ioutil.TempFile(dir.String, prefix.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return &FileObject{File: f, Name: f.Name()}
}

func (i *IOUtilObj) WriteFile(line string, args ...Object) Object {
	if len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "3", len(args)))
	}

	filename, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeFile", "*String", args[0].Type()))
	}

	data, ok := args[1].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "writeFile", "*String", args[1].Type()))
	}

	perm, ok := args[2].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "writeFile", "*String", args[2].Type()))
	}

	err := ioutil.WriteFile(filename.String, []byte(data.String), os.FileMode(int(perm.Int64)))
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

type FileObject struct {
	File    *os.File
	Name    string
	Scanner *bufio.Scanner
	reader  *bufio.Reader
	writer  *bufio.Writer
}

//Implement the 'Closeable' interface
func (f *FileObject) close(line string, args ...Object) Object {
	return f.Close(line, args...)
}

func (f *FileObject) IOWriter() io.Writer { return f.File }
func (f *FileObject) Inspect() string     { return "<file object: " + f.Name + ">" }
func (f *FileObject) Type() ObjectType    { return FILE_OBJ }
func (f *FileObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "close":
		return f.Close(line, args...)
	case "read":
		return f.Read(line, args...)
	case "readAt":
		return f.ReadAt(line, args...)
	case "readRune":
		return f.ReadRune(line, args...)
	case "readLine":
		return f.ReadLine(line, args...)
	case "seek":
		return f.Seek(line, args...)
	case "stat":
		return f.Stat(line, args...)
	case "sync":
		return f.Sync(line, args...)
	case "truncate":
		return f.Truncate(line, args...)
	case "write":
		return f.Write(line, args...)
	case "writeAt":
		return f.WriteAt(line, args...)
	case "writeString":
		return f.WriteString(line, args...)
	case "writeLine":
		return f.WriteLine(line, args...)
	case "name":
		return f.GetName(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, f.Type()))
	}
}

func (f *FileObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := f.File.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//Note: This method will return three different values:
//   1. nil    - with error message    (ERROR)
//   2. nil    - without error message (EOF)
//   3. string - read string
func (f *FileObject) Read(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	readlen, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "read", "*Integer", args[0].Type()))
	}

	buffer := make([]byte, int(readlen.Int64))
	n, err := f.File.Read(buffer)
	if err != io.EOF && err != nil {
		return NewNil(err.Error())
	}

	if n == 0 && err == io.EOF {
		return NIL
	}
	return NewString(string(buffer))
}

func (f *FileObject) ReadAt(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	readlen, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "readAt", "*Integer", args[0].Type()))
	}

	offset, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "readAt", "*Integer", args[1].Type()))
	}

	buffer := make([]byte, int(readlen.Int64))
	_, err := f.File.ReadAt(buffer, offset.Int64)
	if err != io.EOF && err != nil {
		return NewNil(err.Error())
	}

	if err == io.EOF {
		return NIL
	}
	return NewString(string(buffer))
}

func (f *FileObject) ReadRune(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if f.reader == nil {
		f.reader = bufio.NewReader(f.File)
	}

	r, _, err := f.reader.ReadRune()
	if err != io.EOF && err != nil {
		return NewNil(err.Error())
	}

	if err == io.EOF {
		return NIL
	}

	return NewString(string(r))
}

func (f *FileObject) ReadLine(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	if f.Scanner == nil {
		f.Scanner = bufio.NewScanner(f.File)
		f.Scanner.Split(bufio.ScanLines)
	}
	aLine := f.Scanner.Scan()
	if err := f.Scanner.Err(); err != nil {
		return NewNil(err.Error())
	}
	if !aLine {
		return NIL
	}
	return NewString(f.Scanner.Text())
}

func (f *FileObject) Seek(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	offset, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "seek", "*Integer", args[0].Type()))
	}

	whence, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "seek", "*Integer", args[1].Type()))
	}

	ret, err := f.File.Seek(offset.Int64, int(whence.Int64))
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(ret)
}

func (f *FileObject) Stat(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	fileInfo, err := f.File.Stat()
	if err != nil {
		return NewNil(err.Error())
	}

	//return FileInfo2HashObj(fileInfo)
	return &FileInfoObj{Info: fileInfo}
}

func (f *FileObject) Sync(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := f.File.Sync()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (f *FileObject) Truncate(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	size, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "truncate", "*Integer", args[0].Type()))
	}

	err := f.File.Truncate(size.Int64)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (f *FileObject) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "write", "*String", args[0].Type()))
	}

	n, err := f.File.Write([]byte(content.String))
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FileObject) WriteAt(line string, args ...Object) Object {
	if len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "2", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeAt", "*String", args[0].Type()))
	}

	offset, ok := args[1].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "second", "writeAt", "*Integer", args[1].Type()))
	}

	ret, err := f.File.WriteAt([]byte(content.String), offset.Int64)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(ret))
}

func (f *FileObject) WriteString(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeString", "*String", args[0].Type()))
	}

	ret, err := f.File.WriteString(content.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(ret))
}

func (f *FileObject) WriteLine(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	content, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeLine", "*String", args[0].Type()))
	}

	ret, err := f.File.Write([]byte(content.String + "\n"))
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(ret))
}

func (f *FileObject) GetName(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return NewString(f.File.Name())
}
