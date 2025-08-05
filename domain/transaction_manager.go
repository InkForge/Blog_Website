package domain

import "context"

type ITransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
