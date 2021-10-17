package krab

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type HookRunner struct {
	Hooks *Hooks
}

func (h HookRunner) RunBefore(ctx context.Context, db sqlx.ExecerContext) error {
	if h.Hooks.Before != "" {
		_, err := db.ExecContext(ctx, h.Hooks.Before)
		return err
	}

	return nil
}
