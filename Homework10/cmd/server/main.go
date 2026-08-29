package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mod.go/internal/model"
	"mod.go/internal/repository"
)

func main() {
	хранилище := repository.НовоеХранилище()
	if err := хранилище.ЗагрузитьИзФайлов(); err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
	}

	роутер := gin.Default()

	// GET /api/players
	роутер.GET("/api/players", func(c *gin.Context) {
		игроки := хранилище.ПолучитьВсехИгроков()
		c.JSON(http.StatusOK, игроки)
	})

	// POST /api/players
	роутер.POST("/api/players", func(c *gin.Context) {
		var игрок model.Игрок
		if err := c.ShouldBindJSON(&игрок); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		if игрок.Имя() == "" {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "имя обязательно"})
			return
		}

		созданный := хранилище.СоздатьИгрока(игрок)
		c.JSON(http.StatusCreated, созданный)
	})

	// GET /api/players/:id
	роутер.GET("/api/players/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}

		игрок, найден := хранилище.ПолучитьИгрокаПоАйди(id)
		if !найден {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "игрок не найден"})
			return
		}

		c.JSON(http.StatusOK, игрок)
	})

	// PUT /api/players/:id
	роутер.PUT("/api/players/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}

		var игрок model.Игрок
		if err := c.ShouldBindJSON(&игрок); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		if игрок.Имя() == "" {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "имя обязательно"})
			return
		}

		err = хранилище.ОбновитьИгрока(id, игрок)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": err.Error()})
			return
		}

		обновлённый, _ := хранилище.ПолучитьИгрокаПоАйди(id)
		c.JSON(http.StatusOK, обновлённый)
	})

	// DELETE /api/players/:id
	роутер.DELETE("/api/players/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}

		err = хранилище.УдалитьИгрока(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	fmt.Println("Сервер запущен на :8080")
	if err := роутер.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
