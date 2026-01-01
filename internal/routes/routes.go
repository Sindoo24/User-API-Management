package routes

import (
	"github.com/gofiber/fiber/v2"

	"BACKEND/internal/handler"
	"BACKEND/internal/middleware"
)

func Register(app *fiber.App, h *handler.UserHandler) {
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Post("/users", h.Create)
	app.Get("/users/:id", h.GetByID)
	app.Get("/users", h.List)
	app.Put("/users/:id", h.Update)
	app.Delete("/users/:id", h.Delete)
}
