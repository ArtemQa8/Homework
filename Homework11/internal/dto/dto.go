package dto

// CreateGameRequest описывает тело запроса для создания новой игры.
type CreateGameRequest struct {
	Player1 PlayerResponse `json:"игрок1"`
	Player2 PlayerResponse `json:"игрок2"`
	Rows    int            `json:"строки"`
	Cols    int            `json:"столбцы"`
}

// MakeMoveRequest описывает тело запроса для выполнения хода.
type MakeMoveRequest struct {
	FromRow   int    `json:"отСтрока"`
	FromCol   int    `json:"отСтолбец"`
	ToRow     int    `json:"вСтрока"`
	ToCol     int    `json:"вСтолбец"`
	Promotion string `json:"превращение,omitempty"`
}

// PlayerResponse описывает JSON-ответ с данными игрока.
type PlayerResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"имя"`
	Color string `json:"цвет"`
}

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

// GameResponse описывает JSON-ответ с данными игры.
type GameResponse struct {
	ID      int            `json:"id"`
	Player1 PlayerResponse `json:"игрок1"`
	Player2 PlayerResponse `json:"игрок2"`
	Board   *BoardResponse `json:"доска"`
	Current string         `json:"текущий"`
	Moves   []MoveResponse `json:"ходы"`
}

// AutoMoveResponse описывает результат автоматического хода.
type AutoMoveResponse struct {
	Game      GameResponse `json:"игра"`
	Mate      bool         `json:"мат,omitempty"`
	Stalemate bool         `json:"пат,omitempty"`
	Winner    string       `json:"победитель,omitempty"`
}
