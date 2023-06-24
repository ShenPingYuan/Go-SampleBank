-- name: CreateEntry :execresult
INSERT INTO entries
(
    account_id,
    amount
)
VALUES
    (?, ?);

-- name: GetEntry :one
SELECT
    *
FROM
    entries
WHERE
    id = ?;

-- name: ListEntries :many
SELECT
    *
FROM
    entries
ORDER BY
    id ASC
LIMIT ? OFFSET ?;

-- name: DeleteEntry :execresult
DELETE FROM entries
WHERE id = ?;
