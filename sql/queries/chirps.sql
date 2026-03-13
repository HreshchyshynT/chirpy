-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps 
WHERE sqlc.narg(author_id)::UUID IS NULL OR @author_id::UUID = user_id
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE chirps.id = @id;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE chirps.id = @id;
