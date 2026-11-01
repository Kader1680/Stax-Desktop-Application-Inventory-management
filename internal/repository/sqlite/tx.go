package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

// TxManager implements repository.TxManager over a *sql.DB.
type TxManager struct {
	DB *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{DB: db}
}

func (t *TxManager) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("démarrage transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	err = fn(tx)
	return err
}
