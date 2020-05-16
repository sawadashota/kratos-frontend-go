package x

import "github.com/spf13/viper"

func ViperGetString(key, fallback string) string {
	v := viper.GetString(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}

func ViperGetInt(key string, fallback int) int {
	v := viper.GetInt(key)

	if v == 0 {
		return fallback
	}

	return v
}
