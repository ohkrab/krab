package krabapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ohkrab/krab/krab"
)

// Agent exposes API from config.
type Agent struct {
	Registry *krab.CmdRegistry
}

func (a *Agent) Run() error {
	router := gin.Default()
	api := router.Group("/api")
	for _, cmd := range a.Registry.Commands {
		path := fmt.Sprint("/", strings.Join(cmd.Name(), "/"))
		switch cmd.HttpMethod() {
		case http.MethodGet:
			api.GET(path, func(c *gin.Context) {
				resp, err := cmd.Do(c.Request.Context(), krab.CmdOpts{})
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				}
				err = json.NewEncoder(c.Writer).Encode(resp)
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
	return nil
}
