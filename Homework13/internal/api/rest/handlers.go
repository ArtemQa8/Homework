package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"mod.go/internal/auth"
	"mod.go/internal/config"
	"mod.go/internal/dto"
	"mod.go/internal/model"
	"mod.go/internal/repository"
	"mod.go/internal/service"
)

// Handlers - структура со всеми зависимостями RTST-хэндлеров
type Handlers struct {
	storage       *repository.Storage
	cfg           *config.Config
	playerService *service.PlayerService
	gameService   *service.GameService
	moveService   *service.MoveService
}

func New(
	storage *repository.Storage,
	cfg *config.Config,
	playerService *service.PlayerService,
	gameService *service.GameService,
	moveService *service.MoveService,
) *Handlers {
	return &Handlers{
		storage:       storage,
		cfg:           cfg,
		playerService: playerService,
		gameService:   gameService,
		moveService:   moveService,
	}
}

// Login выполняет авторизацию и выдаёт JWT.
// @Summary      Войти
// @Description  Проверяет логин и пароль, при успехе возвращает JWT-токен.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dto.LoginRequest   true  "Логин и пароль"
// @Success      200          {object}  dto.LoginResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /api/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	if req.Login != h.cfg.Login || req.Password != h.cfg.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"Ошибка": "неверный логин или пароль"})
		return
	}

	token, err := auth.GenerateToken(req.Login, h.cfg.JWTSecret, h.cfg.JWTTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "не удалось создать токен"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

// =============PLAYERS=============

// CreatePlayer создаёт нового игрока.
// @Summary      Создать игрока
// @Description  Создаёт нового игрока и возвращает его с присвоенным ID.
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        player  body      dto.CreatePlayerRequest  true  "Данные игрока"
// @Success      201     {object}  dto.PlayerResponse
// @Failure      400     {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/players [post]
func (h *Handlers) CreatePlayer(c *gin.Context) {
	var req dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	created, err := h.playerService.Create(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetPlayers возвращает список всех игроков.
// @Summary      Получить всех игроков
// @Description  Возвращает массив игроков из хранилища.
// @Tags         players
// @Produce      json
// @Success      200 {array}  dto.PlayerResponse
// @Failure      500 {object} map[string]string
// @Router       /api/players [get]
func (h *Handlers) GetPlayers(c *gin.Context) {
	players := h.playerService.List()
	c.JSON(http.StatusOK, players)
}

// GetPlayerByID возвращает игрока по ID.
// @Summary      Получить игрока по ID
// @Description  Возвращает одного игрока по его идентификатору.
// @Tags         players
// @Produce      json
// @Param        id   path      int  true  "ID игрока"
// @Success      200  {object}  dto.PlayerResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/players/{id} [get]
func (h *Handlers) GetPlayerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	player, err := h.playerService.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

// UpdatePlayer обновляет игрока по ID.
// @Summary      Обновить игрока
// @Description  Обновляет имя/цвет игрока с указанным ID.
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id      path      int                 true  "ID игрока"
// @Param        player  body      dto.UpdatePlayerRequest  true  "Новые данные игрока"
// @Success      200     {object}  dto.PlayerResponse
// @Failure      400     {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/players/{id} [put]
func (h *Handlers) UpdatePlayer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	var req dto.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}
	updated, err := h.playerService.Update(id, req.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeletePlayer удаляет игрока по ID.
// @Summary      Удалить игрока
// @Description  Удаляет игрока с указанным ID.
// @Tags         players
// @Param        id   path  int  true  "ID игрока"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/players/{id} [delete]
func (h *Handlers) DeletePlayer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	if err := h.playerService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// =============GAMES=============

// CreateGame создаёт новую игру.
// @Summary      Создать игру
// @Description  Создаёт новую игру с указанными игроками и размерами доски.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        game  body      dto.CreateGameRequest  true  "Данные игры"
// @Success      201   {object}  dto.GameResponse
// @Failure      400   {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/games [post]
func (h *Handlers) CreateGame(c *gin.Context) {
	var req dto.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	p1 := model.NewPlayer(req.Player1.Name)
	if req.Player1.ID != 0 {
		p1.SetID(req.Player1.ID)
	}
	p2 := model.NewPlayer(req.Player2.Name)
	if req.Player2.ID != 0 {
		p2.SetID(req.Player2.ID)
	}

	created, err := h.gameService.Create(*p1, *p2, req.Rows, req.Cols)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetGames возвращает список всех игр.
// @Summary      Получить все игры
// @Description  Возвращает массив игр из хранилища.
// @Tags         games
// @Produce      json
// @Success      200  {array}   dto.GameResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/games [get]
func (h *Handlers) GetGames(c *gin.Context) {
	games := h.gameService.List()
	c.JSON(http.StatusOK, games)
}

// GetGameByID возвращает игру по ID.
// @Summary      Получить игру по ID
// @Description  Возвращает одну игру по её идентификатору.
// @Tags         games
// @Produce      json
// @Param        id   path      int  true  "ID игры"
// @Success      200  {object}  dto.GameResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/games/{id} [get]
func (h *Handlers) GetGameByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	game, err := h.gameService.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "игра не найдена"})
		return
	}
	c.JSON(http.StatusOK, game)
}

// UpdateGame обновляет игру по ID.
// @Summary      Обновить игру
// @Description  Обновляет игроков и/или размеры игры с указанным ID.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "ID игры"
// @Param        game  body      dto.UpdateGameRequest  true  "Новые данные игры"
// @Success      200   {object}  dto.GameResponse
// @Failure      400   {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/games/{id} [put]
func (h *Handlers) UpdateGame(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	var req dto.UpdateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	p1 := model.NewPlayer(req.Player1.Name)
	if req.Player1.ID != 0 {
		p1.SetID(req.Player1.ID)
	}
	p2 := model.NewPlayer(req.Player2.Name)
	if req.Player2.ID != 0 {
		p2.SetID(req.Player2.ID)
	}

	updated, err := h.gameService.Update(id, *p1, *p2)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteGame удаляет игру по ID.
// @Summary      Удалить игру
// @Description  Удаляет игру с указанным ID.
// @Tags         games
// @Param        id   path  int  true  "ID игры"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/games/{id} [delete]
func (h *Handlers) DeleteGame(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	if err := h.gameService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// MakeMove выполняет ход в игре.
// @Summary      Сделать ход
// @Description  Выполняет ход в игре с указанным ID, проверяя правила шахмат.
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "ID игры"
// @Param        move  body      dto.MakeMoveRequest  true  "Координаты хода"
// @Success      200   {object}  dto.GameResponse
// @Failure      400   {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/games/{id}/move [post]
func (h *Handlers) MakeMove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}

	var req dto.MakeMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	move := &model.Move{
		FromRow: req.FromRow,
		FromCol: req.FromCol,
		ToRow:   req.ToRow,
		ToCol:   req.ToCol,
	}
	if req.Promotion != nil {
		move.Promotion = *req.Promotion
	}

	updated, err := h.gameService.MakeMove(id, move)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// AutoMove выполняет случайный ход в игре.
// @Summary      Случайный ход
// @Description  Выбирает и выполняет случайный легальный ход за текущего игрока.
// @Tags         games
// @Produce      json
// @Param        id   path      int  true  "ID игры"
// @Success      200  {object}  dto.AutoMoveResponse
// @Failure      400  {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/games/{id}/auto-move [post]
func (h *Handlers) AutoMove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}

	updated, err := h.gameService.AutoMove(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AutoMoveResponse{Game: toGameResponse(updated)})
}

// =============MOVES=============

// CreateMove создаёт новый ход.
// @Summary      Создать ход
// @Description  Сохраняет ход в хранилище (без применения к доске).
// @Tags         moves
// @Accept       json
// @Produce      json
// @Param        move  body      dto.CreateMoveRequest  true  "Данные хода"
// @Success      201   {object}  dto.MoveResponse
// @Failure      400   {object}  map[string]string
// @Failure		 401 	 {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/moves [post]
func (h *Handlers) CreateMove(c *gin.Context) {
	var req dto.CreateMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	move := model.Move{
		FromRow: req.FromRow,
		FromCol: req.FromCol,
		ToRow:   req.ToRow,
		ToCol:   req.ToCol,
	}
	if req.Promotion != nil {
		move.Promotion = *req.Promotion
	}
	move.SetGameID(req.GameID)

	created, err := h.moveService.Create(move)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetMoves возвращает список всех ходов.
// @Summary      Получить все ходы
// @Description  Возвращает массив всех ходов из хранилища.
// @Tags         moves
// @Produce      json
// @Success      200  {array}   dto.MoveResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/moves [get]
func (h *Handlers) GetMoves(c *gin.Context) {
	moves := h.moveService.List()
	c.JSON(http.StatusOK, moves)
}

// GetMoveByID возвращает ход по ID.
// @Summary      Получить ход по ID
// @Description  Возвращает один ход по его идентификатору.
// @Tags         moves
// @Produce      json
// @Param        id   path      int  true  "ID хода"
// @Success      200  {object}  dto.MoveResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/moves/{id} [get]
func (h *Handlers) GetMoveByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	move, err := h.moveService.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.JSON(http.StatusOK, move)
}

// UpdateMove обновляет ход по ID.
// @Summary      Обновить ход
// @Description  Обновляет данные хода с указанным ID.
// @Tags         moves
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "ID хода"
// @Param        move  body      dto.UpdateMoveRequest  true  "Новые данные хода"
// @Success      200   {object}  dto.MoveResponse
// @Failure      400   {object}  map[string]string
// @Failure		 401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/moves/{id} [put]
func (h *Handlers) UpdateMove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}

	var req dto.UpdateMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	move := model.Move{
		FromRow: req.FromRow,
		FromCol: req.FromCol,
		ToRow:   req.ToRow,
		ToCol:   req.ToCol,
	}
	move.SetGameID(req.GameID)
	move.SetID(id)

	updated, err := h.moveService.Update(id, move)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteMove удаляет ход по ID.
// @Summary      Удалить ход
// @Description  Удаляет ход с указанным ID.
// @Tags         moves
// @Param        id   path  int  true  "ID хода"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure		 401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/moves/{id} [delete]
func (h *Handlers) DeleteMove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный ID"})
		return
	}
	if err := h.moveService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// =============WEB PAGES=============

func (h *Handlers) IndexPage(c *gin.Context) {
	games := h.storage.GetAllGames()
	var builder strings.Builder

	builder.WriteString("<!DOCTYPE html><html><head>")
	builder.WriteString(`<meta charset="utf-8"><title>Текущие игры</title>`)
	builder.WriteString(`<meta http-equiv="refresh" content="3">`)
	builder.WriteString(`<style>body{font-family:Arial;margin:20px}table{border-collapse:collapse;width:auto}td,th{border:1px solid #ccc;padding:4px 8px;white-space:nowrap}th{background:#f0f0f0}a{color:#06c;text-decoration:none}a:hover{text-decoration:underline}</style>`)
	builder.WriteString("</head><body><h1>Текущие игры</h1>")

	if len(games) == 0 {
		builder.WriteString("<p>Нет активных игр.</p>")
	} else {
		builder.WriteString("<table><tr><th>ID игры</th><th>Белые</th><th>Чёрные</th><th>Текущий ход</th><th>Ходов</th><th>Последний ход</th><th>Результат</th></tr>")
		for _, game := range games {
			builder.WriteString("<tr>")
			fmt.Fprintf(&builder, "<td><a href=\"/game?id=%d\">%d</a></td>", game.ID(), game.ID())
			fmt.Fprintf(&builder, "<td>%s</td>", game.Player1().Name())
			fmt.Fprintf(&builder, "<td>%s</td>", game.Player2().Name())
			fmt.Fprintf(&builder, "<td>%s</td>", colorToString(game.CurrentColor()))
			fmt.Fprintf(&builder, "<td>%d</td>", len(game.Moves()))

			if len(game.Moves()) > 0 {
				last := game.Moves()[len(game.Moves())-1]
				fmt.Fprintf(&builder, "<td>%s</td>", model.FormatMove(last))
			} else {
				builder.WriteString("<td>-</td>")
			}

			if game.Board() != nil {
				if game.Checkmate(game.CurrentColor()) {
					winnerColor := "Белые"
					if game.CurrentColor() == model.White {
						winnerColor = "Чёрные"
					}
					fmt.Fprintf(&builder, "<td>Мат. Победили %s</td>", winnerColor)
				} else if game.Stalemate(game.CurrentColor()) {
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
}

func (h *Handlers) GamePage(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный ID игры")
		return
	}
	game, found := h.storage.GetGameByID(id)
	if !found {
		c.String(http.StatusNotFound, "Игра не найдена")
		return
	}
	if game.Board() == nil {
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

	fmt.Fprintf(&builder, "<p>Белые: %s | Чёрные: %s</p>", game.Player1().Name(), game.Player2().Name())
	fmt.Fprintf(&builder, "<p>Текущий ход: %s</p>", colorToString(game.CurrentColor()))

	builder.WriteString("<table class=\"chess\">")
	for r := 0; r < game.Board().Rows(); r++ {
		builder.WriteString("<tr>")
		for c := 0; c < game.Board().Cols(); c++ {
			piece := game.Board().PieceAt(r, c)
			cellColor := "white"
			if (r+c)%2 != 0 {
				cellColor = "black"
			}
			fmt.Fprintf(&builder, `<td class="cell %s">%s</td>`, cellColor, pieceSymbol(piece))
		}
		builder.WriteString("</tr>")
	}
	builder.WriteString("</table>")

	if len(game.Moves()) > 0 {
		last := game.Moves()[len(game.Moves())-1]
		fmt.Fprintf(&builder, "<p>Последний ход: %s</p>", model.FormatMove(last))
	}

	if game.Checkmate(game.CurrentColor()) {
		winnerColor := "Белые"
		if game.CurrentColor() == model.White {
			winnerColor = "Чёрные"
		}
		fmt.Fprintf(&builder, "<p>Результат: Мат. Победили %s</p>", winnerColor)
	} else if game.Stalemate(game.CurrentColor()) {
		builder.WriteString("<p>Результат: Пат</p>")
	} else {
		builder.WriteString("<p>Партия продолжается</p>")
	}

	builder.WriteString(`<p><a href="/">← Назад к списку</a></p>`)
	builder.WriteString("</body></html>")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(builder.String()))
}

// =============HELPERS=============

func colorToString(color model.Color) string {
	if color == model.White {
		return "Белые"
	}
	return "Чёрные"
}

func pieceSymbol(p *model.Piece) string {
	if p == nil {
		return ""
	}
	return p.Symbol()
}

func toGameResponse(game model.Game) dto.GameResponse {
	var board *dto.BoardResponse
	if game.Board() != nil {
		cells := make([][]*dto.PieceResponse, game.Board().Rows())
		for r := 0; r < game.Board().Rows(); r++ {
			cells[r] = make([]*dto.PieceResponse, game.Board().Cols())
			for c := 0; c < game.Board().Cols(); c++ {
				if p := game.Board().PieceAt(r, c); p != nil {
					cells[r][c] = &dto.PieceResponse{
						Color: colorToString(p.Color()),
						Type:  pieceTypeToString(p.Type()),
					}
				}
			}
		}
		board = &dto.BoardResponse{
			Rows:  game.Board().Rows(),
			Cols:  game.Board().Cols(),
			Cells: cells,
		}
	}

	moves := make([]dto.MoveResponse, len(game.Moves()))
	for i, m := range game.Moves() {
		var captured *dto.PieceResponse
		if m.Captured != nil {
			captured = &dto.PieceResponse{
				Color: colorToString(m.Captured.Color()),
				Type:  pieceTypeToString(m.Captured.Type()),
			}
		}
		var movedPiece *dto.PieceResponse
		if m.MovedPiece != nil {
			movedPiece = &dto.PieceResponse{
				Color: colorToString(m.MovedPiece.Color()),
				Type:  pieceTypeToString(m.MovedPiece.Type()),
			}
		}
		moves[i] = dto.MoveResponse{
			GameID:     m.GameID(),
			ID:         m.ID(),
			FromRow:    m.FromRow,
			FromCol:    m.FromCol,
			ToRow:      m.ToRow,
			ToCol:      m.ToCol,
			Captured:   captured,
			Promotion:  pieceTypeToString(m.Promotion),
			MovedPiece: movedPiece,
			Check:      m.Check,
			Mate:       m.Mate,
		}
	}

	return dto.GameResponse{
		ID: game.ID(),
		Player1: dto.GamePlayerResponse{
			ID:    game.Player1().ID(),
			Name:  game.Player1().Name(),
			Color: colorToString(game.Player1Color()),
		},
		Player2: dto.GamePlayerResponse{
			ID:    game.Player2().ID(),
			Name:  game.Player2().Name(),
			Color: colorToString(game.Player2Color()),
		},
		Board:   board,
		Current: colorToString(game.CurrentColor()),
		Moves:   moves,
	}
}

func pieceTypeToString(t model.PieceType) string {
	switch t {
	case model.Pawn:
		return "Пешка"
	case model.Rook:
		return "Ладья"
	case model.Knight:
		return "Конь"
	case model.Bishop:
		return "Слон"
	case model.Queen:
		return "Ферзь"
	case model.King:
		return "Король"
	}
	return ""
}
