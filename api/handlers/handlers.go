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

// Init takes an echo group and registers all the api handlers on it
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
	api.GET("/account/:accountId/templates", GetTemplates, jwtMiddleware)
	api.GET("/account/:accountId/transactions/query", QueryES, jwtMiddleware)
	api.POST("/account/:accountId/transactions", NewTransaction, jwtMiddleware)
	api.POST("/account/:accountId/recurringTransactions", NewRecurringTransaction, jwtMiddleware)
	api.POST("/account/:accountId/templates", NewTemplate, jwtMiddleware)

	api.PUT("/transaction", UpdateTransaction, jwtMiddleware)
	api.PUT("/recurringTransaction", UpdateRecurringTransaction, jwtMiddleware)
	api.PUT("/template", UpdateTemplate, jwtMiddleware)
	api.DELETE("/transaction/:transactionId", DeleteTransaction, jwtMiddleware)
	api.DELETE("/recurringTransaction/:recurringTransactionId", DeleteRecurringTransaction, jwtMiddleware)
	api.DELETE("/template/:templateId", DeleteTemplate, jwtMiddleware)
	api.GET("/transaction/pushAllToES", PushAllToES, jwtMiddleware)
	api.GET("/transaction/genRecurring", GenRecurringTransactions, jwtMiddleware)

	api.GET("/user", GetUser, jwtMiddleware)

	api.POST("/import", Transfer, jwtMiddleware)
}
