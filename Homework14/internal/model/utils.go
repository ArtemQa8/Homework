package model

import (
	"fmt"
	"time"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func ColumnName(index int) string {
	index++
	var result []byte
	for index > 0 {
		index--
		remainder := index % 26
		result = append([]byte{byte('A' + remainder)}, result...)
		index /= 26
	}
	return string(result)
}

func FormatMove(move Move) string {
	// рокировка
	if move.MovedPiece != nil && move.MovedPiece.Type() == King &&
		abs(move.ToCol-move.FromCol) == 2 {
		str := "O-O"
		if move.ToCol < move.FromCol {
			str = "O-O-O"
		}
		if move.Mate {
			str += "#"
		} else if move.Check {
			str += "+"
		}
		return str
	}

	symbol := ""
	if move.MovedPiece != nil && move.MovedPiece.Type() != Pawn {
		symbol = move.MovedPiece.Symbol() + " "
	}

	from := CellToString(move.FromRow, move.FromCol)
	to := CellToString(move.ToRow, move.ToCol)

	separator := "-"
	if move.Captured != nil {
		separator = "x"
	}

	str := symbol + from + separator + to
	if move.Mate {
		str += "#"
	} else if move.Check {
		str += "+"
	}
	return str
}

func CellToString(row, col int) string {
	return ColumnName(col) + fmt.Sprintf("%d", row+1)
}

// GameState хранит информацию о ходе автоматической партии
type GameState struct {
	Game     *Game
	LastMove string
	MoveTime time.Duration
	Finished bool
}

// GenitiveColor возвращает название цвета в родительном падеже.
func GenitiveColor(color Color) string {
	if color == White {
		return "Белых"
	}
	return "Чёрных"
}
