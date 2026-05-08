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

	permissionHandler := NewPermissionHandler(handler.DB)
	router.Get("/permissions/check", permissionHandler.CheckPermission)

	router.Get("/users", handler.ListUsers)
	router.Get("/users/:name", handler.GetUser)
	router.Get("/users/:name/roles", handler.GetUserRoles)
	router.Get("/roles", handler.ListRoles)

	router.Get("/meta/:doctype", handler.GetMeta)

	router.Post("/migration/preview", handler.PreviewMigration)
	router.Post("/migration/apply", handler.ApplyMigration)
	router.Get("/migration/logs", handler.MigrationLogs)

	namingSeriesHandler := NewNamingSeriesHandler(handler.DB)
	router.Get("/naming-series", namingSeriesHandler.ListNamingSeries)
	router.Get("/naming-series/:seriesKey", namingSeriesHandler.GetNamingSeries)
	router.Post("/naming-series/:seriesKey/next", namingSeriesHandler.NextNamingSeries)

	router.Get("/documents/:doctype", handler.ListDocuments)
	router.Get("/documents/:doctype/:name", handler.GetDocument)
	router.Post("/documents/:doctype", handler.CreateDocument)
	router.Put("/documents/:doctype/:name", handler.UpdateDocument)
	router.Delete("/documents/:doctype/:name", handler.DeleteDocument)
}
