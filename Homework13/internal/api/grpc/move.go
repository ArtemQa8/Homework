package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "mod.go/internal/proto"
	"mod.go/internal/service"
)

type MoveServer struct {
	pb.UnimplementedMoveServiceServer
	moveService *service.MoveService
}

func NewMoveServer(moveService *service.MoveService) *MoveServer {
	return &MoveServer{moveService: moveService}
}

func (s *MoveServer) CreateMove(ctx context.Context, req *pb.CreateMoveRequest) (*pb.Move, error) {
	move := moveFromProto(req.FromRow, req.FromCol, req.ToRow, req.ToCol, req.Promotion)
	move.SetGameID(int(req.GameId))

	created, err := s.moveService.Create(move)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoMove(created), nil
}

func (s *MoveServer) GetMove(ctx context.Context, req *pb.GetMoveRequest) (*pb.Move, error) {
	move, err := s.moveService.Get(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoMove(move), nil
}

func (s *MoveServer) ListMoves(ctx context.Context, _ *emptypb.Empty) (*pb.ListMovesResponse, error) {
	moves := s.moveService.List()
	result := make([]*pb.Move, len(moves))
	for i, m := range moves {
		result[i] = toProtoMove(m)
	}
	return &pb.ListMovesResponse{Moves: result}, nil
}

func (s *MoveServer) UpdateMove(ctx context.Context, req *pb.UpdateMoveRequest) (*pb.Move, error) {

	move := moveFromProto(req.FromRow, req.FromCol, req.ToRow, req.ToCol, pb.PieceType_PIECE_TYPE_UNSPECIFIED)
	move.SetGameID(int(req.GameId))
	move.SetID(int(req.Id))

	updated, err := s.moveService.Update(int(req.Id), move)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoMove(updated), nil
}

func (s *MoveServer) DeleteMove(ctx context.Context, req *pb.DeleteMoveRequest) (*emptypb.Empty, error) {
	if err := s.moveService.Delete(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &emptypb.Empty{}, nil
}
