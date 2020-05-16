package driver

import (
	"github.com/gorilla/mux"
	"github.com/ory/kratos-client-go/client"
	"github.com/sawadashota/kratos-gin-frontend/account"
	"github.com/sawadashota/kratos-gin-frontend/authentication"
	"github.com/sawadashota/kratos-gin-frontend/driver/configuration"
	"github.com/sawadashota/kratos-gin-frontend/middleware"
	"github.com/sirupsen/logrus"
)

// Registry .
type Registry interface {
	Logger() logrus.FieldLogger

	Middleware() *middleware.Middleware
	KratosClient() *client.OryKratos

	AccountHandler() *account.Handler
	AuthenticationHandler() *authentication.Handler
	RegisterRoutes(router *mux.Router)
}

// NewRegistry .
func NewRegistry(c configuration.Provider) Registry {
	return NewRegistryDefault(c)
}
