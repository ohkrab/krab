package mustdir

import "os"

// GetPwd gets PWD or panics.
func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
