package api

import "github.com/gofiber/fiber/v2"

// ConfigData is exposed to the frontend at runtime.
type ConfigData struct {
	APIBaseURL string `json:"api_base_url"`
	WSURL      string `json:"ws_url"`
}

// ConfigHandler serves GET /api/config.
func ConfigHandler(c *fiber.Ctx) error {
	return WriteSuccess(c, ConfigData{
		APIBaseURL: "/api",
		WSURL:      "",
	})
}
