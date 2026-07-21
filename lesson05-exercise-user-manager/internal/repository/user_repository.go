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

func (r *InMemoryUserRepository) FindAll() ([]models.User, error) {
	return r.users, nil
}

func (r *InMemoryUserRepository) Create(user models.User) error {
	r.users = append(r.users, user)
	return nil
}

func (r *InMemoryUserRepository) FindByUUID(uuid string) (models.User, bool) {
	for _, user := range r.users {
		if user.UUID == uuid {
			return user, true
		}
	}
	return models.User{}, false
}

func (r *InMemoryUserRepository) Update(uuid string, updatedUser models.User) error {
	for i, user := range r.users {
		if user.UUID == uuid {
			r.users[i] = updatedUser
			return nil
		}
	}
	return nil
}

func (r *InMemoryUserRepository) Delete(uuid string) error {
	for i, user := range r.users {
		if user.UUID == uuid {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (models.User, bool) {
	for _, user := range r.users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}
