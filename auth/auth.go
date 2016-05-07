package auth

import (
	"golang.org/x/net/context"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"user"
)

type Request struct {
	Token string `json:"token" description:"Google ID token provided by sign-in flow"`
}

func AuthUser(c context.Context, req Request) (string, error) {
	log.Debugf(c, "attempting to get google ID from token: %s", req.Token)
	googleId, err := getGoogleIDFromToken(c, req.Token)
	if err != nil {
		return "", err
	}
	log.Debugf(c, "fetched google ID: %s", googleId)

	userId, err := user.FindOrCreateByGoogleId(c, googleId)
	if err != nil {
		return "", err
	}
	log.Debugf(c, "fetched/created user with ID: %s", userId)

	return userId, nil
}

func getGoogleIDFromToken(c context.Context, token string) (string, error) {
	client := urlfetch.Client(c)
	service, err := oauth2.New(client)
	tokenInfoCall := service.Tokeninfo()
	tokenInfoCall.IdToken(token)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return "", err
	}

	log.Debugf(c, "google responded with tokeninfo: %+v", tokenInfo)
	return tokenInfo.UserId, nil
}
