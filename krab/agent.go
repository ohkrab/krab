package krab

import (
	"github.com/gin-gonic/gin"
)

type Agent struct {
	app *gin.Engine
}

func (a *Agent) Run() {
	g := gin.Default()
	a.app = g

	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"name": "Oh! Krab!"})
	})
	g.Run()
}
