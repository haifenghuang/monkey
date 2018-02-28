// Package doc extracts source code documentation from a Monkey AST.
package doc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"monkey/ast"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

var (
	regexpType = regexp.MustCompile(`^\{(.+)\}$`)
	regExample = regexp.MustCompile(`@example([^@]+)@[\r\n]`)

	//table of contents
	toc = `<p><div>
<a id="toc-button" class="toc-button" onclick="toggle_toc()"><span id="btn-text">&#x25BC;</span>&nbsp;Table of Contents</a>
<div id="table-of-contents" style="display:none;">`
	//PlaceHolder line, used only in html output.
	PlaceHolder = "<p>__TOC_PLACEHOLDER_LINE__</p>"
)

// File is the documentation for an entire monkey file.
type File struct {
	Name    string //FileName
	Classes []*Classes
	Enums   []*Value
	Lets    []*Value
	Funcs   []*Function
}

/* Classes is the documention for a class */
type Classes struct {
	Value *Value
	Props []*Value     //Properties
	Lets  []*Value     //Let-statements
	Funcs []*Function //Function
}

/* Classes is the documention for a function */
type Function struct {
	Value *Value
	Params []*FuncInfo
	Returns[]*FuncInfo
}

//function information(@param/@return/@returns part)
type FuncInfo struct {
	Name string //parameter name if @param, or else ""
	Type string //type
	Desc string //description
}

//Value is the documentation for a (possibly grouped) enums, lets, functions, or class declaration.
type Value struct {
	Name string //name
	Doc  string //comment
	Text string //declaration text
}

// Request for github REST API
// URL : https://developer.github.com/v3/markdown/
type Request struct {
	Text    string `json:"text"`
	Mode    string `json:"mode"`
	Context string `json:"context"`
}

func New(name string, program *ast.Program) *File {
	var classes []*ast.ClassStatement
	var enums   []*ast.EnumStatement
	var lets    []*ast.LetStatement
	var funcs   []*ast.FunctionStatement

	for _, statement := range program.Statements {
		switch s := statement.(type) {
		case *ast.ClassStatement:
			if s.Doc != nil {
				classes = append(classes, s)
			}
		case *ast.EnumStatement:
			if s.Doc != nil {
				enums = append(enums, s)
			}
		case *ast.LetStatement:
			if s.Doc != nil {
				lets = append(lets, s)
			}
		case *ast.FunctionStatement:
			if s.Doc != nil {
				funcs = append(funcs, s)
			}
		}
	}

	return &File{
		Name:    filepath.Base(name),
		Classes: sortedClasses(classes),
		Enums:   sortedEnums(enums),
		Lets:    sortedLets(lets),
		Funcs:   sortedFuncs(funcs),
	}
}

// ----------------------------------------------------------------------------
// Markdown document generator

// MdDocGen generates markdown documentation from doc.File.
func MdDocGen(f *File) string {
	var buffer bytes.Buffer
	tmpl, _ := template.New("baseTpl").Funcs(funcs).Parse(templs[0])
	for _, templ := range templs[1:] {
		tmpl, _ = template.Must(tmpl.Clone()).Parse(templ)
	}
	tmpl.Execute(&buffer, f)
	return normalize(buffer.String())
}

func normalize(doc string) string {
	nlReplace := regexp.MustCompile(`\n(\s)+\n`)
	trimCodes := regexp.MustCompile("\n{2,}```")
	doc = nlReplace.ReplaceAllString(doc, "\n\n")
	doc = trimCodes.ReplaceAllString(doc, "\n```")

	return doc
}

// ----------------------------------------------------------------------------
// Html document generator(using github REST API)

// HtmlDocGen generates html documentation from a markdown file.
func HtmlDocGen(content string, file *File) string {
	buf, err := json.Marshal(Request{
		Text:string(content),
		Mode: "gfm",
		Context: "github/gollum",
	})
	if err != nil {
		fmt.Errorf("Marshaling request failed, reason=%v\n", err)
		return ""
	}

	resp, err := http.Post("https://api.github.com/markdown","application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Errorf("Request for github failed, reason:%v\n", err)
		return ""
	}
	defer resp.Body.Close() //must close the 'Body'

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Response read failed, reason:%v\n", err)
		return ""
	}

	var out bytes.Buffer
	//doc type
	out.WriteString("<!DOCTYPE html>")
	//head
	out.WriteString("<head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">")
	out.WriteString(`<script type="text/javascript">
function toggle_toc() {
    var toc=document.getElementById('table-of-contents');
    var btn=document.getElementById('btn-text');
    toc.style.display=(toc.style.display=='none')?'block':'none';
    btn.innerHTML=(toc.style.display=='none')?'&#x25BC;':'&#x25B2;';
}
</script>`)
	out.WriteString("</head>")
	//css style
	out.WriteString("<style>")
	out.WriteString(css)
	out.WriteString("</style>")
	//body
	out.WriteString(`<body><div class="readme"><article class="markdown-body">`)
	out.WriteString(string(body))
	out.WriteString("</article></div></body>")

	html := out.String()
	//The github returned html's inner linking is not working,
	//so we need to fix this.
	for _, enum := range file.Enums {
		enumName := enum.Name
		src  := fmt.Sprintf("<h3>%s</h3>", enumName)
		dest := fmt.Sprintf(`<h3 id="%s">%s</h3>`, SanitizedAnchorName(enumName), enumName)
		html = strings.Replace(html, src, dest, -1)
	}
	for _, let := range file.Lets {
		letName := let.Name
		src  := fmt.Sprintf("<h3>%s</h3>", letName)
		dest := fmt.Sprintf(`<h3 id="%s">%s</h3>`, SanitizedAnchorName(letName), letName)
		html = strings.Replace(html, src, dest, -1)
	}
	for _, fn := range file.Funcs {
		fnName := fn.Value.Name
		src  := fmt.Sprintf("<h3>%s</h3>", fnName)
		dest := fmt.Sprintf(`<h3 id="%s">%s</h3>`, SanitizedAnchorName(fnName), fnName)
		html = strings.Replace(html, src, dest, -1)
	}

	for _, cls := range file.Classes {
		clsName := cls.Value.Name
		src  := fmt.Sprintf("<h3>%s</h3>", clsName)
		dest := fmt.Sprintf(`<h3 id="%s">%s</h3>`, SanitizedAnchorName(clsName), clsName)
		html = strings.Replace(html, src, dest, -1)

		for _, prop := range cls.Props {
			propName := prop.Name
			src  := fmt.Sprintf("<h5>%s</h5>", propName)
			dest := fmt.Sprintf(`<h5 id="%s">%s</h5>`, SanitizedAnchorName(propName), propName)
			html = strings.Replace(html, src, dest, -1)
		}
		for _, let := range cls.Lets {
			letName := let.Name
			src  := fmt.Sprintf("<h5>%s</h5>", letName)
			dest := fmt.Sprintf(`<h5 id="%s">%s</h5>`, SanitizedAnchorName(letName), letName)
			html = strings.Replace(html, src, dest, -1)
		}
		for _, fn := range cls.Funcs {
			fnName := fn.Value.Name
			src  := fmt.Sprintf("<h5>%s</h5>", fnName)
			dest := fmt.Sprintf(`<h5 id="%s">%s</h5>`, SanitizedAnchorName(fnName), fnName)
			html = strings.Replace(html, src, dest, -1)
		}
	}

	html = strings.Replace(html, "<h1>Table of Contents</h1>", toc, 1)
	html = strings.Replace(html, PlaceHolder, "</div>", 1)
	return html
}

// ----------------------------------------------------------------------------
// Sorting

type data struct {
	n    int
	swap func(i, j int)
	less func(i, j int) bool
}

func (d *data) Len() int           { return d.n }
func (d *data) Swap(i, j int)      { d.swap(i, j) }
func (d *data) Less(i, j int) bool { return d.less(i, j) }

// sortBy is a helper function for sorting
func sortBy(less func(i, j int) bool, swap func(i, j int), n int) {
	sort.Sort(&data{n, swap, less})
}

func sortedClasses(classes []*ast.ClassStatement) []*Classes {
	list := make([]*Classes, len(classes))
	i := 0
	for _, c := range classes {

		funcs := make([]*ast.FunctionStatement, 0)
		for _, fn := range c.ClassLiteral.Methods {
			if fn.Doc != nil {
				funcs = append(funcs, fn)
			}
		}

		props := make([]*ast.PropertyDeclStmt, 0)
		for _, prop := range c.ClassLiteral.Properties {
			if prop.Doc != nil {
				props = append(props, prop)
			}
		}

		lets := make([]*ast.LetStatement, 0)
		for _, member := range c.ClassLiteral.Members {
			if member.Doc != nil {
				lets = append(lets, member)
			}
		}

		list[i] = &Classes{
			Value: &Value{
				Name: c.Name.Value,
				Doc:  preProcessCommentExamples(c.Doc.Text()),
				Text: c.Docs(),
			},
			Props: sortedProps(props),
			Lets:  sortedLets(lets),
			Funcs: sortedFuncs(funcs),
		}
		i++
	}

	sortBy(
		func(i, j int) bool { return list[i].Value.Name < list[j].Value.Name },
		func(i, j int) { 
			list[i].Value, list[j].Value = list[j].Value, list[i].Value
			list[i].Props, list[j].Props = list[j].Props, list[i].Props
			list[i].Lets, list[j].Lets = list[j].Lets, list[i].Lets
			list[i].Funcs, list[j].Funcs = list[j].Funcs, list[i].Funcs
		},
		len(list),
	)
	return list
}

func sortedLets(lets []*ast.LetStatement) []*Value {
	list := make([]*Value, len(lets))
	i := 0
	for _, l := range lets {
		list[i] = &Value{
			Name: l.Names[0].Value,
			Doc:  preProcessCommentExamples(l.Doc.Text()),
			Text: l.Docs(),
		}
		i++
	}

	sortBy(
		func(i, j int) bool { return list[i].Name < list[j].Name },
		func(i, j int) { list[i], list[j] = list[j], list[i] },
		len(list),
	)
	return list
}

func sortedEnums(enums []*ast.EnumStatement) []*Value {
	list := make([]*Value, len(enums))
	i := 0
	for _, e := range enums {
		list[i] = &Value{
			Name: e.Name.Value,
			Doc:  preProcessCommentExamples(e.Doc.Text()),
			Text: e.Docs(),
		}
		i++
	}

	sortBy(
		func(i, j int) bool { return list[i].Name < list[j].Name },
		func(i, j int) { list[i], list[j] = list[j], list[i] },
		len(list),
	)
	return list
}

func sortedFuncs(funcs []*ast.FunctionStatement) []*Function {
	list := make([]*Function, len(funcs))
	i := 0
	for _, f := range funcs {
		list[i]= parseFuncComment(f.Name.Value, preProcessCommentExamples(f.Doc.Text()), f.Docs())
		i++
	}

	sortBy(
		func(i, j int) bool { return list[i].Value.Name < list[j].Value.Name },
		func(i, j int) {
			list[i].Value, list[j].Value = list[j].Value, list[i].Value
			list[i].Params, list[j].Params = list[j].Params, list[i].Params
			list[i].Returns, list[j].Returns = list[j].Returns, list[i].Returns
		},
		len(list),
	)
	return list
}

func sortedProps(props []*ast.PropertyDeclStmt) []*Value {
	list := make([]*Value, len(props))
	i := 0
	for _, p := range props {
		list[i] = &Value{
			Name: p.Name.Value,
			Doc:  preProcessCommentExamples(p.Doc.Text()),
			Text: p.Docs(),
		}

		if strings.HasPrefix(p.Name.Value, "this") {
			list[i].Name = "this"
		} else {
			list[i].Name = p.Name.Value
		}
		i++
	}

	sortBy(
		func(i, j int) bool { return list[i].Name < list[j].Name },
		func(i, j int) { list[i], list[j] = list[j], list[i] },
		len(list),
	)
	return list
}

func parseFuncComment(name string, docComments string, text string) (*Function){
	fn := &Function{
		Value:&Value{
			Name: name,
			Text: text,
		},
		Params : make([]*FuncInfo, 0),
		Returns: make([]*FuncInfo, 0),
	}

	var buffer bytes.Buffer
	comments := strings.Split(docComments, "\n")
	for _, comment := range comments {
		if len(comment) > 0 && comment[0] == '@' {
			splitOnSpaces := strings.Split(comment, " ")
			label := splitOnSpaces[0]
			switch label {
			case "@param":
				funcParam := parseValue(splitOnSpaces[1:])
				fn.Params = append(fn.Params, funcParam)
			case "@return", "@returns":
				funcReturn := parseValue(splitOnSpaces[1:])
				fn.Returns = append(fn.Returns, funcReturn)
			}
		} else {
			buffer.WriteString(comment+"\n")
		}
	}
	fn.Value.Doc = buffer.String()

	return fn
}

func parseValue(splitOnSpaces []string) *FuncInfo {
	name  := ""
	types := ""
	var description bytes.Buffer

	description.WriteString("")
	ret := &FuncInfo{Name:"", Type:"", Desc:""}
	for _, item := range splitOnSpaces {
		if m := regexpType.FindStringSubmatch(item); m != nil {
			types = m[1]
		} else if len(name) == 0 {
			if len(item) > 0 && item[0] == '`' {
				name = item[1:len(item)-1]
			} else {
				name = item
			}
		} else {
			description.WriteString(item + " ")
		}
	}

	if (len(name) > 0) { ret.Name = name }
	if (len(types) > 0) { ret.Type = types }
	ret.Desc = description.String()

	return ret
}

// SanitizedAnchorName returns a sanitized anchor name for the given text.
//copied from 'Blackfriday': a markdown processor for Go.
func SanitizedAnchorName(text string) string {
	var anchorName []rune
	futureDash := false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}

/* Process `@example block @` part, and replace this with 
```swift
    block
```
*/
func preProcessCommentExamples(comments string) string {
	retStr := comments
	if m := regExample.FindAllStringSubmatch(comments, -1); m != nil {
		for _, match := range m {
			var buffer bytes.Buffer
			buffer.WriteString("\n```swift")
			buffer.WriteString(match[1])
			buffer.WriteString("```\n")

			retStr = replaceFirstString(regExample, retStr, buffer.String())
		}
		//fmt.Printf("retStr=<%s>\n", retStr) //debugging
	}
	return retStr
}

func replaceFirstString(re *regexp.Regexp, srcStr, replStr string) string {
	src  := []byte(srcStr)
	repl := []byte(replStr)

	if m := re.FindSubmatchIndex(src); m != nil {
		out := make([]byte, m[0])
		copy(out, src[0:m[0]])
		out = re.Expand(out, repl, src, m)
		if m[1] < len(src) {
			out = append(out, src[m[1]:]...)
		}
		return string(out)
	}
	out := make([]byte, len(src))
	copy(out, src)
	return string(out)
}