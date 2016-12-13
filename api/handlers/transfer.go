package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/transfer"
)

// Transfer manages importing files
func Transfer(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("could not get file from context for import")
		return writeError(c, err)
	}

	src, err := file.Open()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("could not open uploaded file for import")
		return writeError(c, err)
	}
	defer src.Close()
	if err := transfer.Import(toContext(c), src); err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
