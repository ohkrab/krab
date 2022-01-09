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

func Env() string {
	return os.Getenv("KRAB_ENV")
}

func Test() bool {
	return Env() == "test"
}
