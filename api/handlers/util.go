package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo"

	"github.com/jchorl/financejc/constants"
)

type Paginated interface {
	Next() string
	Values() []interface{}
}

func writePaginatedEntity(c echo.Context, entity Paginated) error {
	u := url.URL{
		RawPath: c.Request().URL().Path(),
	}
	u.RawQuery = ""
	c.Response().Header().Add("Link", fmt.Sprintf("<%s?start=%s>; rel=\"next\"", u, entity.Next()))
	if entity.Values() != nil {
		return c.JSON(http.StatusOK, entity.Values())
	}

	return c.JSON(http.StatusOK, []interface{}{})
}

func writeError(c echo.Context, err error) error {
	switch err {
	case constants.NotLoggedIn:
		return c.String(http.StatusUnauthorized, err.Error())
	case constants.Forbidden:
		return c.String(http.StatusForbidden, err.Error())
	case constants.BadRequest:
		return c.String(http.StatusBadRequest, err.Error())
	default:
		return c.String(http.StatusInternalServerError, err.Error())
	}
}
