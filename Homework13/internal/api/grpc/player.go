package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "mod.go/internal/proto"
	"mod.go/internal/service"
)

type PlayerServer struct {
	pb.UnimplementedPlayerServiceServer
	playerService *service.PlayerService
}

func NewPlayerServer(playerService *service.PlayerService) *PlayerServer {
	return &PlayerServer{playerService: playerService}
}

// CreatePlayer создаёт нового игрока.
func (s *PlayerServer) CreatePlayer(ctx context.Context, req *pb.CreatePlayerRequest) (*pb.Player, error) {
	player, err := s.playerService.Create(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoPlayer(player), nil
}

// GetPlayer возвращает игрока по ID
func (s *PlayerServer) GetPlayer(ctx context.Context, req *pb.GetPlayerRequest) (*pb.Player, error) {
	player, err := s.playerService.Get(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoPlayer(player), nil
}

// ListPlayers возвращает всех игроков
func (s *PlayerServer) ListPlayers(ctx context.Context, req *pb.ListPlayersRequest) (*pb.ListPlayersResponse, error) {
	players := s.playerService.List()
	result := make([]*pb.Player, len(players))
	for i, p := range players {
		result[i] = toProtoPlayer(p)
	}
	return &pb.ListPlayersResponse{Players: result}, nil
}

// UpdatePlayer обновляет игрока по ID.
func (s *PlayerServer) UpdatePlayer(ctx context.Context, req *pb.UpdatePlayerRequest) (*pb.Player, error) {
	updated, err := s.playerService.Update(int(req.Id), req.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoPlayer(updated), nil
}

// DeletePlayer удаляет игрока по ID.
func (s *PlayerServer) DeletePlayer(ctx context.Context, req *pb.DeletePlayerRequest) (*pb.DeletePlayerResponse, error) {
	if err := s.playerService.Delete(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeletePlayerResponse{}, nil
}
