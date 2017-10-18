package formatter

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

type (
	operatorMode int //operator mode

	Formatter struct {
		input         []rune //the source we need to beautifier
		formattedCode bytes.Buffer
		pos           int
		length        int
	}
)

const (
	op_chars = "+-*/%()[]=<>!&|.,^?:;{}"

	opMode_SPACES    operatorMode = iota //add space to both
	opMode_RSPACES                       //add right space
	opMode_TRIM                          //trim both space
	opMode_RTRIM                         //trim right space
	opMode_AS_SOURCE                     //same as source
)

//New create a new Formatter
func New(input string) *Formatter {
	f := &Formatter{input: []rune(input)}
	f.length = len(f.input)
	return f
}

//Format the source code
func (f *Formatter) Format() {
	for f.pos < f.length {
		current := f.peek(0)
		if current == '"' || current == '`' || current == '~' || current == '\'' {
			f.processText(current)
		} else if current == '#' {
			f.processComment(current)
		} else if strings.Contains(op_chars, string(current)) {
			f.processOperator()
		} else {
			f.formattedCode.WriteRune(f.input[f.pos])
			f.next()
		}
	} //end for

	f.formattedCode.WriteTo(os.Stdout)
	//Note: We could not use `String()` method, it will treat '%' as format. so when we print `{MISSING}`
	//this is not we wanted
	//return f.formattedCode.String()
}

func (f *Formatter) processText(ch rune) error {
	var ret []rune
	ret = append(ret, ch)

	for {
		f.next()
		if f.input[f.pos] == 0 {
			return errors.New("unexpected EOF")
		}

		if f.input[f.pos] == ch {
			f.next()
			break
		}
		ret = append(ret, f.input[f.pos])
	}

	ret = append(ret, ch)
	f.formattedCode.WriteString(string(ret))

	return nil
}

func (f *Formatter) processComment(ch rune) {
	if ch == '#' {
		f.formattedCode.WriteRune(ch)
		f.next()
	} else {
		f.formattedCode.WriteRune(ch)
		f.formattedCode.WriteRune(ch)
		f.next()
		f.next()
	}

	for f.peek(0) != '\n' && f.peek(0) != 0 {
		f.formattedCode.WriteRune(f.input[f.pos])
		f.next()
	}
}

func (f *Formatter) processOperator() {
	opArr := []struct {
		operator string
		opLen    int
		opMode   operatorMode
	}{
		{"+=", 2, opMode_SPACES},
		{"-=", 2, opMode_SPACES},
		{"*=", 2, opMode_SPACES},
		{"/=", 2, opMode_SPACES},
		{"%=", 2, opMode_SPACES},
		{"^=", 2, opMode_SPACES},

		{"++", 2, opMode_AS_SOURCE},
		{"--", 2, opMode_AS_SOURCE},

		{"&&", 2, opMode_SPACES},
		{"||", 2, opMode_SPACES},

		{"<<", 2, opMode_SPACES},
		{">>", 2, opMode_SPACES},

		{"->", 2, opMode_SPACES},
		{"=>", 2, opMode_SPACES},

		{"==", 2, opMode_SPACES},
		{"!=", 2, opMode_SPACES},
		{"<=", 2, opMode_SPACES},
		{">=", 2, opMode_SPACES},

		{"=~", 2, opMode_SPACES},
		{"!~", 2, opMode_SPACES},

		{"+", 1, opMode_SPACES},
		{"-", 1, opMode_SPACES},
		{"*", 1, opMode_SPACES},
		{"/", 1, opMode_SPACES},
		{"%", 1, opMode_SPACES},
		{"(", 1, opMode_AS_SOURCE},
		{")", 1, opMode_AS_SOURCE},
		{"[", 1, opMode_AS_SOURCE},
		{"]", 1, opMode_AS_SOURCE},
		{"=", 1, opMode_SPACES},
		{"<", 1, opMode_SPACES},
		{">", 1, opMode_SPACES},
		{"!", 1, opMode_RTRIM},
		{"&", 1, opMode_SPACES},
		{"|", 1, opMode_SPACES},
		{".", 1, opMode_TRIM},
		{",", 1, opMode_RSPACES},
		{"^", 1, opMode_SPACES},
		{"?", 1, opMode_SPACES},
		{":", 1, opMode_SPACES},
		{";", 1, opMode_RSPACES},
		{"{", 1, opMode_SPACES},
		{"}", 1, opMode_SPACES},
	}

	if f.peek(0) == '/' {
		if f.peek(1) == '/' {
			f.processComment(f.peek(0))
			return
		}
	}

	for i := 0; i < len(opArr); i++ {
		aLen := opArr[i].opLen

		str := string(f.input[f.pos : f.pos+aLen])
		if strings.HasPrefix(str, opArr[i].operator) {
			f.next()
			if aLen == 2 {
				f.next()
			}

			switch opArr[i].opMode {
			case opMode_SPACES:
				f.formatOperator(opArr[i].operator, false, false)
			case opMode_RSPACES:
				f.formatOperator(opArr[i].operator, true, false)
			case opMode_TRIM:
				f.formatOperator(opArr[i].operator, true, true)
			case opMode_RTRIM:
				f.formatOperator(opArr[i].operator, false, true)
			case opMode_AS_SOURCE:
				f.formattedCode.WriteString(opArr[i].operator)
			} //end switch

			break
		} //end if
	} //end for
}

func (f *Formatter) formatOperator(operator string, leftTrim bool, rightTrim bool) {
	tmpBytes := f.formattedCode.Bytes()
	tmpBytes = bytes.TrimRight(tmpBytes, " ")
	f.formattedCode.Reset()
	f.formattedCode.Write(tmpBytes)

	if !leftTrim {
		f.formattedCode.WriteRune(' ')
	}

	for i := 0; i < len(operator); i++ {
		f.formattedCode.WriteRune(rune(operator[i]))
	}

	f.skipWS()
	if !rightTrim {
		f.formattedCode.WriteRune(' ')
	}
}

func (f *Formatter) next() rune {
	f.pos++
	return f.peek(0)
}

func (f *Formatter) peek(relativePos int) rune {
	position := f.pos + relativePos
	if position >= f.length {
		return 0
	}

	return f.input[position]
}

func (f *Formatter) skipWS() {
	if f.pos >= f.length {
		return
	}
	for f.input[f.pos] == ' ' || f.input[f.pos] == '\t' || f.input[f.pos] == '\r' {
		f.next()
		if f.pos >= f.length {
			break
		}
	}
}
