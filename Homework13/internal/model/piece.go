package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Color int

const (
	White Color = iota
	Black
)

type PieceType int

const (
	Pawn PieceType = iota
	Rook
	Knight
	Bishop
	Queen
	King
)

const (
	Reset      = "\033[0m"
	WhitePiece = "\033[31m"
	BlackPiece = "\033[34m"
	WhiteSqBg  = "\033[48;5;250m"
	BlackSqBg  = "\033[48;5;240m"
)

type Piece struct {
	color     Color
	pieceType PieceType
}

func NewPiece(newColor Color, newType PieceType) *Piece {
	return &Piece{color: newColor, pieceType: newType}
}

func (p *Piece) Color() Color    { return p.color }
func (p *Piece) Type() PieceType { return p.pieceType }

func (p *Piece) Render(bg string) string {
	var colorCode string
	if p.color == White {
		colorCode = WhitePiece
	} else {
		colorCode = BlackPiece
	}

	var symbol string
	switch p.pieceType {
	case Pawn:
		if p.color == White {
			symbol = "♙"
		} else {
			symbol = "♟"
		}
	case Rook:
		if p.color == White {
			symbol = "♖"
		} else {
			symbol = "♜"
		}
	case Knight:
		if p.color == White {
			symbol = "♘"
		} else {
			symbol = "♞"
		}
	case Bishop:
		if p.color == White {
			symbol = "♗"
		} else {
			symbol = "♝"
		}
	case Queen:
		if p.color == White {
			symbol = "♕"
		} else {
			symbol = "♛"
		}
	case King:
		if p.color == White {
			symbol = "♔"
		} else {
			symbol = "♚"
		}
	default:
		symbol = " "
	}
	return bg + colorCode + " " + symbol + " " + Reset
}

func (p *Piece) Symbol() string {
	switch p.pieceType {
	case Pawn:
		if p.color == White {
			return "♙"
		} else {
			return "♟"
		}
	case Rook:
		if p.color == White {
			return "♖"
		} else {
			return "♜"
		}
	case Knight:
		if p.color == White {
			return "♘"
		} else {
			return "♞"
		}
	case Bishop:
		if p.color == White {
			return "♗"
		} else {
			return "♝"
		}
	case Queen:
		if p.color == White {
			return "♕"
		} else {
			return "♛"
		}
	case King:
		if p.color == White {
			return "♔"
		} else {
			return "♚"
		}
	}
	return " "
}

// ParsePieceType разбирает строковое название фигуры (RU/EN, любой регистр).
func ParsePieceType(s string) (PieceType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "пешка":
		return Pawn, nil
	case "ладья", "л", "r", "rook":
		return Rook, nil
	case "конь", "к", "n", "knight":
		return Knight, nil
	case "слон", "с", "b", "bishop":
		return Bishop, nil
	case "ферзь", "ф", "q", "queen":
		return Queen, nil
	case "король":
		return King, nil
	}
	return 0, fmt.Errorf("неизвестный тип фигуры: %q", s)
}

// MarshalJSON реализует сериализацию в JSON с сохранением приватных полей.
func (p Piece) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Color     Color     `json:"цвет"`
		PieceType PieceType `json:"тип"`
	}{
		Color:     p.color,
		PieceType: p.pieceType,
	})
}

// UnmarshalJSON реализует десериализацию из JSON с заполнением приватных полей.
func (p *Piece) UnmarshalJSON(data []byte) error {
	var payload struct {
		Color     Color     `json:"цвет"`
		PieceType PieceType `json:"тип"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	p.color = payload.Color
	p.pieceType = payload.PieceType
	return nil
}

// MarshalJSON превращает Color в строку.
func (c Color) MarshalJSON() ([]byte, error) {
	switch c {
	case White:
		return []byte(`"Белые"`), nil
	case Black:
		return []byte(`"Чёрные"`), nil
	}
	return []byte(`"Неизвестно"`), nil
}

// UnmarshalJSON превращает строку обратно в Color.
func (c *Color) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Белые"`, `"белые"`:
		*c = White
	case `"Чёрные"`, `"чёрные"`, `"Черные"`, `"черные"`:
		*c = Black
	default:
		return fmt.Errorf("неизвестный цвет: %s", string(data))
	}
	return nil
}

// MarshalJSON превращает PieceType в строку.
func (t PieceType) MarshalJSON() ([]byte, error) {
	switch t {
	case Pawn:
		return []byte(`"Пешка"`), nil
	case Rook:
		return []byte(`"Ладья"`), nil
	case Knight:
		return []byte(`"Конь"`), nil
	case Bishop:
		return []byte(`"Слон"`), nil
	case Queen:
		return []byte(`"Ферзь"`), nil
	case King:
		return []byte(`"Король"`), nil
	}
	return []byte(`"Неизвестно"`), nil
}

// UnmarshalJSON превращает строку (или пустую строку/null) в PieceType.
func (t *PieceType) UnmarshalJSON(data []byte) error {
	s := string(data)

	// null → считаем, что поле не задано
	if s == "null" {
		*t = Pawn
		return nil
	}

	// снимаем кавычки JSON-строки
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// пустая строка (Swagger может прислать "") → поле не задано
	if s == "" {
		*t = Pawn
		return nil
	}

	pt, err := ParsePieceType(s)
	if err != nil {
		return err
	}
	*t = pt
	return nil
}
