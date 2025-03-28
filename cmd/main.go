package main

import (
	"chenel/passport/app/auth"
	"chenel/passport/app/config"
	"chenel/passport/app/consts"
	"chenel/passport/app/db"
	"chenel/passport/app/security"
	"chenel/passport/pb/pb_auth_service"
	"chenel/passport/service"
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
	jwtManager := auth.NewJWTManager(cfg.GetString("security.jwtSecret"), time.Duration(cfg.GetUint32("security.jwtTokenDuration")))

	// GRPC server with TLS
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))
	reflection.Register(grpcServer)

	//Nombres pochos
	pb_auth_service.RegisterPBAuthServiceServer(grpcServer, service.NewAuthService(*authStore, jwtManager))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GetInt("server.port")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("%s is running on port %d with TLS...", consts.APP_NAME, cfg.GetInt("server.port"))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
