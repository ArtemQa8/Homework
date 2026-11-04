package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"mod.go/internal/model"
	"mod.go/internal/repository/postgres/sqlc"
)

// CreatePlayer создает игрока
func (s *PostgresStorage) CreatePlayer(player model.Player) model.Player {
	ctx := context.Background()

	id, err := s.queries.GetNextSequence(ctx, "player_id")
	if err != nil {
		panic(fmt.Sprintf("ошибка получения ID игрока: %v", err))
	}

	dbPlayer, err := s.queries.CreatePlayer(ctx, sqlc.CreatePlayerParams{
		ID:     id,
		Name:   player.Name(),
		Rating: 0,
	})
	if err != nil {
		panic(fmt.Sprintf("ошибка создания игрока: %v", err))
	}

	player.SetID(int(dbPlayer.ID))
	return player
}

// GetAllPlayers возвращает всех игроков.
func (s *PostgresStorage) GetAllPlayers() []model.Player {
	ctx := context.Background()
	dbPlayers, err := s.queries.ListPlayers(ctx)
	if err != nil {
		return []model.Player{}
	}

	players := make([]model.Player, 0, len(dbPlayers))
	for _, dbP := range dbPlayers {
		p := model.NewPlayer(dbP.Name)
		p.SetID(int(dbP.ID))
		players = append(players, *p)
	}
	return players
}

// GetPlayerByID возвращает игрока по ID
func (s *PostgresStorage) GetPlayerByID(id int) (model.Player, bool) {
	ctx := context.Background()
	dbPlayer, err := s.queries.GetPlayer(ctx, int64(id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Player{}, false
		}
		return model.Player{}, false
	}

	p := model.NewPlayer(dbPlayer.Name)
	p.SetID(int(dbPlayer.ID))
	return *p, true
}

// UpdatePlayer обновляет имя игрока
func (s *PostgresStorage) UpdatePlayer(id int, newPlayer model.Player) error {
	ctx := context.Background()
	return s.queries.UpdatePlayer(ctx, sqlc.UpdatePlayerParams{
		ID:   int64(id),
		Name: newPlayer.Name(),
	})
}

// DeletePlayer удаляет игрока.
func (s *PostgresStorage) DeletePlayer(id int) error {
	ctx := context.Background()
	return s.queries.DeletePlayer(ctx, int64(id))
}
