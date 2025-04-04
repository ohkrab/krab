package parser

import "github.com/ohkrab/krab/ferro/config"

type ParsedConfig struct {
	Files []*ParsedFile
}

type ParsedFile struct {
	Path   string
	Chunks []*ParsedChunk

	Migrations    []*config.Migration
	MigrationSets []*config.MigrationSet
}

type ParsedChunk struct {
	Header *config.Header
	Raw    []byte
}
