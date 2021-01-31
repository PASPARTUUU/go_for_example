package store

import (
	"context"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/config"
	"github.com/PASPARTUUU/go_for_example/service/store/postgres"
	"github.com/PASPARTUUU/go_for_example/tools/errpath"
	"github.com/sirupsen/logrus"
)

// Store - contains all repositories
type Store struct {
	Pg *postgres.Pg

	// ---
	logger *logrus.Logger
	config *config.Config
}

// New - creates new store
func New(ctx context.Context, cfg *config.Config, logger *logrus.Logger) (*Store, error) {
	var store Store

	store.config = cfg
	store.logger = logger

	// connect to postgres
	pgConn, err := postgres.NewConnect(cfg.Postgres)
	if err != nil {
		return &store, errpath.Err(err)
	}
	store.Pg = pgConn

	if pgConn != nil {
		go store.keepAlivePg()
	}

	return &store, nil
}

// KeepAlivePollPeriod - is a Pg keepalive check time period
const KeepAlivePollPeriod = time.Second * 60

// keepAlivePg - makes sure PostgreSQL is alive and reconnects if needed
func (store *Store) keepAlivePg() {
	log := store.logger.WithField("event", "KeepAlivePg")
	var err error
	for {
		// Check if PostgreSQL is alive every 'KeepAlivePollPeriod' seconds
		time.Sleep(KeepAlivePollPeriod)
		lostConnect := false
		if store.Pg == nil {
			lostConnect = true
		} else if _, err = store.Pg.DB.Exec("SELECT 1"); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		log.Warnln(errpath.Infof("Lost PostgreSQL connection. Restoring..."))

		store.Pg, err = postgres.NewConnect(store.config.Postgres)
		if err != nil {
			log.Errorln(errpath.Err(err))
			continue
		}
		log.Infoln(errpath.Infof("PostgreSQL reconnected"))
	}
}
