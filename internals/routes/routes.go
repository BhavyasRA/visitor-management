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

	api.Post("/users", middleware.RequirePermission("create_user"), handlers.CreateUser)
	api.Get("/users", middleware.RequirePermission("view_users"), handlers.GetUsers)
	api.Get("/users/:id", middleware.RequirePermission("view_users"), handlers.GetUser)
	api.Put("/users/:id", middleware.RequirePermission("update_user"), handlers.UpdateUser)
	api.Delete("/users/:id", middleware.RequirePermission("delete_user"), handlers.DeleteUser)
	api.Patch("/users/:id/deactivate", middleware.RequirePermission("update_user"), handlers.DeactivateUser)
	api.Get("/users/:id/visitor-history", middleware.RequirePermission("view_visitors"), handlers.VisitorHistory)

	api.Post("/visitors", middleware.RequirePermission("create_visitor"), handlers.CreateVisitor)
	api.Get("/visitors", middleware.RequirePermission("view_visitors"), handlers.GetVisitors)
	api.Put("/visitors/:id", middleware.RequirePermission("update_visitor"), handlers.UpdateVisitor)
	api.Patch("/visitors/:id/restrict", middleware.RequirePermission("restrict_visitor"), handlers.RestrictVisitor)

	api.Post("/guard/make-entry", middleware.RequirePermission("make_entry"), handlers.MakeEntry)
	api.Get("/guard/entries", middleware.RequirePermission("see_entry"), handlers.SeeEntries)
	api.Patch("/guard/exit/:visitorId", middleware.RequirePermission("make_entry"), handlers.ExitVisitor)
	api.Patch("/guard/restrict-visitor/:visitorId", middleware.RequirePermission("restrict_visitor"), handlers.GuardRestrictVisitor)
	api.Patch("/guard/restrict-employee/:userId", middleware.RequirePermission("update_user"), handlers.GuardRestrictEmployee)

	api.Get(
	"/users/me/visitors",
	middleware.RequirePermission("view_visitors"),
	handlers.GetMyVisitors,
)

}
