package krab

import "github.com/gofiber/fiber"

type Agent struct {
	app *fiber.App
}

func (a *Agent) Run() {
	app := fiber.New()
	a.app = app

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Oh Krab!")
	})
	app.Listen(8888)
}
