package krabapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ohkrab/krab/krab"
)

// Agent exposes API from config.
type Agent struct {
	Registry krab.CmdRegistry
}

func (a *Agent) Run() {
	router := gin.Default()
	api := router.Group("/api")
	for _, cmd := range a.Registry {
		path := fmt.Sprint("/", strings.Join(cmd.Name(), "/"))
		switch cmd.HttpMethod() {
		case http.MethodGet:
			api.GET(path, func(c *gin.Context) {
				err := cmd.Do(c.Request.Context(), krab.CmdOpts{Writer: c.Writer})
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				}
				c.Status(http.StatusOK)
			})
		default:
			panic(fmt.Sprintf("HTTP %s Method not implemented", cmd.HttpMethod()))
		}
	}
	router.Run()
}
