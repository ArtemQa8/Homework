package service_test

import (
	"fmt"

	"mod.go/internal/model"
	"mod.go/internal/service"
)

type mockStorage struct {
	players []model.Player
	games   []model.Game
	moves   []model.Move

	nextPlayerID int
	nextGameID   int
	nextMoveID   int
}

var _ service.Storage = (*mockStorage)(nil)

func newMockStorage() *mockStorage {
	return &mockStorage{
		nextPlayerID: 1,
		nextGameID:   1,
		nextMoveID:   1,
	}
}

// ======== PLAYERS ========

func (m *mockStorage) CreatePlayer(player model.Player) model.Player {
	player.SetID(m.nextPlayerID)
	m.nextPlayerID++
	m.players = append(m.players, player)
	return player
}

func (m *mockStorage) GetAllPlayers() []model.Player {
	result := make([]model.Player, len(m.players))
	copy(result, m.players)
	return result
}

func (m *mockStorage) GetPlayerByID(id int) (model.Player, bool) {
	for _, p := range m.players {
		if p.ID() == id {
			return p, true
		}
	}
	return model.Player{}, false
}

func (m *mockStorage) UpdatePlayer(id int, newPlayer model.Player) error {
	for i := range m.players {
		if m.players[i].ID() == id {
			newPlayer.SetID(id)
			m.players[i] = newPlayer
			return nil
		}
	}
	return fmt.Errorf("Игрок с ID %d не найден", id)
}

func (m *mockStorage) DeletePlayer(id int) error {
	for i := range m.players {
		if m.players[i].ID() == id {
			m.players = append(m.players[:i], m.players[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Игрок с ID %d не найден", id)
}

// ======== GAMES ========

func (m *mockStorage) CreateGame(game model.Game) (model.Game, error) {
	game.SetID(m.nextGameID)
	m.nextGameID++
	m.games = append(m.games, game)
	return game, nil
}

func (m *mockStorage) GetAllGames() []model.Game {
	result := make([]model.Game, len(m.games))
	copy(result, m.games)
	return result
}

func (m *mockStorage) GetGameByID(id int) (model.Game, bool) {
	for _, g := range m.games {
		if g.ID() == id {
			return g, true
		}
	}
	return model.Game{}, false
}

func (m *mockStorage) UpdateGame(id int, newGame model.Game) error {
	for i := range m.games {
		if m.games[i].ID() == id {
			newGame.SetID(id)
			m.games[i] = newGame
			return nil
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

func (m *mockStorage) OverwriteGame(id int, newGame model.Game) error {
	for i := range m.games {
		if m.games[i].ID() == id {
			newGame.SetID(id)
			m.games[i] = newGame
			return nil
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

func (m *mockStorage) DeleteGame(id int) error {
	for i := range m.games {
		if m.games[i].ID() == id {
			m.games = append(m.games[:i], m.games[i+1:]...)

			cleaned := make([]model.Move, 0, len(m.moves))
			for _, mv := range m.moves {
				if mv.GameID() != id {
					cleaned = append(cleaned, mv)
				}
			}
			m.moves = cleaned
			return nil
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

// ======== MOVES ========

func (m *mockStorage) CreateMove(move model.Move) (model.Move, error) {
	if move.GameID() == 0 {
		return model.Move{}, fmt.Errorf("играID обязателен для хода")
	}
	if _, found := m.GetGameByID(move.GameID()); !found {
		return model.Move{}, fmt.Errorf("Игра с ID %d не найдена", move.GameID())
	}

	move.SetID(m.nextMoveID)
	m.nextMoveID++
	m.moves = append(m.moves, move)
	return move, nil
}

func (m *mockStorage) GetAllMoves() []model.Move {
	result := make([]model.Move, len(m.moves))
	copy(result, m.moves)
	return result
}

func (m *mockStorage) GetMoveByID(id int) (model.Move, bool) {
	for _, mv := range m.moves {
		if mv.ID() == id {
			return mv, true
		}
	}
	return model.Move{}, false
}

func (m *mockStorage) UpdateMove(id int, newMove model.Move) error {
	for i := range m.moves {
		if m.moves[i].ID() == id {
			newMove.SetID(id)
			m.moves[i] = newMove
			return nil
		}
	}
	return fmt.Errorf("Ход с ID %d не найден", id)
}

func (m *mockStorage) DeleteMove(id int) error {
	for i := range m.moves {
		if m.moves[i].ID() == id {
			m.moves = append(m.moves[:i], m.moves[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Ход с ID %d не найден", id)
}
