package config

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

const (
	FileExt = ".fyml"
)

type Filesystem struct {
	Dir string
}

func NewFilesystem(dir string) *Filesystem {
	return &Filesystem{
		Dir: dir,
	}
}

func (f *Filesystem) MkdirAll(paths []string) error {
	dir := f.Dir
	for _, path := range paths {
		dir = filepath.Join(dir, path)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", strings.Join(paths, "/"), err)
	}
	return nil
}

func (f *Filesystem) Exists(path string) bool {
	_, err := os.Stat(filepath.Join(f.Dir, path))
	return !os.IsNotExist(err)
}

func (f *Filesystem) TouchFile(path string, data []byte) error {
	if _, err := os.Stat(filepath.Join(f.Dir, path)); os.IsNotExist(err) {
		if err := os.WriteFile(filepath.Join(f.Dir, path), data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}
	}
	return nil
}

func (f *Filesystem) LoadFiles(paths []string) ([]*ParsedFile, error) {
	parsedFiles := make([]*ParsedFile, len(paths))
	for i, path := range paths {
		parsedFiles[i] = &ParsedFile{
			Path:   path,
			Chunks: []*ParsedChunk{},
		}
	}

	eg := errgroup.Group{}
	for i, file := range parsedFiles {
		eg.Go(func() error {
			src, err := os.ReadFile(filepath.Join(f.Dir, file.Path))
			if err != nil {
				return fmt.Errorf("failed to read config file %s: %w", file.Path, err)
			}

			chunks := bytes.Split(src, []byte("\n---"))
			for c, chunk := range chunks {
				var parsedHeader Header
				if err := yaml.Unmarshal(chunk, &parsedHeader); err != nil {
					return fmt.Errorf("failed to unmarshal chunk (%d): %w\n  %s", c, err, string(chunk))
				}

				parsedFiles[i].Chunks = append(parsedFiles[i].Chunks, &ParsedChunk{
					Header: &parsedHeader,
					Raw:    chunk,
				})
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to load files: %w", err)
	}

	return parsedFiles, nil
}

func (f *Filesystem) DirFiles() ([]string, error) {
	paths := []string{}

	err := fs.WalkDir(os.DirFS(f.Dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := f.fileExt(path)
		if ext == "" || f.isIgnoredFile(path) {
			return nil
		}

		paths = append(paths, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to dir files: %w", err)
	}

	return paths, nil
}

func (f *Filesystem) fileExt(path string) string {
	if strings.HasSuffix(path, FileExt) {
		return FileExt
	}

	return ""
}

func (f *Filesystem) isIgnoredFile(name string) bool {
	// ignore dotfiles and emacs/vim backups
	return strings.HasSuffix(name, "~") ||
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#")
}

type ParsedConfig struct {
	Files []*ParsedFile
}

type ParsedFile struct {
	Path   string
	Chunks []*ParsedChunk

	Migrations    []*Migration
	MigrationSets []*MigrationSet
	Drivers       []*Driver
}

type ParsedChunk struct {
	Header *Header
	Raw    []byte
}

func (p *ParsedConfig) BuildConfig() (*Config, *Errors) {
	cfg := New()

	for _, file := range p.Files {
		for _, migration := range file.Migrations {
			if err := cfg.AddMigration(migration); err != nil {
				return nil, Errorf("adding Migration: %w", err)
			}
		}
		for _, migrationSet := range file.MigrationSets {
			if err := cfg.AddMigrationSet(migrationSet); err != nil {
				return nil, Errorf("adding MigrationSet: %w", err)
			}
		}
		for _, driver := range file.Drivers {
			if err := cfg.AddDriver(driver); err != nil {
				return nil, Errorf("adding Driver: %w", err)
			}
		}
	}

	errors := cfg.Validate()
	return cfg, errors
}
