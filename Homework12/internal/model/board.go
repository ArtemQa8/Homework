package model

import "encoding/json"

type Board struct {
	rows  int
	cols  int
	cells [][]*Piece
}

func NewBoard(newRows, newCols int) *Board {
	b := &Board{rows: newRows, cols: newCols}
	b.cells = make([][]*Piece, newRows)
	for i := range b.cells {
		b.cells[i] = make([]*Piece, newCols)
	}

	for col := 0; col < newCols; col++ {
		b.cells[1][col] = NewPiece(White, Pawn)
		b.cells[newRows-2][col] = NewPiece(Black, Pawn)
	}

	b.placePieces(0, White)
	b.placePieces(newRows-1, Black)

	return b
}

func (b *Board) placePieces(row int, color Color) {
	c := b.cols
	b.cells[row][0] = NewPiece(color, Rook)
	b.cells[row][c-1] = NewPiece(color, Rook)
	b.cells[row][1] = NewPiece(color, Knight)
	b.cells[row][c-2] = NewPiece(color, Knight)
	b.cells[row][2] = NewPiece(color, Bishop)
	b.cells[row][c-3] = NewPiece(color, Bishop)
	b.cells[row][c/2-1] = NewPiece(color, Queen)
	b.cells[row][c/2] = NewPiece(color, King)
}

func (b *Board) RenderCell(row, col int) string {
	light := (row+col)%2 != 0
	bg := BlackSqBg
	if light {
		bg = WhiteSqBg
	}
	p := b.cells[row][col]
	if p == nil {
		return bg + "   " + Reset
	}
	return p.Render(bg)
}

func (b *Board) PieceAt(row, col int) *Piece {
	return b.cells[row][col]
}

func (b *Board) SetPiece(row, col int, piece *Piece) {
	b.cells[row][col] = piece
}

func (b *Board) Rows() int { return b.rows }
func (b *Board) Cols() int { return b.cols }

func (b *Board) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rows  int        `json:"строки"`
		Cols  int        `json:"столбцы"`
		Cells [][]*Piece `json:"клетки"`
	}{
		Rows:  b.rows,
		Cols:  b.cols,
		Cells: b.cells,
	})
}

func (b *Board) UnmarshalJSON(data []byte) error {
	var данные struct {
		Rows  int        `json:"строки"`
		Cols  int        `json:"столбцы"`
		Cells [][]*Piece `json:"клетки"`
	}
	if err := json.Unmarshal(data, &данные); err != nil {
		return err
	}
	b.rows = данные.Rows
	b.cols = данные.Cols
	b.cells = данные.Cells
	return nil
}
