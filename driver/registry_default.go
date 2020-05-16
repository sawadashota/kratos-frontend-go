package driver

import (
	"net/url"

	"github.com/gorilla/mux"

	"github.com/sawadashota/kratos-gin-frontend/account"
	"github.com/sawadashota/kratos-gin-frontend/authentication"

	"github.com/sawadashota/kratos-gin-frontend/middleware"

	"github.com/ory/kratos-client-go/client"
	"github.com/sawadashota/kratos-gin-frontend/driver/configuration"
	"github.com/sirupsen/logrus"
)

var _ Registry = new(RegistryDefault)

// RegistryDefault .
type RegistryDefault struct {
	l logrus.FieldLogger
	c configuration.Provider
	r Registry

	mw *middleware.Middleware

	kratosClient *client.OryKratos

	authenticationHandler *authentication.Handler
	accountHandler        *account.Handler
}

func (r *RegistryDefault) AccountHandler() *account.Handler {
	if r.accountHandler == nil {
		r.accountHandler = account.New(r, r.c)
	}
	return r.accountHandler
}

func (r *RegistryDefault) AuthenticationHandler() *authentication.Handler {
	if r.authenticationHandler == nil {
		r.authenticationHandler = authentication.New(r, r.c)
	}
	return r.authenticationHandler
}

func (r *RegistryDefault) RegisterRoutes(router *mux.Router) {
	r.AuthenticationHandler().RegisterRoutes(router)
	r.AccountHandler().RegisterRoutes(router)
}

func NewRegistryDefault(c configuration.Provider) *RegistryDefault {
	return &RegistryDefault{
		c: c,
	}
}

// Logger .
func (r *RegistryDefault) Logger() logrus.FieldLogger {
	if r.l == nil {
		r.l = r.newLogger()
	}
	return r.l
}

func (r *RegistryDefault) newLogger() logrus.FieldLogger {
	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(r.c.LogLevel())
	if err == nil {
		l.SetLevel(level)
	}
	return l
}

func (r *RegistryDefault) Middleware() *middleware.Middleware {
	if r.mw == nil {
		r.mw = middleware.New(r, r.c)
	}
	return r.mw
}

func (r *RegistryDefault) KratosClient() *client.OryKratos {
	if r.kratosClient == nil {
		u, _ := url.Parse(r.c.KratosAdminURL())
		cl := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
			Host:     u.Host,
			BasePath: "/",
			Schemes:  []string{u.Scheme},
		})
		r.kratosClient = cl
	}

	return r.kratosClient
}
