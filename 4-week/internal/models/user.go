package models

import (
	"time"
)

type Role int32

const (
	RoleUser Role = iota
	RoleAdmin
)

type User struct {
	ID        string
	Name      string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt *time.Time
}
