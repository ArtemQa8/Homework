package service

import (
	"errors"
	"fmt"

	"mod.go/internal/model"
)

type GameService struct {
	storage Storage
}

func NewGameService(storage Storage) *GameService {
	return &GameService{storage: storage}
}

// Create создает новую игшру между двумя игроками
func (s *GameService) Create(p1, p2 model.Player, rows, cols int) (model.Game, error) {
	if p1.Name() == "" || p2.Name() == "" {
		return model.Game{}, errors.New("имена обоих игроков обязательны")
	}

	if rows <= 0 || cols <= 0 {
		rows, cols = 8, 8
	}

	game := model.Game{}
	game.SetPlayer1(p1)
	game.SetPlayer2(p2)
	game.SetBoard(model.NewBoard(rows, cols))

	created, err := s.storage.CreateGame(game)
	if err != nil {
		return model.Game{}, err
	}

	created.SetPlayer1Color(model.White)
	created.SetPlayer2Color(model.Black)

	if err := s.storage.OverwriteGame(created.ID(), created); err != nil {
		return model.Game{}, fmt.Errorf("не удалось сохранить игру: %w", err)
	}

	return created, nil
}

// Get возвращает игру по ID.
func (s *GameService) Get(id int) (model.Game, error) {
	game, found := s.storage.GetGameByID(id)
	if !found {
		return model.Game{}, fmt.Errorf("Игра с ID %d не найдена", id)
	}
	return game, nil
}

// List возвроащает все игры
func (s *GameService) List() []model.Game {
	return s.storage.GetAllGames()
}

// Update обновляет игроков игры
func (s *GameService) Update(id int, p1, p2 model.Player) (model.Game, error) {
	game, found := s.storage.GetGameByID(id)
	if !found {
		return model.Game{}, fmt.Errorf("игра с ID %d не найдена", id)
	}

	if p1.Name() != "" {
		game.SetPlayer1(p1)
	}
	if p2.Name() != "" {
		game.SetPlayer2(p2)
	}

	if err := s.storage.UpdateGame(id, game); err != nil {
		return model.Game{}, err
	}

	updated, _ := s.storage.GetGameByID(id)
	return updated, nil
}

// Delete удаляет игру по ID.
func (s *GameService) Delete(id int) error {
	return s.storage.DeleteGame(id)
}

// MakeMove выполняет ход в игре.
func (s *GameService) MakeMove(id int, move *model.Move) (model.Game, error) {
	game, found := s.storage.GetGameByID(id)
	if !found {
		return model.Game{}, fmt.Errorf("Игра с ID %d не найдена", id)
	}

	move.SetGameID(id)

	// Сбрасываем Promotion, если ход не является превращением.
	if piece := game.Board().PieceAt(move.FromRow, move.FromCol); piece != nil {
		isPromotion := false
		if piece.Type() == model.Pawn {
			lastRow := 0
			if piece.Color() == model.White {
				lastRow = game.Board().Rows() - 1
			}
			if move.ToRow == lastRow {
				isPromotion = true
			}
		}
		if !isPromotion {
			move.Promotion = model.Pawn
		}
	}

	if err := game.MakeMove(move); err != nil {
		return model.Game{}, err
	}

	savedMove, err := s.storage.CreateMove(*move)
	if err != nil {
		return model.Game{}, fmt.Errorf("Не удалось сохранить ход: %w", err)
	}
	game.SetLastMoveID(savedMove.ID())

	if err := s.storage.OverwriteGame(id, game); err != nil {
		return model.Game{}, fmt.Errorf("Не удалось сохранить игру: %w", err)
	}

	return game, nil
}

// AutoMove делает случайный ход за текущего игрока
func (s *GameService) AutoMove(id int) (model.Game, error) {
	game, found := s.storage.GetGameByID(id)
	if !found {
		return model.Game{}, fmt.Errorf("Игра с ID %d не найдена", id)
	}

	color := game.CurrentColor()
	if game.Checkmate(color) || game.Stalemate(color) {
		return game, nil
	}

	move, err := ChooseRandomMove(&game)
	if err != nil {
		return model.Game{}, err
	}

	// Автопревращение пешки в ферзя
	if piece := game.Board().PieceAt(move.FromRow, move.FromCol); piece != nil && piece.Type() == model.Pawn {
		lastRow := 0
		if piece.Color() == model.White {
			lastRow = game.Board().Rows() - 1
		}
		if move.ToRow == lastRow {
			move.Promotion = model.Queen
		}
	}

	move.SetGameID(id)
	if err := game.MakeMove(move); err != nil {
		return model.Game{}, err
	}

	savedMove, err := s.storage.CreateMove(*move)
	if err != nil {
		return model.Game{}, fmt.Errorf("Не удалось сохранить ход: %w", err)
	}
	game.SetLastMoveID(savedMove.ID())

	if err := s.storage.OverwriteGame(id, game); err != nil {
		return model.Game{}, fmt.Errorf("Не удалось сохранить игру: %w", err)
	}

	return game, nil
}
