package middleware

import (
	"github.com/sawadashota/kratos-frontend-go/internal/jwt"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	r Registry
	c Configuration
}

type Registry interface {
	Logger() logrus.FieldLogger
	JWTParser() *jwt.Parser
}

type Configuration interface {
	JWKsURL() string
	KratosLoginURL() string
	KratosRegistrationURL() string
	KratosSettingsURL() string
}

func New(r Registry, c Configuration) *Middleware {
	return &Middleware{
		r: r,
		c: c,
	}
}
