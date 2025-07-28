package controllers

import (
	"context"
	"testing"

	pb "gophkeeper/internal/protos/users"
	userv "gophkeeper/internal/server/services/user_service"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewUserController(t *testing.T) {
	service := &userv.UserService{}
	controller := NewUserController(service)

	assert.NotNil(t, controller)
	assert.Equal(t, service, controller.service)
}

func TestUserController_SignUpUser_Validation(t *testing.T) {
	service := &userv.UserService{}
	controller := NewUserController(service)

	tests := []struct {
		name    string
		request *pb.SignUpUserRequest
		wantErr bool
	}{
		{
			name: "empty login",
			request: &pb.SignUpUserRequest{
				User: &pb.User{
					Login:    "",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "empty password",
			request: &pb.SignUpUserRequest{
				User: &pb.User{
					Login:    "testuser",
					Password: "",
				},
			},
			wantErr: true,
		},
		{
			name: "both empty",
			request: &pb.SignUpUserRequest{
				User: &pb.User{
					Login:    "",
					Password: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := controller.SignUpUser(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			} else {
				assert.NotNil(t, response)
			}
		})
	}
}

func TestUserController_SignInUser_Validation(t *testing.T) {
	service := &userv.UserService{}
	controller := NewUserController(service)

	tests := []struct {
		name    string
		request *pb.SignInUserRequest
		wantErr bool
	}{
		{
			name: "empty login",
			request: &pb.SignInUserRequest{
				User: &pb.User{
					Login:    "",
					Password: "password",
				},
			},
			wantErr: true,
		},
		{
			name: "empty password",
			request: &pb.SignInUserRequest{
				User: &pb.User{
					Login:    "testuser",
					Password: "",
				},
			},
			wantErr: true,
		},
		{
			name: "both empty",
			request: &pb.SignInUserRequest{
				User: &pb.User{
					Login:    "",
					Password: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := controller.SignInUser(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			} else {
				assert.NotNil(t, response)
			}
		})
	}
}

func TestUserController_SignUpUser_ValidRequest(t *testing.T) {
	t.Skip("Skipping valid request test - UserService requires repository dependencies")
}

func TestUserController_SignInUser_ValidRequest(t *testing.T) {
	t.Skip("Skipping valid request test - UserService requires repository dependencies")
}
