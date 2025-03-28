package service

import (
	"chenel/passport/app/auth"
	pb "chenel/passport/pb/pb_auth_service"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	pb.UnimplementedPBAuthServiceServer
	authStore  auth.AuthStore
	jwtManager *auth.JWTManager
}

func NewAuthService(authStore auth.AuthStore, jwtManager *auth.JWTManager) pb.PBAuthServiceServer {
	return &AuthService{authStore: authStore, jwtManager: jwtManager}
}

func (server *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.authStore.Find(req.GetUsername())

	if err != nil || user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.PermissionDenied, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.LoginResponse{
		AccessToken: token}
	return res, nil
}

func (server *AuthService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, _ := auth.NewUser(req.GetUsername(), req.GetPassword())
	err := server.authStore.Save(user)

	if err != nil {
		return nil, err
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.CreateUserResponse{
		AccessToken: token}
	return res, nil
}
