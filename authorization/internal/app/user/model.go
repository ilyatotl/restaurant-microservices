package user

import "time"

type UserModel struct {
	ID           int64     `db:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash"`
	Password     string    `db:"-" json:"password"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserDTO struct {
	ID       int64  `db:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Role     string `db:"role" json:"role"`
}
