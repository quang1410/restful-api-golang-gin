package app

import (
	"galvin/lession05-exercise-user-management/internal/handler"
	"galvin/lession05-exercise-user-management/internal/repository"
	"galvin/lession05-exercise-user-management/internal/routes"
	"galvin/lession05-exercise-user-management/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)
	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
