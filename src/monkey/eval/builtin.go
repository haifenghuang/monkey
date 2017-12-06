package eval

import (
	"container/list"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var fileModeTable = map[string]int{
	"r":   os.O_RDONLY,
	"<":   os.O_RDONLY,
	"w":   os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
	">":   os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
	"a":   os.O_APPEND | os.O_CREATE,
	">>":  os.O_APPEND | os.O_CREATE,
	"r+":  os.O_RDWR,
	"+<":  os.O_RDWR,
	"w+":  os.O_RDWR | os.O_CREATE | os.O_TRUNC,
	"+>":  os.O_RDWR | os.O_CREATE | os.O_TRUNC,
	"a+":  os.O_RDWR | os.O_APPEND | os.O_CREATE,
	"+>>": os.O_RDWR | os.O_APPEND | os.O_CREATE,
}

type BuiltinFunc func(line string, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunc
}

var builtins map[string]*Builtin

func absBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			i, ok := args[0].(*Integer)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "abs", "*Integer", args[0].Type()))
			}
			if i.Int64 > -1 {
				return i
			}
			return NewInteger(i.Int64 * -1)
		}, //Here the ',' is a must, it confused me a lot
	}
}

func rangeBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 && len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
			}
			i, ok := args[0].(*Integer)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "range", "*Integer", args[0].Type()))
			}

			iValue := i.Int64
			if iValue <= 0 {
				return &Array{}
			}

			var jValue int64
			if len(args) == 2 {
				j, ok := args[1].(*Integer)
				if !ok {
					panic(NewError(line, PARAMTYPEERROR, "second", "range", "*Integer", args[0].Type()))
				}
				jValue = j.Int64
				if jValue <= 0 {
					panic(NewError(line, GENERICERROR, "second parameter of 'range' should be >=0"))
				}
			}

			var k int64
			methods := &Array{}
			for k = 0; k < iValue; {
				methods.Members = append(methods.Members, NewInteger(k))
				if len(args) == 2 {
					k += jValue
				} else {
					k += 1
				}
			}
			return methods
		},
	}
}

func addmBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 3 {
				panic(NewError(line, ARGUMENTERROR, "3", len(args)))
			}
			st, ok := args[0].(*Struct)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "addm", "*Struct", args[0].Type()))
			}
			name, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "addm", "*String", args[1].Type()))
			}
			fn, ok := args[2].(*Function)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "third", "addm", "*Function", args[2].Type()))
			}
			st.methods[name.String] = fn
			return st
		},
	}
}

func chrBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			i, ok := args[0].(*Integer)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "chr", "*Integer", args[0].Type()))
			}
			if i.Int64 < 0 || i.Int64 > 255 {
				panic(NewError(line, INPUTERROR, i.Inspect(), "chr"))
			}
			return NewString(string(i.Int64))
		},
	}
}

func newFileBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			var fname *String
			var flag int = os.O_RDONLY
			var ok bool
			var perm os.FileMode = os.FileMode(0666)

			argLen := len(args)
			if argLen < 1 {
				panic(NewError(line, ARGUMENTERROR, "at least one", argLen))
			}

			fname, ok = args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "newFile", "*String", args[0].Type()))
			}

			if argLen == 2 {
				m, ok := args[1].(*String)
				if !ok {
					panic(NewError(line, PARAMTYPEERROR, "second", "newFile", "*String", args[1].Type()))
				}

				flag, ok = fileModeTable[m.String]
				if !ok {
					panic(NewError(line, FILEMODEERROR))
				}
			}

			if len(args) == 3 {
				p, ok := args[2].(*Integer)
				if !ok {
					panic(NewError(line, PARAMTYPEERROR, "third", "newFile", "*Integer", args[2].Type()))
				}

				perm = os.FileMode(int(p.Int64))
			}

			f, err := os.OpenFile(fname.String, flag, perm)
			if err != nil {
				return NewNil(err.Error())
			}
			return &FileObject{File: f, Name: "<file object: " + fname.String + ">"}
		},
	}
}

func intBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			switch input := args[0].(type) {
			case *Integer:
				return input
			case *Float:
				return NewInteger(int64(input.Float64))
			case *Boolean:
				if input.Bool {
					return NewInteger(1)
				}
				return NewInteger(0)
			case *String:
				var n int64
				var err error

				if strings.HasPrefix(input.String, "0b") {
					n, err = strconv.ParseInt(input.String[2:], 2, 64)
				} else if strings.HasPrefix(input.String, "0x") {
					n, err = strconv.ParseInt(input.String[2:], 16, 64)
				} else if strings.HasPrefix(input.String, "0c") {
					n, err = strconv.ParseInt(input.String[2:], 8, 64)
				} else {
					n, err = strconv.ParseInt(input.String, 10, 64)
				}
				if err != nil {
					panic(NewError(line, INPUTERROR, "STRING: "+input.String, "int"))
				}
				return NewInteger(int64(n))
			}
			panic(NewError(line, PARAMTYPEERROR, "first", "int", "*String|*Integer|*Boolean|*Float", args[0].Type()))
		},
	}
}

func floatBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			switch input := args[0].(type) {
			case *Integer:
				return NewFloat(float64(input.Int64))
			case *Float:
				return input
			case *Boolean:
				if input.Bool {
					return NewFloat(1)
				}
				return NewFloat(0)
			case *String:
				var n float64
				var err error
				var k int64

				if strings.HasPrefix(input.String, "0b") {
					k, err = strconv.ParseInt(input.String[2:], 2, 64)
					if err == nil {
						n = float64(k)
					}
				} else if strings.HasPrefix(input.String, "0x") {
					k, err = strconv.ParseInt(input.String[2:], 16, 64)
					if err == nil {
						n = float64(k)
					}
				} else if strings.HasPrefix(input.String, "0c") {
					k, err = strconv.ParseInt(input.String[2:], 8, 64)
					if err == nil {
						n = float64(k)
					}
				} else {
					n, err = strconv.ParseFloat(input.String, 64)
				}
				if err != nil {
					panic(NewError(line, INPUTERROR, "STRING: "+input.String, "float"))
				}
				return NewFloat(float64(n))
			}
			panic(NewError(line, PARAMTYPEERROR, "first", "float", "*String|*Integer|*Boolean|*Float", args[0].Type()))
		},
	}
}

func strBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			switch input := args[0].(type) {
			case *String:
				return input
			default:
				return NewString(input.Inspect())
			}
			//panic(NewError(line, INPUTERROR, args[0].Type(), "str"))
		},
	}
}

func arrayBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			switch input := args[0].(type) {
			case *Array:
				return input
			default:
				return &Array{Members: []Object{input}}
			}
		},
	}
}

func lenBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			switch arg := args[0].(type) {
			case *String:
				n := utf8.RuneCountInString(arg.String)
				return NewInteger(int64(n))
			case *Array:
				return NewInteger(int64(len(arg.Members)))
			case *Hash:
				return NewInteger(int64(len(arg.Pairs)))
			case *Nil:
				return NewInteger(0)
			}
			panic(NewError(line, PARAMTYPEERROR, "first", "len", "*String|*Array|*Hash|*Nil", args[0].Type()))
		},
	}
}

func methodsBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			arr := &Array{}
			t := reflect.TypeOf(args[0])
			for i := 0; i < t.NumMethod(); i++ {
				m := t.Method(i).Name
				if !(m == "Type" || m == "CallMethod" || m == "HashKey" || m == "Inspect") {
					arr.Members = append(arr.Members, NewString(strings.ToLower(m)))
				}
			}
			return arr
		},
	}
}

func ordBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			s, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "ord", "*String", args[0].Type()))
			}
			if len(s.String) > 1 {
				panic(NewError(line, INLENERR, "1", len(s.String)))
			}
			return NewInteger(int64(s.String[0]))
		},
	}
}

func printBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) == 0 {
				n, err := fmt.Print()
				if err != nil {
					return NewNil(err.Error())
				}
				return NewInteger(int64(n))
			}

			wrapped := make([]interface{}, len(args))
			for i, v := range args {
				wrapped[i] = &Formatter{Obj: v}
			}

			n, err := fmt.Print(wrapped...)
			if err != nil {
				return NewNil(err.Error())
			}

			return NewInteger(int64(n))
		},
	}
}

func printlnBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) == 0 {
				n, err := fmt.Println()
				if err != nil {
					return NewNil(err.Error())
				}
				return NewInteger(int64(n))
			}

			wrapped := make([]interface{}, len(args))
			for i, v := range args {
				wrapped[i] = &Formatter{Obj: v}
			}

			n, err := fmt.Println(wrapped...)
			if err != nil {
				return NewNil(err.Error())
			}

			return NewInteger(int64(n))
		},
	}
}

func printfBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) < 1 {
				panic(NewError(line, ARGUMENTERROR, ">0", len(args)))
			}

			format, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "printf", "*String", args[0].Type()))
			}

			subArgs := args[1:]
			wrapped := make([]interface{}, len(subArgs))
			for i, v := range subArgs {
				wrapped[i] = &Formatter{Obj: v}
			}

			n, err := fmt.Printf(format.String, wrapped...)
			if err != nil {
				return NewNil(err.Error())
			}

			return NewInteger(int64(n))
		},
	}
}

func typeBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}
			return NewString(fmt.Sprintf("%s", args[0].Type()))
		},
	}
}

func chanBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) == 0 {
				return &ChanObject{make(chan Object)}
			} else if len(args) == 1 {
				v, ok := args[0].(*Integer)
				if !ok {
					panic(NewError(line, PARAMTYPEERROR, "first", "chan", "*Integer", args[0].Type()))
				}
				return &ChanObject{make(chan Object, v.Int64)}
			}
			panic(NewError(line, ARGUMENTERROR, "Not 0|1", len(args)))
		},
	}
}

func assertBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}

			v, ok := args[0].(*Boolean)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "assert", "*Boolean", args[0].Type()))
			}

			if v.Bool == true {
				return NIL
			}

			panic(NewError(line, ASSERTIONERROR))
		},
	}
}


func reverseBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}

			switch input := args[0].(type) {
			case *String:
				word := []rune(input.String)
				reverse := []rune{}
				for i := len(word) - 1; i >= 0; i-- {
					reverse = append(reverse, word[i])
				}

				return NewString(string(reverse))
			case *Array:
				reverse := &Array{}
				for i := len(input.Members) - 1; i >= 0; i-- {
					reverse.Members = append(reverse.Members, input.Members[i])
				}
				return reverse
			case *Hash:
				hash := &Hash{Pairs: make(map[HashKey]HashPair)}
				for _, v := range input.Pairs {
					if hashable, ok := v.Value.(Hashable); ok {
						hash.Pairs[hashable.HashKey()] = HashPair{Key: v.Value, Value: v.Key}
					} else {
						panic(NewError(line, GENERICERROR, fmt.Sprintf("The hash's value(%s) is not hashable", v.Value.Inspect())))
					}
				}
				return hash
			default:
				panic(NewError(line, PARAMTYPEERROR, "first", "reverse", "*Array|*String", args[0].Type()))
			}
			return NIL
		},
	}
}

func dialTCPBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}
			netStr, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "dialTCP", "*String", args[0].Type()))
			}
			addrStr, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "dialTCP", "*String", args[1].Type()))
			}

			tcpAddr, err := net.ResolveTCPAddr(netStr.String, addrStr.String)
			if err != nil {
				return NewNil(err.Error())
			}

			conn, err := net.DialTCP(netStr.String, nil, tcpAddr)
			if err != nil {
				return NewNil(err.Error())
			}

			return &TcpConnObject{Conn: conn, Address: tcpAddr.String()}
		},
	}
}

func listenTCPBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}
			netStr, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "listenTCP", "*String", args[0].Type()))
			}
			addrStr, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "listenTCP", "*String", args[1].Type()))
			}

			tcpAddr, err := net.ResolveTCPAddr(netStr.String, addrStr.String)
			if err != nil {
				return NewNil(err.Error())
			}

			listener, err := net.ListenTCP(netStr.String, tcpAddr)
			if err != nil {
				return NewNil(err.Error())
			}

			return &TCPListenerObject{Listener: listener, Address: tcpAddr.String()}
		},
	}
}

func dialUDPBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}
			netStr, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "dialUDP", "*String", args[0].Type()))
			}
			addrStr, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "dialUDP", "*String", args[1].Type()))
			}

			udpAddr, err := net.ResolveUDPAddr(netStr.String, addrStr.String)
			if err != nil {
				return NewNil(err.Error())
			}

			conn, e := net.DialUDP(netStr.String, nil, udpAddr)
			if e != nil {
				return NewNil(err.Error())
			}

			return &UdpConnObject{Conn: conn, Address: udpAddr.String()}
		},
	}
}

func dialUnixBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}
			netStr, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "dialUnix", "*String", args[0].Type()))
			}
			addrStr, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "dialUnix", "*String", args[1].Type()))
			}

			unixAddr, err := net.ResolveUnixAddr(netStr.String, addrStr.String)
			if err != nil {
				return NewNil(err.Error())
			}

			conn, err := net.DialUnix(netStr.String, nil, unixAddr)
			if err != nil {
				return NewNil(err.Error())
			}

			return &UnixConnObject{Conn: conn, Address: unixAddr.String()}
		},
	}
}

func listenUnixBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}
			netStr, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "listenTCP", "*String", args[0].Type()))
			}
			addrStr, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "listenTCP", "*String", args[1].Type()))
			}

			unixAddr, err := net.ResolveUnixAddr(netStr.String, addrStr.String)
			if err != nil {
				return NewNil(err.Error())
			}

			listener, err := net.ListenUnix(netStr.String, unixAddr)
			if err != nil {
				return NewNil(err.Error())
			}

			return &UnixListenerObject{Listener: listener, Address: unixAddr.String()}
		},
	}
}

func dbOpenBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}
			driverName, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "dbOpen", "*String", args[0].Type()))
			}

			dataSourceName, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "dbOpen", "*String", args[1].Type()))
			}

			db, err := sql.Open(driverName.String, dataSourceName.String)
			if err != nil {
				return NewNil(err.Error())
			}

			return &SqlObject{Db: db, Name: fmt.Sprintf("%s:%s", driverName.String, dataSourceName.String)}
		},
	}
}

func newTimeBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 && len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "0|1", len(args)))
			}

			if len(args) == 0 {
				return &TimeObj{Tm: time.Now(), Valid: true}
			}

			location, ok := args[0].(*Integer)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "newTime", "*Integer", args[0].Type()))
			}
			if location.Int64 == UTC {
				return &TimeObj{Tm: time.Now().UTC(), Valid: true}
			}
			return &TimeObj{Tm: time.Now(), Valid: true}
		},
	}
}

//func Date(year int, month Month, day, hour, min, sec, nsec int, loc int)
func newDateBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			argLen := len(args)
			if argLen != 7 && argLen != 8 {
				panic(NewError(line, ARGUMENTERROR, "7|8", len(args)))
			}

			year, ok1 := args[0].(*Integer)
			if !ok1 {
				panic(NewError(line, PARAMTYPEERROR, "first", "newDate", "*Integer", args[0].Type()))
			}

			month, ok2 := args[1].(*Integer)
			if !ok2 {
				panic(NewError(line, PARAMTYPEERROR, "second", "newDate", "*Integer", args[1].Type()))
			}

			day, ok3 := args[2].(*Integer)
			if !ok3 {
				panic(NewError(line, PARAMTYPEERROR, "third", "newDate", "*Integer", args[2].Type()))
			}

			hour, ok4 := args[3].(*Integer)
			if !ok4 {
				panic(NewError(line, PARAMTYPEERROR, "fourth", "newDate", "*Integer", args[3].Type()))
			}

			min, ok5 := args[4].(*Integer)
			if !ok5 {
				panic(NewError(line, PARAMTYPEERROR, "fifth", "newDate", "*Integer", args[4].Type()))
			}

			sec, ok6 := args[5].(*Integer)
			if !ok6 {
				panic(NewError(line, PARAMTYPEERROR, "sixth", "newDate", "*Integer", args[5].Type()))
			}

			nsec, ok7 := args[6].(*Integer)
			if !ok7 {
				panic(NewError(line, PARAMTYPEERROR, "seventh", "newDate", "*Integer", args[6].Type()))
			}

			var location Object
			var ok8 bool
			if argLen == 8 {
				location, ok8 = args[7].(*Integer)
				if !ok8 {
					panic(NewError(line, PARAMTYPEERROR, "eighth", "newDate", "*Integer", args[7].Type()))
				}
			}

			if argLen == 7 {
				return &TimeObj{Tm: time.Date(int(year.Int64), time.Month(month.Int64), int(day.Int64),
					int(hour.Int64), int(min.Int64), int(sec.Int64), int(nsec.Int64),
					time.Local), Valid: true}
			} else {
				var loc *time.Location
				if location.(*Integer).Int64 == LOCAL {
					loc = time.Local
				} else {
					loc = time.UTC
				}
				return &TimeObj{Tm: time.Date(int(year.Int64), time.Month(month.Int64), int(day.Int64),
					int(hour.Int64), int(min.Int64), int(sec.Int64), int(nsec.Int64),
					loc), Valid: true}
			}

		},
	}
}

func newCondBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", len(args)))
			}

			switch arg := args[0].(type) {
			case *SyncMutexObj:
				return &SyncCondObj{Cond: sync.NewCond(arg.Mutex)}
			case *SyncRWMutexObj:
				return &SyncCondObj{Cond: sync.NewCond(arg.RWMutex)}
			default:
				panic(NewError(line, PARAMTYPEERROR, "first", "newCond", "*SyncMutexObj|*SyncRWMutexObj", args[0].Type()))
			}
		},
	}
}

func newOnceBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0", len(args)))
			}
			return &SyncOnceObj{Once: new(sync.Once)}
		},
	}
}

func newMutexBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0", len(args)))
			}
			return &SyncMutexObj{Mutex: new(sync.Mutex)}
		},
	}
}

func newRWMutexBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0", len(args)))
			}
			return &SyncRWMutexObj{RWMutex: new(sync.RWMutex)}
		},
	}
}

func newWaitGroupBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0", len(args)))
			}
			return &SyncWaitGroupObj{WaitGroup: new(sync.WaitGroup)}
		},
	}
}

func newPipeBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0", len(args)))
			}
			r, w := io.Pipe()
			return &PipeObj{Reader: r, Writer: w}
		},
	}
}

func newLoggerBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 3 && len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0|3", len(args)))
			}

			if len(args) == 0 {
				logger := log.New(os.Stdout, "", log.LstdFlags)
				return &LoggerObj{Logger: logger}
			}

			out, ok := args[0].(Writable)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "newLogger", "Writable", args[0].Type()))
			}

			prefix, ok := args[1].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "second", "newLogger", "*String", args[1].Type()))
			}

			flag, ok := args[2].(*Integer)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "third", "newLogger", "*Integer", args[2].Type()))
			}

			logger := log.New(out.IOWriter(), prefix.String, int(flag.Int64))
			return &LoggerObj{Logger: logger}
		},
	}
}

func newListBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 0 {
				panic(NewError(line, ARGUMENTERROR, "0", len(args)))
			}

			return &ListObject{List: list.New()}
		},
	}
}

func newDeepEqualBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			if len(args) != 2 {
				panic(NewError(line, ARGUMENTERROR, "2", len(args)))
			}

			r := reflect.DeepEqual(args[0], args[1])
			if r {
				return TRUE
			}
			return FALSE
		},
	}
}

func newCsvReaderBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			argLen := len(args)
			if argLen != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", argLen))
			}

			fname, ok := args[0].(*String)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "newCsv", "*String", args[0].Type()))
			}

			f, err := os.Open(fname.String)
			if err != nil {
				return NewNil(err.Error())
			}

			return &CsvObj{Reader: csv.NewReader(f), ReaderFile:f}
		},
	}
}

func newCsvWriterBuiltin() *Builtin {
	return &Builtin{
		Fn: func(line string, args ...Object) Object {
			argLen := len(args)
			if argLen != 1 {
				panic(NewError(line, ARGUMENTERROR, "1", argLen))
			}

			writer, ok := args[0].(Writable)
			if !ok {
				panic(NewError(line, PARAMTYPEERROR, "first", "newCsvWriterBuiltin", "Writable", args[0].Type()))
			}

			return &CsvObj{Writer: csv.NewWriter(writer.IOWriter())}
		},
	}
}

func RegisterBuiltin(name string, f *Builtin) {
	builtins[strings.ToLower(name)] = f
}

func init() {
	builtins = map[string]*Builtin{
		"abs":     absBuiltin(),
		"range":   rangeBuiltin(),
		"addm":    addmBuiltin(),
		"chr":     chrBuiltin(),
		"newFile": newFileBuiltin(),
		"int":     intBuiltin(),
		"float":   floatBuiltin(),
		"str":     strBuiltin(),
		"array":   arrayBuiltin(),
		"len":     lenBuiltin(),
		"methods": methodsBuiltin(),
		"ord":     ordBuiltin(),
		"print":   printBuiltin(),
		"println": printlnBuiltin(),
		"printf":  printfBuiltin(),
		"type":    typeBuiltin(),
		"chan":    chanBuiltin(),
		"assert":  assertBuiltin(),
		"reverse": reverseBuiltin(),

		//net
		"dialTCP":    dialTCPBuiltin(),
		"listenTCP":  listenTCPBuiltin(),
		"dialUDP":    dialUDPBuiltin(),
		"dialUnix":   dialUnixBuiltin(),
		"listenUnix": listenUnixBuiltin(),

		//database
		"dbOpen": dbOpenBuiltin(),

		//time
		"newTime": newTimeBuiltin(),
		"newDate": newDateBuiltin(),

		//sync
		"newCond":      newCondBuiltin(),
		"newOnce":      newOnceBuiltin(),
		"newMutex":     newMutexBuiltin(),
		"newRWMutex":   newRWMutexBuiltin(),
		"newWaitGroup": newWaitGroupBuiltin(),

		//pipe
		"newPipe": newPipeBuiltin(),

		//Logger
		"newLogger": newLoggerBuiltin(),

		//container
		"newList": newListBuiltin(),

		//deepEqual
		"deepEqual": newDeepEqualBuiltin(),

		//csv
		"newCsvReader": newCsvReaderBuiltin(),
		"newCsvWriter": newCsvWriterBuiltin(),
	}
}
