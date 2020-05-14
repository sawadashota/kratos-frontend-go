package jwttoken

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/square/go-jose/v3"
)

func ParseRequest(req *http.Request, jwksURL string) (*jwt.Token, error) {
	return request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, KeyFunc(jwksURL))
}

func KeyFunc(jwksURL string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodRSA); !err {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		th, err := ParseTokenHeader(token.Header)
		if err != nil {
			return nil, err
		}

		jwk, err := findJWK(jwksURL, th.Kid)
		if err != nil {
			return nil, err
		}

		return jwk.Key, nil
	}
}

func findJWK(jwksURL, kid string) (*jose.JSONWebKey, error) {
	res, err := http.Get(jwksURL)
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
