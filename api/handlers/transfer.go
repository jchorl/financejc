package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/transfer/batchTransfer"
	"github.com/jchorl/financejc/api/transfer/userTransfer"
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
	if err := userTransfer.Import(toContext(c), src); err != nil {
		return writeError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Export exports all system data
func Export(c echo.Context) error {
	results, err := batchTransfer.Export(toContext(c))
	if err != nil {
		return writeError(c, err)
	}

	return c.String(http.StatusOK, results)
}
