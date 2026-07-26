package repository

import "lesson08-prepare-connection/internal/models"

type UserRepository interface {
	Create(user *models.User)
	FindById(id int)
}
