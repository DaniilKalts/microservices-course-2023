package user

import "time"

// --- Models ---

type User struct {
	ID        string
	Name      string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// --- Service inputs ---

type CreateInput struct {
	User     *User
	Password string
}

type UpdateInput struct {
	ID           string
	Name         *string
	Email        *string
	Password     *string
	PasswordHash *string
}
