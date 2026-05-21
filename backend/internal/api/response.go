package api

import "github.com/gofiber/fiber/v2"

// Envelope is the standardized API response format.
type Envelope struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// WriteSuccess sends a success envelope.
func WriteSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{Status: "success", Data: data})
}

// WriteError sends an error envelope.
func WriteError(c *fiber.Ctx, status int, message string, errors []string) error {
	return c.Status(status).JSON(Envelope{Status: "error", Message: message, Errors: errors})
}
