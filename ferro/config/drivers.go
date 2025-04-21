package config

// Driver

type Driver struct {
	Path     string     `yaml:"-"`
	Metadata Metadata   `yaml:"metadata"`
	Spec     DriverSpec `yaml:"spec"`
}

type DriverSpec struct {
	Driver string       `yaml:"driver"`
	Config DriverConfig `yaml:"config"`
}

type DriverConfig map[string]any

func (d *Driver) EnforceDefaults() {
	if d.Spec.Config == nil {
		d.Spec.Config = make(DriverConfig)
	}
}

func (d *Driver) Validate() *Errors {
	errors := &Errors{
		Errors: []error{},
	}

	return errors
}
