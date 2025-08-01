package domain

import "context"

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
