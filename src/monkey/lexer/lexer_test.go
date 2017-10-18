package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten_dummy = 10;
	
	let add = fn(x, y) {
		x + y;
	};
	
	let result = add(five, ten);
	5 < 10 > 5;
	
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	5 == 5;
	5 != 6;
	5 >= 6;
	5 <= 6;
	"foobar";
	"foo bar";
	[];
	function.call
	{ "foo" => "bar" }
	[1:3]
	5 % 4
	include tests
	x and y
	x or y
	struct
	do
	if (~\d+(\w)+.*$~.exec("abc def") == 0) {  # this is just a comment
	    return "found"
	}
	# this is another command
	let a234 = ~[ab|cd].*\~efg$~
	let ww = 1.523 + 2    # test for floating point number
	for item in arr
	grep { $_ > 5 }
	if (abc =~ ~\d+~)
	y ? a : b
	52.9..80.7
	52..80
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten_dummy"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.NEQ, "!="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.GE, ">="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LE, "<="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foo bar"},
		{token.SEMICOLON, ";"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "function"},
		{token.DOT, "."},
		{token.IDENT, "call"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.FATARROW, "=>"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COLON, ":"},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.INT, "5"},
		{token.MOD, "%"},
		{token.INT, "4"},
		{token.INCLUDE, "include"},
		{token.IDENT, "tests"},
		{token.IDENT, "x"},
		{token.AND, "and"},
		{token.IDENT, "y"},
		{token.IDENT, "x"},
		{token.OR, "or"},
		{token.IDENT, "y"},
		{token.STRUCT, "struct"},
		{token.DO, "do"},

		//if (~\d+(\w)+.*$~.exec("abc def") == 0) {
		//    return "found"
		//}
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.REGEX, `\d+(\w)+.*$`},
		{token.DOT, "."},
		{token.IDENT, "exec"},
		{token.LPAREN, "("},
		{token.STRING, "abc def"},
		{token.RPAREN, ")"},
		{token.EQ, "=="},
		{token.INT, "0"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.STRING, "found"},
		{token.RBRACE, "}"},

		//let a234 = ~[ab|cd].*\~efg$~
		{token.LET, "let"},
		{token.IDENT, "a234"},
		{token.ASSIGN, "="},
		{token.REGEX, `[ab|cd].*\~efg$`},

		//let ww = 1.523 + 2
		{token.LET, "let"},
		{token.IDENT, "ww"},
		{token.ASSIGN, "="},
		{token.FLOAT, "1.523"},
		{token.PLUS, "+"},
		{token.INT, "2"},

		//for item in arr
		{token.FOR, "for"},
		{token.IDENT, "item"},
		{token.IN, "in"},
		{token.IDENT, "arr"},

		//grep { $_ > 5 }
		{token.GREP, "grep"},
		{token.LBRACE, "{"},
		{token.IDENT, "$_"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.RBRACE, "}"},

		//input := `if (abc =~ ~\d+~)`
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "abc"},
		{token.MATCH, "=~"},
		{token.REGEX, `\d+`},
		{token.RPAREN, ")"},

		//y ? a : b
		{token.IDENT, "y"},
		{token.QUESTIONM, "?"},
		{token.IDENT, "a"},
		{token.COLON, ":"},
		{token.IDENT, "b"},

		//52.9..80.7
		{token.FLOAT, "52.9"},
		{token.DOTDOT, ".."},
		{token.FLOAT, "80.7"},

		//52..80
		{token.INT, "52"},
		{token.DOTDOT, ".."},
		{token.INT, "80"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
