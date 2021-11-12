
-- name: NewRefreshToken :exec
INSERT INTO refresh_tokens (
  user_id, token
) VALUES (
  $1, $2
);

-- name: ClearRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

