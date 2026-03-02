-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  @hashed_password
)
RETURNING *;

-- name: ClearUsers :exec
DELETE FROM users;

-- name: FindUser :one
SELECT * FROM users WHERE email = @email;
