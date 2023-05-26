package user

import (
	"context"
	"errors"
	userRole "github.com/Slintox/user-service/internal/repository/user_role"

	"github.com/Slintox/user-service/internal/model"
	repo "github.com/Slintox/user-service/internal/repository"
	uRepo "github.com/Slintox/user-service/internal/repository/user"
)

type Service interface {
	Create(ctx context.Context, user *model.CreateUser) error
	Get(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, username string, updateData *model.UpdateUser) error
	Delete(ctx context.Context, username string) error
}

type service struct {
	userRepo     uRepo.Repository
	userRoleRepo userRole.Repository
}

func NewService(userRepo uRepo.Repository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) Create(ctx context.Context, user *model.CreateUser) error {
	// Проверка на правильность пароля
	if user.Password != user.ConfirmPassword {
		return errInvalidUserPasswordConfirm
	}

	// Проверка на доступность username
	isUsernameAvailable, err := s.userRepo.IsUsernameAvailable(ctx, user.Username)
	if err != nil {
		return err
	}
	if !isUsernameAvailable {
		return errUsernameIsAlreadyUsed
	}

	// Проверка на существование роли
	isRoleExist, err := s.userRoleRepo.IsRoleExist(ctx, user.RoleID)
	if !isRoleExist {
		return errInvalidUserRole
	}

	// Сохранение нового пользователя
	if err = s.userRepo.Add(ctx, user); err != nil {
		if errors.Is(err, repo.ErrRecordNotFound) {
			return errInvalidUserRole
		}
		return err
	}

	return nil
}

func (s *service) Get(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.Get(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrRecordNotFound) {
			return nil, errUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, username string, updateData *model.UpdateUser) error {
	// Проверка на возможность обновления username
	if updateData.Username != nil {
		isUsernameAvailable, err := s.userRepo.IsUsernameAvailable(ctx, *updateData.Username)
		if err != nil {
			return err
		}

		if !isUsernameAvailable {
			return errUsernameIsAlreadyUsed
		}
	}

	// Обновление пользователя
	if err := s.userRepo.Update(ctx, username, updateData); err != nil {
		if errors.Is(err, repo.ErrRecordNotFound) {
			return errUserNotFound
		}
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, username string) error {
	err := s.userRepo.Delete(ctx, username)
	if err != nil {
		return err
	}

	return nil
}
