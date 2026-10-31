package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mod.go/internal/model"
)

func TestCellToString(t *testing.T) {
	tests := []struct {
		name string
		row  int
		col  int
		want string
	}{
		{name: "a1 - левый нижний", row: 0, col: 0, want: "A1"},
		{name: "h8 - правый верхний", row: 7, col: 7, want: "H8"},
		{name: "e4 - центр", row: 3, col: 4, want: "E4"},
		{name: "d1 - ферзь белых", row: 0, col: 3, want: "D1"},
		{name: "e8 - король чёрных", row: 7, col: 4, want: "E8"},

		{name: "Z - 25-й столбец", row: 0, col: 25, want: "Z1"},
		{name: "AA - 26-й столбец", row: 0, col: 26, want: "AA1"},
		{name: "AB - 27-й столбец", row: 0, col: 27, want: "AB1"},
		{name: "ZZ - 701-й столбец", row: 0, col: 701, want: "ZZ1"},
		{name: "AAA - 702-й столбец", row: 0, col: 702, want: "AAA1"},

		{name: "строка 30", row: 29, col: 0, want: "A30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.CellToString(tt.row, tt.col)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatMove(t *testing.T) {
	tests := []struct {
		name string
		move model.Move
		want string
	}{
		{
			name: "пешка: обычный ход",
			move: model.Move{
				FromRow: 1, FromCol: 4,
				ToRow: 3, ToCol: 4,
				MovedPiece: model.NewPiece(model.White, model.Pawn),
			},
			want: "E2-E4",
		},
		{
			name: "пешка: взятие",
			move: model.Move{
				FromRow: 3, FromCol: 4,
				ToRow: 4, ToCol: 3,
				MovedPiece: model.NewPiece(model.White, model.Pawn),
				Captured:   model.NewPiece(model.Black, model.Pawn),
			},
			want: "E4xD5",
		},

		{
			name: "белая ладья: обычный ход",
			move: model.Move{
				FromRow: 0, FromCol: 0,
				ToRow: 7, ToCol: 0,
				MovedPiece: model.NewPiece(model.White, model.Rook),
			},
			want: "♖ A1-A8",
		},
		{
			name: "чёрный конь: обычный ход",
			move: model.Move{
				FromRow: 7, FromCol: 6,
				ToRow: 5, ToCol: 5,
				MovedPiece: model.NewPiece(model.Black, model.Knight),
			},
			want: "♞ G8-F6",
		},
		{
			name: "белый слон: взятие",
			move: model.Move{
				FromRow: 0, FromCol: 2,
				ToRow: 2, ToCol: 4,
				MovedPiece: model.NewPiece(model.White, model.Bishop),
				Captured:   model.NewPiece(model.Black, model.Pawn),
			},
			want: "♗ C1xE3",
		},

		{
			name: "короткая рокировка",
			move: model.Move{
				FromRow: 0, FromCol: 4,
				ToRow: 0, ToCol: 6,
				MovedPiece: model.NewPiece(model.White, model.King),
			},
			want: "O-O",
		},
		{
			name: "длинная рокировка",
			move: model.Move{
				FromRow: 0, FromCol: 4,
				ToRow: 0, ToCol: 2,
				MovedPiece: model.NewPiece(model.White, model.King),
			},
			want: "O-O-O",
		},
		{
			name: "короткая рокировка с шахом",
			move: model.Move{
				FromRow: 0, FromCol: 4,
				ToRow: 0, ToCol: 6,
				MovedPiece: model.NewPiece(model.White, model.King),
				Check:      true,
			},
			want: "O-O+",
		},
		{
			name: "длинная рокировка с матом",
			move: model.Move{
				FromRow: 0, FromCol: 4,
				ToRow: 0, ToCol: 2,
				MovedPiece: model.NewPiece(model.White, model.King),
				Mate:       true,
			},
			want: "O-O-O#",
		},

		{
			name: "пешка с шахом",
			move: model.Move{
				FromRow: 3, FromCol: 4,
				ToRow: 4, ToCol: 4,
				MovedPiece: model.NewPiece(model.White, model.Pawn),
				Check:      true,
			},
			want: "E4-E5+",
		},
		{
			name: "ферзь с матом",
			move: model.Move{
				FromRow: 3, FromCol: 3,
				ToRow: 7, ToCol: 3,
				MovedPiece: model.NewPiece(model.White, model.Queen),
				Mate:       true,
			},
			want: "♕ D4-D8#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.FormatMove(tt.move)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenitiveColor(t *testing.T) {
	tests := []struct {
		name  string
		color model.Color
		want  string
	}{
		{name: "белые - Белых", color: model.White, want: "Белых"},
		{name: "чёрные - Чёрных", color: model.Black, want: "Чёрных"},

		{name: "неизвестный - Чёрных (default)", color: model.Color(99), want: "Чёрных"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, model.GenitiveColor(tt.color))
		})
	}
}
