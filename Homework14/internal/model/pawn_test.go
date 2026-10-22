package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mod.go/internal/model"
)

func TestPawnRules_CanMove(t *testing.T) {
	board := model.NewBoard(8, 8)

	tests := []struct {
		name    string
		color   model.Color
		fromRow int
		fromCol int
		toRow   int
		toCol   int
		want    bool
	}{
		{name: "белая: 1 вперёд", color: model.White,
			fromRow: 1, fromCol: 4, toRow: 2, toCol: 4, want: true},
		{name: "белая: 2 вперёд со старта", color: model.White,
			fromRow: 1, fromCol: 4, toRow: 3, toCol: 4, want: true},
		{name: "белая: 1 вперёд с середины", color: model.White,
			fromRow: 3, fromCol: 4, toRow: 4, toCol: 4, want: true},
		{name: "белая: 2 вперёд не со старта", color: model.White,
			fromRow: 3, fromCol: 4, toRow: 5, toCol: 4, want: false},
		{name: "белая: назад", color: model.White,
			fromRow: 3, fromCol: 4, toRow: 2, toCol: 4, want: false},
		{name: "белая: на 3 вперёд", color: model.White,
			fromRow: 1, fromCol: 4, toRow: 4, toCol: 4, want: false},
		{name: "белая: по диагонали пусто", color: model.White,
			fromRow: 3, fromCol: 4, toRow: 4, toCol: 5, want: false},

		{name: "чёрная: 1 вперёд", color: model.Black,
			fromRow: 6, fromCol: 4, toRow: 5, toCol: 4, want: true},
		{name: "чёрная: 2 вперёд со старта", color: model.Black,
			fromRow: 6, fromCol: 4, toRow: 4, toCol: 4, want: true},
		{name: "чёрная: 2 вперёд не со старта", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 4, want: false},
		{name: "чёрная: назад", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 4, want: false},
		{name: "чёрная: по диагонали пусто", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 5, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := model.NewPawnRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			got := rules.CanMove(move, board)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPawnRules_BlockedAndCapture(t *testing.T) {
	tests := []struct {
		name       string
		color      model.Color
		fromRow    int
		fromCol    int
		toRow      int
		toCol      int
		extraPiece *struct {
			row   int
			col   int
			color model.Color
			ptype model.PieceType
		}
		want bool
	}{
		// блокировка спереди
		{
			name:    "белая: впереди фигура, на 1 вперёд нельзя",
			color:   model.White,
			fromRow: 3, fromCol: 4, toRow: 4, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 4, model.White, model.Rook},
			want: false,
		},
		// Так же здесь проверяется, что пешка не бьет по прямой
		{
			name:    "белая: впереди на 2-ой фигура, на 2 вперёд нельзя",
			color:   model.White,
			fromRow: 1, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.Black, model.Queen},
			want: false,
		},
		{
			name:    "белая: впереди на 1-ой фигура, на 2 вперёд нельзя",
			color:   model.White,
			fromRow: 1, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{2, 4, model.Black, model.Bishop},
			want: false,
		},

		// взятие по диагонали
		{
			name:    "белая: бьёт врага по диагонали",
			color:   model.White,
			fromRow: 3, fromCol: 4, toRow: 4, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 5, model.Black, model.Knight},
			want: true,
		},
		{
			name:    "белая: не бьёт свою по диагонали",
			color:   model.White,
			fromRow: 3, fromCol: 4, toRow: 4, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 5, model.White, model.Pawn},
			want: false,
		},
		{
			name:    "белая: по диагонали пусто - не идёт",
			color:   model.White,
			fromRow: 3, fromCol: 4, toRow: 4, toCol: 5,
			want: false,
		},
		{
			name:    "чёрная: впереди фигура, на 1 вперёд нельзя",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.White, model.Rook},
			want: false,
		},
		{
			name:    "чёрная: бьёт врага по диагонали",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 5, model.White, model.Knight},
			want: true,
		},
		{
			name:    "чёрная: не бьёт свою по диагонали",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 5, model.Black, model.Pawn},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := emptyBoard(8, 8)
			if tt.extraPiece != nil {
				board.SetPiece(tt.extraPiece.row, tt.extraPiece.col,
					model.NewPiece(tt.extraPiece.color, tt.extraPiece.ptype))
			}
			rules := model.NewPawnRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}
