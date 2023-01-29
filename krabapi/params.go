package krabapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ohkrab/krab/krab"
)

func bindInputs(c *gin.Context) (krab.NamedInputs, error) {
	params := map[string]interface{}{}
	if c.Request.ContentLength == 0 {
		return krab.NamedInputs{}, nil
	}
	err := c.BindJSON(&params)
	if err != nil {
		return krab.NamedInputs{}, fmt.Errorf("Can't bind inputs: %w", err)
	}
	inputs, ok := params["inputs"]
	if !ok {
		return krab.NamedInputs{}, nil
	}

	if inputs, ok := inputs.(map[string]interface{}); ok {
		return krab.NamedInputs(inputs), nil
	}

	return krab.NamedInputs{}, fmt.Errorf("Failed to fetch inputs")
}
