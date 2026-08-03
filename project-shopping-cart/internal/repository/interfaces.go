package repository

import (
	"context"
	"project-shopping-cart/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(uuid string)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(uuid string)
	FindByEmail(email string)
}