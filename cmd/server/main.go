package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"BACKEND/config"
	"BACKEND/db/sqlc/generated"
	"BACKEND/internal/handler"
	"BACKEND/internal/logger"
	"BACKEND/internal/middleware"
	"BACKEND/internal/repository"
	"BACKEND/internal/routes"
	"BACKEND/internal/service"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New()
	defer appLogger.Sync()
	middleware.InitLogger(appLogger)
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	appLogger.Info("Connected to database successfully")

	queries := generated.New(dbPool)
	userRepo := repository.NewUserRepository(queries)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userRepo, userSvc, appLogger)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Register(app, userHandler)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		appLogger.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			appLogger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	serverAddr := ":" + cfg.ServerPort
	appLogger.Info("Starting server on " + serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
