package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"mod.go/internal/model"
	"mod.go/internal/repository/postgres/sqlc"
)

func (s *PostgresStorage) CreateMove(move model.Move) (model.Move, error) {
	ctx := context.Background()

	id, err := s.queries.GetNextSequence(ctx, "move_id")
	if err != nil {
		return model.Move{}, err
	}
	move.SetID(int(id))

	var captureBytes []byte
	if move.Captured != nil {
		captureBytes, _ = json.Marshal(move.Captured)
	}

	var movedPieceBytes []byte
	if move.MovedPiece != nil {
		movedPieceBytes, _ = json.Marshal(move.MovedPiece)
	}

	dbMove, err := s.queries.CreateMove(ctx, sqlc.CreateMoveParams{
		ID:         id,
		GameID:     int64(move.GameID()),
		FromRow:    int32(move.FromRow),
		FromCol:    int32(move.FromCol),
		ToRow:      int32(move.ToRow),
		ToCol:      int32(move.ToCol),
		Captured:   captureBytes,
		Promotion:  int32(move.Promotion),
		MovedPiece: movedPieceBytes,
		Check:      move.Check,
		Mate:       move.Mate,
	})
	if err != nil {
		return model.Move{}, fmt.Errorf("ошибка сохранения хода: %w", err)
	}

	return s.toModelMove(dbMove), nil
}

func (s *PostgresStorage) GetAllMoves() []model.Move {
	ctx := context.Background()
	dbMoves, err := s.queries.ListMoves(ctx)
	if err != nil {
		return []model.Move{}
	}

	moves := make([]model.Move, 0, len(dbMoves))
	for _, dbM := range dbMoves {
		moves = append(moves, s.toModelMove(dbM))
	}
	return moves
}

func (s *PostgresStorage) GetMoveByID(id int) (model.Move, bool) {
	ctx := context.Background()
	dbMove, err := s.queries.GetMove(ctx, int64(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Move{}, false
		}
		return model.Move{}, false
	}

	return s.toModelMove(dbMove), true
}

func (s *PostgresStorage) UpdateMove(id int, newMove model.Move) error {
	ctx := context.Background()

	var capturedBytes []byte
	if newMove.Captured != nil {
		capturedBytes, _ = json.Marshal(newMove.Captured)
	}

	var movedPieceBytes []byte
	if newMove.MovedPiece != nil {
		movedPieceBytes, _ = json.Marshal(newMove.MovedPiece)
	}

	return s.queries.UpdateMove(ctx, sqlc.UpdateMoveParams{
		ID:         int64(id),
		GameID:     int64(newMove.GameID()),
		FromRow:    int32(newMove.FromRow),
		FromCol:    int32(newMove.FromCol),
		ToRow:      int32(newMove.ToRow),
		ToCol:      int32(newMove.ToCol),
		Captured:   capturedBytes,
		Promotion:  int32(newMove.Promotion),
		MovedPiece: movedPieceBytes,
		Check:      newMove.Check,
		Mate:       newMove.Mate,
	})
}

func (s *PostgresStorage) DeleteMove(id int) error {
	ctx := context.Background()
	return s.queries.DeleteMove(ctx, int64(id))
}

func (s *PostgresStorage) toModelMove(dbM sqlc.Move) model.Move {
	m := model.Move{
		FromRow:   int(dbM.FromRow),
		FromCol:   int(dbM.FromCol),
		ToRow:     int(dbM.ToRow),
		ToCol:     int(dbM.ToCol),
		Promotion: model.PieceType(dbM.Promotion),
		Check:     dbM.Check,
		Mate:      dbM.Mate,
	}
	m.SetID(int(dbM.ID))
	m.SetGameID(int(dbM.GameID))

	if len(dbM.Captured) > 0 {
		m.Captured = &model.Piece{}
		_ = json.Unmarshal(dbM.Captured, m.Captured)
	}

	if len(dbM.MovedPiece) > 0 {
		m.MovedPiece = &model.Piece{}
		_ = json.Unmarshal(dbM.MovedPiece, m.MovedPiece)
	}

	return m
}
