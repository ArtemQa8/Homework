package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mod.go/internal/model"
	pb "mod.go/internal/proto"
	"mod.go/internal/repository"
)

type MoveServer struct {
	pb.UnimplementedMoveServiceServer
	storage *repository.Storage
}

func NewMoveServer(storage *repository.Storage) *MoveServer {
	return &MoveServer{storage: storage}
}

func (s *MoveServer) CreateMove(ctx context.Context, req *pb.CreateMoveRequest) (*pb.Move, error) {
	move := model.Move{
		FromRow: int(req.FromRow),
		FromCol: int(req.FromCol),
		ToRow:   int(req.ToRow),
		ToCol:   int(req.ToCol),
	}

	if req.Promotion != pb.PieceType_PIECE_TYPE_UNSPECIFIED {
		move.Promotion = fromProtoPieceType(req.Promotion)
	}

	move.SetGameID(int(req.GameId))

	created, err := s.storage.CreateMove(move)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoMove(created), nil
}

func (s *MoveServer) GetMove(ctx context.Context, req *pb.GetMoveRequest) (*pb.Move, error) {
	move, found := s.storage.GetMoveByID(int(req.Id))
	if !found {
		return nil, status.Errorf(codes.NotFound, "ход с ID %d не найден", req.Id)
	}
	return toProtoMove(move), nil
}

func (s *MoveServer) ListMoves(ctx context.Context, _ *emptypb.Empty) (*pb.ListMovesResponse, error) {
	moves := s.storage.GetAllMoves()
	result := make([]*pb.Move, len(moves))
	for i, m := range moves {
		result[i] = toProtoMove(m)
	}
	return &pb.ListMovesResponse{Moves: result}, nil
}

func (s *MoveServer) UpdateMove(ctx context.Context, req *pb.UpdateMoveRequest) (*pb.Move, error) {
	_, found := s.storage.GetMoveByID(int(req.Id))
	if !found {
		return nil, status.Errorf(codes.NotFound, "ход с ID %d не найден", req.Id)
	}

	move := model.Move{
		FromRow: int(req.FromRow),
		FromCol: int(req.FromCol),
		ToRow:   int(req.ToRow),
		ToCol:   int(req.ToCol),
	}
	move.SetGameID(int(req.GameId))
	move.SetID(int(req.Id))

	if err := s.storage.UpdateMove(int(req.Id), move); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	updated, _ := s.storage.GetMoveByID(int(req.Id))
	return toProtoMove(updated), nil
}

func (s *MoveServer) DeleteMove(ctx context.Context, req *pb.DeleteMoveRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteMove(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &emptypb.Empty{}, nil
}
