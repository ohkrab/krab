package fs

import "os"

// MustGetPwd gets PWD or panics.
func MustGetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
