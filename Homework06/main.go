package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"mod.go/internal/model"
)

func main() {
	var строки, столбцы int
	fmt.Print("Введите кол-во строк: ")
	fmt.Scan(&строки)
	fmt.Print("Введите кол-во столбцов: ")
	fmt.Scan(&столбцы)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // съедаем оставшийся \n

	fmt.Print("Введите имя первого шахматиста: ")
	scanner.Scan()
	имя1 := scanner.Text()

	fmt.Print("Введите имя второго шахматиста: ")
	scanner.Scan()
	имя2 := scanner.Text()

	игра := model.НоваяИгра(имя1, имя2, строки, столбцы)
	fmt.Print(игра.Отобразить())

	for {
		цвет := игра.ТекущийЦвет()
		var имя string
		if цвет == model.Белые {
			имя = игра.Игрок1().Имя()
		} else {
			имя = игра.Игрок2().Имя()
		}

		fmt.Printf("\nХод %s (%s): ", цветРодительный(цвет), имя)
		fmt.Println("Введите ход (например, e2 e4) или 'exit':")

		if !scanner.Scan() {
			break
		}
		ввод := strings.TrimSpace(scanner.Text())
		if ввод == "" {
			continue
		}
		if strings.ToLower(ввод) == "exit" || strings.ToLower(ввод) == "quit" ||
			strings.ToLower(ввод) == "q" || strings.ToLower(ввод) == "й" ||
			strings.ToLower(ввод) == "йгше" || strings.ToLower(ввод) == "учше" ||
			strings.ToLower(ввод) == "выход" {
			fmt.Println("Игра завершена.")
			break
		}

		ход, err := парситьХод(ввод, игра)
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
		if !игра.СделатьХод(ход) {
			fmt.Println("Недопустимый ход!")
			continue
		}
		fmt.Print(игра.Отобразить())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка ввода: %v\n", err)
	}
}

// цветРодительный возвращает цвет в родительном падеже.
func цветРодительный(цвет model.ЦветФигуры) string {
	if цвет == model.Белые {
		return "Белых"
	}
	return "Чёрных"
}

// парситьХод разбирает строку вроде "e2 e4" и возвращает *model.Ход.
func парситьХод(ввод string, игра *model.Игра) (*model.Ход, error) {
	части := strings.Fields(ввод)
	if len(части) != 2 {
		return nil, fmt.Errorf("нужно указать две клетки, например e2 e4")
	}
	отСтр, отСтл, err := разобратьКлетку(части[0], игра)
	if err != nil {
		return nil, fmt.Errorf("начальная клетка: %w", err)
	}
	вСтр, вСтл, err := разобратьКлетку(части[1], игра)
	if err != nil {
		return nil, fmt.Errorf("конечная клетка: %w", err)
	}
	return model.НовыйХод(отСтр, отСтл, вСтр, вСтл), nil
}

// разобратьКлетку преобразует обозначение клетки (например "a1", "aa10", "abc123")
// в индексы строки и столбца (с 0).
func разобратьКлетку(нотация string, игра *model.Игра) (строка, столбец int, err error) {
	нотация = strings.TrimSpace(нотация)
	if len(нотация) < 2 {
		return 0, 0, fmt.Errorf("неверный формат клетки: %q", нотация)
	}

	// Ищем, где заканчиваются буквы и начинаются цифры
	буквеннаяЧасть := ""
	цифроваяЧасть := ""
	for i, r := range нотация {
		if unicode.IsLetter(r) {
			if цифроваяЧасть != "" {
				return 0, 0, fmt.Errorf("перемешаны буквы и цифры в %q", нотация)
			}
			буквеннаяЧасть += string(r)
		} else if unicode.IsDigit(r) {
			цифроваяЧасть = нотация[i:]
			break
		} else {
			return 0, 0, fmt.Errorf("недопустимый символ в %q", нотация)
		}
	}
	if буквеннаяЧасть == "" || цифроваяЧасть == "" {
		return 0, 0, fmt.Errorf("неверный формат клетки: %q", нотация)
	}

	// Преобразуем буквенную часть в столбец (A=0, B=1, ..., AA=26, AB=27...)
	столбец = буквыВИндекс(буквеннаяЧасть)
	if столбец < 0 || столбец >= игра.Доска().Столбцы() {
		return 0, 0, fmt.Errorf("столбец %q за пределами доски", буквеннаяЧасть)
	}

	// Преобразуем цифровую часть в номер строки (1 -> 0)
	var номерСтроки int
	_, err = fmt.Sscanf(цифроваяЧасть, "%d", &номерСтроки)
	if err != nil || номерСтроки < 1 || номерСтроки > игра.Доска().Строки() {
		return 0, 0, fmt.Errorf("неверный номер строки: %q", цифроваяЧасть)
	}
	строка = номерСтроки - 1
	return
}

// буквыВИндекс преобразует буквенное обозначение столбца (A, B, ..., Z, AA, AB...)
// в числовой индекс (0, 1, ..., 25, 26, 27...).
func буквыВИндекс(буквы string) int {
	индекс := 0
	for _, символ := range буквы {
		// Приводим к верхнему регистру
		if символ >= 'a' && символ <= 'z' {
			символ = символ - 'a' + 'A'
		}
		if символ < 'A' || символ > 'Z' {
			return -1
		}
		индекс = индекс*26 + int(символ-'A'+1)
	}
	return индекс - 1
}
