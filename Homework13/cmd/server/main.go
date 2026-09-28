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
	"mod.go/internal/api/rest"
	"mod.go/internal/config"
	"mod.go/internal/repository"
	"mod.go/internal/service"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	storage := repository.NewStorage()
	if err := storage.LoadFromFiles(); err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
	}

	playerService := service.NewPlayerService(storage)
	gameService := service.NewGameService(storage)
	moveService := service.NewMoveService(storage)

	handlers := rest.New(storage, cfg, playerService, gameService, moveService)

	router := gin.Default()

	// ============ Открытые роуты (без авторизации) ============

	// Логин — доступен всем, иначе не залогиниться
	router.POST("/api/login", handlers.Login)

	// Просмотр (GET) — открыт для наблюдателей и клиента
	router.GET("/api/players", handlers.GetPlayers)
	router.GET("/api/players/:id", handlers.GetPlayerByID)
	router.GET("/api/games", handlers.GetGames)
	router.GET("/api/games/:id", handlers.GetGameByID)
	router.GET("/api/moves", handlers.GetMoves)
	router.GET("/api/moves/:id", handlers.GetMoveByID)

	// Веб-страницы наблюдателя
	router.GET("/", handlers.IndexPage)
	router.GET("/game", handlers.GamePage)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ============ Защищённые роуты (требуется JWT) ============

	protected := router.Group("/api")
	protected.Use(handlers.AuthRequired())
	{
		// Игроки
		protected.POST("/players", handlers.CreatePlayer)
		protected.PUT("/players/:id", handlers.UpdatePlayer)
		protected.DELETE("/players/:id", handlers.DeletePlayer)

		// Игры
		protected.POST("/games", handlers.CreateGame)
		protected.PUT("/games/:id", handlers.UpdateGame)
		protected.DELETE("/games/:id", handlers.DeleteGame)
		protected.POST("/games/:id/move", handlers.MakeMove)
		protected.POST("/games/:id/auto-move", handlers.AutoMove)

		// Ходы
		protected.POST("/moves", handlers.CreateMove)
		protected.PUT("/moves/:id", handlers.UpdateMove)
		protected.DELETE("/moves/:id", handlers.DeleteMove)
	}

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
