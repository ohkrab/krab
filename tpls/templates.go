package tpls

import (
	"bytes"
	"text/template"
)

// Templates is used for templates rendering.
type Templates struct {
	values   map[string]interface{}
	template *template.Template
}

type root struct {
	Args map[string]interface{}
}

// New created template renderer with values to replace.
func New(values map[string]interface{}) *Templates {
	t := &Templates{
		values:   values,
		template: template.New(""),
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
	w := bytes.NewBufferString("")
	template, err := t.template.Parse(s)
	if err != nil {
		//TODO: handle error
		panic(err)
	}

	err = template.Execute(w, root{Args: t.values})
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	return w.String()
}
