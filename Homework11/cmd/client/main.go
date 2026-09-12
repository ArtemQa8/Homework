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
	sims      []*Simulation
)

type Simulation struct {
	ID       int
	Game     model.Game
	Finished bool
	LastMove string
	MoveTime time.Duration
}

type AutoMoveResponse struct {
	Game      model.Game  `json:"игра"`
	Mate      bool        `json:"мат,omitempty"`
	Stalemate bool        `json:"пат,omitempty"`
	Winner    model.Color `json:"победитель,omitempty"`
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

	if *gameIDFlag > 0 {
		if err := gameSession(*gameIDFlag, inputChan); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for {
		fmt.Println("Добро пожаловать в шахматный клиент!")
		fmt.Println()
		fmt.Printf("Наблюдение за играми: %s/\n", *serverAddr)
		fmt.Println()
		fmt.Println("1. Создать новую игру")
		fmt.Println("2. Подключиться к существующей")
		fmt.Println("3. Симуляция")
		fmt.Println("0. Выход")
		fmt.Print("Выберите действие: ")

		choice, ok := <-inputChan
		if !ok {
			fmt.Println("Ввод закрыт.")
			return
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := createNewGame(inputChan); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}
		case "2":
			if err := connectToGame(inputChan); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}
		case "3":
			simulationFromMenu(inputChan)
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный ввод, попробуйте ещё раз.")
		}
	}
}

func readString(inputChan chan string) (string, bool) {
	s, ok := <-inputChan
	return strings.TrimSpace(s), ok
}

func readInt(inputChan chan string) (int, bool) {
	s, ok := readString(inputChan)
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func createNewGame(inputChan chan string) error {
	fmt.Print("Имя первого игрока (белые): ")
	name1, ok := readString(inputChan)
	if !ok {
		return fmt.Errorf("ввод закрыт")
	}

	fmt.Print("Имя второго игрока (чёрные): ")
	name2, ok := readString(inputChan)
	if !ok {
		return fmt.Errorf("ввод закрыт")
	}

	fmt.Print("Количество строк (0 для стандартных 8): ")
	rows, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверное число строк, используются 8")
		rows = 0
	}

	fmt.Print("Количество столбцов (0 для стандартных 8): ")
	cols, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверное число столбцов, используются 8")
		cols = 0
	}

	game, err := createGame(name1, name2, rows, cols)
	if err != nil {
		return err
	}
	fmt.Printf("Игра создана, ID: %d\n", game.ID())
	return gameSession(game.ID(), inputChan)
}

func connectToGame(inputChan chan string) error {
	fmt.Print("Введите ID игры (или 'список' для показа всех): ")
	input, ok := readString(inputChan)
	if !ok {
		return fmt.Errorf("ввод закрыт")
	}

	if strings.ToLower(input) == "список" {
		games, err := fetchGames()
		if err != nil {
			return err
		}
		if len(games) == 0 {
			fmt.Println("Нет доступных игр.")
			return nil
		}
		fmt.Println("Доступные игры:")
		for _, game := range games {
			fmt.Printf("ID: %d | %s (белые) - %s (чёрные)\n",
				game.ID(), game.Player1().Name(), game.Player2().Name())
		}
		fmt.Print("Введите ID игры: ")
		input, ok = readString(inputChan)
		if !ok {
			return fmt.Errorf("ввод закрыт")
		}
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("неверный ID")
	}
	return gameSession(id, inputChan)
}

func simulationFromMenu(inputChan chan string) {
	fmt.Print("Количество симуляций: ")
	n, ok := readInt(inputChan)
	if !ok || n <= 0 {
		fmt.Println("Неверное число симуляций.")
		return
	}

	fmt.Print("Количество строк (0 для стандартных 8): ")
	rows, _ := readInt(inputChan)
	fmt.Print("Количество столбцов (0 для стандартных 8): ")
	cols, _ := readInt(inputChan)

	if rows <= 0 || cols <= 0 {
		rows, cols = 8, 8
	}

	startSimulations(n, rows, cols)
	simActive = true
	startSimulationsRenderer()

	fmt.Println("Симуляции запущены. Введите 'стоп' для остановки.")
	for {
		input, ok := readString(inputChan)
		if !ok {
			break
		}
		if strings.ToLower(input) == "стоп" || strings.ToLower(input) == "выход" {
			stopSimulations()
			simActive = false
			fmt.Print("\033[H\033[2J")
			return
		}
	}
}

func renderGame(game model.Game, id int) {
	fmt.Print("\033[H\033[2J")
	rendered := game.Render()
	idx := strings.Index(rendered, "История ходов:")
	if idx == -1 {
		fmt.Printf("ID игры: %d\n%s", id, rendered)
	} else {
		fmt.Print(rendered[:idx])
		fmt.Printf("ID игры: %d\n\n", id)
		fmt.Print(rendered[idx:])
	}
}

func showGameCommands(id int) {
	fmt.Println("Доступные команды:")
	fmt.Println("  автоход <N>     - сделать N автоматических ходов")
	fmt.Println("  симуляция <N>   - запустить N фоновых партий")
	fmt.Println("  стоп            - остановить симуляции (если запущены)")
	fmt.Println("  help / команды  - показать этот список")
	fmt.Println("  выход / exit    - выйти из игры")
	fmt.Printf("  Наблюдать за игрой: %s/game?id=%d\n", *serverAddr, id)
}

func printPrompt(game model.Game) {
	color := game.CurrentColor()
	var name string
	if color == model.White {
		name = game.Player1().Name()
	} else {
		name = game.Player2().Name()
	}
	fmt.Printf("\nХод %s (%s): ", model.GenitiveColor(color), name)
	fmt.Println("Введите ход (например: e2 e4) или 'help':")
}

func askPromotion(inputChan chan string) (model.PieceType, bool) {
	for {
		fmt.Print("Во что превратить пешку? (ферзь, ладья, конь, слон): ")
		input, ok := readString(inputChan)
		if !ok {
			return 0, false
		}
		pt, err := model.ParsePieceType(input)
		if err != nil {
			fmt.Println("Ошибка:", err)
			continue
		}
		switch pt {
		case model.Queen, model.Rook, model.Bishop, model.Knight:
			return pt, true
		default:
			fmt.Println("Можно превратить только в ферзя, ладью, слона или коня.")
		}
	}
}

func gameSession(id int, inputChan chan string) error {
	game, err := fetchGame(id)
	if err != nil {
		return err
	}

	renderGame(game, id)
	showGameCommands(id)
	printPrompt(game)

	var muGame sync.Mutex

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			if simActive {
				continue
			}
			fresh, err := fetchGame(id)
			if err != nil {
				continue
			}

			muGame.Lock()
			changed := len(fresh.Moves()) != len(game.Moves()) || fresh.CurrentColor() != game.CurrentColor()
			if changed {
				game = fresh
				renderGame(game, id)
				printPrompt(game)
			}
			muGame.Unlock()
		}
	}()

	for input := range inputChan {
		if input == "" {
			continue
		}

		muGame.Lock()

		if strings.ToLower(input) == "help" || strings.ToLower(input) == "команды" || strings.ToLower(input) == "помощь" {
			showGameCommands(id)
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" ||
			strings.ToLower(input) == "q" || strings.ToLower(input) == "й" ||
			strings.ToLower(input) == "йгше" || strings.ToLower(input) == "учше" ||
			strings.ToLower(input) == "выход" {
			fmt.Println("Выход из игры.")
			muGame.Unlock()
			return nil
		}

		if simActive {
			if strings.ToLower(input) == "стоп" {
				stopSimulations()
				simActive = false
				renderGame(game, id)
				printPrompt(game)
				muGame.Unlock()
				continue
			}
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		if strings.HasPrefix(strings.ToLower(input), "автоход") {
			parts := strings.Fields(input)
			if len(parts) < 2 {
				fmt.Println("Укажите количество ходов: автоход <число>")
				muGame.Unlock()
				printPrompt(game)
				continue
			}
			n, err := strconv.Atoi(parts[1])
			if err != nil || n <= 0 {
				fmt.Println("Неверное число ходов.")
				muGame.Unlock()
				printPrompt(game)
				continue
			}
			for range n {
				updated, mate, stalemate, winner, err := makeAutoMove(id)
				if err != nil {
					fmt.Println("Ошибка:", err)
					break
				}
				game = updated

				if mate {
					winnerColor := "Белые"
					if winner == model.Black {
						winnerColor = "Чёрные"
					}
					renderGame(game, id)
					printPrompt(game)
					fmt.Printf("\nМат! Победили %s\n", winnerColor)
					muGame.Unlock()
					return nil
				} else if stalemate {
					renderGame(game, id)
					printPrompt(game)
					fmt.Println("\nПат. Ничья")
					muGame.Unlock()
					return nil
				}

				renderGame(game, id)
				printPrompt(game)
				time.Sleep(time.Duration(2000+rand.Intn(2000)) * time.Millisecond)
			}
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		if strings.HasPrefix(strings.ToLower(input), "симуляция") {
			parts := strings.Fields(input)
			if len(parts) < 2 {
				fmt.Println("Укажите количество досок: симуляция <число>")
				muGame.Unlock()
				printPrompt(game)
				continue
			}
			n, err := strconv.Atoi(parts[1])
			if err != nil || n <= 0 {
				fmt.Println("Неверное число.")
				muGame.Unlock()
				printPrompt(game)
				continue
			}
			startSimulations(n, game.Board().Rows(), game.Board().Cols())
			simActive = true
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		parts := strings.Fields(input)
		if len(parts) != 2 {
			fmt.Println("Нужно указать две клетки, например: e2 e4")
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		fromRow, fromCol, err := model.ParseCell(parts[0], game.Board())
		if err != nil {
			fmt.Println("Ошибка:", err)
			muGame.Unlock()
			printPrompt(game)
			continue
		}
		toRow, toCol, err := model.ParseCell(parts[1], game.Board())
		if err != nil {
			fmt.Println("Ошибка:", err)
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		move := model.Move{
			FromRow: fromRow,
			FromCol: fromCol,
			ToRow:   toRow,
			ToCol:   toCol,
		}

		// Если это пешка, которая может дойти до последнего ряда — спросить фигуру.
		// Спрашиваем ТОЛЬКО если ход легален по правилам пешки: при нелегальном
		// (например, e2→e8) сервер вернёт «пешка так не ходит», а не запрос превращения.
		if piece := game.Board().PieceAt(fromRow, fromCol); piece != nil && piece.Type() == model.Pawn {
			lastRow := 0
			if piece.Color() == model.White {
				lastRow = game.Board().Rows() - 1
			}
			if toRow == lastRow {
				if rules := model.RulesFor(piece); rules != nil && rules.CanMove(&move, game.Board()) {
					promo, ok := askPromotion(inputChan)
					if !ok {
						muGame.Unlock()
						printPrompt(game)
						continue
					}
					move.Promotion = promo
				}
			}
		}

		updated, err := makeMove(id, move)
		if err != nil {
			fmt.Println("Ошибка:", err)
			muGame.Unlock()
			printPrompt(game)
			continue
		}

		game = updated
		renderGame(game, id)
		printPrompt(game)

		if game.Checkmate(game.CurrentColor()) {
			winnerColor := "Белые"
			if game.CurrentColor() == model.White {
				winnerColor = "Чёрные"
			}
			fmt.Printf("\nМат! Победили %s\n", winnerColor)
			muGame.Unlock()
			return nil
		} else if game.Stalemate(game.CurrentColor()) {
			fmt.Println("\nПат. Ничья")
			muGame.Unlock()
			return nil
		}

		muGame.Unlock()
	}
	return nil
}

// ---------- HTTP ----------

func fetchGame(id int) (model.Game, error) {
	var game model.Game
	resp, err := httpClient.Get(fmt.Sprintf("%s/api/games/%d", *serverAddr, id))
	if err != nil {
		return game, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return game, fmt.Errorf("Сервер вернул статус %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return game, err
	}
	return game, nil
}

func fetchGames() ([]model.Game, error) {
	var games []model.Game
	resp, err := httpClient.Get(fmt.Sprintf("%s/api/games", *serverAddr))
	if err != nil {
		return games, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return games, fmt.Errorf("Сервер вернул статус %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return games, err
	}
	return games, nil
}

func createGame(name1, name2 string, rows, cols int) (model.Game, error) {
	var game model.Game
	data := struct {
		Player1 model.Player `json:"игрок1"`
		Player2 model.Player `json:"игрок2"`
		Rows    int          `json:"строки"`
		Cols    int          `json:"столбцы"`
	}{
		Player1: *model.NewPlayer(name1),
		Player2: *model.NewPlayer(name2),
		Rows:    rows,
		Cols:    cols,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return game, err
	}

	resp, err := httpClient.Post(fmt.Sprintf("%s/api/games", *serverAddr), "application/json", bytes.NewReader(body))
	if err != nil {
		return game, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return game, fmt.Errorf("Сервер вернул статус %d: %s", resp.StatusCode, errResp["Ошибка"])
	}
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return game, err
	}
	return game, nil
}

func makeMove(id int, move model.Move) (model.Game, error) {
	var game model.Game
	data := struct {
		FromRow   int              `json:"отСтрока"`
		FromCol   int              `json:"отСтолбец"`
		ToRow     int              `json:"вСтрока"`
		ToCol     int              `json:"вСтолбец"`
		Promotion *model.PieceType `json:"превращение,omitempty"`
	}{
		FromRow: move.FromRow,
		FromCol: move.FromCol,
		ToRow:   move.ToRow,
		ToCol:   move.ToCol,
	}
	if move.Promotion != 0 {
		p := move.Promotion
		data.Promotion = &p
	}

	body, err := json.Marshal(data)
	if err != nil {
		return game, err
	}

	resp, err := httpClient.Post(
		fmt.Sprintf("%s/api/games/%d/move", *serverAddr, id),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return game, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return game, fmt.Errorf("Сервер вернул статус %d: %s", resp.StatusCode, errResp["Ошибка"])
	}
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return game, err
	}
	return game, nil
}

func makeAutoMove(id int) (model.Game, bool, bool, model.Color, error) {
	var resp struct {
		Game      model.Game  `json:"игра"`
		Mate      bool        `json:"мат,omitempty"`
		Stalemate bool        `json:"пат,omitempty"`
		Winner    model.Color `json:"победитель,omitempty"`
	}

	httpResp, err := httpClient.Post(
		fmt.Sprintf("%s/api/games/%d/auto-move", *serverAddr, id),
		"application/json",
		nil,
	)
	if err != nil {
		return model.Game{}, false, false, model.White, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(httpResp.Body).Decode(&errResp)
		msg := errResp["Ошибка"]
		if msg == "" {
			msg = errResp["ошибка"]
		}
		return model.Game{}, false, false, model.White, fmt.Errorf("сервер вернул статус %d: %s", httpResp.StatusCode, msg)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return model.Game{}, false, false, model.White, err
	}
	return resp.Game, resp.Mate, resp.Stalemate, resp.Winner, nil
}

// ---------- Симуляции ----------

func startSimulations(n, rows, cols int) {
	muSim.Lock()
	defer muSim.Unlock()

	sims = nil
	simActive = true

	for range n {
		name1 := randomName()
		name2 := randomName()

		game, err := createGame(name1, name2, rows, cols)
		if err != nil {
			fmt.Printf("Ошибка создания игры для симуляции: %v\n", err)
			continue
		}

		sim := &Simulation{
			ID:   game.ID(),
			Game: game,
		}
		sims = append(sims, sim)

		go simulateGame(sim)
	}
}

func stopSimulations() {
	muSim.Lock()
	defer muSim.Unlock()
	for _, sim := range sims {
		sim.Finished = true
	}
}

func randomName() string {
	if rand.Intn(2) == 0 {
		return model.MaleNames[rand.Intn(len(model.MaleNames))] + " " + model.MaleSurnames[rand.Intn(len(model.MaleSurnames))]
	}
	return model.FemaleNames[rand.Intn(len(model.FemaleNames))] + " " + model.FemaleSurnames[rand.Intn(len(model.FemaleSurnames))]
}

func simulateGame(sim *Simulation) {
	for {
		muSim.Lock()
		if sim.Finished {
			muSim.Unlock()
			return
		}
		muSim.Unlock()

		start := time.Now()
		updated, mate, stalemate, winner, err := makeAutoMove(sim.ID)
		if err != nil {
			muSim.Lock()
			sim.LastMove = fmt.Sprintf("Ошибка: %v", err)
			sim.Finished = true
			muSim.Unlock()
			return
		}

		muSim.Lock()
		sim.Game = updated

		if mate {
			winnerColor := "Белые"
			if winner == model.Black {
				winnerColor = "Чёрные"
			}
			sim.LastMove = fmt.Sprintf("Мат! Победили %s", winnerColor)
			sim.Finished = true
			sim.MoveTime = time.Since(start)
		} else if stalemate {
			sim.LastMove = "Пат. Ничья"
			sim.Finished = true
			sim.MoveTime = time.Since(start)
		} else {
			if len(updated.Moves()) > 0 {
				sim.LastMove = model.FormatMove(updated.Moves()[len(updated.Moves())-1])
			}
		}
		muSim.Unlock()

		if mate || stalemate {
			return
		}

		time.Sleep(time.Duration(2000+rand.Intn(2000)) * time.Millisecond)

		muSim.Lock()
		sim.MoveTime = time.Since(start)
		muSim.Unlock()
	}
}

func startSimulationsRenderer() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if !simActive {
				continue
			}
			muSim.Lock()
			fmt.Print("\033[H\033[2J")

			var active []*Simulation
			for _, sim := range sims {
				if !sim.Finished {
					active = append(active, sim)
				}
			}

			if len(active) > 0 {
				sampleLines := active[0].Game.RenderLinesWithoutHistory()
				boardWidth := 0
				if len(sampleLines) > 0 {
					boardWidth = model.VisibleLength(sampleLines[0])
				}
				gap := 3
				perRow := 3
				maxTermWidth := 120
				if boardWidth*3+gap*2 > maxTermWidth {
					perRow = 2
					if boardWidth*2+gap > maxTermWidth {
						perRow = 1
					}
				}
				for start := 0; start < len(active); start += perRow {
					end := start + perRow
					if end > len(active) {
						end = len(active)
					}
					group := active[start:end]
					boards := make([][]string, len(group))
					maxRows := 0
					for i, sim := range group {
						lines := sim.Game.RenderLinesWithoutHistory()
						boards[i] = lines
						if len(lines) > maxRows {
							maxRows = len(lines)
						}
					}
					widths := make([]int, len(group))
					for i, b := range boards {
						for _, line := range b {
							w := model.VisibleLength(line)
							if w > widths[i] {
								widths[i] = w
							}
						}
					}
					for row := 0; row < maxRows; row++ {
						for i := 0; i < len(group); i++ {
							if row < len(boards[i]) {
								text := boards[i][row]
								fmt.Print(model.PadWithSpaces(text, widths[i]) + "   ")
							} else {
								fmt.Print(strings.Repeat(" ", widths[i]) + "   ")
							}
						}
						fmt.Println()
					}
					fmt.Println()
				}
			}

			for i, sim := range sims {
				name1 := sim.Game.Player1().Name()
				name2 := sim.Game.Player2().Name()
				if sim.Finished {
					fmt.Printf("#%d %s vs %s: %s (завершена)\n", i+1, name1, name2, sim.LastMove)
				} else {
					fmt.Printf("#%d %s vs %s: %s за %v\n", i+1, name1, name2, sim.LastMove, sim.MoveTime.Round(time.Millisecond))
				}
			}
			muSim.Unlock()
		}
	}()
}
