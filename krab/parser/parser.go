package parser

import (
	"io/ioutil"
	"os"

	HclAst "github.com/hashicorp/hcl/hcl/ast"
	HclParser "github.com/hashicorp/hcl/hcl/parser"
)

// ParsedFile type.
type ParsedFile struct {
	Ast *HclAst.File
}

// FindFilesToParse recursively find files in current dir.
func FindFilesToParse() []string {
	return make([]string, 0)
}


// ParseFromFile parses HCL files based on path.
func ParseFromFile(path string) (*ParsedFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return parse(b)
}

func parse(source []byte) (*ParsedFile, error) {
	ast, err := HclParser.Parse(source)
	if err != nil {
		return nil, err
	}
	return &ParsedFile{
		Ast: ast,
	}, nil
}
