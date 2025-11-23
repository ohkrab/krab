package main

import (
	_ "embed"
	"os"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/fmtx"
)

var (
	//go:embed res/ferrodbicon.svg
	favicon []byte

	//go:embed res/ferrodb.svg
	logo []byte

	//go:embed tpls/embed/migration.fyml.tpl
	tplMigration []byte

	//go:embed tpls/embed/set.fyml.tpl
	tplSet []byte

	//go:embed tpls/embed/driver.fyml.tpl
	tplDriver []byte
)

func main() {
	logger := fmtx.Default()

	dir, err := config.Dir()
	if err != nil {
		logger.WriteError("can't read config dir: %w", err)
		os.Exit(1)
	}

	app := ferro.App{
		Logger:                   logger,
		Dir:                      dir,
		EmbededMigrationTemplate: tplMigration,
		EmbededDriverTemplate:    tplDriver,
		EmbededSetTemplate:       tplSet,
	}
	os.Exit(app.Run(os.Args))
}
