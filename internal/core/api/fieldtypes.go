package api

import (
	"net/url"
	"sort"

	"github.com/galaxylabs/gogal-studio/internal/core/fieldtype"
	"github.com/gofiber/fiber/v2"
)

func ListFieldTypes(c *fiber.Ctx) error {
	defs := fieldtype.All()

	sort.Slice(defs, func(i, j int) bool {
		return defs[i].Name < defs[j].Name
	})

	return c.JSON(fiber.Map{
		"data": defs,
	})
}

func GetFieldType(c *fiber.Ctx) error {
	name := c.Params("name")

	decoded, err := url.PathUnescape(name)
	if err == nil {
		name = decoded
	}

	def, err := fieldtype.MustGet(name)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": def,
	})
}
