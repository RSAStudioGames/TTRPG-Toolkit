package api

import "github.com/gofiber/fiber/v2"

// Register mounts API routes on the Fiber app.
func Register(app *fiber.App) {
	app.Get("/api/config", ConfigHandler)
}
