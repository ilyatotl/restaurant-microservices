package sessions

import "time"

type SessionModel struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	SessionToken string    `db:"session_token"`
	ExpiresAt    time.Time `db:"expires_at"`
}
