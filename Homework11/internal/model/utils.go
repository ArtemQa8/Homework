package model

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
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

func center(name string, boardWidth int) string {
	length := utf8.RuneCountInString(name)
	totalSpaces := boardWidth - length
	if totalSpaces <= 0 {
		return name
	}
	left := totalSpaces / 2
	right := totalSpaces - left
	return strings.Repeat(" ", left) + name + strings.Repeat(" ", right)
}

func columnName(index int) string {
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
	return columnName(col) + fmt.Sprintf("%d", row+1)
}

// GameState хранит информацию о ходе автоматической партии
type GameState struct {
	Game     *Game
	LastMove string
	MoveTime time.Duration
	Finished bool
}

// VisibleLength считает длину строки без учёта ANSI-кодов.
func VisibleLength(s string) int {
	length := 0
	inEscape := false
	for _, r := range s {
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if r == '\033' {
			inEscape = true
			continue
		}
		length++
	}
	return length
}

// PadWithSpaces дополняет строку пробелами до нужной видимой ширины.
func PadWithSpaces(s string, width int) string {
	current := VisibleLength(s)
	if current >= width {
		return s
	}
	return s + strings.Repeat(" ", width-current)
}

// GenitiveColor возвращает название цвета в родительном падеже.
func GenitiveColor(color Color) string {
	if color == White {
		return "Белых"
	}
	return "Чёрных"
}
