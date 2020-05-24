package driver

import (
	"github.com/gorilla/mux"
	"github.com/ory/kratos-client-go/client"
	"github.com/sawadashota/kratos-frontend-go/account"
	"github.com/sawadashota/kratos-frontend-go/authentication"
	"github.com/sawadashota/kratos-frontend-go/driver/configuration"
	"github.com/sawadashota/kratos-frontend-go/internal/jwt"
	"github.com/sawadashota/kratos-frontend-go/middleware"
	"github.com/sawadashota/kratos-frontend-go/salary"
	"github.com/sirupsen/logrus"
)

// Registry .
type Registry interface {
	Logger() logrus.FieldLogger

	JWTParser() *jwt.Parser
	Middleware() *middleware.Middleware
	KratosClient() *client.OryKratos

	AccountHandler() *account.Handler
	AuthenticationHandler() *authentication.Handler
	SalaryHandler() *salary.Handler
	RegisterRoutes(router *mux.Router)
}

// NewRegistry .
func NewRegistry(c configuration.Provider) Registry {
	return NewRegistryDefault(c)
}
