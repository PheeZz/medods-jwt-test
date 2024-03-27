package service

import (
	"errors"
	"github.com/pheezz/medods-jwt-test/internal/app/database"
	"github.com/pheezz/medods-jwt-test/internal/app/token"
	"strings"
	"time"
)

var errRefreshTokenHashRequired = errors.New("refresh_token cookie not found")
var errRefreshTokenExpired = errors.New("refresh_token expired")
var errRefreshTokenHashInvalid = errors.New("refresh_token hash invalid")

func (s *Service) RefreshKeyPair(RefreshTokenHash string, remoteAddr string, fingerprint string) (TokenPairCookies, error) {
	if RefreshTokenHash == "" {
		return TokenPairCookies{}, errRefreshTokenHashRequired
	}
	ip := strings.Split(remoteAddr, ":")[0]
	user, key, err := database.GetUserByFingerprintAndIP(fingerprint, ip)
	if err != nil {
		return TokenPairCookies{}, err
	}

	if !token.ValidateRefreshToken(RefreshTokenHash, key.RefreshTokenHash) {
		return TokenPairCookies{}, errRefreshTokenHashInvalid
	}

	if key.ExpireAtTimestamp < time.Now().Unix() {
		database.DeleteKey(key)
		return TokenPairCookies{}, errRefreshTokenExpired
	}

	tokenPair := token.New(user.GUID)
	database.AddKeyToUser(user.GUID, database.KeySchema{
		RefreshTokenHash:  tokenPair.RefreshToken.Hash,
		ExpireAtTimestamp: tokenPair.RefreshToken.ExpireAtTimestamp,
		FromIP:            ip,
		Fingerprint:       tokenPair.RefreshToken.Fingerprint,
	})
	database.DeleteKey(key)
	return convertTokensToCookies(*tokenPair), nil
}
