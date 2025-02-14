package helpers

import "github.com/gofiber/fiber/v2"

func ResultJSON(c *fiber.Ctx, code int, error bool, message string, data interface{}) error {

	return c.JSON(fiber.Map{
		"code":    code,
		"error":   error,
		"message": message,
		"data":    data,
	})
}
