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
	vaultStore vault.VaultStore
}

func NewVaultService(vaultStore vault.VaultStore) pb.VaultServiceServer {
	return &VaultService{vaultStore: vaultStore}
}

// func (server VaultService) AddFolder(context.Context, *pb.AddFolderRequest) (*pb.VaultResponse, error)
func (server VaultService) CreateVault(ctx context.Context, req *pb.CreateVaultRequest) (*pb.VaultResponse, error) {

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "vault name cannot be empty")
	}

	vault := &dto.Vault{
		Name: req.GetName(),
		ACL: []dto.ACL{
			{
				UserID: func(u uint64) *uint64 { return &u }(0),
				Role:   "OWNER",
			},
		},
	}

	err := server.vaultStore.Save(vault)
	if err != nil {
		return nil, err
	}
	return &pb.VaultResponse{
		Vault: &pb.Vault{
			Id:   vault.ID,
			Name: vault.Name,
		},
	}, nil
}

// func (server VaultService) UpdateVault(context.Context, *pb.UpdateVaultRequest) (*pb.VaultResponse, error)
