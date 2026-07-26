package repository

import (
	"context"
	"lesson08-prepare-connection/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	FindByUuid(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
}
