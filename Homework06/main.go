package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"mod.go/internal/model"
	"mod.go/internal/repository"
	"mod.go/internal/service"
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
	хранилище := &repository.Хранилище{}

	// Сохраняем игру и игроков один раз
	service.СохранитьИгру(игра, хранилище)
	service.СохранитьИгроков(игра, хранилище)

	// Хранилище (начальное состояние)
	fmt.Printf("[Хранилище] Игр: %d, Игроков: %d, Ходов: %d | Игроки: %s, %s\n",
		len(хранилище.Игры), len(хранилище.Игроки), len(хранилище.Ходы),
		хранилище.Игроки[0].Имя(), хранилище.Игроки[1].Имя()) // Хранилище

	fmt.Print("\033[H\033[2J")
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

		// Сдался
		if strings.ToLower(ввод) == "сдался" {
			fmt.Printf("Игрок %s сдался. Победил(а) %s!\n", имя, имяПобедителя(игра))
			break
		}

		// Ничья
		if strings.ToLower(ввод) == "ничья" {
			fmt.Printf("%s предлагает ничью. %s, вы согласны? (да/нет): ",
				имя, имяСоперника(игра, цвет))
			if !scanner.Scan() {
				break
			}
			ответ := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if ответ == "да" || ответ == "y" || ответ == "yes" || ответ == "д" || ответ == "ага" {
				fmt.Println("Ничья! Игра завершена.")
				break
			} else {
				fmt.Println("Ничья отклонена. Игра продолжается.")
				continue
			}
		}

		// Автоход
		if strings.HasPrefix(strings.ToLower(ввод), "автоход") {
			части := strings.Fields(ввод)
			if len(части) < 2 {
				fmt.Println("Укажите количество ходов: автоход <число>")
				continue
			}
			число, err := strconv.Atoi(части[1])
			if err != nil || число <= 0 {
				fmt.Println("Неверное число ходов. Пример: автоход 5")
				continue
			}
			for i := 0; i < число; i++ {
				ход, err := service.СделатьАвтоход(игра)
				if err != nil {
					fmt.Println("Автоход:", err)
					break
				}
				if err := игра.СделатьХод(ход); err != nil {
					fmt.Println("Ошибка выполнения автохода:", err)
					break
				}
				service.СохранитьХод(игра, хранилище)

				// Хранилище
				последний := хранилище.Ходы[len(хранилище.Ходы)-1]
				fmt.Printf("[Хранилище] Игр: %d, Игроков: %d, Ходов: %d | Игроки: %s, %s | Последний: %s\n",
					len(хранилище.Игры), len(хранилище.Игроки), len(хранилище.Ходы),
					хранилище.Игроки[0].Имя(), хранилище.Игроки[1].Имя(),
					model.ФорматироватьХод(последний)) // Хранилище

				fmt.Print("\033[H\033[2J")
				fmt.Print(игра.Отобразить())
				if игра.Мат(игра.ТекущийЦвет()) {
					fmt.Printf("Мат! Победил(а) %s.\n", имяПобедителя(игра))
					return
				} else if игра.Шах(игра.ТекущийЦвет()) {
					fmt.Println("Шах!")
				}
				time.Sleep(500 * time.Millisecond)
			}
			continue
		}

		ход, err := обработатьВвод(ввод, игра, scanner)
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}

		if err := игра.СделатьХод(ход); err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
		service.СохранитьХод(игра, хранилище)

		// Хранилище
		последний := хранилище.Ходы[len(хранилище.Ходы)-1]
		fmt.Printf("[Хранилище] Игр: %d, Игроков: %d, Ходов: %d | Игроки: %s, %s | Последний: %s\n",
			len(хранилище.Игры), len(хранилище.Игроки), len(хранилище.Ходы),
			хранилище.Игроки[0].Имя(), хранилище.Игроки[1].Имя(),
			model.ФорматироватьХод(последний)) // Хранилище

		fmt.Print("\033[H\033[2J")
		fmt.Print(игра.Отобразить())

		if игра.Мат(игра.ТекущийЦвет()) {
			fmt.Printf("Мат! Победил(а) %s.\n", имяПобедителя(игра))
			break
		} else if игра.Шах(игра.ТекущийЦвет()) {
			fmt.Println("Шах!")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка ввода: %v\n", err)
	}
}

func цветРодительный(цвет model.ЦветФигуры) string {
	if цвет == model.Белые {
		return "Белых"
	}
	return "Чёрных"
}

func имяПобедителя(игра *model.Игра) string {
	if игра.ТекущийЦвет() == model.Белые {
		return игра.Игрок2().Имя()
	}
	return игра.Игрок1().Имя()
}

func имяСоперника(игра *model.Игра, цвет model.ЦветФигуры) string {
	if цвет == model.Белые {
		return игра.Игрок2().Имя()
	}
	return игра.Игрок1().Имя()
}

func обработатьВвод(ввод string, игра *model.Игра, scanner *bufio.Scanner) (*model.Ход, error) {
	части := strings.Fields(ввод)
	if len(части) != 2 {
		return nil, fmt.Errorf("нужно указать две клетки, например e2 e4")
	}

	отСтр, отСтл, err := model.РазобратьКлетку(части[0], игра.Доска())
	if err != nil {
		return nil, err
	}
	вСтр, вСтл, err := model.РазобратьКлетку(части[1], игра.Доска())
	if err != nil {
		return nil, err
	}
	ход := model.НовыйХод(отСтр, отСтл, вСтр, вСтл)

	фигура := игра.Доска().ФигураНа(отСтр, отСтл)
	if фигура != nil && фигура.Тип() == model.Пешка {

		// если правила для фигуры отсуствуют
		// или ход не является разрешенным (геометрически)
		// то сразу выходим и не предлагаем превращение
		правила := model.ПравилаДля(фигура)
		if правила == nil || !правила.МожетХодить(ход, игра.Доска()) {
			return ход, nil
		}

		последняяСтрока := 0
		if фигура.Цвет() == model.Белые {
			последняяСтрока = игра.Доска().Строки() - 1
		}
		if вСтр == последняяСтрока {
			fmt.Print("Во что превратить пешку? (ферзь, ладья, слон, конь): ")
			if !scanner.Scan() {
				return nil, fmt.Errorf("ввод прерван")
			}
			выбор := strings.TrimSpace(scanner.Text())
			switch strings.ToLower(выбор) {
			case "ферзь", "ф", "q":
				ход.Превращение = model.Ферзь
			case "ладья", "л", "r":
				ход.Превращение = model.Ладья
			case "слон", "с", "b":
				ход.Превращение = model.Слон
			case "конь", "к", "n":
				ход.Превращение = model.Конь
			default:
				return nil, fmt.Errorf("недопустимый выбор фигуры")
			}
		}
	}
	return ход, nil
}
