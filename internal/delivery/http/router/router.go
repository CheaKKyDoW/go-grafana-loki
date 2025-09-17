package router

import (
	"go-clean-arch/internal/shared/middleware"
	userHandler "go-clean-arch/internal/user/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App, userHandler *userHandler.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "go-clean-arch-api"})
	})

	// API routes
	api := app.Group("/api/v1")

	// User domain routes
	setupUserRoutes(api, userHandler, authMiddleware)

	// Future domain routes can be added here
	// setupCartRoutes(api, cartHandler, authMiddleware)
	// setupOrderRoutes(api, orderHandler, authMiddleware)
}

func setupUserRoutes(api fiber.Router, userHandler *userHandler.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// Public user routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Protected user routes
	users := api.Group("/users")
	users.Use(authMiddleware.RequireAuth())
	users.Get("/profile", userHandler.GetProfile)
}
