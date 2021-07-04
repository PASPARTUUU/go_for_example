package store

import (
	"context"

	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/PASPARTUUU/go_for_example/service/config"
	"github.com/PASPARTUUU/go_for_example/service/store/mongo"
	"github.com/PASPARTUUU/go_for_example/service/store/postgres"
	"github.com/sirupsen/logrus"
)

// Store - contains all repositories
type Store struct {
	Pg    *postgres.Pg
	Mongo *mongo.Mongo
	// ---
	logger *logrus.Logger
	config *config.Config
}

// New - creates new store
func New(ctx context.Context, cfg *config.Config, logger *logrus.Logger) (*Store, error) {
	var err error
	var store Store

	store.config = cfg
	store.logger = logger

	// connect to postgres
	store.Pg, err = postgres.NewConnect(cfg.Postgres)
	if err != nil {
		return &store, errpath.Err(err)
	}

	// connect to mongo
	store.Mongo, err = mongo.NewConnect(cfg.Mongo)
	if err != nil {
		return &store, errpath.Err(err)
	}

	return &store, nil
}
