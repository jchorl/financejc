package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
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
	c.Response().Header().Add("Link", fmt.Sprintf("<%s?start=%s>; rel=\"next\"", u.String(), entity.Next()))
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

// toContext is supposed to take a context/middleware injected value
// from whatever web framework is being used and convert it to a
// Go context.Context that everything below the handlers can understand.
func toContext(ctx echo.Context) context.Context {
	ret := context.Background()
	for _, ctx_key := range constants.CTX_KEYS {
		switch ctx_key {
		case constants.CTX_USER_ID:
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
			userIdF := claims["sub"].(float64)
			userId := uint(userIdF)
			ret = context.WithValue(ret, ctx_key, userId)
		default:
			ret = context.WithValue(ret, ctx_key, ctx.Get(ctx_key))
		}
	}

	return ret
}
