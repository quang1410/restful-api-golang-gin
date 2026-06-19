package app

import (
	v1handler "project-shopping-cart/internal/handler/v1"
	"project-shopping-cart/internal/repository"
	"project-shopping-cart/internal/routes"
	v1routes "project-shopping-cart/internal/routes/v1"
	v1service "project-shopping-cart/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewSqlUserRepository()
	userService := v1service.NewUserService(userRepo)
	userHandler := v1handler.NewUserHandler(userService)
	userRoutes := v1routes.NewUserRoutes(userHandler)
	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}