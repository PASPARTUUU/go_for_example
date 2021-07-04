package mongo

import (
	"context"

	"github.com/PASPARTUUU/go_for_example/service/models"
)

// для сущностей хранимых исключительно в Mongo

type MongoUser interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
