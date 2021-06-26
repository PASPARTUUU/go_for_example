package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/models"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/store"
	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
)

// UserController -
type UserController struct {
	storage *store.Store
	logger  *logrus.Logger
}

// NewUsers - creates a new user controller.
func NewUsers(storage *store.Store, logger *logrus.Logger) *UserController {
	return &UserController{
		storage: storage,
		logger:  logger,
	}
}

// Create - creates new user
func (ctr *UserController) Create(c echo.Context) error {
	ctx := c.Request().Context()
	log := ctr.logger.WithContext(ctx).WithField("event", "create user")
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		log.Errorln(errpath.Err(err))
		return c.JSON(http.StatusBadRequest, errpath.Err(err, "could not decode user data").Error())
	}

	newUser, err := ctr.storage.Pg.User.CreateUser(ctx, &user)
	if err != nil {
		log.Errorln(errpath.Err(err))
		return c.JSON(http.StatusBadRequest, errpath.Err(err).Error())
	}

	log.Infoln("user created")

	return c.JSON(http.StatusCreated, newUser)
}

// Get -
func (ctr *UserController) Get(c echo.Context) error {
	ctx := c.Request().Context()
	log := ctr.logger.WithContext(ctx)

	uuid := c.Param("uuid")

	user, err := ctr.storage.Pg.User.GetUser(ctx, uuid)
	if err != nil {
		log.Errorln(errpath.Err(err))
		return c.JSON(http.StatusBadRequest, errpath.Err(err).Error())
	}

	log.Infof(errpath.Infof("Done!!!"))
	return c.JSON(http.StatusCreated, user)
}

// Update -
func (ctr *UserController) Update(c echo.Context) error {
	ctx := c.Request().Context()
	log := ctr.logger.WithContext(ctx)

	var user models.User
	err := c.Bind(&user)
	if err != nil {
		log.Errorln(errpath.Err(err))
		return c.JSON(http.StatusBadRequest, errpath.Err(err, "could not decode user data").Error())
	}
	user.ID = c.Param("uuid")

	uUser, err := ctr.storage.Pg.User.UpdateUser(ctx, &user)
	if err != nil {
		log.Errorln(errpath.Err(err))
		return c.JSON(http.StatusBadRequest, errpath.Err(err).Error())
	}

	log.Infof(errpath.Infof("Done!!!"))
	return c.JSON(http.StatusCreated, uUser)
}

// Delete -
func (ctr *UserController) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	log := ctr.logger.WithContext(ctx)

	uuid := c.Param("uuid")

	err := ctr.storage.Pg.User.DeleteUser(ctx, uuid)
	if err != nil {
		log.Errorln(errpath.Err(err))
		return c.JSON(http.StatusBadRequest, errpath.Err(err).Error())
	}

	log.Infof(errpath.Infof("Done!!!"))
	return c.JSON(http.StatusCreated, "User Deleted")
}
