-- name: CreateSession :execresult
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expire_time,
  created_at
) VALUES (
  ?,?,?,?,?,?,?,?
);

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE refresh_token = ? LIMIT 1;