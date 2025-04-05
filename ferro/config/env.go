package config

import "os"

func Dir() (string, error) {
	if dir := os.Getenv("FERRO_DIR"); dir != "" {
		return dir, nil
	}

	return os.Getwd()
}
