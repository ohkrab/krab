package tpls

import (
	"bytes"
	"fmt"
	"text/template"
)

// Templates is used for templates rendering.
type Templates struct {
	template *template.Template
	embedded map[string][]byte
}

// New created template renderer with values to replace.
func New(funcMap template.FuncMap) *Templates {
	t := &Templates{
		template: template.New("").Funcs(funcMap),
		embedded: make(map[string][]byte),
	}
	return t
}

func (t *Templates) AddEmbedded(name string, b []byte) {
	t.embedded[name] = b
}

// Validate verifies if template is correct.
func (t *Templates) Validate(b []byte) error {
	_, err := t.template.Parse(string(b))
	return err
}

// Render applies values and renders final output.
func (t *Templates) Render(b []byte, data any) ([]byte, error) {
	var rendered bytes.Buffer
	template, err := t.template.Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	err = template.Execute(&rendered, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}
	return rendered.Bytes(), nil
}

func (t *Templates) RenderEmbedded(name string, data any) ([]byte, error) {
	b, ok := t.embedded[name]
	if !ok {
		return nil, fmt.Errorf("embedded template %s not found", name)
	}
	return t.Render(b, data)
}
