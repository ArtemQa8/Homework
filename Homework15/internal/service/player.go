package service

import (
	"errors"
	"fmt"

	"mod.go/internal/model"
)

type PlayerService struct {
	storage Storage
}

func NewPlayerService(storage Storage) *PlayerService {
	return &PlayerService{storage: storage}
}

// Create создаёт нового игрока
func (s *PlayerService) Create(name string) (model.Player, error) {
	if name == "" {
		return model.Player{}, errors.New("Имя обязательно")
	}
	player := model.NewPlayer(name)
	return s.storage.CreatePlayer(*player), nil
}

// Get возвращает игрока по ID.
func (s *PlayerService) Get(id int) (model.Player, error) {
	player, found := s.storage.GetPlayerByID(id)
	if !found {
		return model.Player{}, fmt.Errorf("Игрок с ID %d не найден", id)
	}
	return player, nil
}

// List возвращает всех игроков.
func (s *PlayerService) List() []model.Player {
	return s.storage.GetAllPlayers()
}

// Update обновляет имя игрока.
func (s *PlayerService) Update(id int, name string) (model.Player, error) {
	if name == "" {
		return model.Player{}, errors.New("Имя обязательно")
	}
	if err := s.storage.UpdatePlayer(id, *model.NewPlayer(name)); err != nil {
		return model.Player{}, err
	}
	updated, _ := s.storage.GetPlayerByID(id)
	return updated, nil
}

// Delete удаляет игрока по ID.
func (s *PlayerService) Delete(id int) error {
	return s.storage.DeletePlayer(id)
}
