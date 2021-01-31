package postgres

import (
	"context"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/models"
	"github.com/PASPARTUUU/go_for_example/tools/errpath"

	"github.com/go-pg/pg"
	"github.com/gofrs/uuid"
)

// UserPgRepo -
type UserPgRepo struct {
	DB *pg.DB
}

// NewUserRepo -
func NewUserRepo(db *pg.DB) *UserPgRepo {
	return &UserPgRepo{DB: db}
}

// GetUser - retrieves user from Postgres
func (repo *UserPgRepo) GetUser(ctx context.Context, uuid string) (*models.User, error) {
	user := &models.User{}
	err := repo.DB.ModelContext(ctx, user).
		Where("uuid = ?", uuid).
		Select()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, errpath.Err(err)
	}
	return user, nil
}

// CreateUser - creates user in Postgres
func (repo *UserPgRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, errpath.Err(err)
	}
	user.UUID = uuid.String()

	_, err = repo.DB.ModelContext(ctx, user).
		Insert()
	if err != nil {
		return nil, errpath.Err(err)
	}

	return user, nil
}

// UpdateUser - updates user in Postgres
func (repo *UserPgRepo) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()
	_, err := repo.DB.ModelContext(ctx, user).
		WherePK().
		Returning("*").
		UpdateNotNull()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, errpath.Err(err)
	}

	return user, nil
}

// DeleteUser - set deleted time for user in Postgres
// rus: устанавливает время удаления пользователя
func (repo *UserPgRepo) DeleteUser(ctx context.Context, uuid string) error {
	_, err := repo.DB.ModelContext(ctx, (*models.User)(nil)).
		Where("uuid = ?", uuid).
		Set("deleted_at = ?", time.Now()).
		Update()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		return errpath.Err(err)
	}
	return nil
}
