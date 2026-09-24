package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mod.go/internal/model"
	pb "mod.go/internal/proto"
	"mod.go/internal/repository"
)

type PlayerServer struct {
	pb.UnimplementedPlayerServiceServer
	storage *repository.Storage
}

func NewPlayerServer(storage *repository.Storage) *PlayerServer {
	return &PlayerServer{storage: storage}
}

func (s *PlayerServer) CreatePlayer(ctx context.Context, req *pb.CreatePlayerRequest) (*pb.Player, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "имя обязательно")
	}

	player := model.NewPlayer(req.Name)
	created := s.storage.CreatePlayer(*player)
	return toProtoPlayer(created), nil
}

func (s *PlayerServer) GetPlayer(ctx context.Context, req *pb.GetPlayerRequest) (*pb.Player, error) {
	player, found := s.storage.GetPlayerByID(int(req.Id))
	if !found {
		return nil, status.Errorf(codes.NotFound, "игрок с ID %d не найден", req.Id)
	}
	return toProtoPlayer(player), nil
}

func (s *PlayerServer) ListPlayers(ctx context.Context, req *pb.ListPlayersRequest) (*pb.ListPlayersResponse, error) {
	players := s.storage.GetAllPlayers()
	result := make([]*pb.Player, len(players))
	for i, p := range players {
		result[i] = toProtoPlayer(p)
	}
	return &pb.ListPlayersResponse{Players: result}, nil
}

func (s *PlayerServer) UpdatePlayer(ctx context.Context, req *pb.UpdatePlayerRequest) (*pb.Player, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "имя обязательно")
	}

	player := model.NewPlayer(req.Name)
	if err := s.storage.UpdatePlayer(int(req.Id), *player); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	updated, _ := s.storage.GetPlayerByID(int(req.Id))
	return toProtoPlayer(updated), nil
}

func (s *PlayerServer) DeletePlayer(ctx context.Context, req *pb.DeletePlayerRequest) (*pb.DeletePlayerResponse, error) {
	if err := s.storage.DeletePlayer(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeletePlayerResponse{}, nil
}
