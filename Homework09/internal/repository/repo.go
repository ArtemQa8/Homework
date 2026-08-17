package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"mod.go/internal/model"
)

type Хранилище struct {
	mu     sync.Mutex
	Игры   []model.Игра
	Игроки []model.Игрок
	Ходы   []model.Ход

	fileGames   string
	filePlayers string
	fileMoves   string
}

func НовоеХранилище() *Хранилище {
	return &Хранилище{
		fileGames:   "games.json",
		filePlayers: "players.json",
		fileMoves:   "moves.json",
	}
}

// ЗагрузитьИзФайлов наполняет слайсы данными из JSON-файлов.
func (х *Хранилище) ЗагрузитьИзФайлов() error {
	х.mu.Lock()
	defer х.mu.Unlock()

	// Игроки
	данные, err := os.ReadFile(х.filePlayers)
	if err == nil {
		if err := json.Unmarshal(данные, &х.Игроки); err != nil {
			return fmt.Errorf("ошибка загрузки игроков: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Игры
	данные, err = os.ReadFile(х.fileGames)
	if err == nil {
		if err := json.Unmarshal(данные, &х.Игры); err != nil {
			return fmt.Errorf("ошибка загрузки игр: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Ходы
	данные, err = os.ReadFile(х.fileMoves)
	if err == nil {
		if err := json.Unmarshal(данные, &х.Ходы); err != nil {
			return fmt.Errorf("ошибка загрузки ходов: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return nil
}

// сохранитьИгроков записывает слайс Игроки в файл players.json.
func (х *Хранилище) сохранитьИгроков() error {
	данные, err := json.MarshalIndent(х.Игроки, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(х.filePlayers, данные, 0644)
}

// сохранитьИгры записывает слайс Игры в файл games.json.
func (х *Хранилище) сохранитьИгры() error {
	данные, err := json.MarshalIndent(х.Игры, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(х.fileGames, данные, 0644)
}

// сохранитьХоды записывает слайс Ходы в файл moves.json.
func (х *Хранилище) сохранитьХоды() error {
	данные, err := json.MarshalIndent(х.Ходы, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(х.fileMoves, данные, 0644)
}

// СохранитьВсе записывает все слайсы в соответствующие файлы.
func (х *Хранилище) СохранитьВсе() error {
	х.mu.Lock()
	defer х.mu.Unlock()

	if err := х.сохранитьИгры(); err != nil {
		return err
	}
	if err := х.сохранитьИгроков(); err != nil {
		return err
	}
	if err := х.сохранитьХоды(); err != nil {
		return err
	}
	return nil
}

func (х *Хранилище) Добавить(объект model.ОбъектХранилища) {
	х.mu.Lock()
	defer х.mu.Unlock()
	switch v := объект.(type) {
	case *model.Игра:
		х.Игры = append(х.Игры, *v)
		_ = х.сохранитьИгры()
	case model.Игрок:
		х.Игроки = append(х.Игроки, v)
		_ = х.сохранитьИгроков()
	case model.Ход:
		х.Ходы = append(х.Ходы, v)
		_ = х.сохранитьХоды()
	}
}

func (х *Хранилище) Закрыть() { х.mu.Lock() }
func (х *Хранилище) Открыть() { х.mu.Unlock() }
