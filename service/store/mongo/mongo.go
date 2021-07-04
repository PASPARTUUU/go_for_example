package mongo

import (
	"context"
	"log"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/config"
	"github.com/PASPARTUUU/go_for_example/service/store/repo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	DB *mongo.Database
	//--
	User repo.User

	MongoUser MongoUser
	//--
	cfg config.Mongo
}

func NewConnect(cfg config.Mongo) (*Mongo, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.DBName)

	mongo := Mongo{
		DB:        db,
		User:      NewUserRepo(db),
		MongoUser: NewUserRepo(db),
		cfg:       cfg,
	}

	return &mongo, nil
}
