package tpls

import (
	"strings"
	"text/template"
)

// Templates is used for templates rendering.
type Templates struct {
	values   map[string]any
	template *template.Template
}

type root struct {
	Args map[string]any
}

// New created template renderer with values to replace.
func New(values map[string]any, funcMap template.FuncMap) *Templates {
	t := &Templates{
		values:   values,
		template: template.New("").Funcs(funcMap),
	}
	return t
}

// Validate verifies if template is correct.
func (t *Templates) Validate(s string) error {
	_, err := t.template.Parse(s)
	return err
}

// Render applies values and renders final output.
func (t *Templates) Render(s string) string {
	sb := strings.Builder{}
	template, err := t.template.Parse(s)
	if err != nil {
		//TODO: handle error
		panic(err)
	}

	err = template.Execute(&sb, root{Args: t.values})
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	return sb.String()
}
