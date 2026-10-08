package service

import (
	"fmt"

	"mod.go/internal/model"
)

type MoveService struct {
	storage Storage
}

func NewMoveService(storage Storage) *MoveService {
	return &MoveService{storage: storage}
}

// Create сохраняет ход в хранилище.
func (s *MoveService) Create(move model.Move) (model.Move, error) {
	return s.storage.CreateMove(move)
}

// Get возращает ход по ID
func (s *MoveService) Get(id int) (model.Move, error) {
	move, found := s.storage.GetMoveByID(id)
	if !found {
		return model.Move{}, fmt.Errorf("Ход с ID %d не найден", id)
	}
	return move, nil
}

// List возвращает все ходы.
func (s *MoveService) List() []model.Move {
	return s.storage.GetAllMoves()
}

// Update обновляет ход по ID.
func (s *MoveService) Update(id int, move model.Move) (model.Move, error) {
	if err := s.storage.UpdateMove(id, move); err != nil {
		return model.Move{}, err
	}
	updated, _ := s.storage.GetMoveByID(id)
	return updated, nil
}

// Delete удаляет ход по ID
func (s *MoveService) Delete(id int) error {
	return s.storage.DeleteMove(id)
}
