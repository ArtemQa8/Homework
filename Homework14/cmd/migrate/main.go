package main

import (
	"context"
	"fmt"
	"log"

	"mod.go/internal/config"
	"mod.go/internal/repository/filestore"
	"mod.go/internal/repository/mongodb"
)

func main() {
	fmt.Println("=== НАЧАЛО МИГРАЦИИ ===")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка конфига: %v", err)
	}

	// Читаем Джисончики
	fs := filestore.NewStorage("data")
	if err := fs.LoadFromFiles(); err != nil {
		log.Fatalf("Ошибка чтения JSON: %v", err)
	}

	// Подключаемся к Mongo
	db, err := mongodb.NewMongoStorage(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
	defer db.Close(context.Background())

	// Переносим ИГРОКОВ
	players := fs.GetAllPlayers()
	maxPlayerID := 0
	for _, p := range players {
		err := db.ImportPlayer(p.ID(), p.Name())
		if err != nil {
			log.Printf("Пропущен игрок %d: %v\n", p.ID(), err)
		}
		if p.ID() > maxPlayerID {
			maxPlayerID = p.ID()
		}
	}
	fmt.Printf("Перенесено игроков: %d (макс. ID: %d)\n", len(players), maxPlayerID)

	// Переносим ИГРЫ
	games := fs.GetAllGames()
	maxGameID := 0
	for _, g := range games {
		err := db.ImportGame(g)
		if err != nil {
			log.Printf("Пропущена игра %d: %v\n", g.ID(), err)
		}
		if g.ID() > maxGameID {
			maxGameID = g.ID()
		}
	}
	fmt.Printf("Перенесено игр: %d (макс. ID: %d)\n", len(games), maxGameID)

	// Переносим ХОДЫ
	moves := fs.GetAllMoves()
	maxMoveID := 0
	for _, m := range moves {
		err := db.ImportMove(m)
		if err != nil {
			log.Printf("Пропущен ход %d: %v\n", m.ID(), err)
		}
		if m.ID() > maxMoveID {
			maxMoveID = m.ID()
		}
	}
	fmt.Printf("Перенесено ходов: %d (макс. ID: %d)\n", len(moves), maxMoveID)

	// Обновляем счетчики, чтобы новые записи не перезаписали старые
	_ = db.SetCounter("player_id", maxPlayerID)
	_ = db.SetCounter("game_id", maxGameID)
	_ = db.SetCounter("move_id", maxMoveID)

	fmt.Println("=== МИГРАЦИЯ УСПЕШНО ЗАВЕРШЕНА ===")
}
