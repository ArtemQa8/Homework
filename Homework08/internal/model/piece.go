package model

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
