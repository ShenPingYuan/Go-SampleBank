-- name: CreateTransfer :execresult
INSERT INTO transfers
    (
    from_account_id,
    to_account_id,
    amount
    )
VALUES
    (?, ?, ?);

-- name: GetTransfer :one
SELECT
    *
FROM
    transfers
WHERE
    id = ?;

-- name: ListTransfers :many
SELECT
    *
FROM
    transfers
ORDER BY
    id ASC
LIMIT ? OFFSET ?;

-- name: DeleteTransfer :execresult
DELETE FROM transfers
WHERE id = ?;
