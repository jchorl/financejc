package main

import (
	"errors"
	"net/http"
	"os"
	"path"

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

func getUserId(unparsed string) (int, error) {
	token, err := jwt.ParseWithClaims(unparsed, &server.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWT_SIGNING_KEY), nil
	})
	if claims, ok := token.Claims.(*server.JWTClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return -1, err
	}
	return -1, NotLoggedIn
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

func staticHandler(req *restful.Request, resp *restful.Response) {
	actual := path.Join("client/dest", req.PathParameter("subpath"))
	http.ServeFile(resp.ResponseWriter, req.Request, actual)
}

func indexHandler(req *restful.Request, resp *restful.Response) {
	actual := path.Join("client", "index.html")
	http.ServeFile(resp.ResponseWriter, req.Request, actual)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	serv, err := server.NewServer(constants.DB_DRIVER, constants.DB_ADDRESS)
	if err != nil {
		logrus.WithField("Error", err).Fatal("Failed to start server")
	}

	apiWs := new(restful.WebService)

	apiWs.
		Path("/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiWs.Route(apiWs.GET("/auth").Filter(loggedOutFilter).To(serv.CheckAuth).
		Doc("Check if a user is authenticated").
		Operation("CheckAuth"))
	apiWs.Route(apiWs.POST("/auth").Filter(loggedOutFilter).To(serv.AuthUser).
		Doc("Authenticate a user").
		Operation("AuthUser").
		Reads(auth.Request{}))
	apiWs.Route(apiWs.GET("/currencies").To(serv.GetCurrencies).
		Doc("Get all currencies").
		Operation("GetCurrencies").
		Writes(struct{ ISO4217 string }{"Name"}))
	apiWs.Route(apiWs.GET("/account").Filter(loggedInFilter).To(serv.GetAccounts).
		Doc("Get all accounts").
		Operation("GetAccounts").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(account.Account{}))
	apiWs.Route(apiWs.POST("/account").Filter(loggedInFilter).To(serv.NewAccount).
		Doc("Create a new account").
		Operation("NewAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	apiWs.Route(apiWs.PUT("/account").Filter(loggedInFilter).To(serv.UpdateAccount).
		Doc("Update account").
		Operation("UpdateAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	apiWs.Route(apiWs.DELETE("/account/{account-id}").Filter(loggedInFilter).To(serv.DeleteAccount).
		Doc("Delete account").
		Operation("DeleteAccount").
		Param(apiWs.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil))
	apiWs.Route(apiWs.GET("/account/{account-id}/transactions").Filter(loggedInFilter).To(serv.GetTransactions).
		Doc("Get all transactions for an account").
		Operation("GetTransactions").
		Param(apiWs.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(transaction.Transaction{}))
	apiWs.Route(apiWs.POST("/account/{account-id}/transactions").Filter(loggedInFilter).To(serv.NewTransaction).
		Doc("Create a new transaction").
		Operation("NewTransaction").
		Param(apiWs.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Reads(transaction.Transaction{}).
		Writes(transaction.Transaction{}))
	apiWs.Route(apiWs.PUT("/transaction").Filter(loggedInFilter).To(serv.UpdateTransaction).
		Doc("Update transaction").
		Operation("UpdateTransaction").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Reads(transaction.Transaction{}).
		Writes(transaction.Transaction{}))
	apiWs.Route(apiWs.DELETE("/transaction/{transaction-id}").Filter(loggedInFilter).To(serv.DeleteTransaction).
		Doc("Delete transaction").
		Operation("DeleteTransaction").
		Param(apiWs.PathParameter("transaction-id", "id of the transaction").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil))
	apiWs.Route(apiWs.POST("/import").Filter(loggedInFilter).To(serv.Transfer).
		Doc("Import a file from QIF").
		Operation("Transfer"))
	restful.Add(apiWs)

	staticWs := new(restful.WebService)
	staticWs.Route(staticWs.GET("").To(indexHandler))
	staticWs.Route(staticWs.GET("/{subpath}").To(staticHandler))
	restful.Add(staticWs)

	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "client/swagger",
	}
	swagger.InstallSwaggerService(config)

	port := os.Getenv("PORT")
	err = http.ListenAndServeTLS(":"+port, "server.pem", "server.key", nil)
	if err != nil {
		logrus.WithField("Error", err).Fatal("could not serve")
	}
}
