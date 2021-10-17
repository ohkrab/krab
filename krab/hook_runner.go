package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
)

type HookRunner struct {
	Hooks *Hooks
}

func (h HookRunner) SetSearchPath(ctx context.Context, db sqlx.ExecerContext, schema string) error {
	_, err := db.ExecContext(ctx, fmt.Sprint("SET search_path TO ", krabdb.QuoteIdent(schema)))
	return err
}
