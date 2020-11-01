package configs

import "github.com/hashicorp/hcl/v2"

func (p *Parser) LoadConfigFile(path string) (*File, hcl.Diagnostics) {
	body, diags := p.LoadHCLFile(path)
	if body == nil {
		return nil, diags
	}

	file := &File{}
	content, contentDiags := body.Content(configFileSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case "globals":
			globs, globsDiags := decodeGlobalsBlock(block)
			diags = append(diags, globsDiags...)
			file.Globals = append(file.Globals, globs...)
		case "connection":
			conn, connDiags := decodeConnectionBlock(block)
			diags = append(diags, connDiags...)
			file.Connections = append(file.Connections, conn)
		default:
		}
	}

	return file, diags
}
