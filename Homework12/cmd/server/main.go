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

	// ============ Открытые роуты (без авторизации) ============

	// Логин — доступен всем, иначе не залогиниться
	router.POST("/api/login", Login)

	// Просмотр (GET) — открыт для наблюдателей и клиента
	router.GET("/api/players", GetPlayers)
	router.GET("/api/players/:id", GetPlayerByID)
	router.GET("/api/games", GetGames)
	router.GET("/api/games/:id", GetGameByID)
	router.GET("/api/moves", GetMoves)
	router.GET("/api/moves/:id", GetMoveByID)

	// Веб-страницы наблюдателя
	router.GET("/", IndexPage)
	router.GET("/game", GamePage)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ============ Защищённые роуты (требуется JWT) ============

	protected := router.Group("/api")
	protected.Use(AuthRequired())
	{
		// Игроки
		protected.POST("/players", CreatePlayer)
		protected.PUT("/players/:id", UpdatePlayer)
		protected.DELETE("/players/:id", DeletePlayer)

		// Игры
		protected.POST("/games", CreateGame)
		protected.PUT("/games/:id", UpdateGame)
		protected.DELETE("/games/:id", DeleteGame)
		protected.POST("/games/:id/move", MakeMove)
		protected.POST("/games/:id/auto-move", AutoMove)

		// Ходы
		protected.POST("/moves", CreateMove)
		protected.PUT("/moves/:id", UpdateMove)
		protected.DELETE("/moves/:id", DeleteMove)
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
