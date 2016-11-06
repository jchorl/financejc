package handlers

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/jchorl/financejc/constants"
)

var (
	jwtMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(constants.JWT_SIGNING_KEY),
		TokenLookup: "cookie:Authorization",
	})
)

func Init(api *echo.Group) {
	api.POST("/auth", AuthUser)
	api.POST("/auth/logout", Logout)

	api.GET("/currencies", GetCurrencies)

	api.GET("/account", GetAccounts, jwtMiddleware)
	api.POST("/account", NewAccount, jwtMiddleware)
	api.PUT("/account", UpdateAccount, jwtMiddleware)
	api.DELETE("/account/:accountId", DeleteAccount, jwtMiddleware)

	api.GET("/account/:accountId/transactions", GetTransactions, jwtMiddleware)
	api.GET("/account/:accountId/recurringTransactions", GetRecurringTransactions, jwtMiddleware)
	api.GET("/account/:accountId/transactionTemplates", GetTransactionTemplates, jwtMiddleware)
	api.GET("/account/:accountId/transactions/query", QueryES, jwtMiddleware)
	api.POST("/account/:accountId/transactions", NewTransaction, jwtMiddleware)
	api.POST("/account/:accountId/recurringTransactions", NewRecurringTransaction, jwtMiddleware)
	api.POST("/account/:accountId/transactionTemplates", NewTransactionTemplate, jwtMiddleware)

	api.PUT("/transaction", UpdateTransaction, jwtMiddleware)
	api.PUT("/recurringTransaction", UpdateRecurringTransaction, jwtMiddleware)
	api.PUT("/transactionTemplate", UpdateTransactionTemplate, jwtMiddleware)
	api.DELETE("/transaction/:transactionId", DeleteTransaction, jwtMiddleware)
	api.DELETE("/recurringTransaction/:recurringTransactionId", DeleteRecurringTransaction, jwtMiddleware)
	api.DELETE("/transactionTemplate/:transactionTemplateId", DeleteTransactionTemplate, jwtMiddleware)

	api.GET("/user", GetUser, jwtMiddleware)

	api.POST("/import", Transfer, jwtMiddleware)
}
