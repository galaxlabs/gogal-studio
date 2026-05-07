package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *Handler) {
	router.Get("/installed-apps", handler.ListInstalledApps)
	router.Get("/modules", handler.ListModules)

	router.Get("/doctypes", handler.ListDocTypes)
	router.Post("/doctypes", handler.SaveDocType)
	router.Get("/doctypes/:name", handler.GetDocType)
	router.Get("/doctypes/:name/fields", handler.GetDocTypeFields)
	router.Get("/doctypes/:name/permissions", handler.GetDocTypePermissions)

	router.Get("/fieldtypes", ListFieldTypes)
	router.Get("/fieldtypes/:name", GetFieldType)

	router.Get("/users", handler.ListUsers)
	router.Get("/users/:name", handler.GetUser)
	router.Get("/users/:name/roles", handler.GetUserRoles)
	router.Get("/roles", handler.ListRoles)

	namingSeriesHandler := NewNamingSeriesHandler(handler.DB)
	router.Get("/naming-series", namingSeriesHandler.ListNamingSeries)
	router.Get("/naming-series/:seriesKey", namingSeriesHandler.GetNamingSeries)
	router.Post("/naming-series/:seriesKey/next", namingSeriesHandler.NextNamingSeries)
}
