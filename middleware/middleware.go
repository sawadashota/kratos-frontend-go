package middleware

import "github.com/sirupsen/logrus"

type Middleware struct {
	r Registry
	c Configuration
}

type Registry interface {
	Logger() logrus.FieldLogger
}

type Configuration interface {
	JWKsURL() string
	KratosLoginURL() string
	KratosRegistrationURL() string
}

func New(r Registry, c Configuration) *Middleware {
	return &Middleware{
		r: r,
		c: c,
	}
}
