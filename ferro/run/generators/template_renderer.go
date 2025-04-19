package generators

import "strings"

type TemplateRenderer struct {
}

func (r *TemplateRenderer) Render(template string) string {
	return strings.ReplaceAll(template, "\t", "  ")
}
