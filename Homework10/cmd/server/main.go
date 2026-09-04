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

	// =============ИГРОКИ=============
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

	// GET /api/players
	роутер.GET("/api/players", func(c *gin.Context) {
		игроки := хранилище.ПолучитьВсехИгроков()
		c.JSON(http.StatusOK, игроки)
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

	// =============ИГРЫ=============
	// POST /api/games
	роутер.POST("/api/games", func(c *gin.Context) {
		var вход struct {
			Игрок1  model.Игрок `json:"игрок1"`
			Игрок2  model.Игрок `json:"игрок2"`
			Строки  int         `json:"строки"`
			Столбцы int         `json:"столбцы"`
		}
		if err := c.ShouldBindJSON(&вход); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		игра := model.Игра{}
		игра.УстановитьИгрок1(вход.Игрок1)
		игра.УстановитьИгрок2(вход.Игрок2)

		if вход.Строки <= 0 || вход.Столбцы <= 0 {
			вход.Строки, вход.Столбцы = 8, 8
		}
		игра.УстановитьДоску(model.НоваяДоска(вход.Строки, вход.Столбцы))

		созданная, err := хранилище.СоздатьИгру(игра)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		игрок1 := созданная.Игрок1()
		игрок2 := созданная.Игрок2()
		игрок1.УстановитьЦвет(model.Белые)
		игрок2.УстановитьЦвет(model.Чёрные)
		созданная.УстановитьИгрок1(игрок1)
		созданная.УстановитьИгрок2(игрок2)

		if err := хранилище.ПерезаписатьИгру(созданная.ID(), созданная); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "не удалось сохранить игру"})
			return
		}

		c.JSON(http.StatusCreated, созданная)
	})

	// GET /api/games
	роутер.GET("/api/games", func(c *gin.Context) {
		игры := хранилище.ПолучитьВсеИгры()
		c.JSON(http.StatusOK, игры)
	})

	// GET /api/games/:id
	роутер.GET("/api/games/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}
		игра, найден := хранилище.ПолучитьИгруПоАйди(id)
		if !найден {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "игра не найдена"})
			return
		}
		c.JSON(http.StatusOK, игра)
	})

	// PUT /api/games/:id
	роутер.PUT("/api/games/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}

		var игра model.Игра
		if err := c.ShouldBindJSON(&игра); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		if err := хранилище.ОбновитьИгру(id, игра); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": err.Error()})
			return
		}

		обновлённая, _ := хранилище.ПолучитьИгруПоАйди(id)
		c.JSON(http.StatusOK, обновлённая)
	})

	// DELETE /api/games/:id
	роутер.DELETE("/api/games/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}
		if err := хранилище.УдалитьИгру(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	роутер.POST("/api/games/:id/move", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}

		var вход struct {
			ОтСтрока    int             `json:"отСтрока"`
			ОтСтолбец   int             `json:"отСтолбец"`
			ВСтрока     int             `json:"вСтрока"`
			ВСтолбец    int             `json:"вСтолбец"`
			Превращение model.ТипФигуры `json:"превращение,omitempty"`
		}

		if err := c.ShouldBindJSON(&вход); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		игра, найдена := хранилище.ПолучитьИгруПоАйди(id)
		if !найдена {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "игра не найдена"})
			return
		}

		ход := model.Ход{
			ОтСтрока:    вход.ОтСтрока,
			ОтСтолбец:   вход.ОтСтолбец,
			ВСтрока:     вход.ВСтрока,
			ВСтолбец:    вход.ВСтолбец,
			Превращение: вход.Превращение,
		}
		ход.УстановитьИграID(id)

		if err := игра.СделатьХод(&ход); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		сохранённыйХод, err := хранилище.СоздатьХод(ход)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "не удалось сохранить ход"})
			return
		}

		игра.УстановитьIDПоследнегоХода(сохранённыйХод.ID())

		if err := хранилище.ПерезаписатьИгру(id, игра); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ошибка": "не удалось сохранить игру"})
			return
		}

		c.JSON(http.StatusOK, игра)
	})

	// =============ХОДЫ=============
	// POST /api/moves
	роутер.POST("/api/moves", func(c *gin.Context) {
		var ход model.Ход
		if err := c.ShouldBindJSON(&ход); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}
		созданный, err := хранилище.СоздатьХод(ход)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, созданный)
	})

	// GET /api/moves
	роутер.GET("/api/moves", func(c *gin.Context) {
		ходы := хранилище.ПолучитьВсеХоды()
		c.JSON(http.StatusOK, ходы)
	})

	// GET /api/moves/:id
	роутер.GET("/api/moves/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}
		ход, найден := хранилище.ПолучитьХодПоАйди(id)
		if !найден {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "ход не найден"})
			return
		}
		c.JSON(http.StatusOK, ход)
	})

	// PUT /api/moves/:id
	роутер.PUT("/api/moves/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}

		_, найден := хранилище.ПолучитьХодПоАйди(id)
		if !найден {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": "ход не найден"})
			return
		}

		var ход model.Ход
		if err := c.ShouldBindJSON(&ход); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": err.Error()})
			return
		}

		if err := хранилище.ОбновитьХод(id, ход); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"ошибка": err.Error()})
			return
		}

		обновлённый, _ := хранилище.ПолучитьХодПоАйди(id)
		c.JSON(http.StatusOK, обновлённый)
	})

	// DELETE /api/moves/:id
	роутер.DELETE("/api/moves/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ошибка": "неверный ID"})
			return
		}
		if err := хранилище.УдалитьХод(id); err != nil {
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
