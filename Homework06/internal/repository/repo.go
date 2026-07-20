package repository

import "mod.go/internal/model"

type Хранилище struct {
	Игры   []model.Игра
	Игроки []model.Игрок
	Ходы   []model.Ход
}

func (х *Хранилище) Добавить(объект model.ОбъектХранилища) {
	switch v := объект.(type) {
	case *model.Игра:
		х.Игры = append(х.Игры, *v)
	case model.Игрок:
		х.Игроки = append(х.Игроки, v)
	case model.Ход:
		х.Ходы = append(х.Ходы, v)
	}
}
