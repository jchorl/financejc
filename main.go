package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/robfig/cron.v2"

	"github.com/jchorl/financejc/api/handlers"
	"github.com/jchorl/financejc/api/transaction"
	"github.com/jchorl/financejc/constants"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	db, err := sql.Open(constants.DB_DRIVER, constants.DB_ADDRESS)
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to connect to database")
	}

	es, err := elastic.NewClient(elastic.SetURL(constants.ES_ADDRESS))
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to connect to elasticsearch")
	}
	configureEsIndices(es)

	c := cron.New()
	ctx := context.WithValue(context.Background(), constants.CTX_DB, db)
	c.AddFunc("@daily", func() {
		// ignore the error because it should already be logged in GenRecurringTransactions
		transaction.GenRecurringTransactions(ctx)
	})

	e := echo.New()
	e.Pre(middleware.HTTPSRedirect())
	e.Use(
		middleware.Gzip(),
		middleware.Logger(),
		dbMiddleware(db),
		esMiddleware(es),
	)

	e.File("/", "client/index.html")
	e.Static("/static", "client/dest")

	apiRoutes := e.Group("/api")
	handlers.Init(apiRoutes)

	logrus.Debug("starting server")
	e.Run(standard.WithTLS(":"+os.Getenv("PORT"), "/etc/letsencrypt/live/"+os.Getenv("DOMAIN")+"/fullchain.pem", "/etc/letsencrypt/live/"+os.Getenv("DOMAIN")+"/privkey.pem"))
}

func dbMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(constants.CTX_DB, db)
			return next(c)
		}
	}
}

func esMiddleware(client *elastic.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(constants.CTX_ES, client)
			return next(c)
		}
	}
}

func configureEsIndices(client *elastic.Client) {
	exists, err := client.IndexExists(constants.ES_INDEX).Do(context.Background())
	if err != nil {
		logrus.WithField("error", err).Fatal("failed to check if transactions index exists")
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(constants.ES_INDEX).Do(context.Background())
		if err != nil {
			logrus.WithField("error", err).Fatal("failed to create transactions index")
		}
		if !createIndex.Acknowledged {
			logrus.Fatal("creating transactions index in es was not acknowledged")
		}
	}
}
