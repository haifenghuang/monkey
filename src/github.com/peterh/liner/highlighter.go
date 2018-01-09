package liner

import (
	"io"
	"os"
	"errors"
	"strings"
	"unicode"
)

const (
	COLOR_NOCOLOR = ""
	COLOR_RESET  = "\x1b[0m"
	COLOR_BRIGHT = "\x1b[1m"

	COLOR_BLACK   = "\x1b[30m"
	COLOR_RED     = "\x1b[31m"
	COLOR_GREEN   = "\x1b[32m"
	COLOR_YELLOW  = "\x1b[33m"
	COLOR_BLUE    = "\x1b[34m"
	COLOR_MAGENTA = "\x1b[35m"
	COLOR_CYAN    = "\x1b[36m"
	COLOR_WHITE   = "\x1b[37m"
)

type Category int

const (
	NumberType Category = iota //number. e.g. 10, 10.5
	IdentType                  //identifier. e.g. name, age
	KeywordType                //keywords.  e.g. if, else
	StringType                 //string. e.g. "hello", 'hello'
	CommentType                // comment. e.g. #xxxx
	OperatorType               // operators. e.g. ++, --, +-
)

type Highlighter struct {
	input  []rune //the source we need to highlight
	pos    int
	length int

	keywords      map[string]int       //keywords
	category      map[Category]string  //color categories

	//operators
	operatorChars string
	operatorArr   []string
}

func NewHighlighter() *Highlighter {
	h := &Highlighter{}
	h.pos = 0
	h.length = 0

	h.keywords = make(map[string]int)

	h.category = make(map[Category]string)
	//default values(NO COLOR)
	h.category[NumberType]   = COLOR_NOCOLOR
	h.category[IdentType]    = COLOR_NOCOLOR
	h.category[KeywordType]  = COLOR_NOCOLOR
	h.category[StringType]   = COLOR_NOCOLOR
	h.category[CommentType]  = COLOR_NOCOLOR
	h.category[OperatorType] = COLOR_NOCOLOR


	return h
}

//reset the highlighter for reuse.
func (h *Highlighter) Reset(input []rune) {
	h.input = input
	h.pos = 0
	h.length = len(h.input)
}

func (h *Highlighter) RegisterKeywords(keywords []string) {
	for _, v := range keywords {
		h.keywords[v] = 1 //map's value is not important
	}
}

func (h *Highlighter) RegisterOperators(operators []string) {
	for _, v := range operators {
		h.operatorArr = append(h.operatorArr, v)
	}
	h.operatorChars = strings.Join(operators,"")
}

func (h *Highlighter) RegisterColors(category map[Category]string) {
	h.category = category
}

//Highlight the source code
func (h *Highlighter) Highlight() {
	for h.pos < h.length {
		current := h.peek(0)
		if current == '"' || current == '`' || current == '~' || current == '\'' {
			h.processString(current)
		} else if current == '#' {
			h.processComment(current)
		} else if isLetter(current) {
			h.processIdentifier()
		} else if strings.Contains(h.operatorChars, string(current)) {
			h.processOperator()
		} else if isDigit(current) {
			h.processNumber()
		} else {
			if h.input[h.pos] == '\n' {
				//should never happen, because we only support single line.
			} else {
				h.processNormal()
			}
			h.next()
		}
	} //end for
}

//process strings(doule quoted string, single quoted string or raw string)
func (h *Highlighter) processString(ch rune) error {
	var ret []rune

	ret = append(ret, ch)

	for {
		if h.next() == 0 { goto end }
		if h.input[h.pos] == 0 {
			return errors.New("unexpected EOF")
		}

		if h.input[h.pos] == ch {
			if h.next() == 0 {
				ret = append(ret, ch)
				goto end
			}
			break
		}

		ret = append(ret, h.input[h.pos])
	}
	ret = append(ret, ch)

end:
	str := h.category[StringType] + string(ret)
	if h.category[StringType] != COLOR_NOCOLOR {
		str += COLOR_RESET
	}

	io.WriteString(os.Stdout, str)

	return nil
}

func (h *Highlighter) processComment(ch rune) {
	var ret []rune
	if ch == '#' {
		ret = append(ret, ch)
		if h.next() == 0 { goto end }
	} else {
		ret = append(ret, ch)
		ret = append(ret, ch)
		if h.next() == 0 { goto end }
		if h.next() == 0 { goto end }
	}

	for h.peek(0) != '\n' && h.peek(0) != 0 {
		ret = append(ret, h.input[h.pos])
		if h.next() == 0 { goto end }
	}

end:
	str := h.category[CommentType] + string(ret)
	if h.category[CommentType] != COLOR_NOCOLOR {
		str += COLOR_RESET
	}

	io.WriteString(os.Stdout, str)
}

//process operator
func (h *Highlighter) processOperator() {
	if h.peek(0) == '/' {
		if h.peek(1) == '/' {
			h.processComment(h.peek(0))
			return
		}
	}

	for _, operator := range(h.operatorArr) {
		aLen := len(operator)

		var str string
		if (h.pos + aLen < h.length) {
			str = string(h.input[h.pos : h.pos+aLen])
		} else {
			str = string(h.input[h.pos:])
		}

		if strings.HasPrefix(str, operator) {
			if h.next() == 0 { goto end }
			if aLen == 2 {
				if h.next() == 0 { goto end }
			}
end:
			strOut := h.category[OperatorType] + operator
			if h.category[OperatorType] != COLOR_NOCOLOR {
				strOut += COLOR_RESET
			}
			io.WriteString(os.Stdout, strOut)
			break
		} //end if
	} //end for
}

func (h *Highlighter) processIdentifier() {
	pos := h.pos
	for isLetter(h.input[h.pos]) || isDigit(h.input[h.pos]) {
		r := h.next()
		if r == 0 {
			break
		}
	}

	text := string(h.input[pos:h.pos])

	if _, ok := h.keywords[text]; ok {
		str := h.category[KeywordType] + text
		if h.category[KeywordType] != COLOR_NOCOLOR {
			str += COLOR_RESET
		}
		io.WriteString(os.Stdout, str)
	} else {
		str := h.category[IdentType] + text
		if h.category[IdentType] != COLOR_NOCOLOR {
			str += COLOR_RESET
		}
		io.WriteString(os.Stdout, str)
	}
}

func (h *Highlighter) processNumber() error {
	var ret []rune
	var ch rune = h.input[h.pos]

	ret = append(ret, ch)
	if h.next() == 0 { goto end	}

	if ch == '0' && (h.input[h.pos] == 'x' || h.input[h.pos] == 'b' || h.input[h.pos] == 'c') { //support '0x'(hex) and '0b'(bin) and '0c'(octal)
		savedCh := h.input[h.pos]
		ret = append(ret, h.input[h.pos])
		if h.next() == 0 { goto end }
		if savedCh == 'x' {
			for isHex(h.input[h.pos]) || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					ret = append(ret, h.input[h.pos])
					if h.next() == 0 { goto end	}
					continue
				}
				ret = append(ret, h.input[h.pos])
				if h.next() == 0 { goto end	}
			}
		} else if savedCh == 'b' {
			for isBin(h.input[h.pos]) || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					ret = append(ret, h.input[h.pos])
					if h.next() == 0 { goto end	}
					continue
				}
				ret = append(ret, h.input[h.pos])
				if h.next() == 0 { goto end	}
			}
		} else if savedCh == 'c' {
			for isOct(h.input[h.pos]) || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					ret = append(ret, h.input[h.pos])
					if h.next() == 0 { goto end	}
					continue
				}
				ret = append(ret, h.input[h.pos])
				if h.next() == 0 { goto end	}
			}
		}
	} else {
		for isDigit(h.input[h.pos]) || h.input[h.pos] == '.' || h.input[h.pos] == '_' {
			if h.input[h.pos] == '_' {
				ret = append(ret, h.input[h.pos])
				if h.next() == 0 { goto end	}
				continue
			}

			if h.input[h.pos] == '.' {
				if h.peek(1) == '.' { //range operator
					goto end
				} else if !isDigit(h.peek(1)) && h.peek(1) != 'e' && h.peek(1) != 'E' { //should be a method calling, e.g. 10.next()
					//Note: there are some limitation about this, i.e. the method name could not begin with 'e' or 'E'
					goto end
				}
			} //end if

			ret = append(ret, h.input[h.pos])
			if h.next() == 0 { goto end	}
		}

		if h.input[h.pos] == 'e' || h.input[h.pos] == 'E' {
			ret = append(ret, h.input[h.pos])
			if h.next() == 0 { goto end	}
			if isDigit(h.input[h.pos]) || h.input[h.pos] == '+' || h.input[h.pos] == '-' {
				ret = append(ret, h.input[h.pos])
				if h.next() == 0 { goto end	}
				for isDigit(h.input[h.pos]) || h.input[h.pos] == '.' || h.input[h.pos] == '_' {
					if h.input[h.pos] == '_' {
						ret = append(ret, h.input[h.pos])
						if h.next() == 0 { goto end	}
						continue
					}
					ret = append(ret, h.input[h.pos])
					if h.next() == 0 { goto end	}
				}
			}
			for isDigit(h.input[h.pos]) || h.input[h.pos] == '.' || h.input[h.pos] == '_' {
				if h.input[h.pos] == '_' {
					ret = append(ret, h.input[h.pos])
					if h.next() == 0 { goto end	}
					continue
				}
				ret = append(ret, h.input[h.pos])
				if h.next() == 0 { goto end	}
			}
		}
	}

end:
	str := h.category[NumberType] + string(ret)
	if h.category[NumberType] != COLOR_NOCOLOR {
		str += COLOR_RESET
	}

	io.WriteString(os.Stdout, str)

	return nil
}

func (h *Highlighter) processNormal() {
//	str := COLOR_WHITE + string(h.input[h.pos]) + COLOR_RESET
	io.WriteString(os.Stdout, string(h.input[h.pos]))
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || unicode.IsDigit(ch) || ch == '_'
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
