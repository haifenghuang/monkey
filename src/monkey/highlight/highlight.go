package highlight

import (
	"errors"
	_ "fmt"
	"io"
	"strings"
	"unicode"
)

//HighlightIntf is the interface of different Highlighter generators
type HighlightIntf interface {
	Name() string      //The Name of the generator
	Writer() io.Writer //The io.Writer you want the output to write to

	WriteHeader() string
	WriteFooter() string

	WriteLineHead(linNo int) string
	WriteLineTail() string
	WriteNewLine() string

	WriteQuotes(quotes string) string     //quote(single, double, raw)
	WriteComment(comment string) string   //comment
	WriteKeyword(keyword string) string   //keyword
	WriteOperator(operator string) string //operator
	WriteNumber(number string) string     //number
	WriteNormal(text string) string       //normal text
}

//This map is mainly for quick searching
var keywords = map[string]int{
	"fn":       1,
	"let":      1,
	"true":     1,
	"false":    1,
	"if":       1,
	"else":     1,
	"elsif":    1,
	"elseif":   1,
	"elif":     1,
	"return":   1,
	"include":  1,
	"and":      1,
	"or":       1,
	"struct":   1,
	"do":       1,
	"while":    1,
	"break":    1,
	"continue": 1,
	"for":      1,
	"in":       1,
	"where":    1,
	"grep":     1,
	"map":      1,
	"case":     1,
	"is":       1,
	"try":      1,
	"catch":    1,
	"finally":  1,
	"throw":    1,
	"qw":       1,
	"unless":   1,
	"spawn":    1,
	"enum":     1,
	"defer":    1,
	"nil":      1,
	"class":    1,
	"new":      1,
	"this":     1,
	"parent":   1,
	"property": 1,
	"get":      1,
	"set":      1,
	"static":   1,
	"public":   1,
	"private":  1,
	"protected":1,
	"interface":1,
	"default":  1,
}

const (
	op_chars = "+-*/%()[]=<>!&|.,^?:;{}"
)

type Highlighter struct {
	input  []rune //the source we need to highlight
	pos    int
	length int

	lineNo    int
	generator map[string]HighlightIntf
}

func New(input string) *Highlighter {
	h := &Highlighter{input: []rune(input)}
	h.pos = 0
	h.length = len(h.input)
	h.lineNo = 1

	h.generator = make(map[string]HighlightIntf)

	return h
}

//Highlight the source code
func (h *Highlighter) Highlight() {
	//Header
	for _, intf := range h.generator {
		str := intf.WriteHeader()
		if len(str) > 0 {
			io.WriteString(intf.Writer(), str)
		}

		str = intf.WriteLineHead(h.lineNo)
		if len(str) > 0 {
			io.WriteString(intf.Writer(), str)
		}
	}

	for h.pos < h.length {
		current := h.peek(0)
		if current == '"' || current == '`' || current == '~' || current == '\'' {
			h.processText(current)
		} else if current == '#' {
			h.processComment(current)
		} else if isLetter(current) {
			h.processIdentifier()
		} else if strings.Contains(op_chars, string(current)) {
			h.processOperator()
		} else if isDigit(current) {
			h.processNumber()
		} else {
			if h.input[h.pos] == '\n' {
				h.lineNo++
				for _, intf := range h.generator {
					str := intf.WriteNewLine()
					if len(str) > 0 {
						io.WriteString(intf.Writer(), str)
					}

					str = intf.WriteLineTail()
					if len(str) > 0 {
						io.WriteString(intf.Writer(), str)
					}

					str = intf.WriteLineHead(h.lineNo)
					if len(str) > 0 {
						io.WriteString(intf.Writer(), str)
					}
				}
			} else {
				h.processNormal()
			}
			h.next()
		}
	} //end for

	//footer
	for _, intf := range h.generator {
		str := intf.WriteFooter()
		if len(str) > 0 {
			io.WriteString(intf.Writer(), str)
		}
	}
}

//RegisterGenerator register a highlighter
func (h *Highlighter) RegisterGenerator(intf HighlightIntf) {
	h.generator[intf.Name()] = intf
}

func (h *Highlighter) processText(ch rune) error {
	var ret []rune

	ret = append(ret, ch)

	for {
		h.next()
		if h.input[h.pos] == 0 {
			return errors.New("unexpected EOF")
		}

		if h.input[h.pos] == ch {
			h.next()
			break
		}

		ret = append(ret, h.input[h.pos])
	}
	ret = append(ret, ch)

	for _, intf := range h.generator {
		str := intf.WriteQuotes(string(ret))
		if len(str) > 0 {
			io.WriteString(intf.Writer(), str)
		}
	}

	return nil
}

func (h *Highlighter) processComment(ch rune) {
	var ret []rune
	if ch == '#' {
		ret = append(ret, ch)
		h.next()
	} else {
		ret = append(ret, ch)
		ret = append(ret, ch)
		h.next()
		h.next()
	}

	for h.peek(0) != '\n' && h.peek(0) != 0 {
		ret = append(ret, h.input[h.pos])
		h.next()
	}

	for _, intf := range h.generator {
		str := intf.WriteComment(string(ret))
		if len(str) > 0 {
			io.WriteString(intf.Writer(), str)
		}
	}
}

//Really need this function?
func (h *Highlighter) processOperator() {
	opArr := []struct {
		operator string
		opLen    int
	}{
		{"+=", 2},
		{"-=", 2},
		{"*=", 2},
		{"/=", 2},
		{"%=", 2},
		{"^=", 2},

		{"++", 2},
		{"--", 2},

		{"&&", 2},
		{"||", 2},

		{"<<", 2},
		{">>", 2},

		{"->", 2},
		{"=>", 2},

		{"==", 2},
		{"!=", 2},
		{"<=", 2},
		{">=", 2},

		{"=~", 2},
		{"!~", 2},

		{"+", 1},
		{"-", 1},
		{"*", 1},
		{"/", 1},
		{"%", 1},
		{"(", 1},
		{")", 1},
		{"[", 1},
		{"]", 1},
		{"=", 1},
		{"<", 1},
		{">", 1},
		{"!", 1},
		{"&", 1},
		{"|", 1},
		{".", 1},
		{",", 1},
		{"^", 1},
		{"?", 1},
		{":", 1},
		{";", 1},
		{"{", 1},
		{"}", 1},
	}

	if h.peek(0) == '/' {
		if h.peek(1) == '/' {
			h.processComment(h.peek(0))
			return
		}
	}

	for i := 0; i < len(opArr); i++ {
		aLen := opArr[i].opLen

		str := string(h.input[h.pos : h.pos+aLen])

		if strings.HasPrefix(str, opArr[i].operator) {
			h.next()
			if aLen == 2 {
				h.next()
			}

			for _, intf := range h.generator {
				str := intf.WriteOperator(opArr[i].operator)
				if len(str) > 0 {
					io.WriteString(intf.Writer(), str)
				}
			}
			break
		} //end if
	} //end for

}

func (h *Highlighter) processIdentifier() {
	pos := h.pos
	for isLetter(h.input[h.pos]) || isDigit(h.input[h.pos]) {
		h.next()
	}

	text := string(h.input[pos:h.pos])

	if _, ok := keywords[text]; ok {
		for _, intf := range h.generator {
			str := intf.WriteKeyword(text)
			if len(str) > 0 {
				io.WriteString(intf.Writer(), str)
			}
		}
	} else {
		for _, intf := range h.generator {
			str := intf.WriteNormal(text)
			if len(str) > 0 {
				io.WriteString(intf.Writer(), str)
			}
		}
	}
}

func (h *Highlighter) processNumber() error {
	var ret []rune
	var ch rune = h.input[h.pos]

	ret = append(ret, ch)
	h.next()

	if ch == '0' && (h.input[h.pos] == 'x' || h.input[h.pos] == 'b' || h.input[h.pos] == 'c') { //support '0x'(hex) and '0b'(bin) and '0c'(octal)
		savedCh := h.input[h.pos]
		ret = append(ret, h.input[h.pos])
		h.next()

		if savedCh == 'x' {
			for isHex(h.input[h.pos]) || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					h.next()
					continue
				}
				ret = append(ret, h.input[h.pos])
				h.next()
			}
		} else if savedCh == 'b' {
			for isBin(h.input[h.pos]) || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					h.next()
					continue
				}
				ret = append(ret, h.input[h.pos])
				h.next()
			}
		} else if savedCh == 'c' {
			for isOct(h.input[h.pos]) || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					h.next()
					continue
				}
				ret = append(ret, h.input[h.pos])
				h.next()
			}
		}
	} else {
		for isDigit(h.input[h.pos]) || h.input[h.pos] == '.' || h.input[h.pos] == '_' {
			if h.input[h.pos] == '_' {
				h.next()
				continue
			}

			if h.input[h.pos] == '.' {
				if h.peek(1) == '.' { //range operator
					for _, intf := range h.generator {
						str := intf.WriteNumber(string(ret))
						io.WriteString(intf.Writer(), str)
					}
					return nil
				} else if !isDigit(h.peek(1)) && h.peek(1) != 'e' && h.peek(1) != 'E' { //should be a method calling, e.g. 10.next()
					for _, intf := range h.generator {
						str := intf.WriteNumber(string(ret))
						io.WriteString(intf.Writer(), str)
					}
					return nil
				}
			} //end if

			ret = append(ret, h.input[h.pos])
			h.next()
		}

		if h.input[h.pos] == 'e' || h.input[h.pos] == 'E' {
			ret = append(ret, h.input[h.pos])
			h.next()
			if isDigit(h.input[h.pos]) || h.input[h.pos] == '+' || h.input[h.pos] == '-' {
				ret = append(ret, h.input[h.pos])
				h.next()
				for isDigit(h.input[h.pos]) || h.input[h.pos] == '.' || h.input[h.pos] == '_' {
					if h.input[h.pos] == '_' {
						h.next()
						continue
					}
					ret = append(ret, h.input[h.pos])
					h.next()
				}
			}
			for isDigit(h.input[h.pos]) || h.input[h.pos] == '.' || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					h.next()
					continue
				}
				ret = append(ret, h.input[h.pos])
				h.next()
			}
		}
//		if isLetter(h.input[h.pos]) {
//			return errors.New("identifier starts immediately after numeric literal")
//		}
	}

	for _, intf := range h.generator {
		str := intf.WriteNumber(string(ret))
		io.WriteString(intf.Writer(), str)
	}

	return nil
}

func (h *Highlighter) processNormal() {
	for _, intf := range h.generator {
		str := intf.WriteNormal(string(h.input[h.pos]))
		io.WriteString(intf.Writer(), str)
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || unicode.IsDigit(ch)
}

// isHex returns true if the rune is a hex digits.
func isHex(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

// isBin returns true if the rune is a binary digits.
func isBin(ch rune) bool {
	return ('0' == ch || '1' == ch)
}

// isOct returns true if the rune is a octal digits.
func isOct(ch rune) bool {
	return ('0' <= ch && ch <= '7')
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$' || unicode.IsLetter(ch)
}

func (h *Highlighter) next() rune {
	h.pos++
	return h.peek(0)
}

func (h *Highlighter) peek(relativePos int) rune {
	position := h.pos + relativePos
	if position >= h.length {
		return 0
	}

	return h.input[position]
}
