package krabapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ohkrab/krab/krab"
)

func bindInputs(c *gin.Context) (krab.Inputs, error) {
	params := map[string]interface{}{}
	if c.Request.ContentLength == 0 {
		return krab.Inputs{}, nil
	}
	err := c.BindJSON(&params)
	if err != nil {
		return krab.Inputs{}, fmt.Errorf("Can't bind inputs: %w", err)
	}
	inputs, ok := params["inputs"]
	if !ok {
		return krab.Inputs{}, nil
	}

	if inputs, ok := inputs.(map[string]interface{}); ok {
		return krab.Inputs(inputs), nil
	}

	return krab.Inputs{}, fmt.Errorf("Failed to fetch inputs")
}
