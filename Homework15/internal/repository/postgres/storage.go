package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"mod.go/internal/config"
	"mod.go/internal/model"
	"mod.go/internal/repository/postgres/sqlc"
)

type PostgresStorage struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("база данных не отвечает: %w", err)
	}

	fmt.Println("Успешно подключено к PostgreSQL!")

	return &PostgresStorage{
		pool:    pool,
		queries: sqlc.New(pool),
	}, nil
}

func (s *PostgresStorage) Close(ctx context.Context) error {
	s.pool.Close()
	return nil
}

func (s *PostgresStorage) ImportPlayer(id int, name string) error {
	ctx := context.Background()
	_, err := s.queries.CreatePlayer(ctx, sqlc.CreatePlayerParams{
		ID:   int64(id),
		Name: name,
	})
	return err
}

func (s *PostgresStorage) ImportGame(game model.Game) error {
	ctx := context.Background()
	boardBytes, _ := json.Marshal(game.Board())
	_, err := s.queries.CreateGame(ctx, sqlc.CreateGameParams{
		ID:              int64(game.ID()),
		Player1ID:       int64(game.Player1().ID()),
		Player1Name:     game.Player1().Name(),
		Player1Color:    int32(game.Player1Color()),
		Player2ID:       int64(game.Player2().ID()),
		Player2Name:     game.Player2().Name(),
		Player2Color:    int32(game.Player2Color()),
		Board:           boardBytes,
		CurrentColor:    int32(game.CurrentColor()),
		WhiteKingMoved:  game.WhiteKingMoved(),
		BlackKingMoved:  game.BlackKingMoved(),
		WhiteRookAMoved: game.WhiteRookAMoved(),
		WhiteRookHMoved: game.WhiteRookHMoved(),
		BlackRookAMoved: game.BlackRookAMoved(),
		BlackRookHMoved: game.BlackRookHMoved(),
	})
	return err
}

func (s *PostgresStorage) ImportMove(move model.Move) error {
	ctx := context.Background()
	var capturedBytes, movedPieceBytes []byte
	if move.Captured != nil {
		capturedBytes, _ = json.Marshal(move.Captured)
	}
	if move.MovedPiece != nil {
		movedPieceBytes, _ = json.Marshal(move.MovedPiece)
	}
	_, err := s.queries.CreateMove(ctx, sqlc.CreateMoveParams{
		ID:         int64(move.ID()),
		GameID:     int64(move.GameID()),
		FromRow:    int32(move.FromRow),
		FromCol:    int32(move.FromCol),
		ToRow:      int32(move.ToRow),
		ToCol:      int32(move.ToCol),
		Captured:   capturedBytes,
		Promotion:  int32(move.Promotion),
		MovedPiece: movedPieceBytes,
		Check:      move.Check,
		Mate:       move.Mate,
	})
	return err
}

func (s *PostgresStorage) SetCounter(name string, seq int) error {
	ctx := context.Background()
	return s.queries.SetCounter(ctx, sqlc.SetCounterParams{
		Name: name,
		Seq:  int64(seq),
	})
}

func (s *PostgresStorage) GetAllCounters() ([]sqlc.Counter, error) {
	ctx := context.Background()
	return s.queries.ListCounters(ctx)
}
