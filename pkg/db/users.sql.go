// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const confirmUserEmail = `-- name: ConfirmUserEmail :exec
UPDATE users 
SET 
  confirmed_at = NOW()
WHERE confirmation_token = $1 AND confirmed_at is NULL
`

func (q *Queries) ConfirmUserEmail(ctx context.Context, confirmationToken string) error {
	_, err := q.db.ExecContext(ctx, confirmUserEmail, confirmationToken)
	return err
}

const confirmUserEmailByID = `-- name: ConfirmUserEmailByID :exec
UPDATE users 
SET 
  confirmed_at = NOW()
WHERE id = $1 AND confirmed_at is NULL
`

func (q *Queries) ConfirmUserEmailByID(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, confirmUserEmailByID, id)
	return err
}

const doesConfirmationTokenExist = `-- name: DoesConfirmationTokenExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE confirmation_token=$1) 
AS "exists"
`

func (q *Queries) DoesConfirmationTokenExist(ctx context.Context, confirmationToken string) (bool, error) {
	row := q.db.QueryRowContext(ctx, doesConfirmationTokenExist, confirmationToken)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const doesEmailUpdateTokenExist = `-- name: DoesEmailUpdateTokenExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE email_change_token=$1) 
AS "exists"
`

func (q *Queries) DoesEmailUpdateTokenExist(ctx context.Context, emailChangeToken string) (bool, error) {
	row := q.db.QueryRowContext(ctx, doesEmailUpdateTokenExist, emailChangeToken)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const doesPasswordResetTokenExist = `-- name: DoesPasswordResetTokenExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE recovery_token=$1) 
AS "exists"
`

func (q *Queries) DoesPasswordResetTokenExist(ctx context.Context, recoveryToken string) (bool, error) {
	row := q.db.QueryRowContext(ctx, doesPasswordResetTokenExist, recoveryToken)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const doesUserExist = `-- name: DoesUserExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE email=$1) 
AS "exists"
`

func (q *Queries) DoesUserExist(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRowContext(ctx, doesUserExist, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const doesUserIDExist = `-- name: DoesUserIDExist :one
SELECT EXISTS
(SELECT 1 FROM users WHERE id=$1) 
AS "exists"
`

func (q *Queries) DoesUserIDExist(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRowContext(ctx, doesUserIDExist, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getConfirmedUserCount = `-- name: GetConfirmedUserCount :one
SELECT COUNT(*) FROM users WHERE confirmed_at IS NOT NULL
`

func (q *Queries) GetConfirmedUserCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getConfirmedUserCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getResetPasswordOnConfirmation = `-- name: GetResetPasswordOnConfirmation :one
SELECT (coalesce(user_info->>'reset_password_on_confirmation','false'))::boolean
FROM users WHERE ID=$1
`

func (q *Queries) GetResetPasswordOnConfirmation(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRowContext(ctx, getResetPasswordOnConfirmation, id)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const getUser = `-- name: GetUser :one
SELECT id, email, avatar_url, first_name, last_name, encrypted_password, confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token, email_change, email_change_sent_at, last_sign_in_at, user_info, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.AvatarUrl,
		&i.FirstName,
		&i.LastName,
		&i.EncryptedPassword,
		&i.ConfirmedAt,
		&i.InvitedAt,
		&i.ConfirmationToken,
		&i.ConfirmationSentAt,
		&i.RecoveryToken,
		&i.RecoverySentAt,
		&i.EmailChangeToken,
		&i.EmailChange,
		&i.EmailChangeSentAt,
		&i.LastSignInAt,
		&i.UserInfo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByConfirmationToken = `-- name: GetUserByConfirmationToken :one
SELECT id,email, (coalesce(user_info->>'reset_password_on_confirmation','false'))::boolean as ResetPassword
FROM users 
WHERE confirmation_token=$1
`

type GetUserByConfirmationTokenRow struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Resetpassword bool   `json:"resetpassword"`
}

func (q *Queries) GetUserByConfirmationToken(ctx context.Context, confirmationToken string) (GetUserByConfirmationTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByConfirmationToken, confirmationToken)
	var i GetUserByConfirmationTokenRow
	err := row.Scan(&i.ID, &i.Email, &i.Resetpassword)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, avatar_url, first_name, last_name, encrypted_password, confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token, email_change, email_change_sent_at, last_sign_in_at, user_info, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.AvatarUrl,
		&i.FirstName,
		&i.LastName,
		&i.EncryptedPassword,
		&i.ConfirmedAt,
		&i.InvitedAt,
		&i.ConfirmationToken,
		&i.ConfirmationSentAt,
		&i.RecoveryToken,
		&i.RecoverySentAt,
		&i.EmailChangeToken,
		&i.EmailChange,
		&i.EmailChangeSentAt,
		&i.LastSignInAt,
		&i.UserInfo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserCount = `-- name: GetUserCount :one
SELECT Count(*) FROM users
`

func (q *Queries) GetUserCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUserCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUserFromEmailUpdateToken = `-- name: GetUserFromEmailUpdateToken :one
SELECT id, email_change, email_change_sent_at
FROM users
WHERE email_change_token = $1
`

type GetUserFromEmailUpdateTokenRow struct {
	ID                string       `json:"id"`
	EmailChange       string       `json:"email_change"`
	EmailChangeSentAt sql.NullTime `json:"email_change_sent_at"`
}

func (q *Queries) GetUserFromEmailUpdateToken(ctx context.Context, emailChangeToken string) (GetUserFromEmailUpdateTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromEmailUpdateToken, emailChangeToken)
	var i GetUserFromEmailUpdateTokenRow
	err := row.Scan(&i.ID, &i.EmailChange, &i.EmailChangeSentAt)
	return i, err
}

const getUserFromRecoveryToken = `-- name: GetUserFromRecoveryToken :one
SELECT id, recovery_sent_at 
FROM users
WHERE recovery_token = $1
`

type GetUserFromRecoveryTokenRow struct {
	ID             string       `json:"id"`
	RecoverySentAt sql.NullTime `json:"recovery_sent_at"`
}

func (q *Queries) GetUserFromRecoveryToken(ctx context.Context, recoveryToken string) (GetUserFromRecoveryTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRecoveryToken, recoveryToken)
	var i GetUserFromRecoveryTokenRow
	err := row.Scan(&i.ID, &i.RecoverySentAt)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT users.id, users.email, users.avatar_url, users.first_name, users.last_name, users.encrypted_password, users.confirmed_at, users.invited_at, users.confirmation_token, users.confirmation_sent_at, users.recovery_token, users.recovery_sent_at, users.email_change_token, users.email_change, users.email_change_sent_at, users.last_sign_in_at, users.user_info, users.created_at, users.updated_at FROM users,refresh_tokens
WHERE refresh_tokens.token = $1 
AND refresh_tokens.user_id = users.id 
AND revoked=FALSE 
LIMIT 1
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.AvatarUrl,
		&i.FirstName,
		&i.LastName,
		&i.EncryptedPassword,
		&i.ConfirmedAt,
		&i.InvitedAt,
		&i.ConfirmationToken,
		&i.ConfirmationSentAt,
		&i.RecoveryToken,
		&i.RecoverySentAt,
		&i.EmailChangeToken,
		&i.EmailChange,
		&i.EmailChangeSentAt,
		&i.LastSignInAt,
		&i.UserInfo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserInfo = `-- name: GetUserInfo :one
SELECT user_info
FROM users 
WHERE ID=$1
LIMIT 1
`

func (q *Queries) GetUserInfo(ctx context.Context, id string) (json.RawMessage, error) {
	row := q.db.QueryRowContext(ctx, getUserInfo, id)
	var user_info json.RawMessage
	err := row.Scan(&user_info)
	return user_info, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, email, avatar_url, first_name, last_name, last_sign_in_at,created_at, updated_at
FROM users 
ORDER BY created_at DESC
LIMIT $1 
OFFSET $2
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetUsersRow struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	AvatarUrl    string       `json:"avatar_url"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	LastSignInAt sql.NullTime `json:"last_sign_in_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]GetUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersRow
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.AvatarUrl,
			&i.FirstName,
			&i.LastName,
			&i.LastSignInAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByEmail = `-- name: GetUsersByEmail :many
SELECT id, email, avatar_url, first_name, last_name, last_sign_in_at,created_at, updated_at
FROM users 
WHERE email LIKE $1
ORDER BY created_at DESC
LIMIT 20
`

type GetUsersByEmailRow struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	AvatarUrl    string       `json:"avatar_url"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	LastSignInAt sql.NullTime `json:"last_sign_in_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (q *Queries) GetUsersByEmail(ctx context.Context, email string) ([]GetUsersByEmailRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByEmail, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByEmailRow
	for rows.Next() {
		var i GetUsersByEmailRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.AvatarUrl,
			&i.FirstName,
			&i.LastName,
			&i.LastSignInAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersSignedUpNDaysAgo = `-- name: GetUsersSignedUpNDaysAgo :many
SELECT id, email, first_name, last_name, user_info
FROM users
WHERE created_at >= CURRENT_DATE - ($1 || ' day')::INTERVAL
  AND created_at < CURRENT_DATE - ($1 || ' day')::INTERVAL + INTERVAL '1 day'
`

type GetUsersSignedUpNDaysAgoRow struct {
	ID        string          `json:"id"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	UserInfo  json.RawMessage `json:"user_info"`
}

func (q *Queries) GetUsersSignedUpNDaysAgo(ctx context.Context, dollar_1 sql.NullString) ([]GetUsersSignedUpNDaysAgoRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersSignedUpNDaysAgo, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersSignedUpNDaysAgoRow
	for rows.Next() {
		var i GetUsersSignedUpNDaysAgoRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FirstName,
			&i.LastName,
			&i.UserInfo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const initUpdateUserEmail = `-- name: InitUpdateUserEmail :exec
UPDATE users 
SET 
   email_change= $2,
   email_change_token=$3,
   email_change_sent_at=NOW()
WHERE id = $1
`

type InitUpdateUserEmailParams struct {
	ID               string `json:"id"`
	EmailChange      string `json:"email_change"`
	EmailChangeToken string `json:"email_change_token"`
}

func (q *Queries) InitUpdateUserEmail(ctx context.Context, arg InitUpdateUserEmailParams) error {
	_, err := q.db.ExecContext(ctx, initUpdateUserEmail, arg.ID, arg.EmailChange, arg.EmailChangeToken)
	return err
}

const insertUserCustom = `-- name: InsertUserCustom :one
INSERT INTO users (id, email, first_name, last_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, avatar_url, first_name, last_name, encrypted_password, confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token, email_change, email_change_sent_at, last_sign_in_at, user_info, created_at, updated_at
`

type InsertUserCustomParams struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) InsertUserCustom(ctx context.Context, arg InsertUserCustomParams) (User, error) {
	row := q.db.QueryRowContext(ctx, insertUserCustom,
		arg.ID,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.AvatarUrl,
		&i.FirstName,
		&i.LastName,
		&i.EncryptedPassword,
		&i.ConfirmedAt,
		&i.InvitedAt,
		&i.ConfirmationToken,
		&i.ConfirmationSentAt,
		&i.RecoveryToken,
		&i.RecoverySentAt,
		&i.EmailChangeToken,
		&i.EmailChange,
		&i.EmailChangeSentAt,
		&i.LastSignInAt,
		&i.UserInfo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateConfirmationToken = `-- name: UpdateConfirmationToken :exec
UPDATE users 
SET 
  confirmation_token = $2,
  confirmation_sent_at = NOW()
WHERE id = $1
`

type UpdateConfirmationTokenParams struct {
	ID                string `json:"id"`
	ConfirmationToken string `json:"confirmation_token"`
}

func (q *Queries) UpdateConfirmationToken(ctx context.Context, arg UpdateConfirmationTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateConfirmationToken, arg.ID, arg.ConfirmationToken)
	return err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users 
SET 
  encrypted_password = $2
WHERE id = $1
`

type UpdatePasswordParams struct {
	ID                string `json:"id"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.ExecContext(ctx, updatePassword, arg.ID, arg.EncryptedPassword)
	return err
}

const updatePasswordResetToken = `-- name: UpdatePasswordResetToken :exec
UPDATE users 
SET 
  recovery_token = $2,
  recovery_sent_at = NOW()
WHERE id = $1
`

type UpdatePasswordResetTokenParams struct {
	ID            string `json:"id"`
	RecoveryToken string `json:"recovery_token"`
}

func (q *Queries) UpdatePasswordResetToken(ctx context.Context, arg UpdatePasswordResetTokenParams) error {
	_, err := q.db.ExecContext(ctx, updatePasswordResetToken, arg.ID, arg.RecoveryToken)
	return err
}

const updateRecoveryPassword = `-- name: UpdateRecoveryPassword :exec
UPDATE users 
SET 
  encrypted_password = $2,
  recovery_token= '',
  recovery_sent_at = NULL
WHERE id = $1
`

type UpdateRecoveryPasswordParams struct {
	ID                string `json:"id"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (q *Queries) UpdateRecoveryPassword(ctx context.Context, arg UpdateRecoveryPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateRecoveryPassword, arg.ID, arg.EncryptedPassword)
	return err
}

const updateUserEmail = `-- name: UpdateUserEmail :exec
UPDATE users 
SET 
   email= $2,
   email_change ='',
   email_change_token='',
   email_change_sent_at=NULL
WHERE id = $1
`

type UpdateUserEmailParams struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (q *Queries) UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) error {
	_, err := q.db.ExecContext(ctx, updateUserEmail, arg.ID, arg.Email)
	return err
}

const updateUserFirstName = `-- name: UpdateUserFirstName :exec
UPDATE users 
SET 
   first_name= $2
WHERE id = $1
`

type UpdateUserFirstNameParams struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
}

func (q *Queries) UpdateUserFirstName(ctx context.Context, arg UpdateUserFirstNameParams) error {
	_, err := q.db.ExecContext(ctx, updateUserFirstName, arg.ID, arg.FirstName)
	return err
}

const updateUserInfo = `-- name: UpdateUserInfo :exec
UPDATE users
SET
user_info = jsonb_set(user_info, $1, $2)
WHERE id = $3
`

type UpdateUserInfoParams struct {
	Path        interface{}     `json:"path"`
	Replacement json.RawMessage `json:"replacement"`
	ID          string          `json:"id"`
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateUserInfo, arg.Path, arg.Replacement, arg.ID)
	return err
}

const updateUserLastName = `-- name: UpdateUserLastName :exec
UPDATE users 
SET 
   last_name= $2
WHERE id = $1
`

type UpdateUserLastNameParams struct {
	ID       string `json:"id"`
	LastName string `json:"last_name"`
}

func (q *Queries) UpdateUserLastName(ctx context.Context, arg UpdateUserLastNameParams) error {
	_, err := q.db.ExecContext(ctx, updateUserLastName, arg.ID, arg.LastName)
	return err
}

const updateUserProfile = `-- name: UpdateUserProfile :exec
UPDATE users 
SET 
  first_name = $2,
  last_name = $3,
  avatar_url = $4
WHERE id = $1
`

type UpdateUserProfileParams struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarUrl string `json:"avatar_url"`
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) error {
	_, err := q.db.ExecContext(ctx, updateUserProfile,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.AvatarUrl,
	)
	return err
}

const createNewUser = `-- name: createNewUser :one
INSERT INTO users (
  id,email, encrypted_password
) VALUES (
  $1, $2, $3
)
RETURNING id, email, avatar_url, first_name, last_name, encrypted_password, confirmed_at, invited_at, confirmation_token, confirmation_sent_at, recovery_token, recovery_sent_at, email_change_token, email_change, email_change_sent_at, last_sign_in_at, user_info, created_at, updated_at
`

type createNewUserParams struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (q *Queries) createNewUser(ctx context.Context, arg createNewUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createNewUser, arg.ID, arg.Email, arg.EncryptedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.AvatarUrl,
		&i.FirstName,
		&i.LastName,
		&i.EncryptedPassword,
		&i.ConfirmedAt,
		&i.InvitedAt,
		&i.ConfirmationToken,
		&i.ConfirmationSentAt,
		&i.RecoveryToken,
		&i.RecoverySentAt,
		&i.EmailChangeToken,
		&i.EmailChange,
		&i.EmailChangeSentAt,
		&i.LastSignInAt,
		&i.UserInfo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
