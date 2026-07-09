package model

type ФигураСПравилами interface {
	МожетХодить(ход *Ход, доска *Доска) bool
}

func знак(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ПравилаДля(фигура *Фигура) ФигураСПравилами {
	switch фигура.Тип() {
	case Пешка:
		return &Pawn{Цвет: фигура.Цвет()}
	case Ладья:
		return &Rook{Цвет: фигура.Цвет()}
	case Конь:
		return &Knight{Цвет: фигура.Цвет()}
	case Слон:
		return &Bishop{Цвет: фигура.Цвет()}
	case Ферзь:
		return &Queen{Цвет: фигура.Цвет()}
	case Король:
		return &King{Цвет: фигура.Цвет()}
	default:
		return nil
	}
}
