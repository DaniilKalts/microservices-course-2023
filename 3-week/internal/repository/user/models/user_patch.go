package models

type UpdateUserPatch struct {
	Name         *string `db:"name"`
	Email        *string `db:"email"`
	PasswordHash *string `db:"password_hash"`
}
