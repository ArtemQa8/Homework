package model

import (
	"encoding/json"
	"fmt"
)

type ЦветФигуры int

const (
	Белые ЦветФигуры = iota
	Чёрные
)

type ТипФигуры int

const (
	Пешка ТипФигуры = iota
	Ладья
	Конь
	Слон
	Ферзь
	Король
)

const (
	Reset      = "\033[0m"
	WhitePiece = "\033[31m"
	BlackPiece = "\033[34m"
	WhiteSqBg  = "\033[48;5;250m"
	BlackSqBg  = "\033[48;5;240m"
)

type Фигура struct {
	цвет ЦветФигуры
	тип  ТипФигуры
}

func НоваяФигура(новыйЦвет ЦветФигуры, новыйТип ТипФигуры) *Фигура {
	return &Фигура{цвет: новыйЦвет, тип: новыйТип}
}

func (ф *Фигура) Цвет() ЦветФигуры { return ф.цвет }
func (ф *Фигура) Тип() ТипФигуры   { return ф.тип }

func (ф *Фигура) ОтобразитьФигуру(фон string) string {
	var цветКод string
	if ф.цвет == Белые {
		цветКод = WhitePiece
	} else {
		цветКод = BlackPiece
	}

	var символ string
	switch ф.тип {
	case Пешка:
		if ф.цвет == Белые {
			символ = "♙"
		} else {
			символ = "♟"
		}
	case Ладья:
		if ф.цвет == Белые {
			символ = "♖"
		} else {
			символ = "♜"
		}
	case Конь:
		if ф.цвет == Белые {
			символ = "♘"
		} else {
			символ = "♞"
		}
	case Слон:
		if ф.цвет == Белые {
			символ = "♗"
		} else {
			символ = "♝"
		}
	case Ферзь:
		if ф.цвет == Белые {
			символ = "♕"
		} else {
			символ = "♛"
		}
	case Король:
		if ф.цвет == Белые {
			символ = "♔"
		} else {
			символ = "♚"
		}
	default:
		символ = " "
	}
	return фон + цветКод + " " + символ + " " + Reset
}

func (ф *Фигура) Символ() string {
	switch ф.тип {
	case Пешка:
		if ф.цвет == Белые {
			return "♙"
		} else {
			return "♟"
		}
	case Ладья:
		if ф.цвет == Белые {
			return "♖"
		} else {
			return "♜"
		}
	case Конь:
		if ф.цвет == Белые {
			return "♘"
		} else {
			return "♞"
		}
	case Слон:
		if ф.цвет == Белые {
			return "♗"
		} else {
			return "♝"
		}
	case Ферзь:
		if ф.цвет == Белые {
			return "♕"
		} else {
			return "♛"
		}
	case Король:
		if ф.цвет == Белые {
			return "♔"
		} else {
			return "♚"
		}
	}
	return " "
}

// MarshalJSON реализует сериализацию в JSON с сохранением приватных полей.
func (ф Фигура) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Цвет ЦветФигуры `json:"цвет"`
		Тип  ТипФигуры  `json:"тип"`
	}{
		Цвет: ф.цвет,
		Тип:  ф.тип,
	})
}

// UnmarshalJSON реализует десериализацию из JSON с заполнением приватных полей.
func (ф *Фигура) UnmarshalJSON(data []byte) error {
	var данные struct {
		Цвет ЦветФигуры `json:"цвет"`
		Тип  ТипФигуры  `json:"тип"`
	}
	if err := json.Unmarshal(data, &данные); err != nil {
		return err
	}
	ф.цвет = данные.Цвет
	ф.тип = данные.Тип
	return nil
}

// MarshalJSON превращает ЦветФигуры в строку.
func (ц ЦветФигуры) MarshalJSON() ([]byte, error) {
	switch ц {
	case Белые:
		return []byte(`"Белые"`), nil
	case Чёрные:
		return []byte(`"Чёрные"`), nil
	}
	return []byte(`"Неизвестно"`), nil
}

// UnmarshalJSON превращает строку обратно в ЦветФигуры.
func (ц *ЦветФигуры) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Белые"`, `"белые"`:
		*ц = Белые
	case `"Чёрные"`, `"чёрные"`, `"Черные"`, `"черные"`:
		*ц = Чёрные
	default:
		return fmt.Errorf("неизвестный цвет: %s", string(data))
	}
	return nil
}

// MarshalJSON превращает ТипФигуры в строку.
func (т ТипФигуры) MarshalJSON() ([]byte, error) {
	switch т {
	case Пешка:
		return []byte(`"Пешка"`), nil
	case Ладья:
		return []byte(`"Ладья"`), nil
	case Конь:
		return []byte(`"Конь"`), nil
	case Слон:
		return []byte(`"Слон"`), nil
	case Ферзь:
		return []byte(`"Ферзь"`), nil
	case Король:
		return []byte(`"Король"`), nil
	}
	return []byte(`"Неизвестно"`), nil
}

// UnmarshalJSON превращает строку обратно в ТипФигуры.
func (т *ТипФигуры) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Пешка"`, `"пешка"`:
		*т = Пешка
	case `"Ладья"`, `"ладья"`:
		*т = Ладья
	case `"Конь"`, `"конь"`:
		*т = Конь
	case `"Слон"`, `"слон"`:
		*т = Слон
	case `"Ферзь"`, `"ферзь"`:
		*т = Ферзь
	case `"Король"`, `"король"`:
		*т = Король
	default:
		return fmt.Errorf("неизвестный тип фигуры: %s", string(data))
	}
	return nil
}
