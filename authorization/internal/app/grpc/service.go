package grpc_service

import (
	"authorization/internal/app/core"
	"authorization/internal/app/pb"
	"context"
)

type AuthenticationService struct {
	pb.UnimplementedAuthenticationServiceServer

	Core *core.Core
}

func NewAuthenticationService(core *core.Core) *AuthenticationService {
	return &AuthenticationService{Core: core}
}

func (s *AuthenticationService) Authorize(ctx context.Context, req *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
	u, err := s.Core.GetUserInfo(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticationResponse{Id: u.ID, Role: u.Role}, nil
}
