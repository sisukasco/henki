// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

type ApiKey struct {
	Key       string    `json:"key"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        int64     `json:"id"`
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID                 string          `json:"id"`
	Email              string          `json:"email"`
	AvatarUrl          string          `json:"avatar_url"`
	FirstName          string          `json:"first_name"`
	LastName           string          `json:"last_name"`
	EncryptedPassword  string          `json:"encrypted_password"`
	ConfirmedAt        sql.NullTime    `json:"confirmed_at"`
	InvitedAt          sql.NullTime    `json:"invited_at"`
	ConfirmationToken  string          `json:"confirmation_token"`
	ConfirmationSentAt sql.NullTime    `json:"confirmation_sent_at"`
	RecoveryToken      string          `json:"recovery_token"`
	RecoverySentAt     sql.NullTime    `json:"recovery_sent_at"`
	EmailChangeToken   string          `json:"email_change_token"`
	EmailChange        string          `json:"email_change"`
	EmailChangeSentAt  sql.NullTime    `json:"email_change_sent_at"`
	LastSignInAt       sql.NullTime    `json:"last_sign_in_at"`
	UserInfo           json.RawMessage `json:"user_info"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
