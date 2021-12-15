package krabapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ohkrab/krab/krabcmd"
)

// Agent exposes API from config.
type Agent struct {
	Registry krabcmd.Registry
}

func (a *Agent) Run() {
	router := gin.Default()
	api := router.Group("/api")
	for _, action := range a.Registry {
		method, path := action.HttpEndpoint()
		switch method {
		case http.MethodGet:
			api.GET(path, func(c *gin.Context) {
				err := action.Do(c.Request.Context(), c.Writer)
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				}
				c.Status(http.StatusOK)
			})
		default:
			panic(fmt.Sprintf("HTTP %s Method not implemented", method))
		}
	}
	router.Run()
}
