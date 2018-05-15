package main

import (
	"fmt"
	"bufio"
	"io/ioutil"
	"log"
	"time"
	"regexp"
	"runtime"
	"math/rand"
	"monkey/eval"
	"monkey/lexer"
	"monkey/parser"
	"monkey/repl"
	"os"
)

func runProgram(filename string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	f, err := ioutil.ReadFile(wd + "/" + filename)
	if err != nil {
		fmt.Println("monkey: ", err.Error())
		os.Exit(1)
	}
	l := lexer.New(filename, string(f))
	p := parser.New(l, wd)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	scope := eval.NewScope(nil)
	RegisterGoGlobals()
	eval.REPLColor = false
	eval.Eval(program, scope)
//	e := eval.Eval(program, scope)
//	if e.Inspect() != "nil" {
//		fmt.Println(e.Inspect())
//	}
}

// Register go package methods/types
// Note here, we use 'gfmt', 'glog', 'gos' 'gtime', because in monkey
// we already have built in module 'fmt', 'log' 'os', 'time'.
// And Here we demonstrate the use of import go language's methods.
func RegisterGoGlobals() {
	eval.RegisterFunctions("gfmt", []interface{}{
		fmt.Errorf,
		fmt.Println, fmt.Print, fmt.Printf,
		fmt.Fprint, fmt.Fprint, fmt.Fprintln, fmt.Fscan, fmt.Fscanf, fmt.Fscanln,
		fmt.Scan, fmt.Scanf, fmt.Scanln,
		fmt.Sscan, fmt.Sscanf, fmt.Sscanln,
		fmt.Sprint, fmt.Sprintf, fmt.Sprintln,
	})

	eval.RegisterFunctions("glog", []interface{}{
		log.Fatal, log.Fatalf, log.Fatalln, log.Flags, log.Panic, log.Panicf, log.Panicln,
		log.Print, log.Printf, log.Println, log.SetFlags, log.SetOutput, log.SetPrefix,
	})

	eval.RegisterFunctions("gos", []interface{}{
		os.Chdir, os.Chmod, os.Chown, os.Exit, os.Getpid,
		os.Hostname, os.Environ, os.Getenv, os.Setenv,
		os.Create, os.Open,
	})

	argsStart := 1
	if len(os.Args) > 2 {
		argsStart = 2
	}
	eval.RegisterVars("gos", map[string]interface{}{
		"Args": os.Args[argsStart:],
	})

	eval.RegisterVars("runtime", map[string]interface{}{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
	})

	eval.RegisterFunctions("gtime", []interface{}{
		time.Sleep, time.Now, time.Unix,
	})

	eval.RegisterFunctions("math/rand", []interface{}{
		rand.New, rand.NewSource,
		rand.Float64, rand.ExpFloat64, rand.Float32, rand.Int,
		rand.Int31, rand.Int31n, rand.Int63, rand.Int63n, rand.Intn,
		rand.NormFloat64, rand.Perm, rand.Seed, rand.Uint32,
	})

	eval.RegisterFunctions("io/ioutil", []interface{}{
		ioutil.WriteFile, ioutil.ReadFile, ioutil.TempDir, ioutil.TempFile,
		ioutil.ReadAll, ioutil.ReadDir, ioutil.NopCloser,
	})

	eval.RegisterFunctions("bufio", []interface{}{
		bufio.NewWriter, bufio.NewReader, bufio.NewReadWriter, bufio.NewScanner,
	})
	eval.RegisterFunctions("gregex", []interface{}{
		regexp.Match,regexp.MatchReader,regexp.MatchString,regexp.QuoteMeta,
		regexp.Compile,regexp.CompilePOSIX,regexp.MustCompile,regexp.MustCompilePOSIX,
	})
}

func main() {
	args := os.Args[1:]
	//We must reset `os.Args`, or the `flag` module will not functioning correctly
	os.Args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Monkey programming language REPL\n")
		repl.Start(os.Stdout, true)
	} else {
		runProgram(args[0])
	}
}
