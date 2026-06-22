package repository

import "galvin/lession05-exercise-user-management/internal/models"

type InMemoryUserRepository struct {
	users []models.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: []models.User{},
	}
}