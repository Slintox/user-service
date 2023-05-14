package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Slintox/user-service/internal/model"
	desc "github.com/Slintox/user-service/pkg/user_v1"
)

// FromUserDesc converts model.User -> grpc.User
func FromUserDesc(user *model.User) *desc.User {
	return &desc.User{
		Username:  user.Username,
		Email:     user.Email,
		Role:      desc.UserRole(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToCreateUserDesc converts grpc.CreateRequest -> model.User
func ToCreateUserDesc(user *desc.CreateRequest) *model.CreateUser {
	return &model.CreateUser{
		Username:        user.GetUsername(),
		Email:           user.GetEmail(),
		Password:        user.GetPassword(),
		ConfirmPassword: user.GetConfirmPassword(),
		Role:            model.UserRole(user.GetRole()),
	}
}

// ToUpdateUserDesc converts grpc.UpdateUser -> model.UpdateUser
func ToUpdateUserDesc(updateUserFields *desc.UpdateUserFields) *model.UpdateUser {
	updUser := &model.UpdateUser{
		Username: updateUserFields.Username,
		Email:    updateUserFields.Email,
		Password: updateUserFields.Password,
	}

	if updateUserFields.Role != nil {
		updUser.Role = new(model.UserRole)
		*updUser.Role = model.UserRole(*updateUserFields.Role)
	}

	return updUser
}
