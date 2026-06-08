package routes

import (
	"entry-system/internals/handlers"
	"entry-system/internals/middleware"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {
	app.Post("/auth/signup", handlers.Signup)
	app.Post("/auth/login", handlers.Login)
	app.Post("/auth/verify", handlers.VerifyAccount)
	app.Post("/auth/forgot-password", handlers.ForgotPassword)
	app.Post("/auth/reset-password", handlers.ResetPassword)

	api := app.Group("/api", middleware.AuthMiddleware)

	api.Get(
		"/me",
		middleware.RequirePermission("view_own_profile"),
		handlers.GetMe,
	)

	api.Put(
		"/me",
		middleware.RequirePermission("update_own_profile"),
		handlers.UpdateMe,
	)

	api.Post(
		"/users",
		middleware.RequirePermission("create_user"),
		handlers.CreateUser,
	)

	api.Get(
		"/users",
		middleware.RequirePermission("view_users"),
		handlers.GetUsers,
	)
	api.Get(
		"/users/dropdown",
		middleware.RequirePermission("create_visitor"),
		handlers.GetUsersDropdown,
	)

	api.Get(
		"/users/:id",
		middleware.RequirePermission("view_users"),
		handlers.GetUser,
	)

	api.Put(
		"/users/:id",
		middleware.RequirePermission("update_user"),
		handlers.UpdateUser,
	)

	api.Delete(
		"/users/:id",
		middleware.RequirePermission("delete_user"),
		handlers.DeleteUser,
	)

	api.Patch(
		"/users/:id/deactivate",
		middleware.RequirePermission("update_user"),
		handlers.DeactivateUser,
	)

	api.Get(
		"/users/:id/visitor-history",
		middleware.RequirePermission("view_visitors"),
		handlers.VisitorHistory,
	)

	api.Post(
		"/visitors",
		middleware.RequirePermission("create_visitor"),
		handlers.CreateVisitor,
	)

	api.Get(
		"/visitors",
		middleware.RequirePermission("view_visitors"),
		handlers.GetVisitors,
	)

	api.Put(
		"/visitors/:id",
		middleware.RequirePermission("update_visitor"),
		handlers.UpdateVisitor,
	)

	api.Patch(
		"/visitors/:id/restrict",
		middleware.RequirePermission("restrict_visitor"),
		handlers.RestrictVisitor,
	)

	api.Get(
		"/visitor-entries",
		middleware.RequirePermission("view_visitors"),
		handlers.GetVisitorEntriesGrouped,
	)
	api.Patch(
		"/exit/:entryId",
		middleware.RequirePermission("make_entry"),
		handlers.ExitVisitor,
	)
	api.Get(
		"/visitors/mobile/:mobile",
		middleware.RequirePermission("create_visitor"),
		handlers.GetVisitorByMobile,
	)
	api.Get(
		"/visitor-stats",
		middleware.RequirePermission("view_visitors"),
		handlers.GetVisitorStats,
	)

	api.Get(
		"/dropdown/persons",
		middleware.RequirePermission("create_visitor"),
		handlers.GetPersonsDropdown,
	)

	api.Get(
		"/visitor-documents/:documentId",
		middleware.RequirePermission("create_visitor"),
		handlers.GetVisitorDocumentForAI,
	)

	api.Post(
		"/visitor-documents/ai-response",
		middleware.RequirePermission("create_visitor"),
		handlers.UpdateDocumentAIResponse,
	)
	api.Get(
		"/all-onSite",
		middleware.RequirePermission("make_entry"),
		handlers.GetActiveEntries,
	)

	api.Post(
		"/guard-session-start",
		middleware.RequirePermission("make_entry"),
		handlers.StartGuardSession,
	)

	api.Post(
		"/guard-session-end",
		middleware.RequirePermission("make_entry"),
		handlers.EndGuardSession,
	)

	api.Get(
		"/guards",
		middleware.RequirePermission("view_visitors"),
		handlers.GetAllGuardSessions,
	)

	api.Get(
		"/guard/:guardId",
		middleware.RequirePermission("view_visitors"),
		handlers.GetGuardSessionsByGuardID,
	)
}
