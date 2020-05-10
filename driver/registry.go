package driver

import (
	"github.com/ory/kratos-client-go/client"
	"github.com/sawadashota/kratos-gin-frontend/driver/configuration"
	"github.com/sirupsen/logrus"
)

// Registry .
type Registry interface {
	Logger() logrus.FieldLogger
	WithConfig(c configuration.Provider) Registry

	KratosClient() *client.OryKratos
}

// NewRegistry .
func NewRegistry(c configuration.Provider) Registry {
	r := &RegistryBase{}
	r.WithConfig(c)
	return r
}
