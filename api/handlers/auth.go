package handlers

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/jchorl/financejc/api/auth"
	"github.com/jchorl/financejc/constants"
)

func CheckAuth(c echo.Context) error {
	// since it got passed the middleware, the user is authd
	return c.NoContent(http.StatusNoContent)
}

func AuthUser(c echo.Context) error {
	req := new(auth.Request)

	if err := c.Bind(req); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("error decoding auth request")
		return err
	}

	userId, err := auth.AuthUser(c, req.Token)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(constants.JWT_SIGNING_KEY))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"token":   token,
			"context": c,
		}).Error("error getting signed jwt")
		return err
	}

	cookie := new(echo.Cookie)
	cookie.SetName("Authorization")
	cookie.SetValue(tokenStr)
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true)
	cookie.SetPath("/")
	c.SetCookie(cookie)

	return c.NoContent(http.StatusNoContent)
}
