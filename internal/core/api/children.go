package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetDocTypeActions(c *fiber.Ctx) error {
	name := getDocTypeName(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
			action_name,
			label,
			group_name,
			action_type,
			action,
			handler,
			route,
			method,
			permission,
			visible_when,
			hidden,
			custom,
			enabled,
			idx
		FROM core_doctype_action
		WHERE doctype = $1
		ORDER BY idx, id
	`, name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			actionName  string
			label       string
			groupName   *string
			actionType  string
			action      *string
			handler     *string
			route       *string
			method      *string
			permission  *string
			visibleWhen *string
			hidden      bool
			custom      bool
			enabled     bool
			idx         int
		)

		err := rows.Scan(
			&actionName,
			&label,
			&groupName,
			&actionType,
			&action,
			&handler,
			&route,
			&method,
			&permission,
			&visibleWhen,
			&hidden,
			&custom,
			&enabled,
			&idx,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"action_name":  actionName,
			"label":        label,
			"group":        groupName,
			"group_name":   groupName,
			"action_type":  actionType,
			"action":       action,
			"handler":      handler,
			"route":        route,
			"method":       method,
			"permission":   permission,
			"visible_when": visibleWhen,
			"hidden":       hidden,
			"custom":       custom,
			"enabled":      enabled,
			"idx":          idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"doctype": name,
		"data":    data,
	})
}

func (h *Handler) GetDocTypeLinks(c *fiber.Ctx) error {
	name := getDocTypeName(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
			link_doctype,
			link_fieldname,
			parent_doctype,
			table_fieldname,
			group_name,
			hidden,
			is_child_table,
			custom,
			idx
		FROM core_doctype_link
		WHERE doctype = $1
		ORDER BY idx, id
	`, name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			linkDocType    string
			linkFieldname  string
			parentDocType  *string
			tableFieldname *string
			groupName      *string
			hidden         bool
			isChildTable   bool
			custom         bool
			idx            int
		)

		err := rows.Scan(
			&linkDocType,
			&linkFieldname,
			&parentDocType,
			&tableFieldname,
			&groupName,
			&hidden,
			&isChildTable,
			&custom,
			&idx,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"link_doctype":    linkDocType,
			"link_fieldname":  linkFieldname,
			"parent_doctype":  parentDocType,
			"table_fieldname": tableFieldname,
			"group":           groupName,
			"group_name":      groupName,
			"hidden":          hidden,
			"is_child_table":  isChildTable,
			"custom":          custom,
			"idx":             idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"doctype": name,
		"data":    data,
	})
}

func (h *Handler) GetDocTypeStates(c *fiber.Ctx) error {
	name := getDocTypeName(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.DB.Query(ctx, `
		SELECT
			title,
			color,
			custom,
			idx
		FROM core_doctype_state
		WHERE doctype = $1
		ORDER BY idx, id
	`, name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	data := []fiber.Map{}

	for rows.Next() {
		var (
			title  string
			color  string
			custom bool
			idx    int
		)

		err := rows.Scan(
			&title,
			&color,
			&custom,
			&idx,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data = append(data, fiber.Map{
			"title":  title,
			"color":  color,
			"custom": custom,
			"idx":    idx,
		})
	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"doctype": name,
		"data":    data,
	})
}
