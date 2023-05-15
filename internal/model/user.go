package model

import (
	"time"
)

type UserRole int

// User описывает модель пользователя
type User struct {
	Username  string    `db:"username"` // Unique
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      UserRole  `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// CreateUser описывает модель для создания нового пользователя
type CreateUser struct {
	Username        string `db:"username"` // Unique
	Email           string `db:"email"`
	Password        string `db:"password"`
	ConfirmPassword string
	Role            UserRole `db:"role"`
}

// UpdateUser описывает модель для обновления пользователя
type UpdateUser struct {
	Username *string   `db:"username"`
	Email    *string   `db:"email"`
	Password *string   `db:"password"`
	Role     *UserRole `db:"role"`
}
