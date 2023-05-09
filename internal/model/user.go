package model

import (
	"time"
)

type UserRole int

// User описывает модель пользователя
type User struct {
	Username  string // Unique
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUser описывает модель для создания нового пользователя
type CreateUser struct {
	Username        string // Unique
	Email           string
	Password        string
	ConfirmPassword string
	Role            UserRole
}

// UpdateUser описывает модель для обновления пользователя
type UpdateUser struct {
	Username *string
	Email    *string
	Password *string
	Role     *UserRole
}
