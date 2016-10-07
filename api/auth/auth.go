package auth

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	"google.golang.org/api/oauth2/v2"

	"github.com/jchorl/financejc/api/user"
)

type Request struct {
	Token string `json:"token"` // provided by google
}

func AuthUser(c context.Context, token string) (uint, error) {
	googleId, err := getGoogleIDFromToken(token)
	if err != nil {
		return 0, err
	}

	userId, err := user.FindOrCreateByGoogleId(c, googleId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func getGoogleIDFromToken(token string) (string, error) {
	service, err := oauth2.New(http.DefaultClient)
	tokenInfoCall := service.Tokeninfo()
	tokenInfoCall.IdToken(token)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"token": token,
		}).Error("error calling google to validate token")
		return "", err
	}

	return tokenInfo.UserId, nil
}
