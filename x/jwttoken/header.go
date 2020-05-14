package jwttoken

import "github.com/mitchellh/mapstructure"

type TokenHeader struct {
	Alg string `mapstructure:"alg"`
	Kid string `mapstructure:"kid"`
	Typ string `mapstructure:"typ"`
}

func ParseTokenHeader(header interface{}) (*TokenHeader, error) {
	var th TokenHeader
	if err := mapstructure.Decode(header, &th); err != nil {
		return nil, err
	}

	return &th, nil
}
