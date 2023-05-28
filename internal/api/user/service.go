package user

import (
	"context"
	userV1 "github.com/Slintox/user-service/pkg/service/user_v1"

	"google.golang.org/protobuf/types/known/emptypb"

	converter "github.com/Slintox/user-service/internal/converter/user"
	"github.com/Slintox/user-service/internal/service/user"
)

type Implementation struct {
	userV1.UnimplementedUserV1Server

	userService user.Service
}

func NewImplementation(userService user.Service) *Implementation {
	return &Implementation{
		userService: userService,
	}
}

func (i *Implementation) Create(ctx context.Context, req *userV1.CreateRequest) (*emptypb.Empty, error) {
	err := i.userService.Create(ctx, converter.ToCreateUserDesc(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i *Implementation) Get(ctx context.Context, req *userV1.GetRequest) (*userV1.GetResponse, error) {
	userView, err := i.userService.Get(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	return &userV1.GetResponse{
		User: converter.FromUserDesc(userView),
	}, nil
}

func (i *Implementation) Update(ctx context.Context, req *userV1.UpdateRequest) (*emptypb.Empty, error) {
	if req.UpdateData == nil {
		return nil, errNoDataToUpdate
	}

	err := i.userService.Update(ctx, req.GetUsername(), converter.ToUpdateUserDesc(req.GetUpdateData()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i *Implementation) Delete(ctx context.Context, req *userV1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
