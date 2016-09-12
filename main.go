package main

import (
	"net/http"
	"os"
	"path"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server"
	"github.com/jchorl/financejc/server/account"
	"github.com/jchorl/financejc/server/auth"
	"github.com/jchorl/financejc/server/transaction"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

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

	apiWs.Route(apiWs.GET("/auth").Filter(server.LoggedOutFilter).To(serv.CheckAuth).
		Doc("Check if a user is authenticated").
		Operation("CheckAuth"))
	apiWs.Route(apiWs.POST("/auth").Filter(server.LoggedOutFilter).To(serv.AuthUser).
		Doc("Authenticate a user").
		Operation("AuthUser").
		Reads(auth.Request{}))
	apiWs.Route(apiWs.GET("/currencies").To(serv.GetCurrencies).
		Doc("Get all currencies").
		Operation("GetCurrencies").
		Writes(struct{ ISO4217 string }{"Name"}))
	apiWs.Route(apiWs.GET("/account").Filter(server.LoggedInFilter).To(serv.GetAccounts).
		Doc("Get all accounts").
		Operation("GetAccounts").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(account.Account{}))
	apiWs.Route(apiWs.POST("/account").Filter(server.LoggedInFilter).To(serv.NewAccount).
		Doc("Create a new account").
		Operation("NewAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	apiWs.Route(apiWs.PUT("/account").Filter(server.LoggedInFilter).To(serv.UpdateAccount).
		Doc("Update account").
		Operation("UpdateAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	apiWs.Route(apiWs.DELETE("/account/{account-id}").Filter(server.LoggedInFilter).To(serv.DeleteAccount).
		Doc("Delete account").
		Operation("DeleteAccount").
		Param(apiWs.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil))
	apiWs.Route(apiWs.GET("/account/{account-id}/transactions").Filter(server.LoggedInFilter).To(serv.GetTransactions).
		Doc("Get all transactions for an account").
		Operation("GetTransactions").
		Param(apiWs.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(transaction.Transaction{}))
	apiWs.Route(apiWs.POST("/account/{account-id}/transactions").Filter(server.LoggedInFilter).To(serv.NewTransaction).
		Doc("Create a new transaction").
		Operation("NewTransaction").
		Param(apiWs.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Reads(transaction.Transaction{}).
		Writes(transaction.Transaction{}))
	apiWs.Route(apiWs.PUT("/transaction").Filter(server.LoggedInFilter).To(serv.UpdateTransaction).
		Doc("Update transaction").
		Operation("UpdateTransaction").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Reads(transaction.Transaction{}).
		Writes(transaction.Transaction{}))
	apiWs.Route(apiWs.DELETE("/transaction/{transaction-id}").Filter(server.LoggedInFilter).To(serv.DeleteTransaction).
		Doc("Delete transaction").
		Operation("DeleteTransaction").
		Param(apiWs.PathParameter("transaction-id", "id of the transaction").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil))
	apiWs.Route(apiWs.POST("/import").Filter(server.LoggedInFilter).To(serv.Transfer).
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
