package config

import "strconv"

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

func (c DriverConfig) Has(key string) bool {
	_, ok := c[key]
	return ok
}

func (c DriverConfig) String(key string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		if i, ok := v.(int); ok {
			return strconv.Itoa(i)
		}
		if f, ok := v.(float64); ok {
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
		if b, ok := v.(bool); ok {
			return strconv.FormatBool(b)
		}
	}
	return ""
}

func (c DriverConfig) Float(key string) float64 {
	if v, ok := c[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func (c DriverConfig) Int(key string) int {
	if v, ok := c[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

func (c DriverConfig) Bool(key string) bool {
	if v, ok := c[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}
