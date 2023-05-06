package user

import "errors"

// Текст ошибок сделан для отображения "пользователю"
var (
	errInvalidUserPasswordConfirm = errors.New("Пароли не совпадают")
	errUsernameIsAlreadyUsed      = errors.New("Данное имя пользователя уже занято")
)

var (
	errInvalidUserRole = errors.New("Указанная роль пользователя не существует")
	errUserNotFound    = errors.New("Пользователь не найден")
)
