package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/afero"
)

// Parser is the main krab data parser.
type Parser struct {
	fs afero.Afero
	p  *hclparse.Parser
}

// NewParser creates krab parser.
func NewParser() *Parser {
	return &Parser{
		fs: afero.Afero{Fs: afero.OsFs{}},
		p:  hclparse.NewParser(),
	}
}

func (p *Parser) LoadHCLFile(path string) (hcl.Body, hcl.Diagnostics) {
	src, err := p.fs.ReadFile(path)

	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to load file",
				Detail:   fmt.Sprintf("The file %q could not be read.", path),
			},
		}
	}

	file, diags := p.p.ParseHCL(src, path)

	if file == nil || file.Body == nil {
		return hcl.EmptyBody(), diags
	}

	return file.Body, diags
}

var configFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "globals",
		},
		{
			Type:       "connection",
			LabelNames: []string{"name"},
		},
	},
}
