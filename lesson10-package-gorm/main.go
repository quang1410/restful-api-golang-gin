package main

import (
	"lesson08-prepare-connection/internal/db"
	"lesson08-prepare-connection/internal/handlers"
	"lesson08-prepare-connection/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if err := db.InitDB(); err != nil {
		log.Fatal("unable to connnect to db")
	}

	log.Println(db.DB)

	r := gin.Default()

	userRepository := repository.NewSQLUserRepository(db.DB)
	userHandler := handlers.NewUserHandler(userRepository)

	r.GET("/api/v1/users/:id", userHandler.GetUserById)
	r.POST("/api/v1/users", userHandler.CreateUser)

	r.Run(":8080")
}
