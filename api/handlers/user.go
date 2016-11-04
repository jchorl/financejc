package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/user"
)

func GetUser(c echo.Context) error {
	user, err := user.Get(toContext(c))
	if err != nil {
		return writeError(c, err)
	}

	return c.JSON(http.StatusOK, user)
}
