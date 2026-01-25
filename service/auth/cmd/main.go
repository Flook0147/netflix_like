package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Flook0147/netflix_like/service/auth/internal/adapter/bcrypt"
	gormrepo "github.com/Flook0147/netflix_like/service/auth/internal/adapter/gorm"
	"github.com/Flook0147/netflix_like/service/auth/internal/adapter/gorm/model"
	httpadapter "github.com/Flook0147/netflix_like/service/auth/internal/adapter/http"

	"github.com/Flook0147/netflix_like/service/auth/internal/adapter/jwt"
	"github.com/Flook0147/netflix_like/service/auth/internal/core/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialization code here
	db, err := gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	log.Println("migrate database")
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	hasher := bcrypt.New()
	jwtProvider := jwt.New(os.Getenv("JWT_SECRET"), 24*time.Hour)
	userRepo := gormrepo.New(db)
	authService := auth.NewAuthService(jwtProvider, userRepo, hasher)

	handler := httpadapter.NewHTTPHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/register", handler.Register)
	log.Println("Auth service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
