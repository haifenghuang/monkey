package doc

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var (
	nlToSpaces = regexp.MustCompile(`\n`)

	funcs = template.FuncMap{
		"inline": func(txt string) string {
			return fmt.Sprintf("`%s`", txt)
		},
		"codeBlock": func(lang, code string) string {
			return fmt.Sprintf("```%s\n%s\n```", lang, code)
		},
		"lower": strings.ToLower,
		"genDate": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	baseTpl = `# File {{inline .Name}}
Table of Contents
=================
{{block "index" .}}_no index_{{end}}
{{block "lets" .}}_no lets_{{end}}
{{block "enums" .}}_no enums_{{end}}
{{block "functions" .}}_no functions_{{end}}
{{block "classes" .}}_no classes_{{end}}
***
_Last updated {{genDate}}_`

	indexTpl = `{{define "index"}}

{{if gt (len .Lets) 0}}
* Lets{{range $idx, $let := .Lets}}
  * [{{$let.Name}}](#{{lower $let.Name}}){{end}}
{{end}}

{{if gt (len .Enums) 0}}
* Enums{{range $idx, $enum := .Enums}}
  * [{{$enum.Name}}](#{{lower $enum.Name}}){{end}}
{{end}}

{{if gt (len .Funcs) 0}}
* Functions{{range $idx, $fn := .Funcs}}
  * [{{$fn.Name}}](#{{lower $fn.Name}}){{end}}
{{end}}

{{if gt (len .Classes) 0}}
* Classes{{range $idx, $cls := .Classes}}
  * [{{$cls.Value.Name}}](#{{lower $cls.Value.Name}})
{{if gt (len $cls.Lets) 0}}
    * Lets{{range $idx, $let := $cls.Lets}}
      * [{{$let.Name}}](#{{lower $let.Name}}){{end}}
{{end}}
{{if gt (len $cls.Props) 0}}
    * Properties{{range $idx, $prop := $cls.Props}}
      * [{{$prop.Name}}](#{{lower $prop.Name}}){{end}}
{{end}}

{{if gt (len $cls.Funcs) 0}}
    * Functions{{range $idx, $func := $cls.Funcs}}
      * [{{$func.Name}}](#{{lower $func.Name}}){{end}}
{{end}}

{{end}}
{{end}}


{{end}}`

	letsTpl = `{{define "lets"}}
{{if gt (len .Lets) 0}}

## Lets
  {{range $idx, $let := .Lets}}
### {{$let.Name}}
{{codeBlock "monkey" $let.Text}}
{{$let.Doc}}
  {{end}}
{{end}}

{{end}}`

	enumsTpl = `{{define "enums"}}
{{if gt (len .Enums) 0}}

## Enums
  {{range $idx, $enum := .Enums}}
### {{$enum.Name}}
{{codeBlock "monkey" $enum.Text}}
{{$enum.Doc}}
  {{end}}
{{end}}

{{end}}`

	functionsTpl = `{{define "functions"}}
{{if gt (len .Funcs) 0}}

## Functions
  {{range $idx, $fn := .Funcs}}
### {{$fn.Name}}
{{codeBlock "monkey" $fn.Text}}
{{$fn.Doc}}

  {{end}}
{{end}}

{{end}}`

	classesTpl = `{{define "classes"}}
{{if gt (len .Classes) 0}}

## Classes
  {{range $idx, $cls := .Classes}}
### {{$cls.Value.Name}}
{{codeBlock "monkey" $cls.Value.Text}}
{{$cls.Value.Doc}}

{{if gt (len .Lets) 0}}

#### Lets
  {{range $idx, $let := .Lets}}
##### {{$let.Name}}
{{codeBlock "monkey" $let.Text}}
{{$let.Doc}}
  {{end}}
{{end}}

{{if gt (len .Props) 0}}

#### Properties
  {{range $idx, $prop := .Props}}
##### {{$prop.Name}}
{{codeBlock "monkey" $prop.Text}}
{{$prop.Doc}}
  {{end}}
{{end}}

{{if gt (len .Funcs) 0}}

#### Functions
  {{range $idx, $fn := .Funcs}}
##### {{$fn.Name}}
{{codeBlock "monkey" $fn.Text}}
{{$fn.Doc}}
  {{end}}
{{end}}

  {{end}}
  {{end}}
{{end}}`
	templs = []string{baseTpl, indexTpl, letsTpl, enumsTpl, functionsTpl, classesTpl}
)
