package main

import (
	"fmt"
	"log"
	"os"

	authhttp "github.com/Flook0147/netflix_like/internal/auth/adapter/inbound/http"
	authapp "github.com/Flook0147/netflix_like/internal/auth/app"
	userapp "github.com/Flook0147/netflix_like/internal/user/app"
	"github.com/gofiber/fiber/v3"

	authOutbound "github.com/Flook0147/netflix_like/internal/auth/adapter/outbound"
	userRouter "github.com/Flook0147/netflix_like/internal/user/adapter/inbound/http"
	userOutbound "github.com/Flook0147/netflix_like/internal/user/adapter/outbound"
	userdb "github.com/Flook0147/netflix_like/internal/user/adapter/outbound/db"
	"github.com/Flook0147/netflix_like/internal/user/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	serverPort := os.Getenv("SERVER_PORT")

	sslmode := "disable"
	timezone := "Asia/Bangkok"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, dbPort, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	fmt.Println("Database connected 🚀")

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	userRepo := userdb.NewUserRepository(db)

	userService := userapp.NewUserService(userRepo)

	tokenRepo := authOutbound.NewRefreshTokenRepository(db)

	authService := authapp.NewAuthService(userService, tokenRepo)

	authHandler := authhttp.NewAuthHandler(authService)

	tokenHandler := userOutbound.TokenHandler{}.NewTokenHandler(authService)

	userService.SetTokenPort(tokenHandler)

	// authRouter := authhttp.NewRouter(authHandler)

	userHandler := userRouter.NewUserHandler(userService)

	// userRouter := userRouter.NewRouter(userHandler)

	log.Println("Server started at :", serverPort)

	app := fiber.New()

	authhttp.RegisterRoutes(app, authHandler)
	userRouter.RegisterRoutes(app, userHandler)

	log.Println("Server started at :", serverPort)

	log.Fatal(app.Listen(":" + serverPort))
}
