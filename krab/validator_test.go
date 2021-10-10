package krab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRefName(t *testing.T) {
	assert := assert.New(t)

	// allow alphanumeric and underscore
	assert.Nil(ValidateRefName("valid_ref"))
	assert.Nil(ValidateRefName("valid_123"))
	assert.Nil(ValidateRefName("ValidRef"))
	assert.Nil(ValidateRefName("___"))

	// cannot start with number
	assert.NotNil(ValidateRefName("123"))
	assert.NotNil(ValidateRefName("123_abc"))

	// cannot be empty
	assert.NotNil(ValidateRefName(""))

	// no other separators
	assert.NotNil(ValidateRefName("abc-def"))
	assert.NotNil(ValidateRefName("abc def"))
}

func TestValidateStringNonEmpty(t *testing.T) {
	assert := assert.New(t)

	// Length must be > 0
	assert.Nil(ValidateStringNonEmpty("field", "a"))
	assert.NotNil(ValidateStringNonEmpty("field", ""))
}
