package user

import "time"

type Entity struct {
	ID        string
	Name      string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt *time.Time
}
