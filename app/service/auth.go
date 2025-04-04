package service

import (
	"chenel/eride/app/auth"
	"chenel/eride/app/dto"
	pb "chenel/eride/pb"
	"fmt"

	"context"

	"golang.org/x/crypto/bcrypt"
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

	if err != nil || user == nil || !IsCorrectPassword(*user, req.GetPassword()) {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	user := &dto.User{
		Username: req.GetUsername(),
		Password: string(hashedPassword),
	}

	err = server.authStore.Save(user)

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

func IsCorrectPassword(user dto.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
