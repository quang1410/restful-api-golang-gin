package main

import (
	"galvin/lession05-exercise-user-management/internal/config"
	"galvin/lession05-exercise-user-management/internal/handler"
	"galvin/lession05-exercise-user-management/internal/repository"
	"galvin/lession05-exercise-user-management/internal/routes"
	"galvin/lession05-exercise-user-management/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	//initialize config
	config := config.NewConfig()

	//initialize repository
	userRepo := repository.NewInMemoryUserRepository()
	
	//initialize service
	userService := service.NewUserService(userRepo)	

	//initialize handler
	userHandler := handler.NewUserHandler(userService)
	
	//initialize routes
	userRoutes := routes.NewUserRoutes(userHandler)

	r := gin.Default()

	routes.RegisterRoutes(r, userRoutes)

	if err := r.Run(config.ServerAddress); err != nil {
		panic(err)
	}
}
