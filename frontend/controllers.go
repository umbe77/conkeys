package frontend

import "github.com/gofiber/fiber/v2"

func InitControllers(app *fiber.App) error {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	return nil
}
