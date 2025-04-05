package service

import (
	"chenel/eride/app/dto"
	"chenel/eride/app/vault"
	pb "chenel/eride/pb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VaultService struct {
	pb.UnimplementedVaultServiceServer
	authInterceptor *AuthInterceptor
	vaultStore      *vault.VaultStore
}

func NewVaultService(authInterceptor *AuthInterceptor, vaultStore *vault.VaultStore) pb.VaultServiceServer {
	return &VaultService{authInterceptor: authInterceptor, vaultStore: vaultStore}
}

func (server VaultService) CreateVault(ctx context.Context, req *pb.CreateVaultRequest) (*pb.VaultResponse, error) {

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "vault name cannot be empty")
	}

	id, err := server.authInterceptor.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	vault := &dto.Vault{
		Name: req.GetName(),
		ACL: []dto.ACL{
			{
				UserID: id,
				Role:   "OWNER",
			},
		},
	}

	err = server.vaultStore.Save(vault)

	if err != nil {
		return nil, err
	}
	return &pb.VaultResponse{
		Vault: &pb.Vault{
			Id:        vault.ID,
			Name:      vault.Name,
			CreatedAt: vault.CreatedAt.Unix(),
		},
	}, nil
}

// func (server VaultService) AddFolder(context.Context, *pb.AddFolderRequest) (*pb.VaultResponse, error)
// func (server VaultService) UpdateVault(context.Context, *pb.UpdateVaultRequest) (*pb.VaultResponse, error)
