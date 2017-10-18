package main

import (
	"fmt"
	"io/ioutil"
	"monkey/highlight"
	"os"
)

func main() {
	args := os.Args[1:]

	var outFileName string = "output.html"
	var f []byte
	var err error

	if len(args) == 0 {
		f, err = ioutil.ReadAll(os.Stdin)
	} else {
		f, err = ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println("Highlighter: cannot read file", err.Error())
			os.Exit(1)
		}
		outFileName = fmt.Sprintf("%s.html", args[0])
	}

	highlighter := highlight.New(string(f))
	highlighter.RegisterGenerator(&highlight.ConsoleHighlighter{})

	file, err := os.Create(outFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	highlighter.RegisterGenerator(&highlight.HtmlHighlighter{Out: file})
	highlighter.Highlight()
}
