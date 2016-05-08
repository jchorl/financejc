package financejc

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

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
	log.Debugf(c, "starting auth")
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
	log.Debugf(c, "starting auth")
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
	ws.Route(ws.GET("/accountGroup").Filter(loggedInFilter).To(handlers.GetAccountGroups).
		Doc("Get all account groups").
		Operation("GetAccountGroups"))
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
