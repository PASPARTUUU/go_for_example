package repo

import (
	"context"

	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/models"
)

// User - is a store for users
type User interface {
	GetUser(ctx context.Context, uuid string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, uuid string) error
}
