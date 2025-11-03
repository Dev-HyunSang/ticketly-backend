package main

import (
	"context"
	"log"

	"github.com/dev-hyunsang/ticketly-backend/internal/db"
	"github.com/dev-hyunsang/ticketly-backend/internal/handler"
	"github.com/dev-hyunsang/ticketly-backend/internal/middleware"
	"github.com/dev-hyunsang/ticketly-backend/internal/repository/mysql"
	"github.com/dev-hyunsang/ticketly-backend/internal/repository/redis"
	"github.com/dev-hyunsang/ticketly-backend/internal/usecase"
	"github.com/dev-hyunsang/ticketly-backend/internal/util"
	"github.com/gofiber/fiber/v2"
	fiberMiddleware "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// MySQL connection
	client, err := db.ConnectMySQL()
	if err != nil {
		log.Fatalf("failed to mysql connection : %v", err)
	}
	defer client.Close()

	if err = client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed create schema resources : %v", err)
	}

	// Redis connection
	redisClient, err := db.ConnectRedis()
	if err != nil {
		log.Fatalf("failed to redis connection : %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	userRepo := mysql.NewUserRepository(client)
	tokenRepo := redis.NewTokenRepository(redisClient)

	// Initialize utilities
	jwtUtil := util.NewJWTUtil()

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtUtil)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// Apply global middleware
	app.Use(logger.New())
	app.Use(fiberMiddleware.New(fiberMiddleware.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Public routes (no authentication required)
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes (authentication required)
	api := app.Group("/api", authMiddleware.Authenticate)
	api.Get("/me", authHandler.Me)
	api.Post("/logout", authHandler.Logout)

	// User routes
	user := api.Group("/users")
	_ = user // userHandler will be added later

	_ = userUseCase // Use it to avoid unused variable error

	log.Println("Server starting on :3000")
	if err = app.Listen(":3000"); err != nil {
		log.Fatalf("failed starting server %v", err)
	}
}
