-- name: CreateUser :execresult
INSERT INTO users (
    username,
    hashed_password,
    email,
    full_name
) VALUES (
  ?,?,?,?
);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;