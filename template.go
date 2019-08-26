package main

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToTitle": func(s string) string {
		return strings.ToTitle(s[:1]) + s[1:]
	},
}

var codeTemplate = template.Must(template.New("code").Funcs(funcMap).Parse(codeTemplateText))

const codeTemplateText = `
// Code generated by github.com/launchdarkly/go-options.  DO NOT EDIT.

{{ if .imports }}
import (
{{ range .imports }}
{{ if .Alias }}  {{ .Alias }} "{{ .Path }}"{{ else }}  "{{ .Path }}"{{ end -}}
{{ end }}
)
{{ end }}

{{ range .options }}{{ if .DefaultValue -}}
const default{{ $.configTypeName | ToTitle }}{{ .Name | ToTitle }} {{ .Type }} = {{ .DefaultValue }}
{{ end }}{{ end }}

{{ $applyOptionFuncType := or $.applyOptionFuncType (printf "apply%sFunc" (ToTitle $.optionTypeName)) }} 

type {{ $applyOptionFuncType }} func(c *{{ $.configTypeName }}) error

func (f {{ $applyOptionFuncType }}) apply(c *{{ $.configTypeName }}) error {
	return f(c)
}

{{ $applyFuncName := or $.applyFuncName (printf "apply%sOptions" (ToTitle $.configTypeName)) }} 

{{ if $.createNewFunc}} 
func new{{ $.configTypeName | ToTitle}}(options ...{{ $.optionTypeName }}) ({{ $.configTypeName }}, error) {
	var c {{ $.configTypeName }}
	err := {{ $applyFuncName }}(&c, options...)
	return c, err
}
{{ end }}

func {{ $applyFuncName }}(c *{{ $.configTypeName }}, options ...{{ $.optionTypeName }}) error {
{{- range .options -}}{{ if .DefaultValue }}
  c.{{ .Name }} = default{{ $.configTypeName | ToTitle }}{{ .Name | ToTitle }}
{{- end }}{{ end }}
  for _, o := range options {
    if err := o.apply(c); err != nil {
      return err
    }
  }
  return nil
}

type {{ $.optionTypeName }} interface {
  apply(*{{ $.configTypeName }}) error
}

{{ range .options }}
{{ $name := .PublicName | ToTitle | printf "%s%s" $.optionPrefix }} 
{{ if .Docs }}
{{- range $i, $doc := .Docs }}// {{ if eq $i 0 }}{{ $name }} {{ end }}{{ $doc }}{{ end -}}
{{ end -}}
func {{ $name }}(o {{ .Type }}) {{ $applyOptionFuncType }} {
	return func(c *{{ $.configTypeName }}) error {
    c.{{ .Name }} = ({{ .Type }})(o)
    return nil
	}
}
{{ end }}
`
