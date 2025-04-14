package config

import "fmt"

type Errors struct {
	Errors []error
}

func Errorf(format string, args ...any) *Errors {
	return &Errors{
		Errors: []error{fmt.Errorf(format, args...)},
	}
}

func (e *Errors) Append(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *Errors) HasErrors() bool {
	return len(e.Errors) > 0
}
