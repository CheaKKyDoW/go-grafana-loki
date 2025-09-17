package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Go Clean Arch User Service",
		ServerHeader: "Fiber",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	// Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "go-clean-arch-user-service",
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "User service is running",
			"users":   []string{"user1", "user2"},
		})
	})

	// Start server
	go func() {
		log.Println("User service starting on port 8080")
		if err := app.Listen(":8080"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down user service...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("User service exited")
}
