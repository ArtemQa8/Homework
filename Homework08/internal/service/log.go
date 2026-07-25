package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"mod.go/internal/model"
	"mod.go/internal/repository"
)

func ЗапуститьЛоггер(ctx context.Context, хранилище *repository.Хранилище, simActive *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	ожидание := time.NewTicker(200 * time.Millisecond)
	defer ожидание.Stop()

	хранилище.Закрыть()
	игрБыло := len(хранилище.Игры)
	игроковБыло := len(хранилище.Игроки)
	ходовБыло := len(хранилище.Ходы)
	хранилище.Открыть()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ожидание.C:
			if simActive != nil && *simActive {
				continue
			}

			хранилище.Закрыть()
			игрСтало := len(хранилище.Игры)
			игроковСтало := len(хранилище.Игроки)
			ходовСтало := len(хранилище.Ходы)

			if игрСтало > игрБыло {
				for i := игрБыло; i < игрСтало; i++ {
					fmt.Printf("[ЛОГ] Новая игра: %s - %s\n",
						хранилище.Игры[i].Игрок1().Имя(),
						хранилище.Игры[i].Игрок2().Имя())
				}
			}
			if игроковСтало > игроковБыло {
				for i := игроковБыло; i < игроковСтало; i++ {
					fmt.Printf("[ЛОГ] Новый игрок: %s\n",
						хранилище.Игроки[i].Имя())
				}
			}
			if ходовСтало > ходовБыло {
				for i := ходовБыло; i < ходовСтало; i++ {
					fmt.Printf("[ЛОГ] Новый ход: %s\n",
						model.ФорматироватьХод(хранилище.Ходы[i]))
				}
			}

			игрБыло = игрСтало
			игроковБыло = игроковСтало
			ходовБыло = ходовСтало

			хранилище.Открыть()
		}
	}
}
