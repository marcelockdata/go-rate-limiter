package router

import (
	"github.com/go-chi/chi"
	"github.com/marcelockdata/go-rate-limiter/limiter"
)

func Init() {
	router := chi.NewRouter()
	rate_limiter := limiter.InitializeRateLimiters()

	InitializeMiddlewares(router, rate_limiter)
	InitializeRoutes(router)
	InitilizeServer(router)
}
