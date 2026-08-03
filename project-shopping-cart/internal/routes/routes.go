package routes

import (
	"project-shopping-cart/internal/middleware"
	"project-shopping-cart/internal/utils"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {

	httpLogger := utils.NewLoggerWithPath("http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("recovery.log", "warning")
	rateLimterLogger := utils.NewLoggerWithPath("rate_limiter.log", "warning")

	r.Use(
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.TraceMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
		middleware.RateLimiterMiddleware(rateLimterLogger),
)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}