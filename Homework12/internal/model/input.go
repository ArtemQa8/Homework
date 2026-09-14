package model

import (
	"fmt"
	"strings"
	"unicode"
)

func ParseCell(notation string, board *Board) (row, col int, err error) {
	notation = strings.TrimSpace(notation)
	if len(notation) < 2 {
		return 0, 0, fmt.Errorf("неверный формат клетки: %q", notation)
	}

	letterPart := ""
	digitPart := ""
	for i, r := range notation {
		if unicode.IsLetter(r) {
			if digitPart != "" {
				return 0, 0, fmt.Errorf("перемешаны буквы и цифры в %q", notation)
			}
			letterPart += string(r)
		} else if unicode.IsDigit(r) {
			digitPart = notation[i:]
			break
		} else {
			return 0, 0, fmt.Errorf("недопустимый символ в %q", notation)
		}
	}
	if letterPart == "" || digitPart == "" {
		return 0, 0, fmt.Errorf("неверный формат клетки: %q", notation)
	}

	col = LettersToIndex(letterPart)
	if col < 0 || col >= board.Cols() {
		return 0, 0, fmt.Errorf("столбец %q за пределами доски", letterPart)
	}

	var rowNumber int
	_, err = fmt.Sscanf(digitPart, "%d", &rowNumber)
	if err != nil || rowNumber < 1 || rowNumber > board.Rows() {
		return 0, 0, fmt.Errorf("неверный номер строки: %q", digitPart)
	}
	row = rowNumber - 1
	return
}

func LettersToIndex(letters string) int {
	index := 0
	for _, ch := range letters {
		if ch >= 'a' && ch <= 'z' {
			ch = ch - 'a' + 'A'
		}
		if ch < 'A' || ch > 'Z' {
			return -1
		}
		index = index*26 + int(ch-'A'+1)
	}
	return index - 1
}
