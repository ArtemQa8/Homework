package tui

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"mod.go/internal/model"
)

// ============ ANSI-константы ============

const (
	Reset      = "\033[0m"
	WhitePiece = "\033[31m"
	BlackPiece = "\033[34m"
	WhiteSqBg  = "\033[48;5;250m"
	BlackSqBg  = "\033[48;5;240m"
)

// ============ Утилиты ============

// VisibleLength считает длину строки без учёта ANSI-кодов.
func VisibleLength(s string) int {
	length := 0
	inEscape := false
	for _, r := range s {
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if r == '\033' {
			inEscape = true
			continue
		}
		length++
	}
	return length
}

// PadWithSpaces дополняет строку пробелами до нужной видимой ширины.
func PadWithSpaces(s string, width int) string {
	current := VisibleLength(s)
	if current >= width {
		return s
	}
	return s + strings.Repeat(" ", width-current)
}

func center(name string, boardWidth int) string {
	length := utf8.RuneCountInString(name)
	totalSpaces := boardWidth - length
	if totalSpaces <= 0 {
		return name
	}
	left := totalSpaces / 2
	right := totalSpaces - left
	return strings.Repeat(" ", left) + name + strings.Repeat(" ", right)
}

// ============ Фигура ============

// PieceRender возвращает ANSI-строку фигуры с заданным фоном клетки.
func PieceRender(p *model.Piece, bg string) string {
	if p == nil {
		return bg + "   " + Reset
	}
	colorCode := BlackPiece
	if p.Color() == model.White {
		colorCode = WhitePiece
	}

	var symbol string
	switch p.Type() {
	case model.Pawn:
		if p.Color() == model.White {
			symbol = "♙"
		} else {
			symbol = "♟"
		}
	case model.Rook:
		if p.Color() == model.White {
			symbol = "♖"
		} else {
			symbol = "♜"
		}
	case model.Knight:
		if p.Color() == model.White {
			symbol = "♘"
		} else {
			symbol = "♞"
		}
	case model.Bishop:
		if p.Color() == model.White {
			symbol = "♗"
		} else {
			symbol = "♝"
		}
	case model.Queen:
		if p.Color() == model.White {
			symbol = "♕"
		} else {
			symbol = "♛"
		}
	case model.King:
		if p.Color() == model.White {
			symbol = "♔"
		} else {
			symbol = "♚"
		}
	default:
		symbol = " "
	}
	return bg + colorCode + " " + symbol + " " + Reset
}

// ============ Доска ============

// BoardRenderCell возвращает ANSI-строку клетки доски.
func BoardRenderCell(b *model.Board, row, col int) string {
	light := (row+col)%2 != 0
	bg := BlackSqBg
	if light {
		bg = WhiteSqBg
	}
	return PieceRender(b.PieceAt(row, col), bg)
}

// ============ Игра ============

// RenderGame возвращает полное отображение игры в терминал.
func RenderGame(g *model.Game) string {
	var sb strings.Builder
	board := g.Board()
	maxWidth := int(math.Log10(float64(board.Rows()))) + 1
	boardWidth := int(board.Cols())*3 + maxWidth + 1

	sb.WriteString(center(g.Player1().Name(), boardWidth))
	sb.WriteByte('\n')

	sb.WriteString(strings.Repeat(" ", maxWidth+1))
	for c := 0; c < board.Cols(); c++ {
		name := model.ColumnName(c)
		length := len(name)
		if length > 3 {
			name = name[:3]
			length = 3
		}
		left := (3 - length + 1) / 2
		right := 3 - length - left
		sb.WriteString(strings.Repeat(" ", left))
		sb.WriteString(name)
		sb.WriteString(strings.Repeat(" ", right))
	}
	sb.WriteByte('\n')

	for r := 0; r < board.Rows(); r++ {
		fmt.Fprintf(&sb, "%*d ", maxWidth, r+1)
		for c := 0; c < board.Cols(); c++ {
			sb.WriteString(BoardRenderCell(board, r, c))
		}
		sb.WriteByte('\n')
	}

	sb.WriteString(center(g.Player2().Name(), boardWidth))
	sb.WriteByte('\n')

	// История ходов
	sb.WriteByte('\n')
	sb.WriteString("История ходов:\n")

	moves := g.Moves()
	start := len(moves) - 10
	if start < 0 {
		start = 0
	}
	if start%2 != 0 {
		start++
	}

	moveStrings := make([]string, 0, len(moves)-start)
	maxLength := 0
	for i := start; i < len(moves); i++ {
		str := model.FormatMove(moves[i])
		moveStrings = append(moveStrings, str)
		runeLen := utf8.RuneCountInString(str)
		if runeLen > maxLength {
			maxLength = runeLen
		}
	}

	moveNumber := start/2 + 1
	for i, str := range moveStrings {
		missing := maxLength - utf8.RuneCountInString(str)
		str = str + strings.Repeat(" ", missing)

		if (start+i)%2 == 0 {
			sb.WriteString(fmt.Sprintf("%2d. ", moveNumber))
			moveNumber++
			sb.WriteString(str)
			sb.WriteString("  ")
		} else {
			sb.WriteString("    ")
			sb.WriteString(str)
			sb.WriteByte('\n')
		}
	}
	if len(moveStrings) > 0 && len(moveStrings)%2 != 0 {
		sb.WriteByte('\n')
	}

	return sb.String()
}

// RenderGameWithoutHistory возвращает отображение игры без блока "История ходов".
func RenderGameWithoutHistory(g *model.Game) []string {
	full := strings.Split(RenderGame(g), "\n")
	for i, line := range full {
		if strings.Contains(line, "История ходов:") {
			return full[:i]
		}
	}
	return full
}
