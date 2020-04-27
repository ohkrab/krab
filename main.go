package main

import (
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/mustdir"
)

func main() {
	app := krab.New(mustdir.GetPwd())
	app.Run()
}
