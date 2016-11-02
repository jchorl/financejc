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
	api.GET("/auth", CheckAuth, jwtMiddleware)
	api.POST("/auth", AuthUser)
	api.POST("/auth/logout", Logout)

	api.GET("/currencies", GetCurrencies)

	api.GET("/account", GetAccounts, jwtMiddleware)
	api.POST("/account", NewAccount, jwtMiddleware)
	api.PUT("/account", UpdateAccount, jwtMiddleware)
	api.DELETE("/account/:accountId", DeleteAccount, jwtMiddleware)

	api.GET("/account/:accountId/transactions", GetTransactions, jwtMiddleware)
	api.GET("/account/:accountId/transactions/query", QueryES, jwtMiddleware)
	api.POST("/account/:accountId/transactions", NewTransaction, jwtMiddleware)

	api.PUT("/transaction", UpdateTransaction, jwtMiddleware)
	api.DELETE("/transaction/:transactionId", DeleteTransaction, jwtMiddleware)

	api.POST("/import", Transfer, jwtMiddleware)
}
