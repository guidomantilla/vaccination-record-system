package datasource

import (
	"context"

	"gorm.io/gorm"
)

var _ TransactionHandler = (*DefaultTransactionHandler)(nil)

type TransactionHandlerFunction func(ctx context.Context, tx *gorm.DB) error

type TransactionHandler interface {
	HandleTransaction(ctx context.Context, fn TransactionHandlerFunction) error
}
