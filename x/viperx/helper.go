package viperx

import "github.com/spf13/viper"

func GetString(key, fallback string) string {
	v := viper.GetString(key)

	if len(v) == 0 {
		return fallback
	}

	return v
}

func GetInt(key string, fallback int) int {
	v := viper.GetInt(key)

	if v == 0 {
		return fallback
	}

	return v
}
