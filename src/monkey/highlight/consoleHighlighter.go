package highlight

import (
	"io"
	"os"
	"strconv"
)

const (
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

type ConsoleHighlighter struct {
}

func NewConsoleHighlighter() *ConsoleHighlighter {
	return &ConsoleHighlighter{}
}

func (hl *ConsoleHighlighter) Name() string {
	return "Console"
}

func (hl *ConsoleHighlighter) Writer() io.Writer {
	return os.Stdout
}

func (hl *ConsoleHighlighter) WriteQuotes(quotes string) string {
	return COLOR_CYAN + quotes + COLOR_RESET
}

func (hl *ConsoleHighlighter) WriteComment(comment string) string {
	return COLOR_GREEN + comment + COLOR_RESET
}

func (hl *ConsoleHighlighter) WriteKeyword(keyword string) string {
	return COLOR_BRIGHT + COLOR_MAGENTA + keyword + COLOR_RESET
}

func (hl *ConsoleHighlighter) WriteOperator(operator string) string {
	return COLOR_GREEN + operator + COLOR_RESET
}

func (hl *ConsoleHighlighter) WriteNumber(number string) string {
	return COLOR_YELLOW + number + COLOR_RESET
}

func (hl *ConsoleHighlighter) WriteNormal(text string) string {
	return COLOR_WHITE + text + COLOR_RESET

}

func (hl *ConsoleHighlighter) WriteHeader() string {
	return ""
}

func (hl *ConsoleHighlighter) WriteFooter() string {
	return ""
}

func (hl *ConsoleHighlighter) WriteLineHead(lineNo int) string {
	lineNumber := strconv.Itoa(lineNo)
	return COLOR_BRIGHT + COLOR_RED + lineNumber + COLOR_RESET + " "
}

func (hl *ConsoleHighlighter) WriteLineTail() string {
	return ""
}

func (hl *ConsoleHighlighter) WriteNewLine() string {
	return "\n"
}
