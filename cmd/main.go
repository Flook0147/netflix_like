package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// AUTH
	authhttp "github.com/Flook0147/netflix_like/internal/auth/adapter/inbound/http"
	authOutbound "github.com/Flook0147/netflix_like/internal/auth/adapter/outbound"
	authapp "github.com/Flook0147/netflix_like/internal/auth/app"

	// USER
	userRouter "github.com/Flook0147/netflix_like/internal/user/adapter/inbound/http"
	userdb "github.com/Flook0147/netflix_like/internal/user/adapter/outbound/db"
	userapp "github.com/Flook0147/netflix_like/internal/user/app"

	// PAYMENT
	paymentRepo "github.com/Flook0147/netflix_like/internal/payment/adapter/outbound/db"
	paymentapp "github.com/Flook0147/netflix_like/internal/payment/app"

	// SUBSCRIPTION
	subscriptionHandlers "github.com/Flook0147/netflix_like/internal/subscription/adapter/inbound/http"
	subscriptionRepo "github.com/Flook0147/netflix_like/internal/subscription/adapter/outbound/db"
	subscriptionapp "github.com/Flook0147/netflix_like/internal/subscription/app"

	// MOVIE

	movieHandlers "github.com/Flook0147/netflix_like/internal/movie/adapter/inbound/http"
	movieRepos "github.com/Flook0147/netflix_like/internal/movie/adapter/outbound/db"
	movieApp "github.com/Flook0147/netflix_like/internal/movie/app"
	"github.com/Flook0147/netflix_like/internal/router"
)

func main() {
	// ========================
	// ENV
	// ========================
	godotenv.Load()
	serverPort := os.Getenv("SERVER_PORT")

	// ========================
	// DATABASE
	// ========================
	db := initDB()

	// ========================
	// RABBITMQ (SINGLE CONNECTION)
	// ========================

	// ========================
	// USER + AUTH
	// ========================
	userRepo := userdb.NewUserRepository(db)
	userService := userapp.NewUserService(userRepo)

	tokenRepo := authOutbound.NewRefreshTokenRepository(db)
	authService := authapp.NewAuthService(userService, tokenRepo)

	authHandler := authhttp.NewAuthHandler(authService)

	userHandler := userRouter.NewUserHandler(authService)

	// ========================
	// PAYMENT
	// ========================
	paymentRepository := paymentRepo.NewPaymentRepository(db)

	paymentService := paymentapp.NewPaymentService(paymentRepository)

	// ========================
	// SUBSCRIPTION
	// ========================
	subRepo := subscriptionRepo.NewSubscriptionRepo(db)
	planRepo := subscriptionRepo.NewSubscriptionPlanRepo(db)

	subscriptionService := subscriptionapp.NewSubscriptionService(subRepo, planRepo, paymentService)
	subscriptionHandler := subscriptionHandlers.NewSubscriptionHandler(subscriptionService)

	// ========================
	// Movie
	// ========================

	movieRepo := movieRepos.NewMovieRepository(db)
	movieService := movieApp.NewMovieService(movieRepo, subscriptionService)
	movieHandler := movieHandlers.NewMovieHandler(movieService)

	// ========================
	// HTTP SERVER
	// ========================
	app := fiber.New()
	router.SetupRoutes(app, router.Handlers{
		Auth:         authHandler,
		User:         userHandler,
		Subscription: subscriptionHandler,
		Movie:        movieHandler,
	})

	authhttp.RegisterRoutes(app, authHandler)
	userRouter.RegisterRoutes(app, userHandler)

	log.Println("🚀 Server running on port:", serverPort)
	log.Fatal(app.Listen(":" + serverPort))
}

func initDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, user, password, dbname, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}

	log.Println("✅ Database connected")

	return db
}
