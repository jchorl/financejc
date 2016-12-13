package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/constants"
)

// Paginated is an interface for entity lists that support pagination
type Paginated interface {
	Next() string
	Values() []interface{}
}

func writePaginatedEntity(c echo.Context, entity Paginated) error {
	u := url.URL{
		RawPath: c.Request().URL.Path,
	}
	u.RawQuery = ""
	c.Response().Header().Add("Link", fmt.Sprintf("<%s?start=%s>; rel=\"next\"", u.String(), entity.Next()))
	if entity.Values() != nil {
		return c.JSON(http.StatusOK, entity.Values())
	}

	return c.JSON(http.StatusOK, []interface{}{})
}

func writeError(c echo.Context, err error) error {
	switch err {
	case constants.ErrNotLoggedIn:
		return c.String(http.StatusUnauthorized, err.Error())
	case constants.ErrForbidden:
		return c.String(http.StatusForbidden, err.Error())
	case constants.ErrBadRequest:
		return c.String(http.StatusBadRequest, err.Error())
	default:
		return c.String(http.StatusInternalServerError, err.Error())
	}
}

func idFromParam(c echo.Context, paramName string) (int, error) {
	IDStr := c.Param(paramName)
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"context":   c,
			"IDStr":     IDStr,
			"paramName": paramName,
		})
		return 0, constants.ErrBadRequest
	}

	return ID, nil
}

// toContext is supposed to take a context/middleware injected value
// from whatever web framework is being used and convert it to a
// Go context.Context that everything below the handlers can understand.
func toContext(ctx echo.Context) context.Context {
	ret := context.Background()
	for _, ctxKey := range constants.CtxKeys {
		switch ctxKey {
		case constants.CtxUserID:
			// we'll catch errors later when we pull vals from ctx
			// for now just make a valid ctx
			user := ctx.Get("user")
			if user == nil {
				continue
			}

			casted, ok := user.(*jwt.Token)
			if !ok {
				continue
			}

			claims := casted.Claims.(jwt.MapClaims)
			userIDF := claims["sub"].(float64)
			userID := uint(userIDF)
			ret = context.WithValue(ret, ctxKey, userID)
		default:
			ret = context.WithValue(ret, ctxKey, ctx.Get(ctxKey))
		}
	}

	return ret
}
