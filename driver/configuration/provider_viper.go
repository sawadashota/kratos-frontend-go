package configuration

import (
	"strings"

	"github.com/sawadashota/kratos-gin-frontend/x/viperx"
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
	return viperx.GetString(viperAppEnv, "production")
}

// Port .
func (v *ViperProvider) Port() int {
	return viperx.GetInt(viperPort, 8080)
}

// LogLevel .
func (v *ViperProvider) LogLevel() string {
	return viperx.GetString(viperLogLevel, "debug")
}

// LogFormat .
func (v *ViperProvider) LogFormat() string {
	return viperx.GetString(viperLogFormat, "")
}

// CSRFSecret .
func (v *ViperProvider) CSRFSecret() string {
	return viperx.GetString(viperCSRFSecret, "cVWPNmk!9s.bvG_4Aq6Wn-fsF9jTN7jPWDxGnUhPd6!@mmQJoi")
}

// KratosFrontendURL .
func (v *ViperProvider) KratosFrontendURL() string {
	return viperx.GetString(viperKratosFrontendURL, "")
}

// KratosAdminURL .
func (v *ViperProvider) KratosAdminURL() string {
	return viperx.GetString(viperKratosAdminURL, "")
}

// KratosBrowserURL .
func (v *ViperProvider) KratosBrowserURL() string {
	return viperx.GetString(viperKratosBrowserURL, "")
}
