package usecase

import (
	"context"
)

type T[U any] struct {
	set   bool
	value U
}

func (t T[U]) TryGetValue() (U, bool) {
	if t.set {
		return t.value, true
	}
	var zero U
	return zero, false
}

type TransactionManagerRunInTransactionOptions struct {
	runningTransaction T[Transaction]
}

type Transaction interface {
	RepositoryTransaction()
}

type TransactionManager interface {
	RunInTransaction(ctx context.Context, f func(ctx context.Context, tx Transaction) error, options ...TransactionManagerRunInTransactionOption) error
}

type TransactionManagerRunInTransactionOption func(options *TransactionManagerRunInTransactionOptions)

type RepositoryOptions struct {
	transaction Transaction
	forUpdate   bool
}

type RepositoryOption func(options *RepositoryOptions)

func MergeTransactionManagerRunInTransactionOptions(
	options ...TransactionManagerRunInTransactionOption,
) TransactionManagerRunInTransactionOptions {
	var opts TransactionManagerRunInTransactionOptions
	for _, op := range options {
		op(&opts)
	}
	return opts
}

func (op TransactionManagerRunInTransactionOptions) RunningTransaction() T[Transaction] {
	return op.runningTransaction
}

func (r RepositoryOptions) Transaction() Transaction {
	return r.transaction
}

func (r RepositoryOptions) ForUpdate() bool {
	return r.forUpdate
}

func WithTransaction(tx Transaction) RepositoryOption {
	return func(options *RepositoryOptions) {
		options.transaction = tx
	}
}

func ForUpdate() RepositoryOption {
	return func(options *RepositoryOptions) {
		options.forUpdate = true
	}
}

func MergeOptions(options ...RepositoryOption) RepositoryOptions {
	var ops RepositoryOptions
	for _, op := range options {
		op(&ops)
	}
	return ops
}
