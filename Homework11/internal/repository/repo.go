package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"mod.go/internal/model"
)

type Storage struct {
	mu      sync.Mutex
	Games   []model.Game
	Players []model.Player
	Moves   []model.Move

	nextGameID   int
	nextPlayerID int
	nextMoveID   int

	fileGames   string
	filePlayers string
	fileMoves   string
}

func NewStorage() *Storage {
	return &Storage{
		fileGames:   "games.json",
		filePlayers: "players.json",
		fileMoves:   "moves.json",

		nextGameID:   1,
		nextPlayerID: 1,
		nextMoveID:   1,
	}
}

// LoadFromFiles наполняет слайсы данными из JSON-файлов.
func (s *Storage) LoadFromFiles() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Игроки
	data, err := os.ReadFile(s.filePlayers)
	if err == nil {
		if err := json.Unmarshal(data, &s.Players); err != nil {
			return fmt.Errorf("ошибка загрузки игроков: %w", err)
		}

		for _, player := range s.Players {
			if player.ID() >= s.nextPlayerID {
				s.nextPlayerID = player.ID() + 1
			}
		}

	} else if !os.IsNotExist(err) {
		return err
	}

	// Игры
	data, err = os.ReadFile(s.fileGames)
	if err == nil {
		if err := json.Unmarshal(data, &s.Games); err != nil {
			return fmt.Errorf("ошибка загрузки игр: %w", err)
		}

		for _, game := range s.Games {
			if game.ID() >= s.nextGameID {
				s.nextGameID = game.ID() + 1
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Ходы
	data, err = os.ReadFile(s.fileMoves)
	if err == nil {
		if err := json.Unmarshal(data, &s.Moves); err != nil {
			return fmt.Errorf("ошибка загрузки ходов: %w", err)
		}

		for _, move := range s.Moves {
			if move.ID() >= s.nextMoveID {
				s.nextMoveID = move.ID() + 1
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return nil
}

// savePlayers записывает слайс Players в файл players.json.
func (s *Storage) savePlayers() error {
	data, err := json.MarshalIndent(s.Players, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePlayers, data, 0644)
}

// saveGames записывает слайс Games в файл games.json.
func (s *Storage) saveGames() error {
	data, err := json.MarshalIndent(s.Games, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileGames, data, 0644)
}

// saveMoves записывает слайс Moves в файл moves.json.
func (s *Storage) saveMoves() error {
	data, err := json.MarshalIndent(s.Moves, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileMoves, data, 0644)
}

// SaveAll записывает все слайсы в соответствующие файлы.
func (s *Storage) SaveAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.saveGames(); err != nil {
		return err
	}
	if err := s.savePlayers(); err != nil {
		return err
	}
	if err := s.saveMoves(); err != nil {
		return err
	}
	return nil
}

// =============PLAYERS=============
// вызывать только под МЬЮТЕКСОМ!!!
func (s *Storage) findPlayerByID(id int) (model.Player, bool) {
	for _, player := range s.Players {
		if player.ID() == id {
			return player, true
		}
	}
	return model.Player{}, false
}

// вызывать только под МЬЮТЕКСОМ!!!
func (s *Storage) findPlayerByName(name string) (model.Player, int, bool) {
	found := []model.Player{}
	for _, player := range s.Players {
		if player.Name() == name {
			found = append(found, player)
		}
	}
	switch len(found) {
	case 0:
		return model.Player{}, 0, false
	case 1:
		return found[0], 1, true
	default:
		return model.Player{}, len(found), false
	}
}

// вызывать только под МЬЮТЕКСОМ!!!
func (s *Storage) addPlayer(player model.Player) model.Player {
	player.SetID(s.nextPlayerID)
	s.nextPlayerID++
	s.Players = append(s.Players, player)
	_ = s.savePlayers()
	return player
}

// вызывать только под МЬЮТЕКСОМ!!!
func (s *Storage) resolvePlayer(player model.Player) (model.Player, error) {
	if player.ID() != 0 {
		found, ok := s.findPlayerByID(player.ID())
		if !ok {
			return model.Player{}, fmt.Errorf("Игрок с ID %d не найден", player.ID())
		}
		return found, nil
	}

	if player.Name() == "" {
		return model.Player{}, fmt.Errorf("Имя игрока обязательно")
	}

	found, count, ok := s.findPlayerByName(player.Name())
	if ok {
		return found, nil
	}
	if count > 1 {
		return model.Player{}, fmt.Errorf("Несколько игроков с именем %q, уточните ID", player.Name())
	}
	return s.addPlayer(player), nil
}

func (s *Storage) CreatePlayer(player model.Player) model.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addPlayer(player)
}

func (s *Storage) GetAllPlayers() []model.Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]model.Player, len(s.Players))
	copy(result, s.Players)
	return result
}

func (s *Storage) GetPlayerByID(id int) (model.Player, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.findPlayerByID(id)
}

func (s *Storage) UpdatePlayer(id int, newPlayer model.Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Players {
		if s.Players[i].ID() == id {
			newPlayer.SetID(id)
			s.Players[i] = newPlayer
			return s.savePlayers()
		}
	}
	return fmt.Errorf("Игрок с ID %d не найден", id)
}

func (s *Storage) DeletePlayer(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Players {
		if s.Players[i].ID() == id {
			copy(s.Players[i:], s.Players[i+1:])
			s.Players[len(s.Players)-1] = model.Player{}
			s.Players = s.Players[:len(s.Players)-1]
			return s.savePlayers()
		}
	}
	return fmt.Errorf("Игрок с ID %d не найден", id)
}

// =============GAMES=============

// вызывать только под МЬЮТЕКСОМ!!!
func (s *Storage) findGameByID(id int) (model.Game, bool) {
	for _, game := range s.Games {
		if game.ID() == id {
			return game, true
		}
	}
	return model.Game{}, false
}

func (s *Storage) CreateGame(game model.Game) (model.Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	player1, err := s.resolvePlayer(game.Player1())
	if err != nil {
		return model.Game{}, err
	}
	player2, err := s.resolvePlayer(game.Player2())
	if err != nil {
		return model.Game{}, err
	}
	game.SetPlayer1(player1)
	game.SetPlayer2(player2)

	game.SetID(s.nextGameID)
	s.nextGameID++

	s.Games = append(s.Games, game)
	_ = s.saveGames()
	return game, nil
}

func (s *Storage) GetAllGames() []model.Game {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]model.Game, len(s.Games))
	copy(result, s.Games)
	return result
}

func (s *Storage) GetGameByID(id int) (model.Game, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, game := range s.Games {
		if game.ID() == id {
			return game, true
		}
	}
	return model.Game{}, false
}

func (s *Storage) UpdateGame(id int, newGame model.Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player1, err := s.resolvePlayer(newGame.Player1())
	if err != nil {
		return err
	}
	player2, err := s.resolvePlayer(newGame.Player2())
	if err != nil {
		return err
	}
	newGame.SetPlayer1(player1)
	newGame.SetPlayer2(player2)

	for i := range s.Games {
		if s.Games[i].ID() == id {
			newGame.SetID(id)
			s.Games[i] = newGame
			return s.saveGames()
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

// OverwriteGame заменяет игру по ID без повторного разрешения игроков.
func (s *Storage) OverwriteGame(id int, newGame model.Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Games {
		if s.Games[i].ID() == id {
			newGame.SetID(id)
			s.Games[i] = newGame
			return s.saveGames()
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

func (s *Storage) DeleteGame(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Games {
		if s.Games[i].ID() == id {
			copy(s.Games[i:], s.Games[i+1:])
			s.Games[len(s.Games)-1] = model.Game{}
			s.Games = s.Games[:len(s.Games)-1]
			return s.saveGames()
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

// =============MOVES=============
func (s *Storage) CreateMove(move model.Move) (model.Move, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if move.GameID() == 0 {
		return model.Move{}, fmt.Errorf("играID обязателен для хода")
	}
	if _, ok := s.findGameByID(move.GameID()); !ok {
		return model.Move{}, fmt.Errorf("Игра с ID %d не найдена", move.GameID())
	}

	move.SetID(s.nextMoveID)
	s.nextMoveID++

	s.Moves = append(s.Moves, move)
	_ = s.saveMoves()
	return move, nil
}

func (s *Storage) GetAllMoves() []model.Move {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]model.Move, len(s.Moves))
	copy(result, s.Moves)
	return result
}

func (s *Storage) GetMoveByID(id int) (model.Move, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, move := range s.Moves {
		if move.ID() == id {
			return move, true
		}
	}
	return model.Move{}, false
}

func (s *Storage) UpdateMove(id int, newMove model.Move) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if newMove.GameID() == 0 {
		return fmt.Errorf("играID обязателен для хода")
	}
	if _, ok := s.findGameByID(newMove.GameID()); !ok {
		return fmt.Errorf("Игра с ID %d не найдена", newMove.GameID())
	}

	for i := range s.Moves {
		if s.Moves[i].ID() == id {
			newMove.SetID(id)
			s.Moves[i] = newMove
			return s.saveMoves()
		}
	}
	return fmt.Errorf("Ход с ID %d не найден", id)
}

func (s *Storage) DeleteMove(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Moves {
		if s.Moves[i].ID() == id {
			copy(s.Moves[i:], s.Moves[i+1:])
			s.Moves[len(s.Moves)-1] = model.Move{}
			s.Moves = s.Moves[:len(s.Moves)-1]
			return s.saveMoves()
		}
	}
	return fmt.Errorf("Ход с ID %d не найден", id)
}

func (s *Storage) Add(object model.StorageObject) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch v := object.(type) {
	case *model.Game:
		s.Games = append(s.Games, *v)
		_ = s.saveGames()
	case model.Player:
		s.Players = append(s.Players, v)
		_ = s.savePlayers()
	case model.Move:
		s.Moves = append(s.Moves, v)
		_ = s.saveMoves()
	}
}

func (s *Storage) Lock()   { s.mu.Lock() }
func (s *Storage) Unlock() { s.mu.Unlock() }
