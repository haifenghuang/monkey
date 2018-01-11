package main

import (
	"fmt"
	"io/ioutil"
	"monkey/formatter"
	"os"
)

func main() {
	args := os.Args[1:]

	var f []byte
	var err error
	if len(args) == 0 {
		f, err = ioutil.ReadAll(os.Stdin)
	} else {
		f, err = ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println("Formatter: cannot read file", err.Error())
			os.Exit(1)
		}
	}

	formatter := formatter.New(string(f))
	formatter.Format()
}
