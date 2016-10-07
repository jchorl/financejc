package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/transfer"
)

func Transfer(c echo.Context) error {
	if err := transfer.AutoImport(c); err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
