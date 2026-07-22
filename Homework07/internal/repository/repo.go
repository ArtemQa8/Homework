package repository

import (
	"sync"

	"mod.go/internal/model"
)

type Хранилище struct {
	mu     sync.Mutex
	Игры   []model.Игра
	Игроки []model.Игрок
	Ходы   []model.Ход
}

func (х *Хранилище) Добавить(объект model.ОбъектХранилища) {
	х.mu.Lock()         // Закрываем замок (одна горутина забирает с собой ключ и не пускает другие)
	defer х.mu.Unlock() // Открываем замок, но в самом-самом конце (перед последней })
	switch v := объект.(type) {
	case *model.Игра:
		х.Игры = append(х.Игры, *v)
	case model.Игрок:
		х.Игроки = append(х.Игроки, v)
	case model.Ход:
		х.Ходы = append(х.Ходы, v)
	}
}

func (х *Хранилище) Закрыть() { х.mu.Lock() }
func (х *Хранилище) Открыть() { х.mu.Unlock() }
