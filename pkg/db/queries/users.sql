-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: createNewUser :one
INSERT INTO users (
  id,email, encrypted_password
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DoesUserIDExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE id=$1) 
AS "exists";

-- name: DoesUserExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE email=$1) 
AS "exists";

-- name: UpdateUserProfile :exec
UPDATE users 
SET 
  first_name = $2,
  last_name = $3,
  avatar_url = $4
WHERE id = $1;

-- name: ConfirmUserEmailByID :exec
UPDATE users 
SET 
  confirmed_at = NOW()
WHERE id = $1 AND confirmed_at is NULL;

-- name: ConfirmUserEmail :exec
UPDATE users 
SET 
  confirmed_at = NOW()
WHERE confirmation_token = $1 AND confirmed_at is NULL;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users,refresh_tokens
WHERE refresh_tokens.token = $1 
AND refresh_tokens.user_id = users.id 
AND revoked=FALSE 
LIMIT 1;

-- name: DoesConfirmationTokenExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE confirmation_token=$1) 
AS "exists";

-- name: GetUserByConfirmationToken :one
SELECT id,email, (coalesce(user_info->>'reset_password_on_confirmation','false'))::boolean as ResetPassword
FROM users 
WHERE confirmation_token=$1;

-- name: UpdateConfirmationToken :exec
UPDATE users 
SET 
  confirmation_token = $2,
  confirmation_sent_at = NOW()
WHERE id = $1;

-- name: DoesPasswordResetTokenExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE recovery_token=$1) 
AS "exists";

-- name: UpdatePasswordResetToken :exec
UPDATE users 
SET 
  recovery_token = $2,
  recovery_sent_at = NOW()
WHERE id = $1;

-- name: UpdateRecoveryPassword :exec
UPDATE users 
SET 
  encrypted_password = $2,
  recovery_token= '',
  recovery_sent_at = NULL
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE users 
SET 
  encrypted_password = $2
WHERE id = $1;

-- name: GetUserFromRecoveryToken :one
SELECT id, recovery_sent_at 
FROM users
WHERE recovery_token = $1;


-- name: UpdateUserFirstName :exec
UPDATE users 
SET 
   first_name= $2
WHERE id = $1;

-- name: UpdateUserLastName :exec
UPDATE users 
SET 
   last_name= $2
WHERE id = $1;

-- name: InitUpdateUserEmail :exec
UPDATE users 
SET 
   email_change= $2,
   email_change_token=$3,
   email_change_sent_at=NOW()
WHERE id = $1;

-- name: GetUserFromEmailUpdateToken :one
SELECT id, email_change, email_change_sent_at
FROM users
WHERE email_change_token = $1;

-- name: UpdateUserEmail :exec
UPDATE users 
SET 
   email= $2,
   email_change ='',
   email_change_token='',
   email_change_sent_at=NULL
WHERE id = $1;

-- name: DoesEmailUpdateTokenExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE email_change_token=$1) 
AS "exists";

-- name: GetUsers :many
SELECT id, email, avatar_url, first_name, last_name, last_sign_in_at,created_at, updated_at
FROM users 
ORDER BY created_at DESC
LIMIT $1 
OFFSET $2 ;

-- name: GetUsersByEmail :many
SELECT id, email, avatar_url, first_name, last_name, last_sign_in_at,created_at, updated_at
FROM users 
WHERE email LIKE $1
ORDER BY created_at DESC
LIMIT 20 ;

-- name: GetUserCount :one
SELECT Count(*) FROM users;

-- name: GetConfirmedUserCount :one
SELECT COUNT(*) FROM users WHERE confirmed_at IS NOT NULL;

-- name: UpdateUserInfo :exec
UPDATE users
SET
user_info = jsonb_set(user_info, $1, $2)
WHERE id = $3 ;


-- name: GetUserInfo :one
SELECT user_info
FROM users 
WHERE ID=$1
LIMIT 1;


-- name: GetUsersSignedUpNDaysAgo :many
SELECT id, email, first_name, last_name, user_info
FROM users
WHERE created_at >= CURRENT_DATE - ($1 || ' day')::INTERVAL
  AND created_at < CURRENT_DATE - ($1 || ' day')::INTERVAL + INTERVAL '1 day';


-- name: InsertUserCustom :one
INSERT INTO users (id, email, first_name, last_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetResetPasswordOnConfirmation :one
SELECT (coalesce(user_info->>'reset_password_on_confirmation','false'))::boolean
FROM users WHERE ID=$1;


-- name: WasUserBanned :one
SELECT EXISTS
(SELECT 1 FROM users WHERE email=$1 AND banned_at IS NOT NULL ) 
AS "banned";

-- name: BanUser :exec
UPDATE users
SET
  banned_at = NOW()
WHERE id = $1 ;

-- name: LiftBan :exec
UPDATE users
SET
  banned_at = NULL
WHERE id = $1 ;

-- name: UpdateUserPlan :exec
UPDATE users
SET
  plan = $1
WHERE id = $2 ;

