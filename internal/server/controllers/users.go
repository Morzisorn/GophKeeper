package controllers

import (
	"context"
	"errors"
	"gophkeeper/internal/errs"
	pb "gophkeeper/internal/protos/users"
	userv "gophkeeper/internal/server/services/user_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	pb.UnimplementedUserControllerServer
	service *userv.UserService
}

func NewUserController(service *userv.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (us *UserController) SignUpUser(ctx context.Context, in *pb.SignUpUserRequest) (*pb.SignUpUserResponse, error) {
	if in.User.Login == "" || in.User.Password == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	token, salt, err := us.service.SignUpUser(ctx, in.User.Login, in.User.Password)
	switch {
	case errors.Is(err, errs.ErrUserAlreadyRegistered):
		return &pb.SignUpUserResponse{
			Error: errs.ErrUserAlreadyRegistered.Error(),
		}, nil
	case err != nil && !errors.Is(err, errs.ErrUserAlreadyRegistered):
		return nil, status.Error(codes.Internal, errs.ErrInternalServerError.Error())
	}

	return &pb.SignUpUserResponse{
		Token: token,
		Salt: salt,
	}, nil
}

func (us *UserController) SignInUser(ctx context.Context, in *pb.SignInUserRequest) (*pb.SignInUserResponse, error) {
	if in.User.Login == "" || in.User.Password == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	token, salt, err := us.service.SignInUser(ctx, in.User.Login, in.User.Password)
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		return &pb.SignInUserResponse{
			Error: errs.ErrUserNotFound.Error(),
		}, nil
	case errors.Is(err, errs.ErrIncorrectCredentials):
		return &pb.SignInUserResponse{
			Error: errs.ErrIncorrectCredentials.Error(),
		}, nil
	case err != nil:
		return nil, status.Error(codes.Internal, errs.ErrInternalServerError.Error())
	}

	return &pb.SignInUserResponse{
		Token: token,
		Salt: salt,
	}, nil
}