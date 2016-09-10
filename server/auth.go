package server

import (
	"net/http"
	"strconv"

	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/server/auth"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
)

type JWTClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func (s server) CheckAuth(request *restful.Request, response *restful.Response) {
	// if already passed the logged out filter, return 401
	response.WriteErrorString(401, "401: Not Authorized")
}

func (s server) AuthUser(request *restful.Request, response *restful.Response) {
	var req auth.Request

	err := request.ReadEntity(&req)
	if err != nil {
		logrus.WithField("Error", err).Error("error decoding auth request")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	userId, err := auth.AuthUser(s.Context(), req.Token)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	logrus.WithField("User ID", userId).Debug("authd with user id")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{strconv.Itoa(userId), jwt.StandardClaims{}})
	tokenStr, err := token.SignedString([]byte(constants.JWT_SIGNING_KEY))
	if err != nil {
		logrus.WithField("Error", err).Error("error getting signed jwt")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	logrus.WithField("Token", tokenStr).Debug("got token from jwt")

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
