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

	nextGameID   int
	nextPlayerID int
	nextMoveID   int

	fileGames   string
	filePlayers string
	fileMoves   string
}

func НовоеХранилище() *Хранилище {
	return &Хранилище{
		fileGames:   "games.json",
		filePlayers: "players.json",
		fileMoves:   "moves.json",

		nextGameID:   1,
		nextPlayerID: 1,
		nextMoveID:   1,
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

		for _, игрок := range х.Игроки {
			if игрок.ID() >= х.nextPlayerID {
				х.nextPlayerID = игрок.ID() + 1
			}
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

		for _, игра := range х.Игры {
			if игра.ID() >= х.nextGameID {
				х.nextGameID = игра.ID() + 1
			}
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

		for _, ход := range х.Ходы {
			if ход.ID() >= х.nextMoveID {
				х.nextMoveID = ход.ID() + 1
			}
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

// =============ИГРОКИ=============
// вызывать только под МЬЮТЕКСОМ!!!
func (х *Хранилище) findPlayerByID(id int) (model.Игрок, bool) {
	for _, игрок := range х.Игроки {
		if игрок.ID() == id {
			return игрок, true
		}
	}
	return model.Игрок{}, false
}

// вызывать только под МЬЮТЕКСОМ!!!
func (х *Хранилище) findPlayerByName(имя string) (model.Игрок, int, bool) {
	найденные := []model.Игрок{}
	for _, игрок := range х.Игроки {
		if игрок.Имя() == имя {
			найденные = append(найденные, игрок)
		}
	}
	switch len(найденные) {
	case 0:
		return model.Игрок{}, 0, false
	case 1:
		return найденные[0], 1, true
	default:
		return model.Игрок{}, len(найденные), false
	}
}

// вызывать только под МЬЮТЕКСОМ!!!
func (х *Хранилище) addPlayer(игрок model.Игрок) model.Игрок {
	игрок.УстановитьID(х.nextPlayerID)
	х.nextPlayerID++
	х.Игроки = append(х.Игроки, игрок)
	_ = х.сохранитьИгроков()
	return игрок
}

// вызывать только под МЬЮТЕКСОМ!!!
func (х *Хранилище) resolvePlayer(игрок model.Игрок) (model.Игрок, error) {
	if игрок.ID() != 0 {
		найденный, ok := х.findPlayerByID(игрок.ID())
		if !ok {
			return model.Игрок{}, fmt.Errorf("Игрок с ID %d не найден", игрок.ID())
		}
		return найденный, nil
	}

	if игрок.Имя() == "" {
		return model.Игрок{}, fmt.Errorf("Имя игрока обязательно")
	}

	найденный, количество, found := х.findPlayerByName(игрок.Имя())
	if found {
		return найденный, nil
	}
	if количество > 1 {
		return model.Игрок{}, fmt.Errorf("Несколько игроков с именем %q, уточните ID", игрок.Имя())
	}
	return х.addPlayer(игрок), nil
}

func (х *Хранилище) СоздатьИгрока(игрок model.Игрок) model.Игрок {
	х.mu.Lock()
	defer х.mu.Unlock()
	return х.addPlayer(игрок)
}

func (х *Хранилище) ПолучитьВсехИгроков() []model.Игрок {
	х.mu.Lock()
	defer х.mu.Unlock()

	копия := make([]model.Игрок, len(х.Игроки))
	copy(копия, х.Игроки)
	return копия
}

func (х *Хранилище) ПолучитьИгрокаПоАйди(id int) (model.Игрок, bool) {
	х.mu.Lock()
	defer х.mu.Unlock()
	return х.findPlayerByID(id)
}

func (х *Хранилище) ОбновитьИгрока(id int, новыйИгрок model.Игрок) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	for i := range х.Игроки {
		if х.Игроки[i].ID() == id {
			новыйИгрок.УстановитьID(id)
			х.Игроки[i] = новыйИгрок
			return х.сохранитьИгроков()
		}
	}
	return fmt.Errorf("Игрок с ID %d не найден", id)
}

func (х *Хранилище) УдалитьИгрока(id int) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	for i := range х.Игроки {
		if х.Игроки[i].ID() == id {
			copy(х.Игроки[i:], х.Игроки[i+1:])
			х.Игроки[len(х.Игроки)-1] = model.Игрок{}
			х.Игроки = х.Игроки[:len(х.Игроки)-1]
			return х.сохранитьИгроков()
		}
	}
	return fmt.Errorf("Игрок с ID %d не найден", id)
}

// =============ИГРЫ=============
// Вызывать только под МЬЮТЕКСОМ!!!
func (х *Хранилище) findGameByID(id int) (model.Игра, bool) {
	for _, игра := range х.Игры {
		if игра.ID() == id {
			return игра, true
		}
	}
	return model.Игра{}, false
}

func (х *Хранилище) СоздатьИгру(игра model.Игра) (model.Игра, error) {
	х.mu.Lock()
	defer х.mu.Unlock()

	игрок1, err := х.resolvePlayer(игра.Игрок1())
	if err != nil {
		return model.Игра{}, err
	}
	игрок2, err := х.resolvePlayer(игра.Игрок2())
	if err != nil {
		return model.Игра{}, err
	}
	игра.УстановитьИгрок1(игрок1)
	игра.УстановитьИгрок2(игрок2)

	игра.УстановитьID(х.nextGameID)
	х.nextGameID++

	х.Игры = append(х.Игры, игра)
	_ = х.сохранитьИгры()
	return игра, nil
}

func (х *Хранилище) ПолучитьВсеИгры() []model.Игра {
	х.mu.Lock()
	defer х.mu.Unlock()

	копия := make([]model.Игра, len(х.Игры))
	copy(копия, х.Игры)
	return копия
}

func (х *Хранилище) ПолучитьИгруПоАйди(id int) (model.Игра, bool) {
	х.mu.Lock()
	defer х.mu.Unlock()

	for _, игра := range х.Игры {
		if игра.ID() == id {
			return игра, true
		}
	}
	return model.Игра{}, false
}

func (х *Хранилище) ОбновитьИгру(id int, новаяИгра model.Игра) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	игрок1, err := х.resolvePlayer(новаяИгра.Игрок1())
	if err != nil {
		return err
	}
	игрок2, err := х.resolvePlayer(новаяИгра.Игрок2())
	if err != nil {
		return err
	}
	новаяИгра.УстановитьИгрок1(игрок1)
	новаяИгра.УстановитьИгрок2(игрок2)

	for i := range х.Игры {
		if х.Игры[i].ID() == id {
			новаяИгра.УстановитьID(id)
			х.Игры[i] = новаяИгра
			return х.сохранитьИгры()
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

func (х *Хранилище) УдалитьИгру(id int) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	for i := range х.Игры {
		if х.Игры[i].ID() == id {
			copy(х.Игры[i:], х.Игры[i+1:])
			х.Игры[len(х.Игры)-1] = model.Игра{}
			х.Игры = х.Игры[:len(х.Игры)-1]
			return х.сохранитьИгры()
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

func (х *Хранилище) ПерезаписатьИгру(id int, новаяИгра model.Игра) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	for i := range х.Игры {
		if х.Игры[i].ID() == id {
			новаяИгра.УстановитьID(id)
			х.Игры[i] = новаяИгра
			return х.сохранитьИгры()
		}
	}
	return fmt.Errorf("Игра с ID %d не найдена", id)
}

// =============ХОДЫ=============
func (х *Хранилище) СоздатьХод(ход model.Ход) (model.Ход, error) {
	х.mu.Lock()
	defer х.mu.Unlock()

	if ход.ИграID() == 0 {
		return model.Ход{}, fmt.Errorf("играID обязателен для хода")
	}
	if _, ok := х.findGameByID(ход.ИграID()); !ok {
		return model.Ход{}, fmt.Errorf("Игра с ID %d не найдена", ход.ИграID())
	}

	ход.УстановитьID(х.nextMoveID)
	х.nextMoveID++

	х.Ходы = append(х.Ходы, ход)
	_ = х.сохранитьХоды()
	return ход, nil
}

func (х *Хранилище) ПолучитьВсеХоды() []model.Ход {
	х.mu.Lock()
	defer х.mu.Unlock()

	копия := make([]model.Ход, len(х.Ходы))
	copy(копия, х.Ходы)
	return копия
}

func (х *Хранилище) ПолучитьХодПоАйди(id int) (model.Ход, bool) {
	х.mu.Lock()
	defer х.mu.Unlock()

	for _, ход := range х.Ходы {
		if ход.ID() == id {
			return ход, true
		}
	}
	return model.Ход{}, false
}

func (х *Хранилище) ОбновитьХод(id int, новыйХод model.Ход) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	if новыйХод.ИграID() == 0 {
		return fmt.Errorf("играID обязателен для хода")
	}
	if _, ok := х.findGameByID(новыйХод.ИграID()); !ok {
		return fmt.Errorf("Игра с ID %d не найдена", новыйХод.ИграID())
	}

	for i := range х.Ходы {
		if х.Ходы[i].ID() == id {
			новыйХод.УстановитьID(id) // сохраняем прежний ID
			х.Ходы[i] = новыйХод
			return х.сохранитьХоды()
		}
	}
	return fmt.Errorf("Ход с ID %d не найден", id)
}

func (х *Хранилище) УдалитьХод(id int) error {
	х.mu.Lock()
	defer х.mu.Unlock()

	for i := range х.Ходы {
		if х.Ходы[i].ID() == id {
			copy(х.Ходы[i:], х.Ходы[i+1:])
			х.Ходы[len(х.Ходы)-1] = model.Ход{}
			х.Ходы = х.Ходы[:len(х.Ходы)-1]
			return х.сохранитьХоды()
		}
	}
	return fmt.Errorf("Ход с ID %d не найден", id)
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
