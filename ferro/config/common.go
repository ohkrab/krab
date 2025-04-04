package config

type Header struct {
	Kind       string `yaml:"kind"`
	ApiVersion string `yaml:"apiVersion"`
}

type Metadata struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Args        []*Args `yaml:"args"`
}

type Args struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
}
