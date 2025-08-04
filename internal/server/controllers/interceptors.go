package controllers

import (
	"context"
	"fmt"
	"gophkeeper/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(cnfg config.ServerInterceptorsConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return AuthInterceptor(ctx, cnfg, req, info, handler)
	}
}

func AuthInterceptor(ctx context.Context, cnfg config.ServerInterceptorsConfig, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")

	claims, err := validateJWT(cnfg, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	login, ok := claims["login"].(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: missing login")
	}

	ctx = context.WithValue(ctx, "user_claims", claims)
	ctx = context.WithValue(ctx, "login", login)

	return handler(ctx, req)
}

func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/users.UserController/SignUpUser",
		"/users.UserController/SignInUser",
		"/crypto.CryptoController/GetPublicKeyPEM",
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}
	return false
}

type Claims struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

func validateJWT(cnfg config.ServerInterceptorsConfig, tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte(cnfg.GetSecretKey())

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
