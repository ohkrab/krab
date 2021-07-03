package krabenv

import "os"

func GetConfigDir() (string, error) {
	if dir := os.Getenv("KRAB_DIR"); dir != "" {
		return dir, nil
	}

	return os.Getwd()
}
