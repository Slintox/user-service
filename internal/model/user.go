package model

import (
	"time"
)

// User описывает модель пользователя
type User struct {
	Username  string    `db:"username"` // Unique
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	RoleID    int       `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// CreateUser описывает модель для создания нового пользователя
type CreateUser struct {
	Username        string `db:"username"` // Unique
	Email           string `db:"email"`
	Password        string `db:"password"`
	ConfirmPassword string `db:"-"`
	RoleID          int    `db:"role"`
}

// UpdateUser описывает модель для обновления пользователя
type UpdateUser struct {
	Username *string `db:"username"`
	Email    *string `db:"email"`
	Password *string `db:"password"`
	RoleID   *int    `db:"role"`
}
