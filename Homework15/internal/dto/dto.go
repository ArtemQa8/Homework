package dto

import "mod.go/internal/model"

// ============ PLAYERS ============

// PlayerResponse описывает JSON-ответ с данными игрока (вне партии).
type PlayerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"имя"`
}

// GamePlayerResponse описывает игрока в контексте партии — с цветом.
// Используется внутри GameResponse.
type GamePlayerResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"имя"`
	Color string `json:"цвет"`
}

// CreatePlayerRequest описывает тело запроса на создание игрока.
type CreatePlayerRequest struct {
	Name string `json:"имя"`
}

// UpdatePlayerRequest описывает тело запроса на обновление игрока.
type UpdatePlayerRequest struct {
	Name string `json:"имя"`
}

// ============ GAMES ============

// GamePlayerName — имя игрока, которое клиент присылает при создании/обновлении игры.
type GamePlayerName struct {
	ID   int    `json:"id"`
	Name string `json:"имя"`
}

// CreateGameRequest описывает тело запроса для создания новой игры.
type CreateGameRequest struct {
	Player1 GamePlayerName `json:"игрок1" extensions:"x-order=0"`
	Player2 GamePlayerName `json:"игрок2" extensions:"x-order=1"`
	Rows    int            `json:"строки" extensions:"x-order=2"`
	Cols    int            `json:"столбцы" extensions:"x-order=3"`
}

// UpdateGameRequest описывает тело запроса на обновление игры
type UpdateGameRequest struct {
	Player1 GamePlayerName `json:"игрок1" extensions:"x-order=0"`
	Player2 GamePlayerName `json:"игрок2" extensions:"x-order=1"`
}

// GameResponse описывает JSON-ответ с данными игры.
type GameResponse struct {
	ID      int                `json:"id"`
	Player1 GamePlayerResponse `json:"игрок1"`
	Player2 GamePlayerResponse `json:"игрок2"`
	Board   *BoardResponse     `json:"доска"`
	Current string             `json:"текущий"`
	Moves   []MoveResponse     `json:"ходы"`
}

// AutoMoveResponse описывает результат автоматического хода.
type AutoMoveResponse struct {
	Game      GameResponse `json:"игра"`
	Mate      bool         `json:"мат,omitempty"`
	Stalemate bool         `json:"пат,omitempty"`
	Winner    string       `json:"победитель,omitempty"`
}

// ============ MOVES ============

// MoveResponse описывает JSON-ответ с данными хода.
type MoveResponse struct {
	GameID     int            `json:"играID"`
	ID         int            `json:"id"`
	FromRow    int            `json:"отСтрока"`
	FromCol    int            `json:"отСтолбец"`
	ToRow      int            `json:"вСтрока"`
	ToCol      int            `json:"вСтолбец"`
	Captured   *PieceResponse `json:"съедена,omitempty"`
	Promotion  string         `json:"превращение,omitempty"`
	MovedPiece *PieceResponse `json:"ходившаяФигура,omitempty"`
	Check      bool           `json:"шах,omitempty"`
	Mate       bool           `json:"мат,omitempty"`
}

// CreateMoveRequest описывает тело запроса на создание хода.
type CreateMoveRequest struct {
	GameID    int              `json:"играID" extensions:"x-order=0"`
	FromRow   int              `json:"отСтрока" extensions:"x-order=1"`
	FromCol   int              `json:"отСтолбец" extensions:"x-order=2"`
	ToRow     int              `json:"вСтрока" extensions:"x-order=3"`
	ToCol     int              `json:"вСтолбец" extensions:"x-order=4"`
	Promotion *model.PieceType `json:"превращение,omitempty" swaggertype:"string" extensions:"x-order=5"`
}

// UpdateMoveRequest описывает тело запроса на обновление хода.
type UpdateMoveRequest struct {
	GameID  int `json:"играID" extensions:"x-order=0"`
	FromRow int `json:"отСтрока" extensions:"x-order=1"`
	FromCol int `json:"отСтолбец" extensions:"x-order=2"`
	ToRow   int `json:"вСтрока" extensions:"x-order=3"`
	ToCol   int `json:"вСтолбец" extensions:"x-order=4"`
}

// MakeMoveRequest описывает тело запроса для выполнения хода в игре.
type MakeMoveRequest struct {
	FromRow   int              `json:"отСтрока" extensions:"x-order=0"`
	FromCol   int              `json:"отСтолбец" extensions:"x-order=1"`
	ToRow     int              `json:"вСтрока" extensions:"x-order=2"`
	ToCol     int              `json:"вСтолбец" extensions:"x-order=3"`
	Promotion *model.PieceType `json:"превращение,omitempty" swaggertype:"string" extensions:"x-order=4"`
}

// ============ BOARD & PIECES ============

// PieceResponse описывает JSON-ответ с данными фигуры.
type PieceResponse struct {
	Color string `json:"цвет"`
	Type  string `json:"тип"`
}

// BoardResponse описывает JSON-ответ с доской.
type BoardResponse struct {
	Rows  int                `json:"строки"`
	Cols  int                `json:"столбцы"`
	Cells [][]*PieceResponse `json:"клетки"`
}

// ============ AUTH ============

// LoginRequest описывает тело запроса на логин.
type LoginRequest struct {
	Login    string `json:"логин"`
	Password string `json:"пароль"`
}

// LoginResponse описывает ответ при успешном логине.
type LoginResponse struct {
	Token string `json:"токен"`
}
