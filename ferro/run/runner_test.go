package run

import (
	"context"
	"testing"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/qbart/expecto/expecto"
)

func TestRunner_MigrationAuditLog(t *testing.T) {
	db := createTestDB(t, context.Background())
	defer db.clear()
	_, dir, fsCleanup := expecto.TempFS(
		db.fymlFileName,
        db.fymlFileContent,
	)
	defer fsCleanup()

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	expecto.NoErr(t, "parsing config", err)
    expecto.NotNil(t, "parsed", parsed)

	// builder := NewBuilder(fs, parsed, plugins.New())
	// cfg, errs := builder.BuildConfig()
	// expecto.NotNil(t, "build errors", errs)
	// expecto.Eq(t, "number of errors", len(errs.Errors), 0)
}
