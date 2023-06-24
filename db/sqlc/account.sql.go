// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: account.sql

package db

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

const addAccountBalance = `-- name: AddAccountBalance :exec
UPDATE accounts
SET balance = balance + ?
WHERE id = ?
`

type AddAccountBalanceParams struct {
	Amount decimal.Decimal `json:"amount"`
	ID     int64           `json:"id"`
}

// 添加账户余额
func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) error {
	_, err := q.db.ExecContext(ctx, addAccountBalance, arg.Amount, arg.ID)
	return err
}

const createAccount = `-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  ?,?,?
)
`

type CreateAccountParams struct {
	Owner    string          `json:"owner"`
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
}

const deleteAccount = `-- name: DeleteAccount :execresult
DELETE FROM accounts
WHERE id = ?
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteAccount, id)
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = ? LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getLastInsertId = `-- name: GetLastInsertId :one
SELECT LAST_INSERT_ID() as account_id
`

func (q *Queries) GetLastInsertId(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertId)
	var account_id int64
	err := row.Scan(&account_id)
	return account_id, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner, balance, currency, created_at FROM accounts
ORDER BY id ASC
LIMIT ? OFFSET ?
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// 分页查询
func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccountBalance = `-- name: UpdateAccountBalance :execresult
UPDATE accounts
SET balance = ?
WHERE id = ?
`

type UpdateAccountBalanceParams struct {
	Balance decimal.Decimal `json:"balance"`
	ID      int64           `json:"id"`
}

// 更新余额
func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateAccountBalance, arg.Balance, arg.ID)
}
