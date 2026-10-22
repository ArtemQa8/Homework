package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "mod.go/internal/proto"
)

var (
	conn      *grpc.ClientConn
	playerCli pb.PlayerServiceClient
	gameCli   pb.GameServiceClient
	moveCli   pb.MoveServiceClient
)

func main() {
	var err error
	conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Ошибка подключения: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	playerCli = pb.NewPlayerServiceClient(conn)
	gameCli = pb.NewGameServiceClient(conn)
	moveCli = pb.NewMoveServiceClient(conn)

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

	fmt.Println("gRPC-клиент. Подключён к localhost:50051")

	for {
		printMainMenu()
		choice, ok := readString(inputChan)
		if !ok {
			return
		}

		switch choice {
		case "1":
			playersMenu(inputChan)
		case "2":
			gamesMenu(inputChan)
		case "3":
			movesMenu(inputChan)
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный ввод.")
		}
	}
}

// ============ Меню ============

func printMainMenu() {
	fmt.Println()
	fmt.Println("=== Главное меню ===")
	fmt.Println("  1. Игроки")
	fmt.Println("  2. Игры")
	fmt.Println("  3. Ходы")
	fmt.Println("  0. Выход")
	fmt.Print("Выберите раздел: ")
}

func printPlayersMenu() {
	fmt.Println()
	fmt.Println("--- Игроки ---")
	fmt.Println("  1. Создать игрока")
	fmt.Println("  2. Список игроков")
	fmt.Println("  3. Получить игрока по ID")
	fmt.Println("  4. Обновить игрока")
	fmt.Println("  5. Удалить игрока")
	fmt.Println("  0. Назад")
	fmt.Print("Выберите действие: ")
}

func printGamesMenu() {
	fmt.Println()
	fmt.Println("--- Игры ---")
	fmt.Println("  1. Создать игру")
	fmt.Println("  2. Список игр")
	fmt.Println("  3. Получить игру по ID")
	fmt.Println("  4. Обновить игру")
	fmt.Println("  5. Удалить игру")
	fmt.Println("  6. Сделать ход")
	fmt.Println("  7. Автоход")
	fmt.Println("  0. Назад")
	fmt.Print("Выберите действие: ")
}

func printMovesMenu() {
	fmt.Println()
	fmt.Println("--- Ходы ---")
	fmt.Println("  1. Создать ход")
	fmt.Println("  2. Список ходов")
	fmt.Println("  3. Получить ход по ID")
	fmt.Println("  4. Обновить ход")
	fmt.Println("  5. Удалить ход")
	fmt.Println("  0. Назад")
	fmt.Print("Выберите действие: ")
}

// ============ Подменю ============

func playersMenu(inputChan chan string) {
	for {
		printPlayersMenu()
		choice, ok := readString(inputChan)
		if !ok {
			return
		}
		switch choice {
		case "1":
			createPlayer(inputChan)
		case "2":
			listPlayers()
		case "3":
			getPlayer(inputChan)
		case "4":
			updatePlayer(inputChan)
		case "5":
			deletePlayer(inputChan)
		case "0":
			return
		default:
			fmt.Println("Неверный ввод.")
		}
	}
}

func gamesMenu(inputChan chan string) {
	for {
		printGamesMenu()
		choice, ok := readString(inputChan)
		if !ok {
			return
		}
		switch choice {
		case "1":
			createGame(inputChan)
		case "2":
			listGames()
		case "3":
			getGame(inputChan)
		case "4":
			updateGame(inputChan)
		case "5":
			deleteGame(inputChan)
		case "6":
			makeMove(inputChan)
		case "7":
			autoMove(inputChan)
		case "0":
			return
		default:
			fmt.Println("Неверный ввод.")
		}
	}
}

func movesMenu(inputChan chan string) {
	for {
		printMovesMenu()
		choice, ok := readString(inputChan)
		if !ok {
			return
		}
		switch choice {
		case "1":
			createMove(inputChan)
		case "2":
			listMoves()
		case "3":
			getMove(inputChan)
		case "4":
			updateMove(inputChan)
		case "5":
			deleteMove(inputChan)
		case "0":
			return
		default:
			fmt.Println("Неверный ввод.")
		}
	}
}

// ============ Helpers ============

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

func colorString(c pb.Color) string {
	switch c {
	case pb.Color_WHITE:
		return "Белые"
	case pb.Color_BLACK:
		return "Чёрные"
	default:
		return "?"
	}
}

func pieceTypeString(t pb.PieceType) string {
	switch t {
	case pb.PieceType_PAWN:
		return "Пешка"
	case pb.PieceType_ROOK:
		return "Ладья"
	case pb.PieceType_KNIGHT:
		return "Конь"
	case pb.PieceType_BISHOP:
		return "Слон"
	case pb.PieceType_QUEEN:
		return "Ферзь"
	case pb.PieceType_KING:
		return "Король"
	default:
		return ""
	}
}

func printGame(g *pb.Game) {
	fmt.Printf("Игра #%d\n", g.Id)
	fmt.Printf("  %s (белые) - %s (чёрные)\n", g.Player1.Name, g.Player2.Name)
	fmt.Printf("  Ход: %s\n", colorString(g.Current))
	fmt.Printf("  Ходов: %d\n", len(g.Moves))
	if g.Board != nil {
		fmt.Printf("  Доска: %d х %d\n", g.Board.Rows, g.Board.Cols)
	}
	if len(g.Moves) > 0 {
		last := g.Moves[len(g.Moves)-1]
		fmt.Printf("  Последний ход: (%d,%d) -> (%d,%d)\n",
			last.FromRow, last.FromCol, last.ToRow, last.ToCol)
		if last.Promotion != pb.PieceType_PIECE_TYPE_UNSPECIFIED {
			fmt.Printf("  Превращение: %s\n", pieceTypeString(last.Promotion))
		}
	}
}

func printMove(m *pb.Move) {
	fmt.Printf("Ход #%d\n", m.Id)
	fmt.Printf("  Игра: %d\n", m.GameId)
	fmt.Printf("  Откуда: (%d,%d)\n", m.FromRow, m.FromCol)
	fmt.Printf("  Куда:   (%d,%d)\n", m.ToRow, m.ToCol)
	if m.Promotion != pb.PieceType_PIECE_TYPE_UNSPECIFIED {
		fmt.Printf("  Превращение: %s\n", pieceTypeString(m.Promotion))
	}
	if m.Check {
		fmt.Println("  Шах")
	}
	if m.Mate {
		fmt.Println("  Мат")
	}
}

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// ============ Players ============

func createPlayer(inputChan chan string) {
	fmt.Print("Имя игрока: ")
	name, ok := readString(inputChan)
	if !ok {
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := playerCli.CreatePlayer(ctx, &pb.CreatePlayerRequest{Name: name})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Игрок создан: ID=%d, имя=%s\n", resp.Id, resp.Name)
}

func listPlayers() {
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := playerCli.ListPlayers(ctx, &pb.ListPlayersRequest{})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if len(resp.Players) == 0 {
		fmt.Println("Игроков нет.")
		return
	}
	for _, p := range resp.Players {
		fmt.Printf("  ID=%d, имя=%s\n", p.Id, p.Name)
	}
}

func getPlayer(inputChan chan string) {
	fmt.Print("ID игрока: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := playerCli.GetPlayer(ctx, &pb.GetPlayerRequest{Id: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("ID=%d, имя=%s\n", resp.Id, resp.Name)
}

func updatePlayer(inputChan chan string) {
	fmt.Print("ID игрока: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	fmt.Print("Новое имя: ")
	name, ok := readString(inputChan)
	if !ok {
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := playerCli.UpdatePlayer(ctx, &pb.UpdatePlayerRequest{
		Id:   int32(id),
		Name: name,
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Игрок обновлён: ID=%d, имя=%s\n", resp.Id, resp.Name)
}

func deletePlayer(inputChan chan string) {
	fmt.Print("ID игрока: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	_, err := playerCli.DeletePlayer(ctx, &pb.DeletePlayerRequest{Id: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Игрок удалён.")
}

// ============ Games ============

func createGame(inputChan chan string) {
	fmt.Print("Имя первого игрока (белые): ")
	n1, ok := readString(inputChan)
	if !ok {
		return
	}
	fmt.Print("Имя второго игрока (чёрные): ")
	n2, ok := readString(inputChan)
	if !ok {
		return
	}
	fmt.Print("Количество строк (0 для стандартных 8): ")
	rows, _ := readInt(inputChan)
	fmt.Print("Количество столбцов (0 для стандартных 8): ")
	cols, _ := readInt(inputChan)

	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := gameCli.CreateGame(ctx, &pb.CreateGameRequest{
		Player1Name: n1,
		Player2Name: n2,
		Rows:        int32(rows),
		Cols:        int32(cols),
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printGame(resp)
}

func listGames() {
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := gameCli.ListGames(ctx, &pb.ListGamesRequest{})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if len(resp.Games) == 0 {
		fmt.Println("Игр нет.")
		return
	}
	for _, g := range resp.Games {
		fmt.Printf("  ID=%d | %s (белые) - %s (чёрные) | ходов: %d | ход: %s\n",
			g.Id, g.Player1.Name, g.Player2.Name, len(g.Moves), colorString(g.Current))
	}
}

func getGame(inputChan chan string) {
	fmt.Print("ID игры: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := gameCli.GetGame(ctx, &pb.GetGameRequest{Id: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printGame(resp)
}

func updateGame(inputChan chan string) {
	fmt.Print("ID игры: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	fmt.Print("Новое имя первого игрока (Enter — не менять): ")
	n1, _ := readString(inputChan)
	fmt.Print("Новое имя второго игрока (Enter — не менять): ")
	n2, _ := readString(inputChan)

	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := gameCli.UpdateGame(ctx, &pb.UpdateGameRequest{
		Id:          int32(id),
		Player1Name: n1,
		Player2Name: n2,
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printGame(resp)
}

func deleteGame(inputChan chan string) {
	fmt.Print("ID игры: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	_, err := gameCli.DeleteGame(ctx, &pb.DeleteGameRequest{Id: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Игра удалена.")
}

// ============ Moves ============

func makeMove(inputChan chan string) {
	fmt.Print("ID игры: ")
	gameID, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID игры.")
		return
	}
	fmt.Print("Откуда (строка, столбец): ")
	fromRow, _ := readInt(inputChan)
	fromCol, _ := readInt(inputChan)
	fmt.Print("Куда (строка, столбец): ")
	toRow, _ := readInt(inputChan)
	toCol, _ := readInt(inputChan)

	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := gameCli.MakeMove(ctx, &pb.MakeMoveRequest{
		GameId:  int32(gameID),
		FromRow: int32(fromRow),
		FromCol: int32(fromCol),
		ToRow:   int32(toRow),
		ToCol:   int32(toCol),
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printGame(resp)
}

func autoMove(inputChan chan string) {
	fmt.Print("ID игры: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := gameCli.AutoMove(ctx, &pb.AutoMoveRequest{GameId: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printGame(resp)
}

func createMove(inputChan chan string) {
	fmt.Print("ID игры: ")
	gameID, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID игры.")
		return
	}
	fmt.Print("Откуда (строка, столбец): ")
	fromRow, _ := readInt(inputChan)
	fromCol, _ := readInt(inputChan)
	fmt.Print("Куда (строка, столбец): ")
	toRow, _ := readInt(inputChan)
	toCol, _ := readInt(inputChan)

	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := moveCli.CreateMove(ctx, &pb.CreateMoveRequest{
		GameId:  int32(gameID),
		FromRow: int32(fromRow),
		FromCol: int32(fromCol),
		ToRow:   int32(toRow),
		ToCol:   int32(toCol),
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printMove(resp)
}

func listMoves() {
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := moveCli.ListMoves(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	if len(resp.Moves) == 0 {
		fmt.Println("Ходов нет.")
		return
	}
	for _, m := range resp.Moves {
		fmt.Printf("  ID=%d | игра=%d | (%d,%d) -> (%d,%d)\n",
			m.Id, m.GameId, m.FromRow, m.FromCol, m.ToRow, m.ToCol)
	}
}

func getMove(inputChan chan string) {
	fmt.Print("ID хода: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := moveCli.GetMove(ctx, &pb.GetMoveRequest{Id: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printMove(resp)
}

func updateMove(inputChan chan string) {
	fmt.Print("ID хода: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	fmt.Print("ID игры: ")
	gameID, _ := readInt(inputChan)
	fmt.Print("Откуда (строка, столбец): ")
	fromRow, _ := readInt(inputChan)
	fromCol, _ := readInt(inputChan)
	fmt.Print("Куда (строка, столбец): ")
	toRow, _ := readInt(inputChan)
	toCol, _ := readInt(inputChan)

	ctx, cancel := withTimeout()
	defer cancel()
	resp, err := moveCli.UpdateMove(ctx, &pb.UpdateMoveRequest{
		Id:      int32(id),
		GameId:  int32(gameID),
		FromRow: int32(fromRow),
		FromCol: int32(fromCol),
		ToRow:   int32(toRow),
		ToCol:   int32(toCol),
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	printMove(resp)
}

func deleteMove(inputChan chan string) {
	fmt.Print("ID хода: ")
	id, ok := readInt(inputChan)
	if !ok {
		fmt.Println("Неверный ID.")
		return
	}
	ctx, cancel := withTimeout()
	defer cancel()
	_, err := moveCli.DeleteMove(ctx, &pb.DeleteMoveRequest{Id: int32(id)})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Ход удалён.")
}
