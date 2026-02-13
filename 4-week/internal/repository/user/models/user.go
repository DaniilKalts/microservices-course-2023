package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string       `db:"id"`
	Name         string       `db:"name"`
	Email        string       `db:"email"`
	PasswordHash string       `db:"password_hash"`
	Role         int32        `db:"role"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}
