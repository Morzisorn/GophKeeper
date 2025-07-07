package client

import (
	"context"
	"errors"
	pb "gophkeeper/internal/protos/users"
)

func (g *GRPCClient) SignUpUser(ctx context.Context, login, password string) (token string, err error) {
	req := &pb.SignUpUserRequest{
		User: &pb.User{
			Login:    login,
			Password: password,
		},
	}

	resp, err := g.User.SignUpUser(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}

func (g *GRPCClient) SignInUser(ctx context.Context, login, password string) (token string, err error) {
	req := &pb.SignInUserRequest{
		User: &pb.User{
			Login:    login,
			Password: password,
		},
	}

	resp, err := g.User.SignInUser(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}

func (c *GRPCClient) SetJWTToken(token string) {
	c.token = token
}

func (c *GRPCClient) GetJWTToken() string {
	return c.token
}