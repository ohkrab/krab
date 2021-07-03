package krab

import (
	"fmt"
	"regexp"
)

type Validator interface {
	Validate() error
}

// ErrorCoalesce returns first non empty error.
func ErrorCoalesce(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateStringNonEmpty checks if string is not empty.
func ValidateStringNonEmpty(what, s string) error {
	if len(s) > 0 {
		return nil
	}

	return fmt.Errorf("%s cannot be empty", what)
}

// ValidateRefName checks if reference name matches allowed format.
func ValidateRefName(refName string) error {
	matched, err := regexp.Match("^[a-zA-Z_][a-zA-Z0-9_]*$", []byte(refName))
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("Reference `%s` has invalid format", refName)
	}
	return nil
}
