package financejc

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"account"
	"auth"
	"credentials"
	"handlers"
)

var NotLoggedIn = errors.New("User is not logged in")

func getGaeURL() string {
	if appengine.IsDevAppServer() {
		return "http://localhost:8080"
	} else {
		return "https://financejc.appspot.com"
	}
}

func getUserId(unparsed string) (string, error) {
	token, err := jwt.Parse(unparsed, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(credentials.JWT_SIGNING_KEY), nil
	})
	if err == nil && token.Valid {
		return token.Claims["userId"].(string), nil
	} else if err != nil {
		return "", err
	}
	return "", NotLoggedIn
}

// only accept logged in users
func loggedInFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(request.Request)
	cookie, err := request.Request.Cookie("auth")
	if err != nil {
		log.Errorf(c, "could not get auth cookie: %+v", err)
	} else {
		userId, err := getUserId(cookie.Value)
		if err == nil {
			request.SetAttribute("userId", userId)
			chain.ProcessFilter(request, response)
			return
		} else if err != NotLoggedIn {
			log.Errorf(c, "error parsing jwt: %+v", err)
		}
	}
	response.WriteErrorString(401, "401: Not Authorized")
}

// only accept logged out users
func loggedOutFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	c := appengine.NewContext(request.Request)
	cookie, err := request.Request.Cookie("auth")
	if err == nil {
		_, err := getUserId(cookie.Value)
		if err == nil {
			return
		}
	}
	log.Debugf(c, "passed through logged out filter and proceeding")
	chain.ProcessFilter(request, response)
}

func init() {
	ws := new(restful.WebService)

	ws.
		Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/auth").Filter(loggedOutFilter).To(handlers.AuthUser).
		Doc("Authenticate a user").
		Operation("AuthUser").
		Reads(auth.Request{}))
	ws.Route(ws.GET("/currencies").To(handlers.GetCurrencies).
		Doc("Get all currencies").
		Operation("GetCurrencies").
		Writes(struct{ ISO4217 string }{"Name"}))
	ws.Route(ws.GET("/account").Filter(loggedInFilter).To(handlers.GetAccounts).
		Doc("Get all accounts").
		Operation("GetAccounts").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Writes(account.Account{}))
	ws.Route(ws.POST("/account").Filter(loggedInFilter).To(handlers.NewAccount).
		Doc("Create a new account").
		Operation("NewAccount").
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	ws.Route(ws.PUT("/account/{account-id}").Filter(loggedInFilter).To(handlers.UpdateAccount).
		Doc("Update account").
		Operation("UpdateAccount").
		Param(ws.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil).
		Reads(account.Account{}).
		Writes(account.Account{}))
	ws.Route(ws.DELETE("/account/{account-id}").Filter(loggedInFilter).To(handlers.DeleteAccount).
		Doc("Delete account").
		Operation("DeleteAccount").
		Param(ws.PathParameter("account-id", "id of the account").DataType("string")).
		Returns(http.StatusUnauthorized, "User is not authorized", nil).
		Returns(http.StatusForbidden, "Invalid currency", nil))
	restful.Add(ws)

	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		WebServicesUrl:  getGaeURL(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "client/swagger",
	}
	swagger.InstallSwaggerService(config)
}
