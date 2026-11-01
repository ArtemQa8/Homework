package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"mod.go/internal/model"
	"mod.go/internal/repository/postgres/sqlc"
)

// resolvePlayer ищет игрока по ID или имени, если нет - создает нового
func (s *PostgresStorage) resolvePlayer(ctx context.Context, player model.Player) (model.Player, error) {
	if player.ID() != 0 {
		p, found := s.GetPlayerByID(player.ID())
		if !found {
			return model.Player{}, fmt.Errorf("Игрок с ID %d не найден", player.ID())
		}
		return p, nil
	}

	if player.Name() == "" {
		return model.Player{}, fmt.Errorf("Имя игрока обязательно")
	}

	dbPlayers, err := s.queries.GetPlayerByName(ctx, player.Name())
	if err != nil {
		return model.Player{}, err
	}

	if len(dbPlayers) > 1 {
		return model.Player{}, fmt.Errorf("Несколько игроков с именем %q, уточните ID", player.Name())
	}

	if len(dbPlayers) == 1 {
		p := model.NewPlayer(dbPlayers[0].Name)
		p.SetID(int(dbPlayers[0].ID))
		return *p, nil
	}

	id, err := s.queries.GetNextSequence(ctx, "player_id")
	if err != nil {
		return model.Player{}, err
	}

	dbPlayer, err := s.queries.CreatePlayer(ctx, sqlc.CreatePlayerParams{
		ID:     id,
		Name:   player.Name(),
		Rating: 0,
	})
	if err != nil {
		return model.Player{}, err
	}

	p := model.NewPlayer(dbPlayer.Name)
	p.SetID(int(dbPlayer.ID))
	return *p, nil
}

// CreateGame создает игру, проверяя/создавая игроков
func (s *PostgresStorage) CreateGame(game model.Game) (model.Game, error) {
	ctx := context.Background()

	p1, err := s.resolvePlayer(ctx, game.Player1())
	if err != nil {
		return model.Game{}, err
	}
	p2, err := s.resolvePlayer(ctx, game.Player2())
	if err != nil {
		return model.Game{}, err
	}
	game.SetPlayer1(p1)
	game.SetPlayer2(p2)

	id, err := s.queries.GetNextSequence(ctx, "game_id")
	if err != nil {
		return model.Game{}, err
	}
	game.SetID(int(id))

	boardBytes, err := json.Marshal(game.Board())
	if err != nil {
		return model.Game{}, fmt.Errorf("ошибка кодирования доски: %w", err)
	}

	dbGame, err := s.queries.CreateGame(ctx, sqlc.CreateGameParams{
		ID:              id,
		Player1ID:       int64(p1.ID()),
		Player1Name:     p1.Name(),
		Player1Color:    int32(game.Player1Color()),
		Player2ID:       int64(p2.ID()),
		Player2Name:     p2.Name(),
		Player2Color:    int32(game.Player2Color()),
		Board:           boardBytes,
		CurrentColor:    int32(game.CurrentColor()),
		WhiteKingMoved:  false,
		BlackKingMoved:  false,
		WhiteRookAMoved: false,
		WhiteRookHMoved: false,
		BlackRookAMoved: false,
		BlackRookHMoved: false,
		LastMoveID:      0,
	})
	if err != nil {
		return model.Game{}, fmt.Errorf("ошибка сохранения игры: %w", err)
	}

	return s.toModelGame(dbGame, nil), nil
}

// GetAllGames возвращает все игры со всеми ходами
func (s *PostgresStorage) GetAllGames() []model.Game {
	ctx := context.Background()
	dbGames, err := s.queries.ListGames(ctx)
	if err != nil {
		return []model.Game{}
	}

	games := make([]model.Game, 0, len(dbGames))
	for _, dbG := range dbGames {
		dbMoves, _ := s.queries.ListMovesByGameID(ctx, dbG.ID)
		games = append(games, s.toModelGame(dbG, dbMoves))
	}
	return games
}

// GetGameByID возвращает игру со всеми ходами
func (s *PostgresStorage) GetGameByID(id int) (model.Game, bool) {
	ctx := context.Background()
	dbGame, err := s.queries.GetGame(ctx, int64(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Game{}, false
		}
		return model.Game{}, false
	}

	dbMoves, err := s.queries.ListMovesByGameID(ctx, int64(id))
	if err != nil {
		return model.Game{}, false
	}

	return s.toModelGame(dbGame, dbMoves), true
}

// UpdateGame обновляет игру (только имена игроков)
func (s *PostgresStorage) UpdateGame(id int, newGame model.Game) error {
	ctx := context.Background()

	p1, err := s.resolvePlayer(ctx, newGame.Player1())
	if err != nil {
		return err
	}
	p2, err := s.resolvePlayer(ctx, newGame.Player2())
	if err != nil {
		return err
	}

	return s.queries.UpdateGame(ctx, sqlc.UpdateGameParams{
		ID:          int64(id),
		Player1ID:   int64(p1.ID()),
		Player1Name: p1.Name(),
		Player2ID:   int64(p2.ID()),
		Player2Name: p2.Name(),
	})
}

// OverwriteGame перезаписывает ВСЁ состояние игры (доску, ход, флаги).
func (s *PostgresStorage) OverwriteGame(id int, newGame model.Game) error {
	ctx := context.Background()

	boardBytes, err := json.Marshal(newGame.Board())
	if err != nil {
		return fmt.Errorf("ошибка кодирования доски: %w", err)
	}

	var lastMoveID int64
	if moves := newGame.Moves(); len(moves) > 0 {
		lastMoveID = int64(moves[len(moves)-1].ID())
	}

	return s.queries.OverwriteGame(ctx, sqlc.OverwriteGameParams{
		ID:              int64(id),
		Player1ID:       int64(newGame.Player1().ID()),
		Player1Name:     newGame.Player1().Name(),
		Player1Color:    int32(newGame.Player1Color()),
		Player2ID:       int64(newGame.Player2().ID()),
		Player2Name:     newGame.Player2().Name(),
		Player2Color:    int32(newGame.Player2Color()),
		Board:           boardBytes,
		CurrentColor:    int32(newGame.CurrentColor()),
		WhiteKingMoved:  newGame.WhiteKingMoved(),
		BlackKingMoved:  newGame.BlackKingMoved(),
		WhiteRookAMoved: newGame.WhiteRookAMoved(),
		WhiteRookHMoved: newGame.WhiteRookHMoved(),
		BlackRookAMoved: newGame.BlackRookAMoved(),
		BlackRookHMoved: newGame.BlackRookHMoved(),
		LastMoveID:      lastMoveID,
	})
}

// DeleteGame удаляет игру по ID.
func (s *PostgresStorage) DeleteGame(id int) error {
	ctx := context.Background()
	return s.queries.DeleteGame(ctx, int64(id))
}

// toModelGame - конвертер. Превращает структуры из бд (sqlc.Game, sqlc.Move) в доменные модели (model.Game, model.Move).
func (s *PostgresStorage) toModelGame(dbGame sqlc.Game, dbMoves []sqlc.Move) model.Game {
	p1 := model.NewPlayer(dbGame.Player1Name)
	p1.SetID(int(dbGame.Player1ID))
	p2 := model.NewPlayer(dbGame.Player2Name)
	p2.SetID(int(dbGame.Player2ID))

	var board *model.Board
	if len(dbGame.Board) > 0 {
		board = &model.Board{}
		_ = json.Unmarshal(dbGame.Board, board)
	}

	game := model.NewGameWithBoard(*p1, *p2, board)
	game.SetID(int(dbGame.ID))
	game.SetPlayer1Color(model.Color(dbGame.Player1Color))
	game.SetPlayer2Color(model.Color(dbGame.Player2Color))
	game.SetCurrentColor(model.Color(dbGame.CurrentColor))
	game.SetWhiteKingMoved(dbGame.WhiteKingMoved)
	game.SetBlackKingMoved(dbGame.BlackKingMoved)
	game.SetWhiteRookAMoved(dbGame.WhiteRookAMoved)
	game.SetWhiteRookHMoved(dbGame.WhiteRookHMoved)
	game.SetBlackRookAMoved(dbGame.BlackRookAMoved)
	game.SetBlackRookHMoved(dbGame.BlackRookHMoved)

	moves := make([]model.Move, 0, len(dbMoves))
	for _, dbM := range dbMoves {
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
		moves = append(moves, m)
	}
	game.SetMoves(moves)

	return *game
}

func (s *PostgresStorage) MakeMove(ctx context.Context, gameID int, move *model.Move, game *model.Game) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	id, err := qtx.GetNextSequence(ctx, "move_id")
	if err != nil {
		return fmt.Errorf("ошибка получения ID хода: %w", err)
	}
	move.SetID(int(id))

	var capturedBytes []byte
	if move.Captured != nil {
		capturedBytes, _ = json.Marshal(move.Captured)
	}
	var movedPieceBytes []byte
	if move.MovedPiece != nil {
		movedPieceBytes, _ = json.Marshal(move.MovedPiece)
	}

	_, err = qtx.CreateMove(ctx, sqlc.CreateMoveParams{
		ID:         id,
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
	if err != nil {
		return fmt.Errorf("ошибка сохранения хода: %w", err)
	}

	game.SetLastMoveID(int(id))

	// Вычисляем lastMoveID для передачи в OverwriteGame
	var lastMoveID int64 = int64(id)

	boardBytes, err := json.Marshal(game.Board())
	if err != nil {
		return fmt.Errorf("ошибка кодирования доски: %w", err)
	}

	err = qtx.OverwriteGame(ctx, sqlc.OverwriteGameParams{
		ID:              int64(gameID),
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
		LastMoveID:      lastMoveID,
	})
	if err != nil {
		return fmt.Errorf("ошибка обновления игры: %w", err)
	}

	return tx.Commit(ctx)
}
