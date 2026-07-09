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
