package financejc

import (
	"errors"
	"net/http"
	"os"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server"
	"github.com/jchorl/financejc/server/account"
	"github.com/jchorl/financejc/server/auth"
	"github.com/jchorl/financejc/server/transaction"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

var NotLoggedIn = errors.New("User is not logged in")

func getUserId(unparsed string) (string, error) {
	token, err := jwt.ParseWithClaims(unparsed, &server.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWT_SIGNING_KEY), nil
	})
	if claims, ok := token.Claims.(*server.JWTClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return "", err
	}
	return "", NotLoggedIn
}

// only accept logged in users
func loggedInFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	cookie, err := request.Request.Cookie("auth")
	if err != nil {
		logrus.Errorf("could not get auth cookie: %+v", err)
	} else {
		userId, err := getUserId(cookie.Value)
		if err == nil {
			request.SetAttribute("userId", userId)
			chain.ProcessFilter(request, response)
			return
		} else if err != NotLoggedIn {
			logrus.Errorf("error parsing jwt: %+v", err)
		}
	}
	response.WriteErrorString(401, "401: Not Authorized")
}

// only accept logged out users
func loggedOutFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	cookie, err := request.Request.Cookie("auth")
	if err == nil {
		_, err := getUserId(cookie.Value)
		if err == nil {
			return
		}
	}
	logrus.Debugf("passed through logged out filter and proceeding")
	chain.ProcessFilter(request, response)
}

func main() {
	serv, err := server.NewServer("postgres", "postgresurl")
	if err != nil {
		logrus.WithField("Error", err).Fatal("Failed to start server")
	}

	ws := new(restful.WebService)

	ws.
		Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/auth").Filter(loggedOutFilter).To(serv.CheckAuth).
		Doc("Check if a user is authenticated").
		Operation("CheckAuth"))
	ws.Route(ws.POST("/auth").Filter(loggedOutFilter).To(serv.AuthUser).
		Doc("Authenticate a user").
		Operation("AuthUser").
		Reads(auth.Request{}))
	ws.Route(ws.GET("/currencies").To(serv.GetCurrencies).
		Doc("Get all currencies").
		Operation("GetCurrencies").
		Writes(struct{ ISO4217 string }{"Name"}))
	ws.Route(ws.GET("/account").Filter(loggedInFilter).To(serv.GetAccounts).
		Doc("Get all accounts").
		Operation("GetAccounts").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(account.Account{}))
	ws.Route(ws.POST("/account").Filter(loggedInFilter).To(serv.NewAccount).
		Doc("Create a new account").
		Operation("NewAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	ws.Route(ws.PUT("/account").Filter(loggedInFilter).To(serv.UpdateAccount).
		Doc("Update account").
		Operation("UpdateAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	ws.Route(ws.DELETE("/account/{account-id}").Filter(loggedInFilter).To(serv.DeleteAccount).
		Doc("Delete account").
		Operation("DeleteAccount").
		Param(ws.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil))
	ws.Route(ws.GET("/account/{account-id}/transactions").Filter(loggedInFilter).To(serv.GetTransactions).
		Doc("Get all transactions for an account").
		Operation("GetTransactions").
		Param(ws.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(transaction.Transaction{}))
	ws.Route(ws.POST("/account/{account-id}/transactions").Filter(loggedInFilter).To(serv.NewTransaction).
		Doc("Create a new transaction").
		Operation("NewTransaction").
		Param(ws.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Reads(transaction.Transaction{}).
		Writes(transaction.Transaction{}))
	ws.Route(ws.PUT("/transaction").Filter(loggedInFilter).To(serv.UpdateTransaction).
		Doc("Update transaction").
		Operation("UpdateTransaction").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Reads(transaction.Transaction{}).
		Writes(transaction.Transaction{}))
	ws.Route(ws.DELETE("/transaction/{transaction-id}").Filter(loggedInFilter).To(serv.DeleteTransaction).
		Doc("Delete transaction").
		Operation("DeleteTransaction").
		Param(ws.PathParameter("transaction-id", "id of the transaction").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil))
	ws.Route(ws.POST("/import").Filter(loggedInFilter).To(serv.Transfer).
		Doc("Import a file from QIF").
		Operation("Transfer"))
	restful.Add(ws)

	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "client/swagger",
	}
	swagger.InstallSwaggerService(config)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
}
