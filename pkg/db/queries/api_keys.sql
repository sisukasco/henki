-- name: NewApiKey :one
INSERT INTO api_keys (
  key,user_id
) VALUES (
  $1,$2
)
RETURNING *;

-- name: GetAPIKeys :many
SELECT
   key, 
   created_at
FROM api_keys
WHERE user_id = $1 
ORDER BY created_at;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys 
WHERE key=$1 AND 
user_id=$2;

-- name: DoesAPIKeyExist :one
SELECT EXISTS
(SELECT 1 FROM api_keys WHERE key=$1 ) 
AS "exists";

-- name: GetUserFromAPIKey :one
SELECT
   users.*
FROM users,api_keys
WHERE api_keys.key = $1 AND 
api_keys.user_id = users.id;
