package service

import (
	"mod.go/internal/model"
	"mod.go/internal/repository"
)

type Storage interface {
	// Игроки
	CreatePlayer(player model.Player) model.Player
	GetAllPlayers() []model.Player
	GetPlayerByID(id int) (model.Player, bool)
	UpdatePlayer(id int, newPlayer model.Player) error
	DeletePlayer(id int) error

	// Игры
	CreateGame(game model.Game) (model.Game, error)
	GetAllGames() []model.Game
	GetGameByID(id int) (model.Game, bool)
	UpdateGame(id int, newGame model.Game) error
	OverwriteGame(id int, newGame model.Game) error
	DeleteGame(id int) error

	// Ходы
	CreateMove(move model.Move) (model.Move, error)
	GetAllMoves() []model.Move
	GetMoveByID(id int) (model.Move, bool)
	UpdateMove(id int, newMove model.Move) error
	DeleteMove(id int) error
}

var _ Storage = (*repository.Storage)(nil)
