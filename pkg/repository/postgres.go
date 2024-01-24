package repository

import (
	"context"
	"database/sql"

	"github.com/taroxii/vote-api/pkg/usecase"
)

type pgTransactionWrapper struct {
	tx *sql.DB
}

func newPgTransactionWrapper(tx *sql.DB) pgTransactionWrapper {

	return pgTransactionWrapper{
		tx: tx,
	}
}

func (w pgTransactionWrapper) RepositoryTransaction() {}

type PgTransactionManager struct {
	db *sql.DB
}

func NewPgTransactionManager(db *sql.DB) PgTransactionManager {
	return PgTransactionManager{
		db: db,
	}
}

func (m PgTransactionManager) RunInTx(ctx context.Context, f func(ctx context.Context, tx usecase.Transaction) error, options ...usecase.TransactionManagerRunInTransactionOption) error {
	opts := usecase.MergeTransactionManagerRunInTransactionOptions(options...)
	if tx, ok := opts.RunningTransaction().TryGetValue(); ok {
		return f(ctx, tx)
	} else {
		return f(ctx, newPgTransactionWrapper(m.db))
	}
}

func ResolveDB(db *sql.DB, tx usecase.Transaction) (*sql.DB, bool) {
	if tx == nil {
		return db, false
	}
	return tx.(pgTransactionWrapper).tx, true
}

// ResolveBaseQuery constructs base query from `use_case.RepositoryOptions`.
// Please note that this function *must* be used in query operation only.
func ResolveBaseQuery(db *sql.DB, opts usecase.RepositoryOptions) (tx *sql.Tx) {
	if opts.ForUpdate() {
		tx, err := db.Begin()
		if err != nil {
			return nil
		}
		return tx

	}
	return
}
