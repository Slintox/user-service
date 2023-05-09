package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	converter "github.com/Slintox/user-service/internal/converter/user"
	"github.com/Slintox/user-service/internal/service/user"
	desc "github.com/Slintox/user-service/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server

	userService user.Service
}

func NewImplementation(userService user.Service) *Implementation {
	return &Implementation{
		userService: userService,
	}
}

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*emptypb.Empty, error) {
	err := i.userService.Create(ctx, converter.ToCreateUserDesc(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userView, err := i.userService.Get(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.FromUserDesc(userView),
	}, nil
}

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.UpdateData == nil {
		return nil, errNoDataToUpdate
	}

	err := i.userService.Update(ctx, req.GetUsername(), converter.ToUpdateUserDesc(req.GetUpdateData()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
