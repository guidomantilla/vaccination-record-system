package datasource

import (
	"context"

	"gorm.io/gorm"
)

type DefaultTransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *DefaultTransactionHandler {
	return &DefaultTransactionHandler{db: db}
}

func (handler *DefaultTransactionHandler) HandleTransaction(ctx context.Context, fn TransactionHandlerFunction) error {
	return handler.db.Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}
