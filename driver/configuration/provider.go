package configuration

type Provider interface {
	AppEnv() string
	Port() int
	LogLevel() string
	LogFormat() string

	CSRFSecret() string

	KratosFrontendURL() string
	KratosAdminURL() string
	KratosBrowserURL() string
	JWKsURL() string
}
