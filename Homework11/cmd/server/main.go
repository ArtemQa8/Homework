package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "mod.go/docs"

	"github.com/gin-gonic/gin"
	"mod.go/internal/config"
	"mod.go/internal/repository"
)

// @title           Шахматный API
// @version         1.0
// @description     REST API для шахматного сервера: игроки, игры, ходы и выполнение ходов.
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT-токен в формате: Bearer <token>
func main() {
	loadedCfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}
	cfg = loadedCfg

	storage = repository.NewStorage()
	if err := storage.LoadFromFiles(); err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
	}

	router := gin.Default()

	// Авторизация
	router.POST("/api/login", Login)

	// Игроки
	router.POST("/api/players", CreatePlayer)
	router.GET("/api/players", GetPlayers)
	router.GET("/api/players/:id", GetPlayerByID)
	router.PUT("/api/players/:id", UpdatePlayer)
	router.DELETE("/api/players/:id", DeletePlayer)

	// Игры
	router.POST("/api/games", CreateGame)
	router.GET("/api/games", GetGames)
	router.GET("/api/games/:id", GetGameByID)
	router.PUT("/api/games/:id", UpdateGame)
	router.DELETE("/api/games/:id", DeleteGame)
	router.POST("/api/games/:id/move", MakeMove)
	router.POST("/api/games/:id/auto-move", AutoMove)

	// Ходы
	router.POST("/api/moves", CreateMove)
	router.GET("/api/moves", GetMoves)
	router.GET("/api/moves/:id", GetMoveByID)
	router.PUT("/api/moves/:id", UpdateMove)
	router.DELETE("/api/moves/:id", DeleteMove)

	// Веб-страницы наблюдателя
	router.GET("/", IndexPage)
	router.GET("/game", GamePage)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Сервер запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	fmt.Println("\nПолучен сигнал остановки, завершаем работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при graceful shutdown: %v", err)
	}

	if err := storage.SaveAll(); err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
	}

	fmt.Println("Сервер остановлен, данные сохранены.")
}
