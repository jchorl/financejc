package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/robfig/cron.v2"

	"github.com/jchorl/financejc/api/handlers"
	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/api/transfer/batchTransfer"
	"github.com/jchorl/financejc/constants"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	db, err := sql.Open(constants.DbDriver, constants.DbAddress)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to connect to database")
	}

	es, err := elastic.NewClient(elastic.SetURL(constants.EsAddress))
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to connect to elasticsearch")
	}
	configureEsIndices(es)

	c := cron.New()
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxDB, db)
	ctx = context.WithValue(ctx, constants.CtxES, es)
	ctx = context.WithValue(ctx, constants.CtxInternalReq, true)
	c.AddFunc("@daily", func() {
		// ignore the error because it should already be logged in GenRecurringTransactions
		transaction.GenRecurringTransactions(ctx)
	})
	c.AddFunc("@daily", func() {
		// ignore the error because it should already be logged in BackupToGCS
		batchTransfer.BackupToGCS(ctx)
	})
	c.Start()

	e := echo.New()
	e.Use(
		middleware.Gzip(),
		middleware.Logger(),
		dbMiddleware(db),
		esMiddleware(es),
	)

	apiRoutes := e.Group("/api")
	handlers.Init(apiRoutes)

	logrus.Debug("starting server")
	e.Start(":" + os.Getenv("PORT"))
}

func dbMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(constants.CtxDB, db)
			return next(c)
		}
	}
}

func esMiddleware(client *elastic.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(constants.CtxES, client)
			return next(c)
		}
	}
}

func configureEsIndices(client *elastic.Client) {
	exists, err := client.IndexExists(constants.ESIndex).Do(context.Background())
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to check if transactions index exists")
	}
	if !exists {
		// Create a new index.
		if err := transaction.InitES(client); err != nil {
			logrus.WithError(err).Fatal("could not configure ES")
		}
	}
}
