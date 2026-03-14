package main

import (
	"fmt"
	"kyaho-space/api/database"
	"kyaho-space/api/routes"
	"kyaho-space/pkg/config"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/job"
	"kyaho-space/pkg/meet_plan"
	"kyaho-space/pkg/middleware"
	"kyaho-space/pkg/user"
	"kyaho-space/pkg/wishlist"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	database.OpenConnection(cfg)

	database.MigrateSchema(&entities.Job{})
	database.MigrateSchema(&entities.MeetPlan{})
	database.MigrateSchema(&entities.Wishlist{})
	database.MigrateSchema(&entities.User{})

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
	})

	// Configure CORS based on Environment
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}

	if cfg.Env == "development" {
		corsConfig.AllowOrigins = []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			cfg.FrontendURL,
		}
	} else {
		corsConfig.AllowOrigins = []string{cfg.FrontendURL}
	}

	app.Use(cors.New(corsConfig))

	api := app.Group("/api")

	// --- PUBLIC ROUTES ---

	// Auth
	userRepo := user.NewRepo(database.DB)
	userService := user.NewService(userRepo)
	routes.AuthRouter(api, userService, cfg.JWTSecret, cfg.Env)

	// --- AUTH MIDDLEWARE ---
	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)

	// MeetPlan
	meetPlanRepo := meet_plan.NewRepo(database.DB)
	meetPlanService := meet_plan.NewService(meetPlanRepo)
	routes.MeetPlanRouter(api, meetPlanService, authMiddleware)

	// --- PROTECTED ROUTES ---
	api.Use(authMiddleware)

	// Job
	jobRepo := job.NewRepo(database.DB)
	jobService := job.NewService(jobRepo)
	routes.JobRouter(api, jobService)

	// Wishlist
	wishlistRepo := wishlist.NewRepo(database.DB)
	wishlistService := wishlist.NewService(wishlistRepo)
	routes.WishlistRouter(api, wishlistService)

	// Start server
	listenAddr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("Server starting on %s\n", listenAddr)
	if err := app.Listen(listenAddr); err != nil {
		panic(err)
	}
}
