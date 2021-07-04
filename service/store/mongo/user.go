package mongo

import (
	"context"

	"github.com/PASPARTUUU/go_for_example/service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Age      int                `bson:"password"`
}

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db.Collection("user_collection"),
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	model := toMongoUser(user)
	res, err := r.db.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (r UserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	user := new(User)
	err := r.db.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(user)

	if err != nil {
		return nil, err
	}

	return toModel(user), nil
}

func (r UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// TODO: implement
	return nil, nil
}
func (r UserRepository) DeleteUser(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}

// -------------------------------------------------

func toMongoUser(u *models.User) *User {
	return &User{
		Username: u.Name,
		Age:      u.Age,
	}
}

func toModel(u *User) *models.User {
	return &models.User{
		ID:   u.ID.Hex(),
		Name: u.Username,
		Age:  u.Age,
	}
}
