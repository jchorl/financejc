package server

import (
	"net/http"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server/auth"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
)

var NotLoggedIn = errors.New("User is not logged in")

type JWTClaims struct {
	UserId int `json:"userId"`
	*jwt.StandardClaims
}

func (s server) CheckAuth(request *restful.Request, response *restful.Response) {
	// if already passed the logged out filter, return 401
	response.WriteErrorString(401, "401: Not Authorized")
}

func (s server) AuthUser(request *restful.Request, response *restful.Response) {
	var req auth.Request

	err := request.ReadEntity(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":   err,
			"Request": request,
		}).Error("error decoding auth request")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	userId, err := auth.AuthUser(s.Context(), req.Token)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{userId, jwt.StandardClaims{}})
	tokenStr, err := token.SignedString([]byte(constants.JWT_SIGNING_KEY))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"Token": token,
		}).Error("error getting signed jwt")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	cookie := http.Cookie{
		Name:     "auth",
		Value:    tokenStr,
		MaxAge:   604800,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(response, &cookie)
}

func getUserId(unparsedJWT string) (int, error) {
	token, err := jwt.ParseWithClaims(unparsedJWT, &server.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
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
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Error("could not get auth cookie")
	} else {
		userId, err := getUserId(cookie.Value)
		if err == nil {
			request.SetAttribute("userId", userId)
			chain.ProcessFilter(request, response)
			return
		} else if err != NotLoggedIn {
			logrus.WithFields(logrus.Fields{
				"Error":  err,
				"Cookie": cookie,
			}).Error("error parsing jwt")
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
	chain.ProcessFilter(request, response)
}
