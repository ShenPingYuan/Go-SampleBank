package db

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error, isolationLevel sql.IsolationLevel) error {
	var tx, err = store.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolationLevel,
	})

	if err != nil {
		return err
	}
	query := New(tx)

	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer tx
type TransferTxParams struct {
	FromAccountID int64           `json:"from_account_id"`
	ToAccountID   int64           `json:"to_account_id"`
	Amount        decimal.Decimal `json:"amount"`
	//备注
	Mark string `json:"mark"`
}

// TransferTxResult is the result of the transfer tx
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction.
// 转账
func (store *Store) TransferTx(ctx context.Context, param TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(query *Queries) error {
		createTransferResult, err := query.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: param.FromAccountID,
			ToAccountID:   param.ToAccountID,
			Amount:        param.Amount,
		})
		if err != nil {
			return err
		}
		lastInsertId, err := createTransferResult.LastInsertId()
		if err != nil {
			return err
		}
		transfer, err := query.GetTransfer(ctx, lastInsertId)
		if err != nil {
			return err
		}
		result.Transfer = transfer

		createFromEntryResult, err := query.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.FromAccountID,
			Amount:    param.Amount.Neg(),
		})
		if err != nil {
			return err
		}
		lastInsertId, err = createFromEntryResult.LastInsertId()
		if err != nil {
			return err
		}
		fromEntry, err := query.GetEntry(ctx, lastInsertId)
		if err != nil {
			return err
		}
		result.FromEntry = fromEntry

		createToEntryResult, err := query.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.ToAccountID,
			Amount:    param.Amount,
		})
		if err != nil {
			return err
		}
		lastInsertId, err = createToEntryResult.LastInsertId()
		if err != nil {
			return err
		}
		toEntry, err := query.GetEntry(ctx, lastInsertId)
		if err != nil {
			return err
		}
		result.ToEntry = toEntry

		//更新账户余额
		err = query.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     param.FromAccountID,
			Amount: param.Amount.Neg(),
		})
		if err != nil {
			return err
		}
		err = query.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     param.ToAccountID,
			Amount: param.Amount,
		})
		if err != nil {
			return err
		}
		//查询更新后的账户
		result.FromAccount, err = query.GetAccount(ctx, param.FromAccountID)
		if err != nil {
			return err
		}
		//查询更新后的账户
		result.ToAccount, err = query.GetAccount(ctx, param.ToAccountID)
		if err != nil {
			return err
		}
		return nil
	}, sql.LevelRepeatableRead)
	return result, err
}
