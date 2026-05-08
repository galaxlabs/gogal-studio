package api

import (
	"context"
	"net/url"
	"time"

	"github.com/galaxylabs/gogal-studio/internal/core/meta"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetMeta(c *fiber.Ctx) error {
	doctype := decodeMetaDocTypeName(c.Params("doctype"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loader := meta.NewLoader(h.DB)

	result, err := loader.GetDocTypeMeta(ctx, doctype)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "DocType metadata not found",
			"doctype": doctype,
			"detail":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func decodeMetaDocTypeName(raw string) string {
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return raw
	}

	return decoded
}
