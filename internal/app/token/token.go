package token

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pheezz/medods-jwt-test/internal/app/config"
	"github.com/pheezz/medods-jwt-test/internal/app/randomizer"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var conf = config.Conf

type RefreshToken struct {
	Token             string
	Hash              []byte
	ExpireAtTimestamp int64
	Fingerprint       string
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     RefreshToken
	RefreshTokenHash []byte
}

func New(GUID string) *TokenPair {
	at := createAccessJWT(GUID)
	rt := createRefreshToken()
	fingerprint := randomizer.String(50)
	rtHash, err := bcrypt.GenerateFromPassword([]byte(rt), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error generating hash: ", err)
	}

	return &TokenPair{
		AccessToken: at,
		RefreshToken: RefreshToken{
			Token:             rt,
			Hash:              rtHash,
			ExpireAtTimestamp: time.Now().Add(conf.RefreshTokenDuration).Unix(),
			Fingerprint:       fingerprint,
		},
	}
}

func createAccessJWT(GUID string) string {
	expiration := time.Now().Add(conf.AccessTokenDuration)
	iat := time.Now().Unix()
	payload := jwt.MapClaims{
		"sub":    GUID,
		"iat":    iat,
		"maxAge": expiration.Unix()}
	token := jwt.NewWithClaims(conf.JwtAlgorithm, payload)
	t, err := token.SignedString(conf.JwtSecretKey)

	if err != nil {
		log.Fatal("Error signing token: ", err)
	}
	return t

}

func createRefreshToken() string {
	const length = 50 // in case of that bcrypt hash can't be longer than 72 bytes
	token := randomizer.String(length)
	b64t := base64.StdEncoding.EncodeToString([]byte(token))
	return b64t
}

func ValidateRefreshToken(token string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(token))
	return err == nil
}
