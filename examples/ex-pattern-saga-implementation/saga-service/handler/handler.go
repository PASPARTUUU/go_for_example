package handler

import (
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/store"
	"github.com/sirupsen/logrus"
)

// Handler -
type Handler struct {
	UserHandler *UserController
}

// New -
func New(storage *store.Store, log *logrus.Logger) *Handler {
	return &Handler{
		UserHandler: NewUsers(storage, log),
	}
}
