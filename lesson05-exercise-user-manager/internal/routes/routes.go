package routes

import (
	"galvin/lession05-exercise-user-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	r.Use(
		middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
		middleware.RateLimiterMiddleware(),
	)
	apiGroup := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(apiGroup)
	}
}
