-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner_id,
  balance,
  currency
) VALUES (
  ?,?,?
);
-- name: GetLastInsertId :one
SELECT LAST_INSERT_ID() as account_id;
-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;
-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: ListAccounts :many
-- 分页查询
SELECT * FROM accounts
ORDER BY id ASC
LIMIT ? OFFSET ?;

-- name: DeleteAccount :execresult
DELETE FROM accounts
WHERE id = ?;

-- 更新余额
-- name: UpdateAccountBalance :execresult
UPDATE accounts
SET balance = ?
WHERE id = ?; 

-- 添加账户余额
-- name: AddAccountBalance :exec
UPDATE accounts
SET balance = balance + sqlc.arg(Amount)
WHERE id = sqlc.arg(Id);