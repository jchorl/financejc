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

// AuthUser authenticates a user based on a auth.Request JSON body
func AuthUser(c echo.Context) error {
	req := new(auth.Request)

	if err := c.Bind(req); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"context": c,
		}).Error("error decoding auth request")
		return err
	}

	user, err := auth.LoginByGoogleToken(toContext(c), req.Token)
	if err != nil {
		return constants.ErrNotLoggedIn
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(constants.JwtSigningKey))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"token":   token,
			"context": c,
		}).Error("error getting signed jwt")
		return constants.ErrNotLoggedIn
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = tokenStr
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

// Logout logs the authenticated user out
func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.Expires = time.Unix(1, 0)
	c.SetCookie(cookie)

	return c.NoContent(http.StatusNoContent)
}
