package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/square/go-jose/v3"
)

type Parser struct {
	r Registry
	c Configuration
}

type Registry interface {
	Logger() logrus.FieldLogger
}

type Configuration interface {
	JWKsURL() string
}

func New(r Registry, c Configuration) *Parser {
	return &Parser{
		r: r,
		c: c,
	}
}

func (p *Parser) ParseTokenClaims(claims interface{}) (*Claims, error) {
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
	Exp     int64   `json:"exp" mapstructure:"exp"`
	Iat     int64   `json:"iat" mapstructure:"iat"`
	Nbf     int64   `json:"nbf" mapstructure:"nbf"`
	Iss     string  `json:"iss" mapstructure:"iss"`
	Jti     string  `json:"jti" mapstructure:"jti"`
	Sub     string  `json:"sub" mapstructure:"sub"`
	Session Session `json:"session" mapstructure:"session"`
}

type Session struct {
	Sid             string    `json:"sid" mapstructure:"sid"`
	AuthenticatedAt time.Time `json:"authenticated_at" mapstructure:"authenticated_at"`
	ExpiresAt       time.Time `json:"expires_at" mapstructure:"expires_at"`
	IssuedAt        time.Time `json:"issued_at" mapstructure:"issued_at"`
	Identity        Identity  `json:"identity" mapstructure:"identity"`
}

type Identity struct {
	ID              string    `json:"id" mapstructure:"id"`
	Addresses       []Address `json:"addresses" mapstructure:"addresses"`
	TraitsSchemaID  string    `json:"traits_schema_id" mapstructure:"traits_schema_id"`
	TraitsSchemaURL string    `json:"traits_schema_url" mapstructure:"traits_schema_url"`
}

type Traits struct {
	Email string `json:"email" mapstructure:"email"`
}

type Address struct {
	ID        string    `json:"id" mapstructure:"id"`
	ExpiresAt time.Time `json:"expires_at" mapstructure:"expires_at"`
	Value     string    `json:"value" mapstructure:"value"`
	Verified  bool      `json:"verified" mapstructure:"verified"`
	Via       string    `json:"via" mapstructure:"via"`
}

func (c *Claims) IsExpired() bool {
	now := time.Now()
	return now.Unix() >= c.Exp || now.After(c.Session.ExpiresAt)
}

func (p *Parser) ParseRequest(req *http.Request) (*jwt.Token, error) {
	return request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, p.keyFunc())
}

func (p *Parser) keyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodRSA); !err {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		th, err := p.ParseTokenHeader(token.Header)
		if err != nil {
			return nil, err
		}

		jwk, err := p.findJWK(th.Kid)
		if err != nil {
			return nil, err
		}

		return jwk.Key, nil
	}
}

func (p *Parser) findJWK(kid string) (*jose.JSONWebKey, error) {
	p.r.Logger().Infof("fetch JWKs from %s", p.c.JWKsURL())
	res, err := http.Get(p.c.JWKsURL())
	if err != nil {
		return nil, err
	}

	var jwks jose.JSONWebKeySet
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	jwk := jwks.Key(kid)
	if len(jwk) == 0 {
		return nil, fmt.Errorf("JWLs does't have the kid: %s", kid)
	}

	return &jwk[0], nil
}

type TokenHeader struct {
	Alg string `json:"alg" mapstructure:"alg"`
	Kid string `json:"kid" mapstructure:"kid"`
	Typ string `json:"typ" mapstructure:"typ"`
}

func (p *Parser) ParseTokenHeader(header interface{}) (*TokenHeader, error) {
	var th TokenHeader
	if err := mapstructure.Decode(header, &th); err != nil {
		return nil, err
	}

	return &th, nil
}
