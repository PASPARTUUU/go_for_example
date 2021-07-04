package server

import (
	"net/http"

	"github.com/PASPARTUUU/go_for_example/service/handler"
	"github.com/labstack/echo/v4"
)

const (
	gatePrefix = "/gate/v1"
	rpcPrefix  = "/rpc/v1"
	apiPrefix  = "/api/v1"
)

// Rest -
type Rest struct {
	Router *echo.Echo
}

// RestInit -
func RestInit(e *echo.Echo, hndl *handler.Handler) {
	var rest = Rest{
		Router: e,
	}

	rest.Route(hndl)
}

// Route -
func (r *Rest) Route(hndl *handler.Handler) {

	open := r.Router.Group(apiPrefix)

	open.POST("/bung", bung)

	open.GET("/user/:id", hndl.UserHandler.Get)
	open.POST("/user", hndl.UserHandler.Create)
	open.PUT("/user/:id", hndl.UserHandler.Update)
	open.DELETE("/user/:id", hndl.UserHandler.Delete)

}

func bung(c echo.Context) error {
	return c.JSON(http.StatusOK, "normalin normalin!!!")
}
