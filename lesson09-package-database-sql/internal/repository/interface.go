package repository

import "lesson08-prepare-connection/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindById(id int, user *models.User) error
}
