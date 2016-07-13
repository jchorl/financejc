package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"

	"github.com/jchorl/financejc/auth"
	"github.com/jchorl/financejc/credentials"
)

type JWTClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func CheckAuth(request *restful.Request, response *restful.Response) {
	// if already passed the logged out filter, return 401
	response.WriteErrorString(401, "401: Not Authorized")
}

func AuthUser(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	var req auth.Request

	err := request.ReadEntity(&req)
	if err != nil {
		log.Errorf(c, "error decoding auth request: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	userId, err := auth.AuthUser(c, req)
	if err != nil {
		log.Errorf(c, "error authenticating user: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	log.Debugf(c, "authed with userId: %s", userId)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{userId, jwt.StandardClaims{}})
	tokenStr, err := token.SignedString([]byte(credentials.JWT_SIGNING_KEY))
	if err != nil {
		log.Errorf(c, "error getting signed jwt: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	log.Debugf(c, "Token: %s", tokenStr)

	cookie := http.Cookie{
		Name:   "auth",
		Value:  tokenStr,
		MaxAge: 604800,
		// TODO: uncomment for prod
		// Secure:   true,
		// HttpOnly: true,
	}
	http.SetCookie(response, &cookie)
}
