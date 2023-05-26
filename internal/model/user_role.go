package model

// UserRoleID используется для хранения ID роли пользователя
type UserRoleID int

// UserRole описывает сущность "роль пользователя"
type UserRole struct {
	ID   UserRoleID `db:"id"`
	Name string     `db:"name"`
}
