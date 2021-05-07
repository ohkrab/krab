package krabdb

import (
	"context"
	"fmt"
)

type AdvisoryLock struct {
	Errs     chan<- error
	acquired bool
}

func (db *AdvisoryLock) Lock(ctx context.Context) {
	var err error
	fmt.Println("ACQUIRE LOCK")

	if err != nil {
		db.Errs <- err
		return
	}

	db.acquired = true
}

func (db *AdvisoryLock) Unlock(ctx context.Context) {
	var err error

	if db.acquired {
		fmt.Println("RELEASE LOCK")
		if err != nil {
			db.Errs <- err
		}
	}
}
