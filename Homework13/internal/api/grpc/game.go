package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mod.go/internal/model"
	pb "mod.go/internal/proto"
	"mod.go/internal/service"
)

type GameServer struct {
	pb.UnimplementedGameServiceServer
	gameService *service.GameService
}

func NewGameServer(gameService *service.GameService) *GameServer {
	return &GameServer{gameService: gameService}
}

func (s *GameServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.Game, error) {
	p1 := model.NewPlayer(req.Player1Name)
	p2 := model.NewPlayer(req.Player2Name)

	game, err := s.gameService.Create(*p1, *p2, int(req.Rows), int(req.Cols))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoGame(game), nil
}

func (s *GameServer) GetGame(ctx context.Context, req *pb.GetGameRequest) (*pb.Game, error) {
	game, err := s.gameService.Get(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoGame(game), nil
}

func (s *GameServer) ListGames(ctx context.Context, req *pb.ListGamesRequest) (*pb.ListGamesResponse, error) {
	games := s.gameService.List()
	result := make([]*pb.Game, len(games))
	for i, g := range games {
		result[i] = toProtoGame(g)
	}
	return &pb.ListGamesResponse{Games: result}, nil
}

func (s *GameServer) UpdateGame(ctx context.Context, req *pb.UpdateGameRequest) (*pb.Game, error) {
	p1 := model.NewPlayer(req.Player1Name)
	p2 := model.NewPlayer(req.Player2Name)

	game, err := s.gameService.Update(int(req.Id), *p1, *p2)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoGame(game), nil
}

func (s *GameServer) DeleteGame(ctx context.Context, req *pb.DeleteGameRequest) (*pb.DeleteGameResponse, error) {
	if err := s.gameService.Delete(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteGameResponse{}, nil
}

func (s *GameServer) MakeMove(ctx context.Context, req *pb.MakeMoveRequest) (*pb.Game, error) {
	move := &model.Move{
		FromRow: int(req.FromRow),
		FromCol: int(req.FromCol),
		ToRow:   int(req.ToRow),
		ToCol:   int(req.ToCol),
	}

	if req.Promotion != pb.PieceType_PIECE_TYPE_UNSPECIFIED {
		move.Promotion = fromProtoPieceType(req.Promotion)
	}

	game, err := s.gameService.MakeMove(int(req.GameId), move)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoGame(game), nil
}

func (s *GameServer) AutoMove(ctx context.Context, req *pb.AutoMoveRequest) (*pb.Game, error) {
	game, err := s.gameService.AutoMove(int(req.GameId))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoGame(game), nil
}
