package jwttoken

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

func ParseTokenClaims(claims interface{}) (*Claims, error) {
	var c Claims
	claimDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeHookFunc("2006-01-02T15:04:05.999999999Z"),
		Result:     &c,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to initialize token claim decoder. %s", err)
	}

	if err := claimDecoder.Decode(claims); err != nil {
		return nil, fmt.Errorf("fail to decode token claim. error: %s value: %v", err, claims)
	}
	return &c, nil
}

type Claims struct {
	Exp     uint    `mapstructure:"exp"`
	Iat     uint    `mapstructure:"iat"`
	Nbf     uint    `mapstructure:"nbf"`
	Iss     string  `mapstructure:"iss"`
	Jti     string  `mapstructure:"jti"`
	Sub     string  `mapstructure:"sub"`
	Session Session `mapstructure:"session"`
}

type Session struct {
	Sid             string    `mapstructure:"sid"`
	AuthenticatedAt time.Time `mapstructure:"authenticated_at"`
	ExpiresAt       time.Time `mapstructure:"expires_at"`
	IssuedAt        time.Time `mapstructure:"issued_at"`
	Identity        Identity  `mapstructure:"identity"`
}

type Identity struct {
	ID              string    `mapstructure:"id"`
	Addresses       []Address `mapstructure:"addresses"`
	TraitsSchemaID  string    `mapstructure:"traits_schema_id"`
	TraitsSchemaURL string    `mapstructure:"traits_schema_url"`
}

type Traits struct {
	Email string `mapstructure:"email"`
}

type Address struct {
	ID        string    `mapstructure:"id"`
	ExpiresAt time.Time `mapstructure:"expires_at"`
	Value     string    `mapstructure:"value"`
	Verified  bool      `mapstructure:"verified"`
	Via       string    `mapstructure:"via"`
}
