package repository

import "errors"

var (
	// ErrRecordNotFound ошибка возвращаемая из уровня repository
	// для обработки на уровне usecase.
	ErrRecordNotFound = errors.New("Запись не найдена")
)
