package auth

import (
	"golang.org/x/net/context"
	"server/user"
)

type Request struct {
	Token string `json:"token"`
}

func AuthUser(c context.Context, req Request) (string, error) {
	// validate token with google
	googleId := "1"
	userId, err := user.FindOrCreateByGoogleId(c, googleId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
