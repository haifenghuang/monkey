package repl

import (
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

const PROMPT = ">> "

func Start(out io.Writer, color bool) {
	history := filepath.Join(os.TempDir(), ".monkey_history")
	l := liner.NewLiner()
	defer l.Close()

	l.SetCtrlCAborts(true)

	if f, err := os.Open(history); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	if color {
		eval.REPLColor = true
	}
	scope := eval.NewScope(nil)
	wd, err := os.Getwd()
	if err != nil {
		io.WriteString(out, err.Error())
		os.Exit(1)
	}

	var tmplines []string
	for {
		if line, err := l.Prompt(PROMPT); err == nil {
			if line == "exit" || line == "quit" {
				if f, err := os.Create(history); err == nil {
					l.WriteHistory(f)
					f.Close()
				}
				break
			}

			tmpline := strings.TrimSpace(line)
			if len(tmpline) == 0 { //empty line
				continue
			}
			//check if the `line` variable is ended with '\'
			if tmpline[len(tmpline)-1:] =="\\" { //the expression/statement has remaining part
				tmplines = append(tmplines, strings.TrimRight(tmpline, "\\"))
				continue
			} else {
				tmplines = append(tmplines, line)
			}

			resultLine := strings.Join(tmplines, "")
			l.AppendHistory(resultLine)
			tmplines = nil // clear the array

			lex := lexer.New("", resultLine)
			p := parser.New(lex, wd)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			eval.Eval(program, scope)
			//e := eval.Eval(program, scope)
			//io.WriteString(out, e.Inspect())
			//io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
