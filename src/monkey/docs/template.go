package doc

import (
	"fmt"
	"regexp"
	"text/template"
	"time"
)

var (
	nlToSpaces = regexp.MustCompile(`\n`)

	funcs = template.FuncMap{
		"inline": func(txt string) string {
			if len(txt) == 0 {
				return txt
			}
			return fmt.Sprintf("`%s`", txt)
		},
		"codeBlock": func(code string) string {
			return fmt.Sprintf("```swift\n%s\n```", code)
		},
		"sanitizedAnchorName": SanitizedAnchorName,
		"genDate": func() string {
			return time.Now().Format("2006-01-02")
		},
	}

	baseTpl = `# File {{inline .Name}}
Table of Contents
=================
{{block "index" .}}_no index_{{end}}
{{if eq 1 .GenHTML}}<p>__TOC_PLACEHOLDER_LINE_END__</p>{{end}}
{{block "lets" .}}_no lets_{{end}}
{{block "enums" .}}_no enums_{{end}}
{{block "functions" .}}_no functions_{{end}}
{{block "classes" .}}_no classes_{{end}}
***
_Last updated {{genDate}}_`

	indexTpl = `{{define "index"}}

{{if gt (len .Lets) 0}}
* Lets{{range $idx, $let := .Lets}}
  * [{{$let.Name}}](#{{sanitizedAnchorName $let.Name}}){{end}}
{{end}}

{{if gt (len .Enums) 0}}
* Enums{{range $idx, $enum := .Enums}}
  * [{{$enum.Name}}](#{{sanitizedAnchorName $enum.Name}}){{end}}
{{end}}

{{if gt (len .Funcs) 0}}
* Functions{{range $idx, $fn := .Funcs}}
  * [{{$fn.Value.Name}}](#{{sanitizedAnchorName $fn.Value.Name}}){{end}}
{{end}}

{{if gt (len .Classes) 0}}
* Classes{{range $idx, $cls := .Classes}}
  * [{{$cls.Value.Name}}](#{{sanitizedAnchorName $cls.Value.Name}})
{{if gt (len $cls.Lets) 0}}
    * Lets{{range $idx, $let := $cls.Lets}}
      * [{{$let.Name}}](#{{sanitizedAnchorName $let.Name}}){{end}}
{{end}}
{{if gt (len $cls.Props) 0}}
    * Properties{{range $idx, $prop := $cls.Props}}
      * [{{$prop.Name}}](#{{sanitizedAnchorName $prop.Name}}){{end}}
{{end}}

{{if gt (len $cls.Funcs) 0}}
    * Functions{{range $idx, $func := $cls.Funcs}}
      * [{{$func.Value.Name}}](#{{sanitizedAnchorName $func.Value.Name}}){{end}}
{{end}}

{{end}}
{{end}}


{{end}}`

	letsTpl = `{{define "lets"}}
{{if gt (len .Lets) 0}}

## Lets
  {{range $idx, $let := .Lets}}
### {{$let.Name}}
{{codeBlock $let.Src}}
{{if eq 1 .GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__{{$let.SrcLines}}</p>{{end}}
{{$let.Doc}}
  {{end}}
{{end}}

{{end}}`

	enumsTpl = `{{define "enums"}}
{{if gt (len .Enums) 0}}

## Enums
  {{range $idx, $enum := .Enums}}
### {{$enum.Name}}
{{codeBlock $enum.Src}}
{{if eq 1 .GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__{{$enum.SrcLines}}</p>{{end}}
{{$enum.Doc}}
  {{end}}
{{end}}

{{end}}`

	functionsTpl = `{{define "functions"}}
{{if gt (len .Funcs) 0}}

## Functions
  {{range $idx, $fn := .Funcs}}
### {{$fn.Value.Name}}
{{codeBlock $fn.Value.Text}}
{{if eq 1 $fn.Value.GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__1</p>{{end}}
{{$fn.Value.Doc}}

    {{if gt (len $fn.Params) 0}}
#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |{{range $idx, $param := $fn.Params}}
|{{$param.Name}}|{{inline $param.Type}}|{{$param.Desc}}|{{end}}
    {{end}}

    {{if gt (len $fn.Returns) 0}}
#### Returns
        {{range $idx, $ret := $fn.Returns}}
- {{inline $ret.Type}} {{$ret.Desc}}
        {{end}}
    {{end}}{{if eq .Value.ShowSrc 1}}
{{if eq 1 .Value.GenHTML}}SHOWSOURCE_PLACEHOLDER_LINE_BEGIN{{$fn.Value.Name}}{{else}}#### Source{{end}}
{{codeBlock $fn.Value.Src}}
{{if eq 1 .Value.GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__{{$fn.Value.SrcLines}}</p>{{end}}
{{if eq 1 .Value.GenHTML}}<p>__SHOWSOURCE_PLACEHOLDER_LINE_END__</p>{{end}}
{{end}}

  {{end}}
{{end}}

{{end}}`

	classesTpl = `{{define "classes"}}
{{if gt (len .Classes) 0}}

## Classes
  {{range $idx, $cls := .Classes}}
### {{$cls.Value.Name}}
{{codeBlock $cls.Value.Text}}
{{if eq 1 $cls.Value.GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__1</p>{{end}}
{{$cls.Value.Doc}}

{{if gt (len .Lets) 0}}

#### Lets
  {{range $idx, $let := .Lets}}
##### {{$let.Name}}
{{codeBlock $let.Text}}
{{if eq 1 $let.GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__1</p>{{end}}
{{$let.Doc}}
  {{end}}
{{end}}

{{if gt (len .Props) 0}}

#### Properties
  {{range $idx, $prop := .Props}}
##### {{$prop.Name}}
{{codeBlock $prop.Text}}
{{$prop.Doc}}
  {{end}}
{{end}}

{{if gt (len .Funcs) 0}}

#### Functions
  {{range $idx, $fn := .Funcs}}
##### {{$fn.Value.Name}}
{{codeBlock $fn.Value.Text}}
{{if eq 1 $fn.Value.GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__1</p>{{end}}
{{$fn.Value.Doc}}

    {{if gt (len $fn.Params) 0}}
#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |{{range $idx, $param := $fn.Params}}
|{{$param.Name}}|{{inline $param.Type}}|{{$param.Desc}}|{{end}}
    {{end}}

    {{if gt (len $fn.Returns) 0}}
#### Returns
        {{range $idx, $ret := $fn.Returns}}
- {{inline $ret.Type}} {{$ret.Desc}}
        {{end}}
    {{end}}
  {{end}}
{{end}}{{if eq .Value.ShowSrc 1}}
{{if eq 1 .Value.GenHTML}}SHOWSOURCE_PLACEHOLDER_LINE_BEGIN{{$cls.Value.Name}}{{else}}#### Source{{end}}
{{codeBlock $cls.Value.Src}}
{{if eq 1 .Value.GenHTML}}<p>__LINENUMBER_PLACEHOLDER_LINE__{{$cls.Value.SrcLines}}</p>{{end}}
{{if eq 1 .Value.GenHTML}}<p>__SHOWSOURCE_PLACEHOLDER_LINE_END__</p>{{end}}
  {{end}}
  {{end}}
  {{end}}
{{end}}`
	templs = []string{baseTpl, indexTpl, letsTpl, enumsTpl, functionsTpl, classesTpl}
)
