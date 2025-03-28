package main

import (
	"chenel/passport/app/db"
	"chenel/passport/app/security"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, stream)
}

func main() {

	creds, err := security.LoadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	dbstore, err := db.NewDBProvider("password_manager.db")

	dbstore.Create(&db.User{FirstName: "Alice", LastName: "Smith", Username: "alice_smith", Password: "hashed_password_1", TOTPSecret: "TOTP_SECRET_1"})

	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	// authStore := auth.NewUserStore(dbstore.DB)
	// jwtManager := auth.NewJWTManager("secret", 150000)

	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))

	// Postman
	reflection.Register(grpcServer)

	//Nombres pochos
	//pb_auth_service.RegisterPBAuthServiceServer(grpcServer, service.NewAuthService(*authStore, jwtManager))

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server is running on port 50051 with TLS...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
