package frontend

import "github.com/gofiber/fiber/v2"

func index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func InitControllers(app *fiber.App) error {
	app.Get("/", index)
	return nil
}
