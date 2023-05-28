package user

import (
	userV1 "github.com/Slintox/user-service/pkg/service/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Slintox/user-service/internal/model"
)

// FromUserDesc converts model.User -> grpc.User
func FromUserDesc(user *model.User) *userV1.User {
	return &userV1.User{
		Username:  user.Username,
		Email:     user.Email,
		RoleId:    userV1.UserRoleID(user.RoleID),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToCreateUserDesc converts grpc.CreateRequest -> model.User
func ToCreateUserDesc(user *userV1.CreateRequest) *model.CreateUser {
	return &model.CreateUser{
		Username:        user.GetUsername(),
		Email:           user.GetEmail(),
		Password:        user.GetPassword(),
		ConfirmPassword: user.GetConfirmPassword(),
		RoleID:          int(user.GetRoleId()),
	}
}

// ToUpdateUserDesc converts grpc.UpdateUser -> model.UpdateUser
func ToUpdateUserDesc(updateUserFields *userV1.UpdateUserFields) *model.UpdateUser {
	updUser := &model.UpdateUser{
		Username: updateUserFields.Username,
		Email:    updateUserFields.Email,
		Password: updateUserFields.Password,
	}

	if updateUserFields.Role != nil {
		updUser.RoleID = new(int)
		*updUser.RoleID = int(*updateUserFields.Role)
	}

	return updUser
}
