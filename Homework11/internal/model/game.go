package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type Game struct {
	player1 Player
	player2 Player
	board   *Board
	current Color
	moves   []Move
	id      int

	whiteKingMoved  bool
	blackKingMoved  bool
	whiteRookAMoved bool
	whiteRookHMoved bool
	blackRookAMoved bool
	blackRookHMoved bool
}

func NewGame(name1, name2 string, rows, cols int) *Game {
	return &Game{
		player1: *NewPlayer(name1, White),
		player2: *NewPlayer(name2, Black),
		board:   NewBoard(rows, cols),
		current: White,
	}
}

func (g *Game) Player1() Player     { return g.player1 }
func (g *Game) Player2() Player     { return g.player2 }
func (g *Game) Board() *Board       { return g.board }
func (g *Game) CurrentColor() Color { return g.current }
func (g *Game) ID() int             { return g.id }
func (g *Game) SetID(id int)        { g.id = id }

func (g *Game) SetPlayer1(player Player) { g.player1 = player }
func (g *Game) SetPlayer2(player Player) { g.player2 = player }

func (g *Game) Moves() []Move {
	result := make([]Move, len(g.moves))
	copy(result, g.moves)
	return result
}

func (g *Game) InCheck(color Color) bool {
	return kingInCheck(g, color)
}

// ── Отображение ────────────────────────────────────────────

func (g *Game) Render() string {
	var sb strings.Builder
	maxWidth := int(math.Log10(float64(g.board.Rows()))) + 1
	boardWidth := int(g.board.Cols())*3 + maxWidth + 1

	sb.WriteString(center(g.player1.Name(), boardWidth))
	sb.WriteByte('\n')

	sb.WriteString(strings.Repeat(" ", maxWidth+1))
	for c := 0; c < g.board.Cols(); c++ {
		name := columnName(c)
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

	for r := 0; r < g.board.Rows(); r++ {
		fmt.Fprintf(&sb, "%*d ", maxWidth, r+1)
		for c := 0; c < g.board.Cols(); c++ {
			sb.WriteString(g.board.RenderCell(r, c))
		}
		sb.WriteByte('\n')
	}

	sb.WriteString(center(g.player2.Name(), boardWidth))
	sb.WriteByte('\n')

	// ── История ходов ──────────────────────────────────────
	sb.WriteByte('\n')
	sb.WriteString("История ходов:\n")

	start := len(g.moves) - 10
	if start < 0 {
		start = 0
	}
	if start%2 != 0 {
		start++
	}

	moveStrings := make([]string, 0, len(g.moves)-start)
	maxLength := 0
	for i := start; i < len(g.moves); i++ {
		str := FormatMove(g.moves[i])
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

func (g *Game) BoardRow(row int) string {
	maxWidth := int(math.Log10(float64(g.board.Rows()))) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%*d ", maxWidth, row+1)
	for c := 0; c < g.board.Cols(); c++ {
		sb.WriteString(g.board.RenderCell(row, c))
	}
	return sb.String()
}

// ── Приватные проверки (шах, легальность, мат) ────────────

func kingInCheck(game *Game, color Color) bool {
	var kingRow, kingCol int
	found := false
	for r := 0; r < game.board.Rows(); r++ {
		for c := 0; c < game.board.Cols(); c++ {
			p := game.board.PieceAt(r, c)
			if p != nil && p.Color() == color && p.Type() == King {
				kingRow, kingCol = r, c
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return false
	}

	oppositeColor := White
	if color == White {
		oppositeColor = Black
	}
	for r := 0; r < game.board.Rows(); r++ {
		for c := 0; c < game.board.Cols(); c++ {
			p := game.board.PieceAt(r, c)
			if p == nil || p.Color() != oppositeColor {
				continue
			}
			checkMove := &Move{
				FromRow: r, FromCol: c,
				ToRow: kingRow, ToCol: kingCol,
			}
			rules := RulesFor(p)
			if rules != nil && rules.CanMove(checkMove, game.board) {
				return true
			}
		}
	}
	return false
}

func (g *Game) LegalMove(move *Move, color Color) bool {
	pieceFrom := g.board.PieceAt(move.FromRow, move.FromCol)
	pieceTo := g.board.PieceAt(move.ToRow, move.ToCol)

	g.board.SetPiece(move.ToRow, move.ToCol, pieceFrom)
	g.board.SetPiece(move.FromRow, move.FromCol, nil)

	inCheck := kingInCheck(g, color)

	g.board.SetPiece(move.FromRow, move.FromCol, pieceFrom)
	g.board.SetPiece(move.ToRow, move.ToCol, pieceTo)

	return !inCheck
}

func (g *Game) Checkmate(color Color) bool {
	if !kingInCheck(g, color) {
		return false
	}

	for fromRow := 0; fromRow < g.board.Rows(); fromRow++ {
		for fromCol := 0; fromCol < g.board.Cols(); fromCol++ {
			piece := g.board.PieceAt(fromRow, fromCol)
			if piece == nil || piece.Color() != color {
				continue
			}
			for toRow := 0; toRow < g.board.Rows(); toRow++ {
				for toCol := 0; toCol < g.board.Cols(); toCol++ {
					if fromRow == toRow && fromCol == toCol {
						continue
					}
					testMove := &Move{FromRow: fromRow, FromCol: fromCol, ToRow: toRow, ToCol: toCol}
					rules := RulesFor(piece)
					if rules == nil || !rules.CanMove(testMove, g.board) {
						continue
					}
					if g.LegalMove(testMove, color) {
						return false
					}
				}
			}
		}
	}
	return true
}

func (g *Game) Stalemate(color Color) bool {
	if kingInCheck(g, color) {
		return false
	}

	for fromRow := 0; fromRow < g.board.Rows(); fromRow++ {
		for fromCol := 0; fromCol < g.board.Cols(); fromCol++ {
			piece := g.board.PieceAt(fromRow, fromCol)
			if piece == nil || piece.Color() != color {
				continue
			}
			rules := RulesFor(piece)
			if rules == nil {
				continue
			}
			for toRow := 0; toRow < g.board.Rows(); toRow++ {
				for toCol := 0; toCol < g.board.Cols(); toCol++ {
					if fromRow == toRow && fromCol == toCol {
						continue
					}
					testMove := NewMove(fromRow, fromCol, toRow, toCol)
					if rules.CanMove(testMove, g.board) && g.LegalMove(testMove, color) {
						return false
					}
				}
			}
		}
	}
	return true
}

// ── Выполнение хода ───────────────────────────────────────

func (g *Game) MakeMove(move *Move) error {
	if move.FromRow < 0 || move.FromRow >= g.board.Rows() ||
		move.FromCol < 0 || move.FromCol >= g.board.Cols() ||
		move.ToRow < 0 || move.ToRow >= g.board.Rows() ||
		move.ToCol < 0 || move.ToCol >= g.board.Cols() {
		return errors.New("выход за границы доски")
	}

	piece := g.board.PieceAt(move.FromRow, move.FromCol)
	if piece == nil {
		return errors.New("на начальной клетке нет фигуры")
	}
	if piece.Color() != g.current {
		return errors.New("сейчас ход другого игрока")
	}

	if piece.Type() == King && abs(move.ToCol-move.FromCol) == 2 {
		move.MovedPiece = piece
		if g.tryCastling(move, piece) {
			if kingInCheck(g, g.current) {
				if g.Checkmate(g.current) {
					move.Mate = true
				} else {
					move.Check = true
				}
			}
			g.moves = append(g.moves, *move)
			if g.current == White {
				g.current = Black
			} else {
				g.current = White
			}
			return nil
		}
		return errors.New("невозможно выполнить рокировку")
	}

	rules := RulesFor(piece)
	if rules == nil {
		return errors.New("для этой фигуры нет правил хода")
	}

	moveAllowed := rules.CanMove(move, g.board)

	enPassant := false
	if !moveAllowed && piece.Type() == Pawn {
		enPassant = isEnPassant(move, piece, g)
	}
	if !moveAllowed && !enPassant {
		return fmt.Errorf("%s так не ходит", pieceName(piece.Type()))
	}

	enemy := g.board.PieceAt(move.ToRow, move.ToCol)
	if !enPassant && enemy != nil && enemy.Color() == g.current {
		return errors.New("нельзя бить свою фигуру")
	}

	oldPieceFrom := piece
	oldPieceTo := enemy
	enPassantPawn := (*Piece)(nil)
	lastMove := (*Move)(nil)

	g.board.SetPiece(move.ToRow, move.ToCol, piece)
	g.board.SetPiece(move.FromRow, move.FromCol, nil)

	if enPassant {
		lastMove = &g.moves[len(g.moves)-1]
		enPassantPawn = g.board.PieceAt(lastMove.ToRow, lastMove.ToCol)
		g.board.SetPiece(lastMove.ToRow, lastMove.ToCol, nil)
	}

	if kingInCheck(g, piece.Color()) {
		g.board.SetPiece(move.FromRow, move.FromCol, oldPieceFrom)
		g.board.SetPiece(move.ToRow, move.ToCol, oldPieceTo)
		if enPassant && enPassantPawn != nil {
			g.board.SetPiece(lastMove.ToRow, lastMove.ToCol, enPassantPawn)
		}
		return errors.New("этот ход оставляет вашего короля под шахом")
	}

	switch {
	case piece.Type() == King:
		if piece.Color() == White {
			g.whiteKingMoved = true
		} else {
			g.blackKingMoved = true
		}
	case piece.Type() == Rook:
		if piece.Color() == White {
			if move.FromCol == 0 {
				g.whiteRookAMoved = true
			} else if move.FromCol == g.board.Cols()-1 {
				g.whiteRookHMoved = true
			}
		} else {
			if move.FromCol == 0 {
				g.blackRookAMoved = true
			} else if move.FromCol == g.board.Cols()-1 {
				g.blackRookHMoved = true
			}
		}
	}

	move.Captured = enemy
	if enPassant {
		move.Captured = enPassantPawn
	}
	move.MovedPiece = piece

	if piece.Type() == Pawn {
		lastRow := 0
		if piece.Color() == White {
			lastRow = g.board.Rows() - 1
		}
		if move.ToRow == lastRow && move.Promotion != 0 {
			newPiece := NewPiece(piece.Color(), move.Promotion)
			g.board.SetPiece(move.ToRow, move.ToCol, newPiece)
		}
	}

	if g.current == White {
		g.current = Black
	} else {
		g.current = White
	}

	if kingInCheck(g, g.current) {
		if g.Checkmate(g.current) {
			move.Mate = true
		} else {
			move.Check = true
		}
	}

	g.moves = append(g.moves, *move)

	return nil
}

// ── Специальные ходы ──────────────────────────────────────

func (g *Game) tryCastling(move *Move, king *Piece) bool {
	color := king.Color()
	row := move.FromRow
	fromCol := move.FromCol
	toCol := move.ToCol

	direction := 1
	if toCol < fromCol {
		direction = -1
	}
	if abs(toCol-fromCol) != 2 {
		return false
	}

	rookCol := 0
	rookColAfter := 0
	if direction == 1 {
		rookCol = g.board.Cols() - 1
		rookColAfter = toCol - 1
	} else {
		rookCol = 0
		rookColAfter = toCol + 1
	}

	kingMoved := false
	rookMoved := false
	switch {
	case color == White && direction == 1:
		kingMoved = g.whiteKingMoved
		rookMoved = g.whiteRookHMoved
	case color == White && direction == -1:
		kingMoved = g.whiteKingMoved
		rookMoved = g.whiteRookAMoved
	case color == Black && direction == 1:
		kingMoved = g.blackKingMoved
		rookMoved = g.blackRookHMoved
	case color == Black && direction == -1:
		kingMoved = g.blackKingMoved
		rookMoved = g.blackRookAMoved
	}
	if kingMoved || rookMoved {
		return false
	}

	step := 1
	if direction == -1 {
		step = -1
	}
	for col := fromCol + step; col != rookCol; col += step {
		if g.board.PieceAt(row, col) != nil {
			return false
		}
	}

	if kingInCheck(g, color) {
		return false
	}

	intermediateCol := fromCol + step
	tempMove := &Move{FromRow: row, FromCol: fromCol, ToRow: row, ToCol: intermediateCol}
	if !g.LegalMove(tempMove, color) {
		return false
	}

	rook := g.board.PieceAt(row, rookCol)
	if rook == nil || rook.Type() != Rook || rook.Color() != color {
		return false
	}

	g.board.SetPiece(row, toCol, king)
	g.board.SetPiece(row, fromCol, nil)
	g.board.SetPiece(row, rookColAfter, rook)
	g.board.SetPiece(row, rookCol, nil)

	switch {
	case color == White:
		g.whiteKingMoved = true
		if direction == 1 {
			g.whiteRookHMoved = true
		} else {
			g.whiteRookAMoved = true
		}
	case color == Black:
		g.blackKingMoved = true
		if direction == 1 {
			g.blackRookHMoved = true
		} else {
			g.blackRookAMoved = true
		}
	}
	return true
}

func isEnPassant(move *Move, piece *Piece, game *Game) bool {
	if len(game.moves) == 0 {
		return false
	}
	last := &game.moves[len(game.moves)-1]
	enemy := game.board.PieceAt(last.ToRow, last.ToCol)
	if enemy == nil || enemy.Type() != Pawn {
		return false
	}
	if abs(last.ToRow-last.FromRow) != 2 {
		return false
	}
	if abs(last.ToCol-move.FromCol) != 1 {
		return false
	}
	shift := 1
	if piece.Color() == Black {
		shift = -1
	}
	return move.ToRow == last.ToRow+shift && move.ToCol == last.ToCol
}

func pieceName(pieceType PieceType) string {
	switch pieceType {
	case Pawn:
		return "Пешка"
	case Rook:
		return "Ладья"
	case Knight:
		return "Конь"
	case Bishop:
		return "Слон"
	case Queen:
		return "Ферзь"
	case King:
		return "Король"
	default:
		return "Фигура"
	}
}

func (g *Game) ObjectType() string { return "игра" }

func (g *Game) RenderLines() []string {
	return strings.Split(g.Render(), "\n")
}

// RenderLinesWithoutHistory возвращает отображение игры без блока "История ходов".
func (g *Game) RenderLinesWithoutHistory() []string {
	full := strings.Split(g.Render(), "\n")
	for i, line := range full {
		if strings.Contains(line, "История ходов:") {
			return full[:i]
		}
	}
	return full
}

// MarshalJSON реализует сериализацию в JSON с сохранением приватных полей.
func (g Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID      int    `json:"id"`
		Player1 Player `json:"игрок1"`
		Player2 Player `json:"игрок2"`
		Board   *Board `json:"доска"`
		Current Color  `json:"текущий"`
		Moves   []Move `json:"ходы"`

		WhiteKingMoved  bool `json:"белыйКорольДвигался,omitempty"`
		BlackKingMoved  bool `json:"чёрныйКорольДвигался,omitempty"`
		WhiteRookAMoved bool `json:"белаяЛадьяAДвигалась,omitempty"`
		WhiteRookHMoved bool `json:"белаяЛадьяHДвигалась,omitempty"`
		BlackRookAMoved bool `json:"чёрнаяЛадьяAДвигалась,omitempty"`
		BlackRookHMoved bool `json:"чёрнаяЛадьяHДвигалась,omitempty"`
	}{
		ID:      g.id,
		Player1: g.player1,
		Player2: g.player2,
		Board:   g.board,
		Current: g.current,
		Moves:   g.moves,

		WhiteKingMoved:  g.whiteKingMoved,
		BlackKingMoved:  g.blackKingMoved,
		WhiteRookAMoved: g.whiteRookAMoved,
		WhiteRookHMoved: g.whiteRookHMoved,
		BlackRookAMoved: g.blackRookAMoved,
		BlackRookHMoved: g.blackRookHMoved,
	})
}

// UnmarshalJSON реализует десериализацию из JSON с заполнением приватных полей.
func (g *Game) UnmarshalJSON(data []byte) error {
	var payload struct {
		ID      int    `json:"id"`
		Player1 Player `json:"игрок1"`
		Player2 Player `json:"игрок2"`
		Board   *Board `json:"доска"`
		Current Color  `json:"текущий"`
		Moves   []Move `json:"ходы"`

		WhiteKingMoved  bool `json:"белыйКорольДвигался,omitempty"`
		BlackKingMoved  bool `json:"чёрныйКорольДвигался,omitempty"`
		WhiteRookAMoved bool `json:"белаяЛадьяAДвигалась,omitempty"`
		WhiteRookHMoved bool `json:"белаяЛадьяHДвигалась,omitempty"`
		BlackRookAMoved bool `json:"чёрнаяЛадьяAДвигалась,omitempty"`
		BlackRookHMoved bool `json:"чёрнаяЛадьяHДвигалась,omitempty"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	g.id = payload.ID
	g.player1 = payload.Player1
	g.player2 = payload.Player2
	g.board = payload.Board
	g.current = payload.Current
	g.moves = payload.Moves

	g.whiteKingMoved = payload.WhiteKingMoved
	g.blackKingMoved = payload.BlackKingMoved
	g.whiteRookAMoved = payload.WhiteRookAMoved
	g.whiteRookHMoved = payload.WhiteRookHMoved
	g.blackRookAMoved = payload.BlackRookAMoved
	g.blackRookHMoved = payload.BlackRookHMoved

	return nil
}

func (g *Game) SetBoard(board *Board) {
	g.board = board
}

func (g *Game) SetLastMoveID(id int) {
	if len(g.moves) == 0 {
		return
	}
	g.moves[len(g.moves)-1].id = id
}
