package router

import (
	"github.com/gofiber/fiber/v3"

	// handlers
	authhttp "github.com/Flook0147/netflix_like/internal/auth/adapter/inbound/http"
	moviehttp "github.com/Flook0147/netflix_like/internal/movie/adapter/inbound/http"
	subhttp "github.com/Flook0147/netflix_like/internal/subscription/adapter/inbound/http"
	userhttp "github.com/Flook0147/netflix_like/internal/user/adapter/inbound/http"

	// middleware
	authMiddleware "github.com/Flook0147/netflix_like/internal/auth/middleware"
	movieMiddleware "github.com/Flook0147/netflix_like/internal/movie/middleware"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type Handlers struct {
	Auth         *authhttp.AuthHandler
	User         *userhttp.UserHandler
	Subscription *subhttp.SubscriptionHandler
	Movie        *moviehttp.MovieHandler
}

func SetupRoutes(app *fiber.App, h Handlers) {

	// ========================
	// GLOBAL
	// ========================
	app.Use("/videos", movieMiddleware.VideoAuthMiddleware())
	app.Use("/videos", static.New("./videos"))

	// ========================
	// API GROUP
	// ========================
	api := app.Group("/api")

	// 🔓 public
	authhttp.RegisterRoutes(api, h.Auth)

	// 🔐 protected
	protected := api.Group("/", authMiddleware.JWTMiddleware())

	userhttp.RegisterRoutes(protected, h.User)
	subhttp.RegisterSubscriptionRoutes(protected, h.Subscription)

	// 🎬 movie (split public/protected)
	moviehttp.RegisterMovieRoutes(api, protected, h.Movie)
}
