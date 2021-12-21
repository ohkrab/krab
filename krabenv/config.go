package krabenv

import "os"

func ConfigDir() (string, error) {
	if dir := os.Getenv("KRAB_DIR"); dir != "" {
		return dir, nil
	}

	return os.Getwd()
}

func DatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}
