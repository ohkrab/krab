package main

import (
	"github.com/ohkrab/krab/fs"
	"github.com/ohkrab/krab/krab"
)

func main() {
	app := krab.New(fs.MustGetPwd())
	app.Run()
}
