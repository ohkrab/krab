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
		cmd := cmd
		path := fmt.Sprint("/", strings.Join(cmd.Name(), "/"))
		switch cmd.HttpMethod() {
		case http.MethodPost:
			api.POST(path, func(c *gin.Context) {

				inputs, err := bindInputs(c)
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}

				resp, err := cmd.Do(c.Request.Context(), krab.CmdOpts{Inputs: inputs})
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}
				err = json.NewEncoder(c.Writer).Encode(resp)
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}
				c.Status(http.StatusOK)
			})
		case http.MethodGet:
			api.GET(path, func(c *gin.Context) {
				c.Header("Cache-Control", "no-store")

				inputs, err := bindInputs(c)
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}

				resp, err := cmd.Do(c.Request.Context(), krab.CmdOpts{Inputs: inputs})
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
				}
				err = json.NewEncoder(c.Writer).Encode(resp)
				if err != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
					return
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
