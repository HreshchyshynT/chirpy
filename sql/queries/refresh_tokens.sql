-- name: SaveRefreshToken :one
INSERT INTO refresh_tokens (
  token,
  created_at,
  updated_at,
  user_id,
  expires_at
) VALUES (
  @token,
  NOW(),
  NOW(),
  @user_id,
  @expires_at
) RETURNING *;

-- name: FindRefreshToken :one
SELECT * from refresh_tokens 
WHERE token = @token;
