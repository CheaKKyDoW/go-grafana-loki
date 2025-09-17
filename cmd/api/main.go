package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-clean-arch/internal/delivery/http/router"
	"go-clean-arch/internal/shared/config"
	"go-clean-arch/internal/shared/infrastructure/database"
	"go-clean-arch/internal/shared/infrastructure/logger"
	"go-clean-arch/internal/shared/middleware"

	// User domain

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger with Loki integration
	logger := logger.NewWithLoki(cfg.Log.Level, cfg.Loki)
	defer logger.Sync()

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	userRepo := userRepository.NewUserRepository(db)

	// Initialize use cases
	userUC := userUseCase.NewUserUseCase(userRepo, logger)

	// Initialize handlers
	userH := userHandler.NewUserHandler(userUC, logger)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Go Clean Arch API",
		ServerHeader: "Fiber",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": "Internal server error",
				"error":   err.Error(),
			})
		},
	})

	// Setup routes
	router.SetupRoutes(app, userH, authMiddleware)

	// Start server in goroutine
	go func() {
		logger.Info("Server starting on port " + cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
