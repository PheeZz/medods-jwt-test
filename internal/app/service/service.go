package service

import (
	"github.com/pheezz/medods-jwt-test/internal/app/config"
	"net/http"
)

var conf = config.Conf

type TokenPairCookies struct {
	AccessToken  http.Cookie
	RefreshToken http.Cookie
	Fingerprint  http.Cookie
}

type Service struct {
}

func New() *Service {
	return &Service{}
}
