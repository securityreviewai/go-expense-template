package models

import "time"

type Role string

const (
	RoleUser           Role = "user"
	RoleProjectManager Role = "project_manager"
	RoleSuperAdmin     Role = "super_admin"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         Role      `json:"role"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
