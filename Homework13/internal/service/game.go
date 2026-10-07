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

	if rows <= 0 {
		rows = 8
	}
	if cols <= 0 {
		cols = 8
	}

	if rows < 4 || cols < 4 {
		return model.Game{}, fmt.Errorf("минимальный размер доски 4x4, получено %dx%d", rows, cols)
	}

	game := model.NewGame(p1, p2, rows, cols)
	created, err := s.storage.CreateGame(*game)
	if err != nil {
		return model.Game{}, err
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
	if !move.IsPromotion(game.Board()) {
		move.Promotion = model.Pawn
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
	if move.IsPromotion(game.Board()) {
		move.Promotion = model.Queen
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
