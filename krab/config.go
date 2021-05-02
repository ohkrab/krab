package krab

// Config represents all configuration loaded from directory.
//
type Config struct {
	MigrationSets map[string]*MigrationSet
	Migrations    map[string]*Migration
}

func NewConfig(files []*File) (*Config, error) {
	c := &Config{
		MigrationSets: make(map[string]*MigrationSet),
		Migrations:    make(map[string]*Migration),
	}

	for _, f := range files {
		if err := c.appendFile(f); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) appendFile(file *File) error {
	return nil
}
