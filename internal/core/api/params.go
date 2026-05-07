package api

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func getDocTypeName(c *fiber.Ctx) string {
	name := c.Params("name")

	decoded, err := url.PathUnescape(name)
	if err != nil {
		return name
	}

	return decoded
}
