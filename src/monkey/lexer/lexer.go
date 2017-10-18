package lexer

import (
	"bytes"
	"errors"
	"monkey/token"
	"strings"
	"unicode"
	//	"fmt"
)

const bom = 0xFEFF // byte order mark, only permitted as very first character
var prevToken token.Token

type Lexer struct {
	filename     string
	input        []rune
	ch           rune //current character
	position     int  //character offset
	readPosition int  //reading offset

	lineOffset     int // current line offset
	tokStartOffset int // character offset of the start of the current token
	line           int
}

func New(filename, input string) *Lexer {
	l := &Lexer{filename: filename, input: []rune(input)}
	l.ch = ' '
	l.position = 0
	l.readPosition = 0
	l.lineOffset = 0
	l.tokStartOffset = 0
	l.line = 1

	l.readNext()
	if l.ch == bom {
		l.readNext() //ignore BOM at file beginning
	}

	return l
}

func (l *Lexer) readNext() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		if l.ch == '\n' {
			l.lineOffset = l.position + 1
			l.line++
		}
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peek() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekn(n int) rune {
	if l.readPosition+n >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+n]
}

var tokenMap = map[rune]token.TokenType{
	'=': token.ASSIGN,
	'.': token.DOT,
	';': token.SEMICOLON,
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	'[': token.LBRACKET,
	']': token.RBRACKET,
	'+': token.PLUS,
	',': token.COMMA,
	'-': token.MINUS,
	'!': token.BANG,
	'*': token.ASTERISK,
	'/': token.SLASH,
	'<': token.LT,
	'>': token.GT,
	':': token.COLON,
	'%': token.MOD,
	'#': token.COMMENT,
	'?': token.QUESTIONM,
	'&': token.BITAND,
	'|': token.BITOR,
	'^': token.BITXOR,
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	pos := l.getPos()
	l.skipWhitespace()
	l.tokStartOffset = l.position
	if t, ok := tokenMap[l.ch]; ok {
		switch t {
		case token.ASSIGN:
			if l.peek() == '=' {
				tok = token.Token{Type: token.EQ, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '~' {
				tok = token.Token{Type: token.MATCH, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '>' {
				tok = token.Token{Type: token.FATARROW, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.ASSIGN, l.ch)
			}
		case token.PLUS:
			if l.peek() == '+' {
				tok = token.Token{Type: token.INCREMENT, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '=' {
				tok = token.Token{Type: token.PLUS_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.PLUS, l.ch)
			}
		case token.MINUS:
			if l.peek() == '-' {
				tok = token.Token{Type: token.DECREMENT, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '=' {
				tok = token.Token{Type: token.MINUS_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.MINUS, l.ch)
			}
		case token.GT:
			if l.peek() == '=' {
				tok = token.Token{Type: token.GE, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '>' {
				tok = token.Token{Type: token.SHIFT_R, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.GT, l.ch)
			}
		case token.LT:
			if l.peek() == '=' {
				tok = token.Token{Type: token.LE, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '<' {
				tok = token.Token{Type: token.SHIFT_L, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.LT, l.ch)
			}
		case token.BANG:
			if l.peek() == '=' {
				tok = token.Token{Type: token.NEQ, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '~' {
				tok = token.Token{Type: token.NOTMATCH, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.BANG, l.ch)
			}
		case token.SLASH:
			if l.peek() == '/' {
				l.skipComment()
				return l.NextToken()
			} else {
				if prevToken.Type == token.RBRACE || // impossible?
					prevToken.Type == token.RPAREN || // (a+c) / b
					prevToken.Type == token.RBRACKET || // a[3] / b
					prevToken.Type == token.IDENT || // a / b
					prevToken.Type == token.INT || // 3 / b
					prevToken.Type == token.FLOAT { // 3.5 / b
					if l.peek() == '=' {
						tok = token.Token{Type: token.SLASH_A, Literal: string(l.ch) + string(l.peek())}
						l.readNext()
					} else {
						tok = newToken(token.SLASH, l.ch)
					}
				} else { //regexp
					tok.Literal = l.readRegExLiteral()
					tok.Type = token.REGEX
					tok.Pos = l.getPos()
					return tok
				}
			}
		case token.MOD:
			if l.peek() == '=' {
				tok = token.Token{Type: token.MOD_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.MOD, l.ch)
			}
		case token.ASTERISK:
			if l.peek() == '=' {
				tok = token.Token{Type: token.ASTERISK_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '*' {
				tok = token.Token{Type: token.POWER, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.ASTERISK, l.ch)
			}
		case token.DOT:
			if l.peek() == '.' {
				l.readNext()
				if l.peek() == '.' {
					tok = token.Token{Type: token.ELLIPSIS, Literal: "..."}
					l.readNext()
				} else {
					tok = token.Token{Type: token.DOTDOT, Literal: ".."}
				}
			} else {
				tok = newToken(token.DOT, l.ch)
			}
		case token.COMMENT:
			l.skipComment()
			return l.NextToken()
		case token.BITAND:
			if l.peek() == '=' {
				tok = token.Token{Type: token.BITAND_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '&' {
				tok = token.Token{Type: token.CONDAND, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.BITAND, l.ch)
			}
		case token.BITOR:
			if l.peek() == '=' {
				tok = token.Token{Type: token.BITOR_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '>' {
				tok = token.Token{Type: token.PIPE, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else if l.peek() == '|' {
				tok = token.Token{Type: token.CONDOR, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.BITOR, l.ch)
			}
		case token.BITXOR:
			if l.peek() == '=' {
				tok = token.Token{Type: token.BITXOR_A, Literal: string(l.ch) + string(l.peek())}
				l.readNext()
			} else {
				tok = newToken(token.BITXOR, l.ch)
			}
		default:
			tok = newToken(t, l.ch)
		}

		l.readNext()

		tok.Pos = pos
		prevToken = tok
		return tok
	}

	newTok := l.readRunesToken()
	newTok.Pos = pos
	prevToken = newTok
	return newTok
}

func (l *Lexer) readRunesToken() token.Token {
	var tok token.Token
	switch {
	case l.ch == 0:
		tok.Literal = ""
		tok.Type = token.EOF
		return tok
	case isLetter(l.ch):
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		return tok
	case isDigit(l.ch):
		literal, _ := l.readNumber()
		if strings.Contains(literal, ".") {
			tok.Type = token.FLOAT
		} else {
			tok.Type = token.INT
		}
		tok.Literal = literal
		return tok
	case isQuote(l.ch):
		if l.ch == 34 { //double quotes
			if s, err := l.readString(l.ch); err == nil {
				tok.Type = token.STRING
				tok.Literal = s
				return tok
			}
		} else if l.ch == 96 { //raw string
			if s, err := l.readRawString(); err == nil {
				tok.Type = token.STRING
				tok.Literal = s
				return tok
			}
		}
	case isSingleQuote(l.ch):
		if s, err := l.readInterpString(); err == nil {
			tok.Type = token.ISTRING
			tok.Literal = s
			return tok
		}
	}
	l.readNext()
	return newToken(token.ILLEGAL, l.ch)
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readRawString() (string, error) {
	var ret []rune
	for {
		l.readNext()
		if l.ch == 0 {
			return "", errors.New("unexpected EOF")
		}

		if l.ch == 96 {
			l.readNext()
			break
		}
		ret = append(ret, l.ch)
	}
	return string(ret), nil
}

func (l *Lexer) readString(r rune) (string, error) {
	var ret []rune
eos:
	for {
		l.readNext()
		switch l.ch {
		case '\n':
			return "", errors.New("unexpected EOL")
		case 0:
			return "", errors.New("unexpected EOF")
		case r:
			l.readNext()
			break eos //eos:end of string
		case '\\':
			l.readNext()
			switch l.ch {
			case 'b':
				ret = append(ret, '\b')
				continue
			case 'f':
				ret = append(ret, '\f')
				continue
			case 'r':
				ret = append(ret, '\r')
				continue
			case 'n':
				ret = append(ret, '\n')
				continue
			case 't':
				ret = append(ret, '\t')
				continue
			}
			ret = append(ret, l.ch)
			continue
		default:
			ret = append(ret, l.ch)
		}
	}

	return string(ret), nil
}

func (l *Lexer) readInterpString() (string, error) {
	start := l.position + 1
	var out bytes.Buffer
	pos := "0"[0]
	for {
		l.readNext()
		if isSingleQuote(l.ch) {
			l.readNext()
			break
		}
		if l.ch == 0 {
			err := errors.New("")
			return "", err
		}
		if l.ch == 123 {
			if l.peek() != 125 {
				out.WriteRune(l.ch)
				for l.ch != 125 || l.ch == 0 {
					l.readNext()
				}
				if l.ch != 0 {
					out.WriteRune(rune(pos))
					pos++
				}
			}
		}
		out.WriteRune(l.ch)
	}
	l.position = start - 1
	l.readPosition = start
	l.ch = l.input[start]
	return out.String(), nil
}

func (l *Lexer) NextInterpToken() token.Token {
	var tok token.Token
	for {
		if l.ch == '{' {
			if l.peek() == '}' {
				continue
			}
			tok = newToken(token.LBRACE, l.ch)
			break
		}
		if l.ch == 0 {
			tok.Type = token.EOF
			tok.Literal = ""
			break
		}
		if isSingleQuote(l.ch) {
			tok = newToken(token.ISTRING, l.ch)
			break
		}
		l.readNext()
	}
	l.readNext()
	tok.Pos = l.getPos()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readNext()
	}
	return string(l.input[position:l.position])

}

func (l *Lexer) readRegExLiteral() (literal string) {
	position := l.position
	/* read until closing slash */
	for {
		l.readNext()
		if l.ch == '\\' {
			// Skip escape sequence
			l.readNext()
		} else if l.ch == '/' {
			// This is the closing
			literal = string(l.input[position+1 : l.position])
			l.readNext() //skip the '/'

			return
		}
	}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$' || ch == '@' || unicode.IsLetter(ch)
}

// scanNumber returns number begining at current position.
func (l *Lexer) readNumber() (string, error) {
	var ret []rune
	ch := l.ch
	ret = append(ret, ch)
	l.readNext()

	if ch == '0' && (l.ch == 'x' || l.ch == 'b' || l.ch == 'c') { //support '0x'(hex) and '0b'(bin) and '0c'(octal)
		savedCh := l.ch
		ret = append(ret, l.ch)
		l.readNext()
		if savedCh == 'x' {
			for isHex(l.ch) || l.ch == '_' {
				if l.ch == '_' {
					l.readNext()
					continue
				}
				ret = append(ret, l.ch)
				l.readNext()
			}
		} else if savedCh == 'b' {
			for isBin(l.ch) || l.ch == '_' {
				if l.ch == '_' {
					l.readNext()
					continue
				}
				ret = append(ret, l.ch)
				l.readNext()
			}
		} else if savedCh == 'c' {
			for isOct(l.ch) || l.ch == '_' {
				if l.ch == '_' {
					l.readNext()
					continue
				}
				ret = append(ret, l.ch)
				l.readNext()
			}
		}
	} else {
		for isDigit(l.ch) || l.ch == '.' || l.ch == '_' {
			if l.ch == '_' {
				l.readNext()
				continue
			}

			if l.ch == '.' {
				if l.peek() == '.' { //range operator
					return string(ret), nil
				}
			} //end if

			ret = append(ret, l.ch)
			l.readNext()
		}

		if l.ch == 'e' || l.ch == 'E' {
			ret = append(ret, l.ch)
			l.readNext()
			if isDigit(l.ch) || l.ch == '+' || l.ch == '-' {
				ret = append(ret, l.ch)
				l.readNext()
				for isDigit(l.ch) || l.ch == '.' || l.ch == '_' {
					if l.ch == '_' {
						l.readNext()
						continue
					}
					ret = append(ret, l.ch)
					l.readNext()
				}
			}
			for isDigit(l.ch) || l.ch == '.' || l.ch == '_' {
				if l.ch == '_' {
					l.readNext()
					continue
				}
				ret = append(ret, l.ch)
				l.readNext()
			}
		}
		if isLetter(l.ch) {
			return "", errors.New("identifier starts immediately after numeric literal")
		}
	}
	return string(ret), nil
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

func isQuote(ch rune) bool {
	return ch == 34 || ch == 96
}

func isSingleQuote(ch rune) bool {
	return ch == 39
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readNext()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readNext()
	}
}

func (l *Lexer) getPos() token.Position {
	return token.Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Col:      l.tokStartOffset - l.lineOffset + 1,
		//Col:     l.lineOffset - l.tokStartOffset + 1,
	}
}
