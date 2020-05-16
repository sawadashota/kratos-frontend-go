package driver

import "github.com/sawadashota/kratos-frontend-go/driver/configuration"

type Driver interface {
	Configuration() configuration.Provider
	Registry() Registry
}

type DefaultDriver struct {
	c configuration.Provider
	r Registry
}

// NewDefaultDriver .
func NewDefaultDriver() Driver {
	c := configuration.NewViperProvider()

	return &DefaultDriver{
		c: c,
		r: NewRegistry(c),
	}
}

// Configuration .
func (d *DefaultDriver) Configuration() configuration.Provider {
	return d.c
}

// Registry .
func (d *DefaultDriver) Registry() Registry {
	return d.r
}
