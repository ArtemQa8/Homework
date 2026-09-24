package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"mod.go/internal/model"
	"mod.go/internal/repository"
)

func RunLogger(ctx context.Context, storage *repository.Storage, simActive *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	storage.Lock()
	gamesBefore := len(storage.Games)
	playersBefore := len(storage.Players)
	movesBefore := len(storage.Moves)
	storage.Unlock()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if simActive != nil && *simActive {
				continue
			}

			storage.Lock()
			gamesAfter := len(storage.Games)
			playersAfter := len(storage.Players)
			movesAfter := len(storage.Moves)

			if gamesAfter > gamesBefore {
				for i := gamesBefore; i < gamesAfter; i++ {
					fmt.Printf("[ЛОГ] Новая игра: %s - %s\n",
						storage.Games[i].Player1().Name(),
						storage.Games[i].Player2().Name())
				}
			}
			if playersAfter > playersBefore {
				for i := playersBefore; i < playersAfter; i++ {
					fmt.Printf("[ЛОГ] Новый игрок: %s\n",
						storage.Players[i].Name())
				}
			}
			if movesAfter > movesBefore {
				for i := movesBefore; i < movesAfter; i++ {
					fmt.Printf("[ЛОГ] Новый ход: %s\n",
						model.FormatMove(storage.Moves[i]))
				}
			}

			gamesBefore = gamesAfter
			playersBefore = playersAfter
			movesBefore = movesAfter

			storage.Unlock()
		}
	}
}
