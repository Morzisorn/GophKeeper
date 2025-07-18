package client

import (
	"context"
	"errors"
	pb "gophkeeper/internal/protos/users"
	"gophkeeper/models"
)

func (g *GRPCClient) SignUpUser(ctx context.Context, user *models.User) (token string, salt string, err error) {
	req := &pb.SignUpUserRequest{
		User: &pb.User{
			Login:    user.Login,
			Password: string(user.Password),
		},
	}

	resp, err := g.User.SignUpUser(ctx, req)
	if err != nil {
		return "", "", err
	}
	if resp.Error != "" {
		return "", "", errors.New(resp.Error)
	}

	return resp.Token, resp.Salt, nil
}

func (g *GRPCClient) SignInUser(ctx context.Context, user *models.User) (token string, salt string, err error) {
	req := &pb.SignInUserRequest{
		User: &pb.User{
			Login:    user.Login,
			Password: string(user.Password),
		},
	}

	resp, err := g.User.SignInUser(ctx, req)
	if err != nil {
		return "", "", err
	}
	if resp.Error != "" {
		return "", "", errors.New(resp.Error)
	}

	return resp.Token, resp.Salt, nil
}

func (c *GRPCClient) SetJWTToken(token string) {
	c.token = token
}

func (c *GRPCClient) GetJWTToken() string {
	return c.token
}
