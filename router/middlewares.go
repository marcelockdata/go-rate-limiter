package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/marcelockdata/go-rate-limiter/limiter"
)

func InitializeMiddlewares(router *chi.Mux, limiter *limiter.RateLimiter) {
	router.Use(limiter.Middleware)

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
}
