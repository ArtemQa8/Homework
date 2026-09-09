package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"mod.go/internal/model"
)

var serverAddr = flag.String("addr", "http://localhost:8080", "адрес сервера")
var gameIDFlag = flag.Int("game", 0, "ID игры для подключения")

var httpClient = &http.Client{Timeout: 10 * time.Second}

var (
	simActive bool
	muSim     sync.Mutex
	sims      []*Симуляция
)

type Симуляция struct {
	ID           int
	Игра         model.Игра
	Завершена    bool
	ПоследнийХод string
	ВремяХода    time.Duration
}

type ОтветАвтоХода struct {
	Игра       model.Игра       `json:"игра"`
	Мат        bool             `json:"мат,omitempty"`
	Пат        bool             `json:"пат,omitempty"`
	Победитель model.ЦветФигуры `json:"победитель,omitempty"`
}

func main() {
	flag.Parse()

	inputChan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Ошибка чтения ввода: %v\n", err)
		}
		close(inputChan)
	}()

	запуститьОтрисовщикСимуляций()

	if *gameIDFlag > 0 {
		if err := игроваяСессия(*gameIDFlag, inputChan); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for {
		fmt.Println()
		fmt.Println("Добро пожаловать в шахматный клиент!")
		fmt.Printf("Наблюдение за играми: %s/\n", *serverAddr)
		fmt.Print("Выберите действие: ")
		fmt.Println()
		fmt.Println("1. Создать новую игру")
		fmt.Println("2. Подключиться к существующей")
		fmt.Println("3. Симуляция")
		fmt.Println("0. Выход")

		выбор, ok := <-inputChan
		if !ok {
			fmt.Println("Ввод закрыт.")
			return
		}
		выбор = strings.TrimSpace(выбор)

		switch выбор {
		case "1":
			if err := создатьНовуюИгру(inputChan); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}
		case "2":
			if err := подключитьсяКИгре(inputChan); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}
		case "3":
			симуляцияИзМеню(inputChan)
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный ввод, попробуйте ещё раз.")
		}
	}
}

func прочитатьСтроку(inputChan chan string) (string, bool) {
	строка, ok := <-inputChan
	return strings.TrimSpace(строка), ok
}

func прочитатьЧисло(inputChan chan string) (int, bool) {
	строка, ok := прочитатьСтроку(inputChan)
	if !ok {
		return 0, false
	}
	число, err := strconv.Atoi(строка)
	if err != nil {
		return 0, false
	}
	return число, true
}

func создатьНовуюИгру(inputChan chan string) error {
	fmt.Print("Имя первого игрока (белые): ")
	имя1, ok := прочитатьСтроку(inputChan)
	if !ok {
		return fmt.Errorf("ввод закрыт")
	}

	fmt.Print("Имя второго игрока (чёрные): ")
	имя2, ok := прочитатьСтроку(inputChan)
	if !ok {
		return fmt.Errorf("ввод закрыт")
	}

	fmt.Print("Количество строк (0 для стандартных 8): ")
	строки, ok := прочитатьЧисло(inputChan)
	if !ok {
		fmt.Println("Неверное число строк, используются 8")
		строки = 0
	}

	fmt.Print("Количество столбцов (0 для стандартных 8): ")
	столбцы, ok := прочитатьЧисло(inputChan)
	if !ok {
		fmt.Println("Неверное число столбцов, используются 8")
		столбцы = 0
	}

	игра, err := создатьИгру(имя1, имя2, строки, столбцы)
	if err != nil {
		return err
	}
	fmt.Printf("Игра создана, ID: %d\n", игра.ID())
	return игроваяСессия(игра.ID(), inputChan)
}

func подключитьсяКИгре(inputChan chan string) error {
	fmt.Print("Введите ID игры (или 'список' для показа всех): ")
	ввод, ok := прочитатьСтроку(inputChan)
	if !ok {
		return fmt.Errorf("ввод закрыт")
	}

	if strings.ToLower(ввод) == "список" {
		игры, err := получитьСписокИгр()
		if err != nil {
			return err
		}
		if len(игры) == 0 {
			fmt.Println("Нет доступных игр.")
			return nil
		}
		fmt.Println("Доступные игры:")
		for _, игра := range игры {
			fmt.Printf("ID: %d | %s (белые) - %s (чёрные)\n",
				игра.ID(), игра.Игрок1().Имя(), игра.Игрок2().Имя())
		}
		fmt.Print("Введите ID игры: ")
		ввод, ok = прочитатьСтроку(inputChan)
		if !ok {
			return fmt.Errorf("ввод закрыт")
		}
	}

	id, err := strconv.Atoi(ввод)
	if err != nil {
		return fmt.Errorf("неверный ID")
	}
	return игроваяСессия(id, inputChan)
}

func симуляцияИзМеню(inputChan chan string) {
	fmt.Print("Количество симуляций: ")
	n, ok := прочитатьЧисло(inputChan)
	if !ok || n <= 0 {
		fmt.Println("Неверное число симуляций.")
		return
	}

	fmt.Print("Количество строк (0 для стандартных 8): ")
	строки, _ := прочитатьЧисло(inputChan)
	fmt.Print("Количество столбцов (0 для стандартных 8): ")
	столбцы, _ := прочитатьЧисло(inputChan)

	if строки <= 0 || столбцы <= 0 {
		строки, столбцы = 8, 8
	}

	запуститьСимуляции(n, строки, столбцы)
	simActive = true

	fmt.Println("Симуляции запущены. Введите 'стоп' для остановки.")
	for {
		ввод, ok := прочитатьСтроку(inputChan)
		if !ok {
			break
		}
		if strings.ToLower(ввод) == "стоп" || strings.ToLower(ввод) == "выход" {
			остановитьСимуляции()
			simActive = false
			fmt.Print("\033[H\033[2J")
			return
		}
	}
}

func отобразитьИгру(игра model.Игра, id int) {
	fmt.Print("\033[H\033[2J")
	fmt.Print(игра.ОтобразитьСИД(id))
}

func показатьКомандыИгры(id int) {
	fmt.Println("Доступные команды:")
	fmt.Println("  автоход <N>     - сделать N автоматических ходов")
	fmt.Println("  симуляция <N>   - запустить N фоновых партий")
	fmt.Println("  стоп            - остановить симуляции (если запущены)")
	fmt.Println("  help / команды  - показать этот список")
	fmt.Println("  выход / exit    - выйти из игры")
	fmt.Printf("  Наблюдать за игрой: %s/game?id=%d\n", *serverAddr, id)
}

func печататьПриглашение(игра model.Игра) {
	цвет := игра.ТекущийЦвет()
	var имя string
	if цвет == model.Белые {
		имя = игра.Игрок1().Имя()
	} else {
		имя = игра.Игрок2().Имя()
	}
	fmt.Printf("\nХод %s (%s): ", model.РодительныйЦвет(цвет), имя)
	fmt.Println("Введите ход (например: e2 e4) или 'help':")
}

func игроваяСессия(id int, inputChan chan string) error {
	игра, err := получитьИгру(id)
	if err != nil {
		return err
	}

	отобразитьИгру(игра, id)
	fmt.Printf("Наблюдать за партией: %s/game?id=%d\n", *serverAddr, id)
	показатьКомандыИгры(id)
	печататьПриглашение(игра)

	var muGame sync.Mutex

	// Горутина-наблюдатель за обновлениями
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			if simActive {
				continue
			}
			свежая, err := получитьИгру(id)
			if err != nil {
				continue
			}

			muGame.Lock()
			изменилась := len(свежая.Ходы()) != len(игра.Ходы()) || свежая.ТекущийЦвет() != игра.ТекущийЦвет()
			if изменилась {
				игра = свежая
				отобразитьИгру(игра, id)
				печататьПриглашение(игра)
			}
			muGame.Unlock()
		}
	}()

	for ввод := range inputChan {
		if ввод == "" {
			continue
		}

		muGame.Lock()

		if strings.ToLower(ввод) == "help" || strings.ToLower(ввод) == "команды" || strings.ToLower(ввод) == "помощь" {
			показатьКомандыИгры(id)
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		if strings.ToLower(ввод) == "exit" || strings.ToLower(ввод) == "quit" ||
			strings.ToLower(ввод) == "q" || strings.ToLower(ввод) == "й" ||
			strings.ToLower(ввод) == "йгше" || strings.ToLower(ввод) == "учше" ||
			strings.ToLower(ввод) == "выход" {
			fmt.Println("Выход из игры.")
			muGame.Unlock()
			return nil
		}

		if simActive {
			if strings.ToLower(ввод) == "стоп" {
				остановитьСимуляции()
				simActive = false
				отобразитьИгру(игра, id)
				печататьПриглашение(игра)
				muGame.Unlock()
				continue
			}
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		if strings.HasPrefix(strings.ToLower(ввод), "автоход") {
			части := strings.Fields(ввод)
			if len(части) < 2 {
				fmt.Println("Укажите количество ходов: автоход <число>")
				muGame.Unlock()
				печататьПриглашение(игра)
				continue
			}
			n, err := strconv.Atoi(части[1])
			if err != nil || n <= 0 {
				fmt.Println("Неверное число ходов.")
				muGame.Unlock()
				печататьПриглашение(игра)
				continue
			}
			for range n {
				обновлённая, мат, пат, победитель, err := сделатьАвтоХод(id)
				if err != nil {
					fmt.Println("Ошибка:", err)
					break
				}
				игра = обновлённая

				if мат {
					победительЦвет := "Белые"
					if победитель == model.Чёрные {
						победительЦвет = "Чёрные"
					}
					отобразитьИгру(игра, id)
					печататьПриглашение(игра)
					fmt.Printf("\nМат! Победили %s\n", победительЦвет)
					muGame.Unlock()
					return nil
				} else if пат {
					отобразитьИгру(игра, id)
					печататьПриглашение(игра)
					fmt.Println("\nПат. Ничья")
					muGame.Unlock()
					return nil
				}

				отобразитьИгру(игра, id)
				печататьПриглашение(игра)
				time.Sleep(time.Duration(2000+rand.Intn(2000)) * time.Millisecond)
			}
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		if strings.HasPrefix(strings.ToLower(ввод), "симуляция") {
			части := strings.Fields(ввод)
			if len(части) < 2 {
				fmt.Println("Укажите количество досок: симуляция <число>")
				muGame.Unlock()
				печататьПриглашение(игра)
				continue
			}
			n, err := strconv.Atoi(части[1])
			if err != nil || n <= 0 {
				fmt.Println("Неверное число.")
				muGame.Unlock()
				печататьПриглашение(игра)
				continue
			}
			запуститьСимуляции(n, игра.Доска().Строки(), игра.Доска().Столбцы())
			simActive = true
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		части := strings.Fields(ввод)
		if len(части) != 2 {
			fmt.Println("Нужно указать две клетки, например: e2 e4")
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		отСтр, отСтл, err := model.РазобратьКлетку(части[0], игра.Доска())
		if err != nil {
			fmt.Println("Ошибка:", err)
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}
		вСтр, вСтл, err := model.РазобратьКлетку(части[1], игра.Доска())
		if err != nil {
			fmt.Println("Ошибка:", err)
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		ход := model.Ход{
			ОтСтрока:  отСтр,
			ОтСтолбец: отСтл,
			ВСтрока:   вСтр,
			ВСтолбец:  вСтл,
		}

		обновлённая, err := сделатьХод(id, ход)
		if err != nil {
			fmt.Println("Ошибка:", err)
			muGame.Unlock()
			печататьПриглашение(игра)
			continue
		}

		игра = обновлённая
		отобразитьИгру(игра, id)
		печататьПриглашение(игра)

		if игра.Мат(игра.ТекущийЦвет()) {
			победительЦвет := "Белые"
			if игра.ТекущийЦвет() == model.Белые {
				победительЦвет = "Чёрные"
			}
			fmt.Printf("\nМат! Победили %s\n", победительЦвет)
			muGame.Unlock()
			return nil
		} else if игра.Пат(игра.ТекущийЦвет()) {
			fmt.Println("\nПат. Ничья")
			muGame.Unlock()
			return nil
		}

		muGame.Unlock()
	}
	return nil
}

// ---------- HTTP-функции ----------

func получитьИгру(id int) (model.Игра, error) {
	var игра model.Игра
	resp, err := httpClient.Get(fmt.Sprintf("%s/api/games/%d", *serverAddr, id))
	if err != nil {
		return игра, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return игра, fmt.Errorf("Сервер вернул статус %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&игра); err != nil {
		return игра, err
	}
	return игра, nil
}

func получитьСписокИгр() ([]model.Игра, error) {
	var игры []model.Игра
	resp, err := httpClient.Get(fmt.Sprintf("%s/api/games", *serverAddr))
	if err != nil {
		return игры, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return игры, fmt.Errorf("Сервер вернул статус %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&игры); err != nil {
		return игры, err
	}
	return игры, nil
}

func создатьИгру(имя1, имя2 string, строки, столбцы int) (model.Игра, error) {
	var игра model.Игра
	данные := struct {
		Игрок1  model.Игрок `json:"игрок1"`
		Игрок2  model.Игрок `json:"игрок2"`
		Строки  int         `json:"строки"`
		Столбцы int         `json:"столбцы"`
	}{
		Игрок1:  *model.НовыйИгрок(имя1, model.Белые),
		Игрок2:  *model.НовыйИгрок(имя2, model.Чёрные),
		Строки:  строки,
		Столбцы: столбцы,
	}

	тело, err := json.Marshal(данные)
	if err != nil {
		return игра, err
	}

	resp, err := httpClient.Post(fmt.Sprintf("%s/api/games", *serverAddr), "application/json", bytes.NewReader(тело))
	if err != nil {
		return игра, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var ответ map[string]string
		json.NewDecoder(resp.Body).Decode(&ответ)
		return игра, fmt.Errorf("Сервер вернул статус %d: %s", resp.StatusCode, ответ["ошибка"])
	}
	if err := json.NewDecoder(resp.Body).Decode(&игра); err != nil {
		return игра, err
	}
	return игра, nil
}

func сделатьХод(id int, ход model.Ход) (model.Игра, error) {
	var игра model.Игра
	данные := struct {
		ОтСтрока    int             `json:"отСтрока"`
		ОтСтолбец   int             `json:"отСтолбец"`
		ВСтрока     int             `json:"вСтрока"`
		ВСтолбец    int             `json:"вСтолбец"`
		Превращение model.ТипФигуры `json:"превращение,omitempty"`
	}{
		ОтСтрока:    ход.ОтСтрока,
		ОтСтолбец:   ход.ОтСтолбец,
		ВСтрока:     ход.ВСтрока,
		ВСтолбец:    ход.ВСтолбец,
		Превращение: ход.Превращение,
	}

	тело, err := json.Marshal(данные)
	if err != nil {
		return игра, err
	}

	resp, err := httpClient.Post(
		fmt.Sprintf("%s/api/games/%d/move", *serverAddr, id),
		"application/json",
		bytes.NewReader(тело),
	)
	if err != nil {
		return игра, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var ответ map[string]string
		json.NewDecoder(resp.Body).Decode(&ответ)
		return игра, fmt.Errorf("Сервер вернул статус %d: %s", resp.StatusCode, ответ["ошибка"])
	}
	if err := json.NewDecoder(resp.Body).Decode(&игра); err != nil {
		return игра, err
	}
	return игра, nil
}

func сделатьАвтоХод(id int) (model.Игра, bool, bool, model.ЦветФигуры, error) {
	var ответ ОтветАвтоХода

	resp, err := httpClient.Post(
		fmt.Sprintf("%s/api/games/%d/auto-move", *serverAddr, id),
		"application/json",
		nil,
	)
	if err != nil {
		return model.Игра{}, false, false, model.Белые, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ошибка map[string]string
		json.NewDecoder(resp.Body).Decode(&ошибка)
		сообщение := ошибка["Ошибка"]
		if сообщение == "" {
			сообщение = ошибка["ошибка"]
		}
		return model.Игра{}, false, false, model.Белые, fmt.Errorf("сервер вернул статус %d: %s", resp.StatusCode, сообщение)
	}

	if err := json.NewDecoder(resp.Body).Decode(&ответ); err != nil {
		return model.Игра{}, false, false, model.Белые, err
	}
	return ответ.Игра, ответ.Мат, ответ.Пат, ответ.Победитель, nil
}

// ---------- Симуляции ----------

func запуститьСимуляции(n, строки, столбцы int) {
	muSim.Lock()
	defer muSim.Unlock()

	sims = nil
	simActive = true

	for range n {
		имя1 := случайноеИмя()
		имя2 := случайноеИмя()

		игра, err := создатьИгру(имя1, имя2, строки, столбцы)
		if err != nil {
			fmt.Printf("Ошибка создания игры для симуляции: %v\n", err)
			continue
		}

		сим := &Симуляция{
			ID:   игра.ID(),
			Игра: игра,
		}
		sims = append(sims, сим)

		go симуляцияПартии(сим)
	}
}

func остановитьСимуляции() {
	muSim.Lock()
	defer muSim.Unlock()
	for _, сим := range sims {
		сим.Завершена = true
	}
}

func случайноеИмя() string {
	if rand.Intn(2) == 0 {
		return model.ИменаМуж[rand.Intn(len(model.ИменаМуж))] + " " + model.ФамилииМуж[rand.Intn(len(model.ФамилииМуж))]
	}
	return model.ИменаЖен[rand.Intn(len(model.ИменаЖен))] + " " + model.ФамилииЖен[rand.Intn(len(model.ФамилииЖен))]
}

func симуляцияПартии(сим *Симуляция) {
	for {
		muSim.Lock()
		if сим.Завершена {
			muSim.Unlock()
			return
		}
		muSim.Unlock()

		начало := time.Now()
		обновлённая, мат, пат, победитель, err := сделатьАвтоХод(сим.ID)
		if err != nil {
			muSim.Lock()
			сим.ПоследнийХод = fmt.Sprintf("Ошибка: %v", err)
			сим.Завершена = true
			muSim.Unlock()
			return
		}

		muSim.Lock()
		сим.Игра = обновлённая
		сим.ВремяХода = time.Since(начало)

		if мат {
			победительЦвет := "Белые"
			if победитель == model.Чёрные {
				победительЦвет = "Чёрные"
			}
			сим.ПоследнийХод = fmt.Sprintf("Мат! Победили %s", победительЦвет)
			сим.Завершена = true
			сим.ВремяХода = time.Since(начало)
		} else if пат {
			сим.ПоследнийХод = "Пат. Ничья"
			сим.Завершена = true
			сим.ВремяХода = time.Since(начало)
		} else {
			if len(обновлённая.Ходы()) > 0 {
				сим.ПоследнийХод = model.ФорматироватьХод(обновлённая.Ходы()[len(обновлённая.Ходы())-1])
			}
		}
		muSim.Unlock()

		if мат || пат {
			return
		}

		time.Sleep(time.Duration(2000+rand.Intn(2000)) * time.Millisecond)

		muSim.Lock()
		сим.ВремяХода = time.Since(начало)
		muSim.Unlock()
	}
}

func запуститьОтрисовщикСимуляций() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if !simActive {
				continue
			}
			muSim.Lock()
			fmt.Print("\033[H\033[2J")

			var активные []*Симуляция
			for _, сим := range sims {
				if !сим.Завершена {
					активные = append(активные, сим)
				}
			}

			if len(активные) > 0 {
				примерСтроки := активные[0].Игра.СтрокиОтображенияБезИстории()
				ширинаДоски := 0
				if len(примерСтроки) > 0 {
					ширинаДоски = model.ВидимаяДлина(примерСтроки[0])
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
							w := model.ВидимаяДлина(стр)
							if w > ширины[i] {
								ширины[i] = w
							}
						}
					}
					for стр := 0; стр < максСтрок; стр++ {
						for i := 0; i < len(группа); i++ {
							if стр < len(доски[i]) {
								текст := доски[i][стр]
								fmt.Print(model.ДополнитьПробелами(текст, ширины[i]) + "   ")
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
	}()
}
