package repl

import (
	"io"
	"monkey/eval"
	"monkey/lexer"
	"monkey/parser"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

const PROMPT = ">> "

func Start(out io.Writer) {
	history := filepath.Join(os.TempDir(), ".monkey_history")
	l := liner.NewLiner()
	defer l.Close()

	l.SetCtrlCAborts(true)

	if f, err := os.Open(history); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	scope := eval.NewScope(nil)
	wd, err := os.Getwd()
	if err != nil {
		io.WriteString(out, err.Error())
		os.Exit(1)
	}
	for {
		if line, err := l.Prompt(PROMPT); err == nil {
			if line == "exit" {
				if f, err := os.Create(history); err == nil {
					l.WriteHistory(f)
					f.Close()
				}
				break
			}
			l.AppendHistory(line)
			lex := lexer.New("", line)
			p := parser.New(lex, wd)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}
			e := eval.Eval(program, scope)
			io.WriteString(out, e.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
