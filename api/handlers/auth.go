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

func AuthUser(c echo.Context) error {
	req := new(auth.Request)

	if err := c.Bind(req); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("error decoding auth request")
		return err
	}

	user, err := auth.AuthUser(toContext(c), req.Token)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
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

	return c.JSON(http.StatusOK, user)
}

func Logout(c echo.Context) error {
	cookie := new(echo.Cookie)
	cookie.SetName("Authorization")
	cookie.SetValue("")
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true)
	cookie.SetPath("/")
	cookie.SetExpires(time.Unix(1, 0))
	c.SetCookie(cookie)

	return c.NoContent(http.StatusNoContent)
}
