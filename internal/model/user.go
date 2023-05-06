package model

import (
	"time"
)

type UserRole int

// User описывает модель пользователя
type User struct {
	Username  string    // Unique
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateUser описывает модель для создания нового пользователя
type CreateUser struct {
	Username        string   `json:"username"` // Unique
	Email           string   `json:"email"`
	Password        string   `json:"password"`
	ConfirmPassword string   `json:"confirmPassword"`
	Role            UserRole `json:"role"`
}

// UpdateUser описывает модель для обновления пользователя
type UpdateUser struct {
	Username *string   `json:"username"`
	Email    *string   `json:"email"`
	Password *string   `json:"password"`
	Role     *UserRole `json:"role"`
}
