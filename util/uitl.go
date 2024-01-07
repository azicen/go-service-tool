package util

import "os"

func Getenv(key, defVal string) string {
	if envValue, exists := os.LookupEnv(key); exists {
		return envValue
	}
	return defVal
}
