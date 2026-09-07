package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"mod.go/internal/model"
	"mod.go/internal/repository"
	"mod.go/internal/service"
)

var (
	sims      []*model.СостояниеПартии
	muSim     sync.Mutex
	simActive bool
	wgSim     sync.WaitGroup // только для симуляций
	wgMain    sync.WaitGroup // для читателя, логгера, отрисовщика
)

func main() {
	var строки, столбцы int
	fmt.Print("Введите кол-во строк: ")
	fmt.Scan(&строки)
	fmt.Print("Введите кол-во столбцов: ")
	fmt.Scan(&столбцы)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	fmt.Print("Введите имя первого шахматиста: ")
	scanner.Scan()
	имя1 := scanner.Text()

	fmt.Print("Введите имя второго шахматиста: ")
	scanner.Scan()
	имя2 := scanner.Text()

	ctxRoot, cancelRoot := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelRoot()

	игра := model.НоваяИгра(имя1, имя2, строки, столбцы)
	хранилище := repository.НовоеХранилище()
	if err := хранилище.ЗагрузитьИзФайлов(); err != nil {
		fmt.Printf("Ошибка загрузки данных: %v\n", err)
	}
	fmt.Println("Загружено игр:", len(хранилище.Игры))
	fmt.Println("Загружено игроков:", len(хранилище.Игроки))
	fmt.Println("Загружено ходов:", len(хранилище.Ходы))

	шина := make(chan model.ОбъектХранилища, 100)

	// Горутина-читатель
	wgMain.Add(1)
	go func() {
		defer wgMain.Done()
		for {
			select {
			case <-ctxRoot.Done():
				return
			case объект, ok := <-шина:
				if !ok {
					return
				}
				хранилище.Добавить(объект)
			}
		}
	}()

	// Логгер
	wgMain.Add(1)
	go service.ЗапуститьЛоггер(ctxRoot, хранилище, &simActive, &wgMain)

	service.СохранитьИгру(игра, шина)
	service.СохранитьИгроков(игра, шина)

	fmt.Print("\033[H\033[2J")
	fmt.Print(игра.Отобразить())
	показатьПомощь()

	// Горутина-отрисовщик
	wgMain.Add(1)
	go func() {
		defer wgMain.Done()
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctxRoot.Done():
				return
			case <-ticker.C:
				if !simActive {
					continue
				}
				fmt.Print("\033[H\033[2J")
				muSim.Lock()
				var активные []*model.СостояниеПартии
				for _, сим := range sims {
					if !сим.Завершена {
						активные = append(активные, сим)
					}
				}
				if len(активные) > 0 {
					примерСтроки := активные[0].Игра.СтрокиОтображенияБезИстории()
					ширинаДоски := 0
					if len(примерСтроки) > 0 {
						ширинаДоски = visibleLength(примерСтроки[0])
					}
					отступ := 3
					вРяд := 3
					maxTermWidth := 120
					if ширинаДоски*3+отступ*2 > maxTermWidth {
						вРяд = 2
						if ширинаДоски*2+отступ > maxTermWidth {
							вРяд = 1
						}
					}
					for начало := 0; начало < len(активные); начало += вРяд {
						конец := начало + вРяд
						if конец > len(активные) {
							конец = len(активные)
						}
						группа := активные[начало:конец]
						доски := make([][]string, len(группа))
						максСтрок := 0
						for i, сим := range группа {
							строки := сим.Игра.СтрокиОтображенияБезИстории()
							доски[i] = строки
							if len(строки) > максСтрок {
								максСтрок = len(строки)
							}
						}
						ширины := make([]int, len(группа))
						for i, д := range доски {
							for _, стр := range д {
								w := visibleLength(стр)
								if w > ширины[i] {
									ширины[i] = w
								}
							}
						}
						for стр := 0; стр < максСтрок; стр++ {
							for i := 0; i < len(группа); i++ {
								if стр < len(доски[i]) {
									текст := доски[i][стр]
									fmt.Print(alignVisible(текст, ширины[i]) + "   ")
								} else {
									fmt.Print(strings.Repeat(" ", ширины[i]) + "   ")
								}
							}
							fmt.Println()
						}
						fmt.Println()
					}
				}
				for i, сим := range sims {
					имя1 := сим.Игра.Игрок1().Имя()
					имя2 := сим.Игра.Игрок2().Имя()
					if сим.Завершена {
						fmt.Printf("#%d %s vs %s: %s (завершена)\n", i+1, имя1, имя2, сим.ПоследнийХод)
					} else {
						fmt.Printf("#%d %s vs %s: %s за %v\n", i+1, имя1, имя2, сим.ПоследнийХод, сим.ВремяХода.Round(time.Millisecond))
					}
				}
				muSim.Unlock()
			}
		}
	}()

	var cancelSim context.CancelFunc

	// Главный игровой цикл
	for {
		цвет := игра.ТекущийЦвет()
		var имя string
		if цвет == model.Белые {
			имя = игра.Игрок1().Имя()
		} else {
			имя = игра.Игрок2().Имя()
		}

		if simActive {
			fmt.Println("\nСимуляции идут. Введите 'стоп' для остановки, 'exit' для выхода.")
		} else {
			fmt.Printf("\nХод %s (%s): ", model.РодительныйЦвет(цвет), имя)
			fmt.Println("Введите ход (например, e2 e4) или 'exit':")
		}

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
			cancelRoot()
			wgSim.Wait()
			wgMain.Wait()
			fmt.Println("Программа завершена.")
			return
		}

		if strings.ToLower(ввод) == "стоп" {
			if simActive && cancelSim != nil {
				cancelSim()
				wgSim.Wait()
				simActive = false
				cancelSim = nil
				fmt.Print("\033[H\033[2J")
				fmt.Print(игра.Отобразить())
			}
			continue
		}

		if simActive {
			fmt.Println("Ручные ходы недоступны во время симуляций.")
			continue
		}

		if strings.ToLower(ввод) == "сдался" {
			fmt.Printf("Игрок %s сдался. Победил(а) %s!\n", имя, имяПобедителя(игра))
			cancelRoot()
			break
		}

		if strings.ToLower(ввод) == "ничья" {
			fmt.Printf("%s предлагает ничью. %s, вы согласны? (да/нет): ",
				имя, имяСоперника(игра, цвет))
			if !scanner.Scan() {
				break
			}
			ответ := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if ответ == "да" || ответ == "y" || ответ == "yes" || ответ == "д" || ответ == "ага" {
				fmt.Println("Ничья! Игра завершена.")
				cancelRoot()
				break
			} else {
				fmt.Println("Ничья отклонена. Игра продолжается.")
				continue
			}
		}

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
				service.СохранитьХод(игра, шина)
				time.Sleep(500 * time.Millisecond)
			}
			continue
		}

		if strings.HasPrefix(strings.ToLower(ввод), "симуляция") {
			части := strings.Fields(ввод)
			if len(части) < 2 {
				fmt.Println("Укажите количество досок: симуляция <число>")
				continue
			}
			число, err := strconv.Atoi(части[1])
			if err != nil || число <= 0 {
				fmt.Println("Неверное число. Пример: симуляция 3")
				continue
			}
			simActive = true
			ctxSim, cancel := context.WithCancel(ctxRoot)
			cancelSim = cancel

			for i := 0; i < число; i++ {
				var имя1, имя2 string
				if rand.Intn(2) == 0 {
					имя1 = model.ИменаМуж[rand.Intn(len(model.ИменаМуж))] + " " + model.ФамилииМуж[rand.Intn(len(model.ФамилииМуж))]
				} else {
					имя1 = model.ИменаЖен[rand.Intn(len(model.ИменаЖен))] + " " + model.ФамилииЖен[rand.Intn(len(model.ФамилииЖен))]
				}
				if rand.Intn(2) == 0 {
					имя2 = model.ИменаМуж[rand.Intn(len(model.ИменаМуж))] + " " + model.ФамилииМуж[rand.Intn(len(model.ФамилииМуж))]
				} else {
					имя2 = model.ИменаЖен[rand.Intn(len(model.ИменаЖен))] + " " + model.ФамилииЖен[rand.Intn(len(model.ФамилииЖен))]
				}

				новаяИгра := model.НоваяИгра(имя1, имя2, строки, столбцы)
				service.СохранитьИгру(новаяИгра, шина)
				service.СохранитьИгроков(новаяИгра, шина)

				сим := &model.СостояниеПартии{Игра: новаяИгра}
				muSim.Lock()
				sims = append(sims, сим)
				muSim.Unlock()
				wgSim.Add(1)
				go func(с *model.СостояниеПартии) {
					defer wgSim.Done()
					service.СимулироватьПартию(ctxSim, с, шина)
				}(сим)
				fmt.Printf("Запущена симуляция #%d: %s vs %s\n", len(sims), имя1, имя2)
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
		service.СохранитьХод(игра, шина)
		fmt.Print("\033[H\033[2J")
		fmt.Print(игра.Отобразить())

		if игра.Мат(игра.ТекущийЦвет()) {
			fmt.Printf("Мат! Победил(а) %s.\n", имяПобедителя(игра))
			cancelRoot()
			break
		} else if игра.Шах(игра.ТекущийЦвет()) {
			fmt.Println("Шах!")
		}
	}

	// Если вышли из цикла по `break` (сдался, ничья, мат, ошибка сканера), а не по `exit`
	if cancelSim != nil {
		cancelSim()
	}
	cancelRoot()
	wgSim.Wait()
	wgMain.Wait()
	if err := хранилище.СохранитьВсе(); err != nil {
		fmt.Printf("Ошибка сохранения данных: %v\n", err)
	}
	fmt.Println("Программа завершена.")
}

// ---------- Вспомогательные функции ----------

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
	// Превращение пешки предлагаем только если фигура есть и её цвет совпадает с текущим
	if фигура != nil && фигура.Цвет() == игра.ТекущийЦвет() && фигура.Тип() == model.Пешка {
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

func visibleLength(s string) int {
	count := 0
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
		count++
	}
	return count
}

func alignVisible(s string, width int) string {
	current := visibleLength(s)
	if current >= width {
		return s
	}
	return s + strings.Repeat(" ", width-current)
}

func показатьПомощь() {
	fmt.Println()
	fmt.Println("Доступные команды:")
	fmt.Println("  сдался            — признать поражение")
	fmt.Println("  ничья             — предложить ничью")
	fmt.Println("  автоход <N>       — сделать N случайных ходов")
	fmt.Println("  симуляция <N>     — запустить N автоматических партий в фоне")
	fmt.Println("  стоп              — остановить симуляции")
	fmt.Println("  exit / выход      — завершить программу")
	fmt.Println()
}
