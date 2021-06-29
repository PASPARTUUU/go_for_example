package handler

import (
	"github.com/PASPARTUUU/go_for_example/service/rabbitmq/rabpub"
	"github.com/PASPARTUUU/go_for_example/service/store"
	"github.com/sirupsen/logrus"
)

// Handler -
type Handler struct {
	Storage *store.Store

	UserHandler *UserController

	Rabbit *rabpub.Publisher
}

// New -
func New(storage *store.Store, rabb *rabpub.Publisher, log *logrus.Logger) *Handler {
	return &Handler{
		Storage:     storage,
		UserHandler: NewUsers(storage, log),
		Rabbit:      rabb,
	}
}
