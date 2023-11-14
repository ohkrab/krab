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

func Ext() string {
	return ".krab.hcl"
}

func Auth() string {
	switch os.Getenv("KRAB_AUTH") {
	case "basic":
		return "basic"

	default:
		return "none"
	}
}

func HttpBasicAuthData() map[string]string {
	users := map[string]string{}
	name := os.Getenv("KRAB_AUTH_BASIC_USERNAME")
	pass := os.Getenv("KRAB_AUTH_BASIC_PASSWORD")
	if name == "" || pass == "" {
		panic("KRAB_AUTH_BASIC_USER or KRAB_AUTH_BASIC_PASSWORD is not set")
	}
	users[name] = pass
	return users
}
