package user

type Credentials struct {
	ID           string
	PasswordHash string
	Role         Role
}
