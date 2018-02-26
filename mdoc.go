/* Documentor generator for monkey. */
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"monkey/docs"
	"monkey/lexer"
	"monkey/parser"
	"path/filepath"
	"strings"
	"os"
)

func genDocs(path string, htmlFlag bool, isDir bool) {
	if !isDir { //single file
		genDoc(path, htmlFlag)
		return
	}

	//processing directory
	fd, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Open directory '%s' failed, reason:%v\n", path, err)
		return
	}
	defer fd.Close()

	list, err := fd.Readdir(-1)
	if err != nil {
		fmt.Errorf("Read directory '%s' failed, reason:%v\n", path, err)
		return
	}

	for _, d := range list {
		if strings.HasSuffix(d.Name(), ".my") {
			filename := filepath.Join(path, d.Name())
			genDoc(filename, htmlFlag)
		}
	}
}

func genDoc(filename string, htmlFlag bool) {
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
	md := doc.MdDocGen(file)

	//remove the '.my' extension
	genFilename := strings.TrimSuffix(filename, filepath.Ext(filename))

	//create markdown file
	mdFile := genFilename + ".md"
	outMd, err := os.Create(mdFile)
	if err != nil {
		fmt.Printf("Error creating '%s' file, reason:%v\n", mdFile, err)
		os.Exit(1)
	}


	//generate markdown file
	fmt.Fprintln(outMd, md)
	outMd.Close()

	if !htmlFlag {
		return
	}
	
	//create html file
	htmlFile := genFilename + ".html"
	outHtml, err := os.Create(htmlFile)
	if err != nil {
		fmt.Printf("Error creating '%s' file, reason:%v\n", htmlFile, err)
		os.Exit(1)
	}
	defer outHtml.Close()

	html := doc.HtmlDocGen(md, file)
	fmt.Fprintln(outHtml, html)

	err = os.Remove(mdFile)
	if err != nil {
		fmt.Printf("Error remove file '%s', reason : %v\n", mdFile, err)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [monkey file]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	var htmlFlag bool
	flag.BoolVar(&htmlFlag, "html", false, "Generate html file using github REST API.")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Invalid number of arguments!")
		flag.Usage()
	}

	path := flag.Arg(0)
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		genDocs(path, htmlFlag, true)
	case mode.IsRegular():
		genDocs(path, htmlFlag, false)
	}
	
}
