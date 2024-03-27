package endpoint

import (
	"github.com/pheezz/medods-jwt-test/internal/app/service"
)

type Service interface {
	GetKeyPair(GUID string, remoteAddr string) (service.TokenPairCookies, error)
	RefreshKeyPair(RefreshTokenHash string, remoteAddr string, fingerprint string) (service.TokenPairCookies, error)
}

type Endpoint struct {
	service Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		service: s,
	}
}
