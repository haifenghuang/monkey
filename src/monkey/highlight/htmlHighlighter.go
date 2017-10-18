package highlight

import (
	"io"
	"strconv"
	"strings"
)

type HtmlHighlighter struct {
	Out io.Writer
}

func NewHtmlHighlighter(writer io.Writer) *HtmlHighlighter {
	return &HtmlHighlighter{Out: writer}
}

func (hl *HtmlHighlighter) Name() string {
	return "Html"
}

func (hl *HtmlHighlighter) Writer() io.Writer {
	return hl.Out
}

func (hl *HtmlHighlighter) WriteQuotes(quotes string) string {
	quotes = escape(quotes)
	return `<span style="color:#032F62">` + quotes + `</span>`
}

func (hl *HtmlHighlighter) WriteComment(comment string) string {
	comment = escape(comment)
	return `<span style="color:#7A737D">` + comment + `</span>`
}

func (hl *HtmlHighlighter) WriteKeyword(keyword string) string {
	return `<span style="color:#D73A49">` + keyword + `</span>`
}

func (hl *HtmlHighlighter) WriteOperator(operator string) string {
	return `<span style="color:#3D2F62">` + operator + `</span>`
}

func (hl *HtmlHighlighter) WriteNumber(number string) string {
	return `<span style="color:#3D2F62">` + number + `</span>`
}

func (hl *HtmlHighlighter) WriteNormal(text string) string {
	text = escape(text)
	return `<span style="color:#000000">` + text + `</span>`

}

func (hl *HtmlHighlighter) WriteHeader() string {
	return `
<html xmlns="http://www.w3.org/1999/xhtml">
    <head>
	 <meta http-equiv="content-type" content="text/html;charset=utf-8">
        <style>
            <!--
.lineNumber {
    font-size:10.0pt;
    font-family:"Consolas","sans-serif";
    text-align: right;
    background-color: Beige; 
    color: darkgray;
    width: 20pt;
}
.code {    
    font-size:10.0pt;
    font-family:"Consolas","sans-serif";
    background-color: #ffffff; 
}
.code td { border-bottom:1px dotted #BDB76B; }
-->
        </style>
    </head>
    <body bgcolor="white" lang="EN-US" link="blue" vlink="purple">
        <table class="code" style="width:100%;cellpadding="0"; cellspacing="0">`
}

func (hl *HtmlHighlighter) WriteFooter() string {
	return `
</table>

    </body>
</html>`
}

func (hl *HtmlHighlighter) WriteLineHead(lineNo int) string {
	lineNumber := strconv.Itoa(lineNo)
	return `<tr><td class="lineNumber">&nbsp;` + lineNumber + `&nbsp;</td><td>`
}

func (hl *HtmlHighlighter) WriteLineTail() string {
	return `</td></tr>`
}

func (hl *HtmlHighlighter) WriteNewLine() string {
	return `<span style="color:#008000">&nbsp;</span>`
}

func escape(text string) string {
	text = strings.Replace(text, " ", "&nbsp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)

	return text
}
