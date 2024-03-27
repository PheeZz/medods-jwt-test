package service

import (
	"errors"
	"github.com/pheezz/medods-jwt-test/internal/app/database"
	"github.com/pheezz/medods-jwt-test/internal/app/token"
	"net/http"
	"strings"
)

var errUserNotFound = errors.New("user not found")
var errGuidRequired = errors.New("GUID is required in query params")

func (s *Service) GetKeyPair(GUID string, remoteAddr string) (TokenPairCookies, error) {
	if GUID == "" {
		return TokenPairCookies{}, errGuidRequired
	}
	ip := strings.Split(remoteAddr, ":")[0]
	user, err := database.GetUserByGUID(GUID)
	if err != nil {
		return TokenPairCookies{}, errUserNotFound
	}
	tokens := token.New(user.GUID)
	database.AddKeyToUser(user.GUID, database.KeySchema{
		RefreshTokenHash:  tokens.RefreshToken.Hash,
		ExpireAtTimestamp: tokens.RefreshToken.ExpireAtTimestamp,
		FromIP:            ip,
		Fingerprint:       tokens.RefreshToken.Fingerprint,
	})

	return convertTokensToCookies(*tokens), nil

}

func convertTokensToCookies(tokens token.TokenPair) TokenPairCookies {
	return TokenPairCookies{
		AccessToken: http.Cookie{
			Name:     "access_token",
			Value:    tokens.AccessToken,
			MaxAge:   int(conf.AccessTokenDuration.Seconds()),
			Secure:   true,
			HttpOnly: true,
		},
		RefreshToken: http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken.Token,
			MaxAge:   int(conf.RefreshTokenDuration.Seconds()),
			Secure:   true,
			HttpOnly: true,
		},
		Fingerprint: http.Cookie{
			Name:     "fingerprint",
			Value:    tokens.RefreshToken.Fingerprint,
			MaxAge:   int(conf.RefreshTokenDuration.Seconds()),
			Secure:   true,
			HttpOnly: true,
		},
	}
}
