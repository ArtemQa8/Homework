package main

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

// цвета
const (
	Reset       = "\033[0m"
	FgWhite     = "\033[31m" // красный
	FgLightGray = "\033[34m" // синий
)

// центрирование имени
func центрирование(имя string, ширинаДоски int) string {
	имяСимволов := utf8.RuneCountInString(имя)
	всегоПробелов := ширинаДоски - имяСимволов
	if всегоПробелов <= 0 {
		return имя // Если имя длинное - вернуть имя
	}
	левыйОтступ := всегоПробелов / 2
	правыйОтступ := всегоПробелов - левыйОтступ // так мы гарантируем, что сложив левый и правый получим всего
	return strings.Repeat(" ", левыйОтступ) + имя + strings.Repeat(" ", правыйОтступ)
}

// выравнивание номеров строк, если числа двухзначные, трехзначные и т.д.
func выравнивание(номерСтроки, ширина int) string {
	return fmt.Sprintf("%*d ", ширина, номерСтроки)
}

// расстановка и покраска фигур
func расстановка(строка, столбец, a, b uint) string {
	// 1. Пешки
	if строка == 1 { // белые пешки во второй строке
		return FgWhite + "♙" + Reset
	}
	if строка == a-2 { // черные пешки в предпоследней строке
		return FgLightGray + "♟" + Reset
	}

	// 2. Белые фигуры
	if строка == 0 { // первая строка
		switch столбец {
		case 0, b - 1: // первая и последняя клетки в строке
			return FgWhite + "♖" + Reset
		case 1, b - 2: // вторая и предпоследняя клетки в строке
			return FgWhite + "♘" + Reset
		case 2, b - 3:
			return FgWhite + "♗" + Reset
		case b / 2:
			return FgWhite + "♕" + Reset
		case b/2 - 1:
			return FgWhite + "♔" + Reset
		default: // если нет фигуры, вернуть "цвет"
			if (строка+столбец)%2 == 0 {
				return " "
			} else {
				return "#"
			}
		}
	}

	// 3. Чёрные фигуры
	if строка == a-1 { // последняя строка
		switch столбец {
		case 0, b - 1: // первая и последняя клетки в строке
			return FgLightGray + "♜" + Reset
		case 1, b - 2: // вторая и предпоследняя клетки в строке
			return FgLightGray + "♞" + Reset
		case 2, b - 3:
			return FgLightGray + "♝" + Reset
		case b / 2:
			return FgLightGray + "♛" + Reset
		case b/2 - 1:
			return FgLightGray + "♚" + Reset
		default: // если нет фигуры, вернуть "цвет"
			if (строка+столбец)%2 == 0 {
				return " "
			} else {
				return "#"
			}
		}
	}

	// 4. Середина доски (пустые клетки)
	// возвращаем "цвет"
	if (строка+столбец)%2 == 0 {
		return " "
	} else {
		return "#"
	}
}

func main() {
	var chessGrid, name1, name2 string
	var a, b uint

	fmt.Print("Введите кол-во строк: ")
	fmt.Scan(&a)

	fmt.Print("Введите кол-во столбцов: ")
	fmt.Scan(&b)

	fmt.Print("Введите имя первого шахматиста: ")
	fmt.Scan(&name1)

	fmt.Print("Введите имя второго шахматиста: ")
	fmt.Scan(&name2)

	// максимальная ширина номеров строк
	максШирина := int(math.Log10(float64(a))) + 1

	// ширина доски со всеми пробелами
	ширинаДоски := int(b)*2 + максШирина + 1

	// отцентрованное имя первого шахматиста
	chessGrid += центрирование(name1, ширинаДоски) + "\n"

	// отступ слева, чтобы буквы не налезали на цифры
	chessGrid += strings.Repeat(" ", максШирина+1)

	for столбец := uint(0); столбец < b; столбец++ {
		chessGrid += string('A'+столбец) + " "
	}
	chessGrid += "\n"

	for строка := uint(0); строка < a; строка++ {
		chessGrid += выравнивание(int(строка)+1, максШирина)

		for столбец := uint(0); столбец < b; столбец++ {
			поле := расстановка(строка, столбец, a, b)

			if столбец == b-1 {
				chessGrid += поле + " "
			} else {
				chessGrid += поле + "|"
			}
		}
		chessGrid += "\n"
	}

	// отцентрованное имя второго шахматиста
	chessGrid += центрирование(name2, ширинаДоски) + "\n"

	fmt.Print(chessGrid)
}
