package main

import (
	"chenel/eride/app/auth"
	"chenel/eride/app/config"
	"chenel/eride/app/consts"
	"chenel/eride/app/db"
	"chenel/eride/app/security"
	"chenel/eride/app/service"
	"chenel/eride/app/vault"
	pb "chenel/eride/pb"
	"context"
	"fmt"
	"log"
	"net"
	"time"

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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	creds, err := security.LoadTLSCredentials(
		&security.TLSConfigData{
			CA:         cfg.GetString("tls.ca"),
			ServerCert: cfg.GetString("tls.cert"),
			ServerKey:  cfg.GetString("tls.key"),
		},
	)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	dbstore, err := db.NewSQLiteDBProvider("password_manager.db")
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	authStore := auth.NewUserStore(dbstore)
	fmt.Printf("Token duration: %d ms\n", time.Duration(cfg.GetUint64("security.jwtTokenDuration")))

	jwtManager := auth.NewJWTManager(cfg.GetString("security.jwtSecret"), time.Duration(cfg.GetUint64("security.jwtTokenDuration"))*time.Second)
	vaultStore := vault.NewVaultStore(dbstore)

	// GRPC server with TLS

	authInterceptor := service.NewAuthInterceptor(jwtManager)
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	reflection.Register(grpcServer)

	pb.RegisterAuthServiceServer(grpcServer, service.NewAuthService(authStore, jwtManager))
	pb.RegisterVaultServiceServer(grpcServer, service.NewVaultService(authInterceptor, vaultStore))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GetInt("server.port")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("%s is running on port %d with TLS...", consts.APP_NAME, cfg.GetInt("server.port"))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
