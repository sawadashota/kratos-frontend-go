package driver

import (
	"net/url"

	"github.com/ory/kratos-client-go/client"
	"github.com/sawadashota/kratos-gin-frontend/driver/configuration"
	"github.com/sirupsen/logrus"
)

// RegistryBase .
type RegistryBase struct {
	l logrus.FieldLogger
	c configuration.Provider
	r Registry

	kratosClient *client.OryKratos
}

// WithConfig .
func (r *RegistryBase) WithConfig(c configuration.Provider) Registry {
	r.c = c
	return r.r
}

// Logger .
func (r *RegistryBase) Logger() logrus.FieldLogger {
	if r.l == nil {
		r.l = r.newLogger()
	}
	return r.l
}

func (r *RegistryBase) newLogger() logrus.FieldLogger {
	l := logrus.New()

	l.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(r.c.LogLevel())
	if err == nil {
		l.SetLevel(level)
	}
	return l
}

func (r *RegistryBase) KratosClient() *client.OryKratos {
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
