package configuration

import (
	"strings"

	"github.com/sawadashota/kratos-frontend-go/x"
	"github.com/spf13/viper"
)

type ViperProvider struct{}

const (
	viperAppEnv            = "app.env"
	viperPort              = "port"
	viperLogLevel          = "log.level"
	viperLogFormat         = "log.format"
	viperCSRFSecret        = "secret.csrf"
	viperKratosFrontendURL = "kratos.frontend_url"
	viperKratosAdminURL    = "kratos.admin_url"
	viperKratosBrowserURL  = "kratos.browser_url"
	viperJWKsURL           = "jwks_url"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func NewViperProvider() Provider {
	return new(ViperProvider)
}

// AppEnv .
func (v *ViperProvider) AppEnv() string {
	return x.ViperGetString(viperAppEnv, "production")
}

// Port .
func (v *ViperProvider) Port() int {
	return x.ViperGetInt(viperPort, 8080)
}

// LogLevel .
func (v *ViperProvider) LogLevel() string {
	return x.ViperGetString(viperLogLevel, "debug")
}

// LogFormat .
func (v *ViperProvider) LogFormat() string {
	return x.ViperGetString(viperLogFormat, "")
}

// CSRFSecret .
func (v *ViperProvider) CSRFSecret() string {
	return x.ViperGetString(viperCSRFSecret, "cVWPNmk!9s.bvG_4Aq6Wn-fsF9jTN7jPWDxGnUhPd6!@mmQJoi")
}

// KratosFrontendURL .
func (v *ViperProvider) KratosFrontendURL() string {
	return x.ViperGetString(viperKratosFrontendURL, "")
}

// KratosAdminURL .
func (v *ViperProvider) KratosAdminURL() string {
	return x.ViperGetString(viperKratosAdminURL, "")
}

// KratosBrowserURL .
func (v *ViperProvider) KratosBrowserURL() string {
	return x.ViperGetString(viperKratosBrowserURL, "")
}

// KratosLogoutURL .
func (v *ViperProvider) KratosLogoutURL() string {
	return v.KratosBrowserURL() + "/self-service/browser/flows/logout"
}

// KratosLoginURL .
func (v *ViperProvider) KratosLoginURL() string {
	return v.KratosBrowserURL() + "/self-service/browser/flows/login"
}

// KratosRegistrationURL .
func (v *ViperProvider) KratosRegistrationURL() string {
	return v.KratosBrowserURL() + "/self-service/browser/flows/registration"
}

// JWKsURL .
func (v *ViperProvider) JWKsURL() string {
	return x.ViperGetString(viperJWKsURL, "")
}
