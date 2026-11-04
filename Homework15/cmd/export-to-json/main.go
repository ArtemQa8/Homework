package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mod.go/internal/config"
	"mod.go/internal/repository/postgres"
)

func main() {
	fmt.Println("=== НАЧАЛО ЭКСПОРТА ===")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка конфига: %v", err)
	}

	db, err := postgres.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close(context.Background())

	players := db.GetAllPlayers()
	games := db.GetAllGames()
	moves := db.GetAllMoves()
	counters, _ := db.GetAllCounters()

	exportDir := "export_data"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Fatalf("Ошибка создания папки: %v", err)
	}

	// Сохраняем игроков
	err = saveToJSON(exportDir+"/players.json", players)
	if err != nil {
		log.Fatalf("Ошибка экспорта игроков: %v", err)
	}
	fmt.Printf("Экспортировано игроков: %d\n", len(players))

	// Сохраняем игры
	err = saveToJSON(exportDir+"/games.json", games)
	if err != nil {
		log.Fatalf("Ошибка экспорта игр: %v", err)
	}
	fmt.Printf("Экспортировано игр: %d\n", len(games))

	// Сохраняем ходы
	err = saveToJSON(exportDir+"/moves.json", moves)
	if err != nil {
		log.Fatalf("Ошибка экспорта ходов: %v", err)
	}
	fmt.Printf("Экспортировано ходов: %d\n", len(moves))

	// Сохраняем счетчики в формате filestore (next*ID),
	nextPlayerID, nextGameID, nextMoveID := 1, 1, 1
	for _, c := range counters {
		switch c.Name {
		case "player_id":
			nextPlayerID = int(c.Seq) + 1
		case "game_id":
			nextGameID = int(c.Seq) + 1
		case "move_id":
			nextMoveID = int(c.Seq) + 1
		}
	}
	countersPayload := map[string]int{
		"nextPlayerID": nextPlayerID,
		"nextGameID":   nextGameID,
		"nextMoveID":   nextMoveID,
	}
	err = saveToJSON(exportDir+"/counters.json", countersPayload)
	if err != nil {
		log.Fatalf("Ошибка экспорта счетчиков: %v", err)
	}
	fmt.Println("Экспортировано счетчиков: 3")

	fmt.Println("=== ЭКСПОРТ УСПЕШНО ЗАВЕРШЕН ===")
	fmt.Printf("Данные сохранены в папку '%s'\n", exportDir)
}

func saveToJSON(filename string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}
