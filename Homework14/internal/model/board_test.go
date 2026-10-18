package model_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/model"
)

func TestNewBoard_Size(t *testing.T) {
	tests := []struct {
		name string
		rows int
		cols int
	}{
		{name: "стандарт 8х8", rows: 8, cols: 8},
		{name: "минимум 4х4", rows: 4, cols: 4},
		{name: "большая 30х30", rows: 30, cols: 30},
		{name: "неквадратная 4х8", rows: 4, cols: 8},
		{name: "неквадратная 10х4", rows: 10, cols: 4},
		{name: "неквадратная 14х5", rows: 14, cols: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := model.NewBoard(tt.rows, tt.cols)
			assert.Equal(t, tt.rows, b.Rows())
			assert.Equal(t, tt.cols, b.Cols())
		})
	}
}

func TestNewBoard_WhitePieces(t *testing.T) {
	board := model.NewBoard(8, 8)

	tests := []struct {
		name      string
		row       int
		col       int
		wantType  model.PieceType
		wantColor model.Color
	}{
		{name: "ладья А", row: 0, col: 0, wantType: model.Rook, wantColor: model.White},
		{name: "конь B", row: 0, col: 1, wantType: model.Knight, wantColor: model.White},
		{name: "слон C", row: 0, col: 2, wantType: model.Bishop, wantColor: model.White},
		{name: "ферзь D", row: 0, col: 3, wantType: model.Queen, wantColor: model.White},
		{name: "король E", row: 0, col: 4, wantType: model.King, wantColor: model.White},
		{name: "слон F", row: 0, col: 5, wantType: model.Bishop, wantColor: model.White},
		{name: "конь G", row: 0, col: 6, wantType: model.Knight, wantColor: model.White},
		{name: "ладья H", row: 0, col: 7, wantType: model.Rook, wantColor: model.White},

		{name: "пешка А", row: 1, col: 0, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка B", row: 1, col: 1, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка C", row: 1, col: 2, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка D", row: 1, col: 3, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка E", row: 1, col: 4, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка F", row: 1, col: 5, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка G", row: 1, col: 6, wantType: model.Pawn, wantColor: model.White},
		{name: "пешка H", row: 1, col: 7, wantType: model.Pawn, wantColor: model.White},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			piece := board.PieceAt(tt.row, tt.col)
			require.NotNil(t, piece, "фигура должна быть на клетке")
			assert.Equal(t, tt.wantType, piece.Type())
			assert.Equal(t, tt.wantColor, piece.Color())
		})
	}
}

func TestNewBoard_BlackPieces(t *testing.T) {
	board := model.NewBoard(8, 8)

	tests := []struct {
		name      string
		row       int
		col       int
		wantType  model.PieceType
		wantColor model.Color
	}{
		{name: "ладья A", row: 7, col: 0, wantType: model.Rook, wantColor: model.Black},
		{name: "конь B", row: 7, col: 1, wantType: model.Knight, wantColor: model.Black},
		{name: "слон C", row: 7, col: 2, wantType: model.Bishop, wantColor: model.Black},
		{name: "ферзь D", row: 7, col: 3, wantType: model.Queen, wantColor: model.Black},
		{name: "король E", row: 7, col: 4, wantType: model.King, wantColor: model.Black},
		{name: "слон F", row: 7, col: 5, wantType: model.Bishop, wantColor: model.Black},
		{name: "конь G", row: 7, col: 6, wantType: model.Knight, wantColor: model.Black},
		{name: "ладья H", row: 7, col: 7, wantType: model.Rook, wantColor: model.Black},

		{name: "пешка A", row: 6, col: 0, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка B", row: 6, col: 1, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка C", row: 6, col: 2, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка D", row: 6, col: 3, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка E", row: 6, col: 4, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка F", row: 6, col: 5, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка G", row: 6, col: 6, wantType: model.Pawn, wantColor: model.Black},
		{name: "пешка H", row: 6, col: 7, wantType: model.Pawn, wantColor: model.Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			piece := board.PieceAt(tt.row, tt.col)
			require.NotNil(t, piece, "фигура должна быть на клетке")
			assert.Equal(t, tt.wantType, piece.Type())
			assert.Equal(t, tt.wantColor, piece.Color())
		})
	}
}

func TestNewBoard_EmptyMiddle(t *testing.T) {
	sizes := []struct {
		rows int
		cols int
	}{
		{rows: 4, cols: 4},
		{rows: 5, cols: 5},
		{rows: 6, cols: 6},
		{rows: 7, cols: 7},
		{rows: 8, cols: 8},
		{rows: 12, cols: 12},
		{rows: 4, cols: 8},
		{rows: 20, cols: 6},
	}
	for _, s := range sizes {
		t.Run(fmt.Sprintf("%dx%d", s.rows, s.cols), func(t *testing.T) {
			board := model.NewBoard(s.rows, s.cols)

			for row := 2; row <= s.rows-3; row++ {
				for col := range s.cols {
					assert.Nil(t, board.PieceAt(row, col),
						"клетка (%d,%d) должна быть пустой", row, col)
				}
			}
		})
	}
}

func TestNewBoard_Row0(t *testing.T) {
	sizes := []struct {
		rows int
		cols int
	}{
		{rows: 8, cols: 8},
		{rows: 8, cols: 9},
		{rows: 8, cols: 10},
		{rows: 8, cols: 11},
		{rows: 8, cols: 12},
		{rows: 8, cols: 20},
	}

	for _, s := range sizes {
		t.Run(fmt.Sprintf("%dx%d", s.rows, s.cols), func(t *testing.T) {
			board := model.NewBoard(s.rows, s.cols)
			c := s.cols

			assert.Equal(t, model.Rook, board.PieceAt(0, 0).Type(), "ладья A")
			assert.Equal(t, model.Knight, board.PieceAt(0, 1).Type(), "конь B")
			assert.Equal(t, model.Bishop, board.PieceAt(0, 2).Type(), "слон C")
			assert.Equal(t, model.Queen, board.PieceAt(0, c/2-1).Type(), "ферзь")
			assert.Equal(t, model.King, board.PieceAt(0, c/2).Type(), "король")
			assert.Equal(t, model.Bishop, board.PieceAt(0, c-3).Type(), "слон F")
			assert.Equal(t, model.Knight, board.PieceAt(0, c-2).Type(), "конь G")
			assert.Equal(t, model.Rook, board.PieceAt(0, c-1).Type(), "ладья H")

			for col := 3; col <= c/2-2; col++ {
				assert.Nil(t, board.PieceAt(0, col),
					"левый промежуток на (%d,%d)", 0, col)
			}

			for col := c/2 + 1; col <= c-4; col++ {
				assert.Nil(t, board.PieceAt(0, col),
					"правый промежуток на (%d,%d)", 0, col)
			}
		})
	}
}

func TestBoard_SetGetPiece(t *testing.T) {
	board := model.NewBoard(8, 8)

	// На пустой клетке — nil
	assert.Nil(t, board.PieceAt(3, 3))

	// Ставим фигуру
	queen := model.NewPiece(model.White, model.Queen)
	board.SetPiece(3, 3, queen)

	// Читаем обратно — та же фигура
	got := board.PieceAt(3, 3)
	require.NotNil(t, got)
	assert.Equal(t, model.Queen, got.Type())
	assert.Equal(t, model.White, got.Color())

	// Убираем фигуру
	board.SetPiece(3, 3, nil)

	// Клетка снова пустая
	assert.Nil(t, board.PieceAt(3, 3))
}

func TestBoard_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		rows int
		cols int
	}{
		{name: "8x8", rows: 8, cols: 8},
		{name: "12x12", rows: 12, cols: 12},
		{name: "4x8 (неквадратная)", rows: 4, cols: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := model.NewBoard(tt.rows, tt.cols)
			got, err := json.Marshal(board)
			require.NoError(t, err)

			var parsed struct {
				Rows  int     `json:"строки"`
				Cols  int     `json:"столбцы"`
				Cells [][]any `json:"клетки"`
			}
			err = json.Unmarshal(got, &parsed)
			require.NoError(t, err)

			assert.Equal(t, tt.rows, parsed.Rows)
			assert.Equal(t, tt.cols, parsed.Cols)
			assert.Len(t, parsed.Cells, tt.rows, "должно быть %d рядов", tt.rows)
			assert.Len(t, parsed.Cells[0], tt.cols, "в каждом ряду %d клеток", tt.cols)
		})
	}
}

func TestBoard_UnmarshalJSON(t *testing.T) {
	t.Run("валидная 8x8", func(t *testing.T) {
		input := `{
			"строки": 8,
			"столбцы": 8,
			"клетки": [
				[{"цвет":"Белые","тип":"Ладья"},{"цвет":"Белые","тип":"Конь"},null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null]
			]
		}`

		var board model.Board
		err := json.Unmarshal([]byte(input), &board)
		require.NoError(t, err)

		assert.Equal(t, 8, board.Rows())
		assert.Equal(t, 8, board.Cols())

		require.NotNil(t, board.PieceAt(0, 0))
		assert.Equal(t, model.Rook, board.PieceAt(0, 0).Type())
		assert.Equal(t, model.White, board.PieceAt(0, 0).Color())

		require.NotNil(t, board.PieceAt(0, 1))
		assert.Equal(t, model.Knight, board.PieceAt(0, 1).Type())
		assert.Equal(t, model.White, board.PieceAt(0, 1).Color())

		assert.Nil(t, board.PieceAt(0, 2))
		assert.Nil(t, board.PieceAt(7, 7))
	})

	t.Run("сломанный JSON", func(t *testing.T) {
		var board model.Board
		err := json.Unmarshal([]byte(`{`), &board)
		assert.Error(t, err)
	})

	t.Run("неверный тип поля", func(t *testing.T) {
		// "строки" ожидает int, а пришла строка — ошибка внутри UnmarshalJSON
		// Покрываем return в анмаршале
		input := `{"строки": "abc", "столбцы": 8, "клетки": []}`
		var board model.Board
		err := json.Unmarshal([]byte(input), &board)
		assert.Error(t, err)
	})
}

func TestBoard_RoundTrip(t *testing.T) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{name: "8x8", rows: 8, cols: 8},
		{name: "4x4", rows: 4, cols: 4},
		{name: "12x12", rows: 12, cols: 12},
		{name: "4x8 неквадратная", rows: 4, cols: 8},
	}

	for _, tt := range sizes {
		t.Run(tt.name, func(t *testing.T) {
			original := model.NewBoard(tt.rows, tt.cols)
			data, err := json.Marshal(original)
			require.NoError(t, err)

			var restored model.Board
			err = json.Unmarshal(data, &restored)
			require.NoError(t, err)

			assert.Equal(t, original.Rows(), restored.Rows())
			assert.Equal(t, original.Cols(), restored.Cols())

			for row := range original.Rows() {
				for col := range original.Cols() {
					origPiece := original.PieceAt(row, col)
					restPiece := restored.PieceAt(row, col)

					if origPiece == nil {
						assert.Nil(t, restPiece, "(%d,%d) должно быть пусто", row, col)
						continue
					}
					require.NotNil(t, restPiece, "(%d,%d) не должно быть пусто", row, col)
					assert.Equal(t, origPiece.Type(), restPiece.Type(), "тип (%d,%d)", row, col)
					assert.Equal(t, origPiece.Color(), restPiece.Color(), "цвет (%d,%d)", row, col)
				}
			}
		})
	}
}
