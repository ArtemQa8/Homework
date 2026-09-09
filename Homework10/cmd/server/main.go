package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"mod.go/internal/model"
	"mod.go/internal/repository"
	"mod.go/internal/service"
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

	// POST /api/games/:id/move
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

	// POST /api/games/:id/auto-move
	роутер.POST("/api/games/:id/auto-move", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
			return
		}

		игра, найдена := хранилище.ПолучитьИгруПоАйди(id)
		if !найдена {
			c.JSON(http.StatusNotFound, gin.H{"Ошибка": "игра не найдена"})
			return
		}

		цвет := игра.ТекущийЦвет()

		// Проверка Мата
		if игра.Мат(цвет) {
			победитель := model.Белые
			if цвет == model.Белые {
				победитель = model.Чёрные
			}
			c.JSON(http.StatusOK, gin.H{
				"игра":       игра,
				"мат":        true,
				"победитель": победитель,
			})
			return
		}

		// Проверка Пата
		if игра.Пат(цвет) {
			c.JSON(http.StatusOK, gin.H{
				"игра": игра,
				"пат":  true,
			})
			return
		}

		ход, err := service.ВыбратьСлучайныйХод(&игра)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
			return
		}

		// Превращение пешки в ферзя
		if фигура := игра.Доска().ФигураНа(ход.ОтСтрока, ход.ОтСтолбец); фигура != nil &&
			фигура.Тип() == model.Пешка {
			последняяСтрока := 0
			if фигура.Цвет() == model.Белые {
				последняяСтрока = игра.Доска().Строки() - 1
			}
			if ход.ВСтрока == последняяСтрока {
				ход.Превращение = model.Ферзь
			}
		}

		ход.УстановитьИграID(id)
		if err := игра.СделатьХод(ход); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
			return
		}

		сохранённыйХод, err := хранилище.СоздатьХод(*ход)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "не удалось сохранить ход"})
			return
		}
		игра.УстановитьIDПоследнегоХода(сохранённыйХод.ID())

		if err := хранилище.ПерезаписатьИгру(id, игра); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "не удалось сохранить игру"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"игра": игра})
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

	// Отображение страницы наблюдателя
	роутер.GET("/", func(c *gin.Context) {
		игры := хранилище.ПолучитьВсеИгры()
		var builder strings.Builder

		builder.WriteString("<!DOCTYPE html><html><head>")
		builder.WriteString(`<meta charset="utf-8"><title>Текущие игры</title>`)
		builder.WriteString(`<meta http-equiv="refresh" content="3">`)
		builder.WriteString(`<style>body{font-family:Arial;margin:20px}table{border-collapse:collapse;width:auto}td,th{border:1px solid #ccc;padding:4px 8px;white-space:nowrap}th{background:#f0f0f0}a{color:#06c;text-decoration:none}a:hover{text-decoration:underline}</style>`)
		builder.WriteString("</head><body><h1>Текущие игры</h1>")

		if len(игры) == 0 {
			builder.WriteString("<p>Нет активных игр.</p>")
		} else {
			builder.WriteString("<table><tr><th>ID игры</th><th>Белые</th><th>Чёрные</th><th>Текущий ход</th><th>Ходов</th><th>Последний ход</th><th>Результат</th></tr>")
			for _, игра := range игры {
				fmt.Fprintf(&builder, "<tr>")
				fmt.Fprintf(&builder, "<td><a href=\"/game?id=%d\">%d</a></td>", игра.ID(), игра.ID())
				fmt.Fprintf(&builder, "<td>%s</td>", игра.Игрок1().Имя())
				fmt.Fprintf(&builder, "<td>%s</td>", игра.Игрок2().Имя())
				fmt.Fprintf(&builder, "<td>%s</td>", модельЦветаСтрокой(игра.ТекущийЦвет()))
				fmt.Fprintf(&builder, "<td>%d</td>", len(игра.Ходы()))

				if len(игра.Ходы()) > 0 {
					последний := игра.Ходы()[len(игра.Ходы())-1]
					fmt.Fprintf(&builder, "<td>%s</td>", model.ФорматироватьХод(последний))
				} else {
					builder.WriteString("<td>-</td>")
				}

				if игра.Доска() != nil {
					if игра.Мат(игра.ТекущийЦвет()) {
						победительЦвет := "Белые"
						if игра.ТекущийЦвет() == model.Белые {
							победительЦвет = "Чёрные"
						}
						fmt.Fprintf(&builder, "<td>Мат. Победили %s</td>", победительЦвет)
					} else if игра.Пат(игра.ТекущийЦвет()) {
						builder.WriteString("<td>Пат</td>")
					} else {
						builder.WriteString("<td>Идёт</td>")
					}
				} else {
					builder.WriteString("<td>Нет доски</td>")
				}
				builder.WriteString("</tr>")
			}
			builder.WriteString("</table>")
		}

		builder.WriteString("</body></html>")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(builder.String()))
	})

	// Страница игры
	роутер.GET("/game", func(c *gin.Context) {
		idStr := c.Query("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Неверный ID игры")
			return
		}
		игра, найдена := хранилище.ПолучитьИгруПоАйди(id)
		if !найдена {
			c.String(http.StatusNotFound, "Игра не найдена")
			return
		}

		if игра.Доска() == nil {
			c.String(http.StatusInternalServerError, "У игры отсутствует доска")
			return
		}

		var builder strings.Builder
		builder.WriteString("<!DOCTYPE html><html><head>")
		fmt.Fprintf(&builder, `<meta charset="utf-8"><title>Игра #%d</title>`, id)
		builder.WriteString(`<meta http-equiv="refresh" content="2">`)
		builder.WriteString(`<style>body{font-family:Arial;margin:20px}table.chess{border-collapse:collapse;margin:10px 0}td.cell{width:40px;height:40px;text-align:center;font-size:24px;border:1px solid #999}.white{background:#f0d9b5}.black{background:#b58863}</style>`)
		builder.WriteString("</head><body>")
		fmt.Fprintf(&builder, "<h1>Игра #%d</h1>", id)

		fmt.Fprintf(&builder, "<p>Белые: %s | Чёрные: %s</p>", игра.Игрок1().Имя(), игра.Игрок2().Имя())
		fmt.Fprintf(&builder, "<p>Текущий ход: %s</p>", модельЦветаСтрокой(игра.ТекущийЦвет()))

		builder.WriteString("<table class=\"chess\">")
		for стр := 0; стр < игра.Доска().Строки(); стр++ {
			builder.WriteString("<tr>")
			for стл := 0; стл < игра.Доска().Столбцы(); стл++ {
				фигура := игра.Доска().ФигураНа(стр, стл)
				клеткаЦвет := "white"
				if (стр+стл)%2 != 0 {
					клеткаЦвет = "black"
				}
				fmt.Fprintf(&builder, `<td class="cell %s">%s</td>`, клеткаЦвет, символФигуры(фигура))
			}
			builder.WriteString("</tr>")
		}
		builder.WriteString("</table>")

		if len(игра.Ходы()) > 0 {
			последний := игра.Ходы()[len(игра.Ходы())-1]
			fmt.Fprintf(&builder, "<p>Последний ход: %s</p>", model.ФорматироватьХод(последний))
		}

		if игра.Мат(игра.ТекущийЦвет()) {
			победительЦвет := "Белые"
			if игра.ТекущийЦвет() == model.Белые {
				победительЦвет = "Чёрные"
			}
			fmt.Fprintf(&builder, "<p>Результат: Мат. Победили %s</p>", победительЦвет)
		} else if игра.Пат(игра.ТекущийЦвет()) {
			builder.WriteString("<p>Результат: Пат</p>")
		} else {
			builder.WriteString("<p>Партия продолжается</p>")
		}

		builder.WriteString(`<p><a href="/">← Назад к списку</a></p>`)
		builder.WriteString("</body></html>")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(builder.String()))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: роутер,
	}

	go func() {
		fmt.Println("Сервер запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	сигналы := make(chan os.Signal, 1)
	signal.Notify(сигналы, syscall.SIGINT, syscall.SIGTERM)
	<-сигналы

	fmt.Println("\nПолучен сигнал остановки, завершаем работу...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при graceful shutdown: %v", err)
	}

	if err := хранилище.СохранитьВсе(); err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
	}

	fmt.Println("Сервер остановлен, данные сохранены.")
}

func модельЦветаСтрокой(цвет model.ЦветФигуры) string {
	if цвет == model.Белые {
		return "Белые"
	}
	return "Чёрные"
}

func символФигуры(ф *model.Фигура) string {
	if ф == nil {
		return ""
	}
	return ф.Символ()
}
