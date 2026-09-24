package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mod.go/internal/model"
	pb "mod.go/internal/proto"
	"mod.go/internal/repository"
	"mod.go/internal/service"
)

type GameServer struct {
	pb.UnimplementedGameServiceServer
	storage *repository.Storage
}

func NewGameServer(storage *repository.Storage) *GameServer {
	return &GameServer{storage: storage}
}

func (s *GameServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.Game, error) {
	if req.Player1Name == "" || req.Player2Name == "" {
		return nil, status.Error(codes.InvalidArgument, "имена обоих игроков обязательны")
	}

	rows := int(req.Rows)
	cols := int(req.Cols)
	if rows <= 0 || cols <= 0 {
		rows, cols = 8, 8
	}

	p1 := model.NewPlayer(req.Player1Name)
	p2 := model.NewPlayer(req.Player2Name)

	game := model.Game{}
	game.SetPlayer1(*p1)
	game.SetPlayer2(*p2)
	game.SetBoard(model.NewBoard(rows, cols))

	created, err := s.storage.CreateGame(game)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	created.SetPlayer1Color(model.White)
	created.SetPlayer2Color(model.Black)

	if err := s.storage.OverwriteGame(created.ID(), created); err != nil {
		return nil, status.Error(codes.Internal, "не удалось сохранить игру")
	}

	return toProtoGame(created), nil
}

func (s *GameServer) GetGame(ctx context.Context, req *pb.GetGameRequest) (*pb.Game, error) {
	game, found := s.storage.GetGameByID(int(req.Id))
	if !found {
		return nil, status.Errorf(codes.NotFound, "игра с ID %d не найдена", req.Id)
	}
	return toProtoGame(game), nil
}

func (s *GameServer) ListGames(ctx context.Context, req *pb.ListGamesRequest) (*pb.ListGamesResponse, error) {
	games := s.storage.GetAllGames()
	result := make([]*pb.Game, len(games))
	for i, g := range games {
		result[i] = toProtoGame(g)
	}
	return &pb.ListGamesResponse{Games: result}, nil
}

func (s *GameServer) UpdateGame(ctx context.Context, req *pb.UpdateGameRequest) (*pb.Game, error) {
	game, found := s.storage.GetGameByID(int(req.Id))
	if !found {
		return nil, status.Errorf(codes.NotFound, "игра с ID %d не найдена", req.Id)
	}

	if req.Player1Name != "" {
		game.SetPlayer1(*model.NewPlayer(req.Player1Name))
	}
	if req.Player2Name != "" {
		game.SetPlayer2(*model.NewPlayer(req.Player2Name))
	}

	if err := s.storage.OverwriteGame(int(req.Id), game); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	updated, _ := s.storage.GetGameByID(int(req.Id))
	return toProtoGame(updated), nil
}

func (s *GameServer) DeleteGame(ctx context.Context, req *pb.DeleteGameRequest) (*pb.DeleteGameResponse, error) {
	if err := s.storage.DeleteGame(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteGameResponse{}, nil
}

func (s *GameServer) MakeMove(ctx context.Context, req *pb.MakeMoveRequest) (*pb.Game, error) {
	game, found := s.storage.GetGameByID(int(req.GameId))
	if !found {
		return nil, status.Errorf(codes.NotFound, "игра с ID %d не найдена", req.GameId)
	}

	move := &model.Move{
		FromRow: int(req.FromRow),
		FromCol: int(req.FromCol),
		ToRow:   int(req.ToRow),
		ToCol:   int(req.ToCol),
	}

	if req.Promotion != pb.PieceType_PIECE_TYPE_UNSPECIFIED {
		move.Promotion = fromProtoPieceType(req.Promotion)
	}

	move.SetGameID(int(req.GameId))

	if err := game.MakeMove(move); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	savedMove, err := s.storage.CreateMove(*move)
	if err != nil {
		return nil, status.Error(codes.Internal, "не удалось сохранить ход")
	}
	game.SetLastMoveID(savedMove.ID())

	if err := s.storage.OverwriteGame(int(req.GameId), game); err != nil {
		return nil, status.Error(codes.Internal, "не удалось сохзранить игру")
	}

	return toProtoGame(game), nil
}

func (s *GameServer) AutoMove(ctx context.Context, req *pb.AutoMoveRequest) (*pb.Game, error) {
	game, found := s.storage.GetGameByID(int(req.GameId))
	if !found {
		return nil, status.Errorf(codes.NotFound, "игра с ID %d не найдена", req.GameId)
	}

	color := game.CurrentColor()
	if game.Checkmate(color) || game.Stalemate(color) {
		return toProtoGame(game), nil
	}

	move, err := service.ChooseRandomMove(&game)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if piece := game.Board().PieceAt(move.FromRow, move.FromCol); piece != nil && piece.Type() == model.Pawn {
		lastRow := 0
		if piece.Color() == model.White {
			lastRow = game.Board().Rows() - 1
		}
		if move.ToRow == lastRow {
			move.Promotion = model.Queen
		}
	}

	move.SetGameID(int(req.GameId))
	if err := game.MakeMove(move); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	savedMove, err := s.storage.CreateMove(*move)
	if err != nil {
		return nil, status.Error(codes.Internal, "не удалось сохранить ход")
	}
	game.SetLastMoveID(savedMove.ID())

	if err := s.storage.OverwriteGame(int(req.GameId), game); err != nil {
		return nil, status.Error(codes.Internal, "не удалось сохранить игру")
	}

	return toProtoGame(game), nil
}
