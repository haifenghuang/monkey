/* Documentor generator for monkey. */
package main

import (
	"fmt"
	"io/ioutil"
	"monkey/docs"
	"monkey/lexer"
	"monkey/parser"
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
	p := parser.NewWithDoc(l, wd)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	//generate markdown docs
	file := doc.New(filename, program)
	fmt.Println(doc.MdDocGen(file))
//	for _, cls := range file.Classes {
//		fmt.Printf("class name=[%s]\n", cls.Value.Name)
//		fmt.Printf("class Doc=[%s]\n", cls.Value.Doc)
//		fmt.Printf("class Text=[%s]\n", cls.Value.Text)
//		fmt.Printf("------------------------------------------------\n")
//
//		for _, fn := range cls.Funcs {
//			fmt.Printf("func name=[%s]\n", fn.Name)
//			fmt.Printf("func Doc=[%s]\n", fn.Doc)
//			fmt.Printf("func Text=[%s]\n", fn.Text)
//		}
//		fmt.Printf("------------------------------------------------\n")
//
//		for _, let := range cls.Lets {
//			fmt.Printf("let name=[%s]\n", let.Name)
//			fmt.Printf("let Doc=[%s]\n", let.Doc)
//			fmt.Printf("let Text=[%s]\n", let.Text)
//		}
//		fmt.Printf("------------------------------------------------\n")
//
//		for _, prop := range cls.Props {
//			fmt.Printf("prop name=[%s]\n", prop.Name)
//			fmt.Printf("prop Doc=[%s]\n", prop.Doc)
//			fmt.Printf("prop Text=[%s]\n", prop.Text)
//		}
//
//		fmt.Printf("\n\n\n\n\n")
//	}
}

func main() {
	args := os.Args[1:]
	//We must reset `os.Args`, or the `flag` module will not functioning correctly
	os.Args = os.Args[1:]
	if len(args) > 0 {
		runProgram(args[0])
	}
}
