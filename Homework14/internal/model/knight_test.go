package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mod.go/internal/model"
)

func TestKnightRules(t *testing.T) {
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
		{name: "Г: 2 вверх, 1 влево", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3, want: true},
		{name: "Г: 2 вверх, 1 вправо", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 5, want: true},
		{name: "Г: 2 вниз, 1 влево", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 3, want: true},
		{name: "Г: 2 вниз, 1 вправо", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 5, want: true},
		{name: "Г: 1 вверх, 2 влево", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 2, want: true},
		{name: "Г: 1 вверх, 2 вправо", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 6, want: true},
		{name: "Г: 1 вниз, 2 влево", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 2, want: true},
		{name: "Г: 1 вниз, 2 вправо", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 6, want: true},

		{name: "прямо вперёд", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 4, want: false},
		{name: "прямо вбок", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 6, want: false},
		{name: "по диагонали", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 6, want: false},
		{name: "не двинулась", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 4, want: false},
		{name: "на 3 вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 1, toCol: 4, want: false},
		{name: "2x2", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 2, want: false},

		{
			name:    "прыжок через свои фигуры",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 3, model.White, model.Pawn},
			want: true,
		},
		{
			name:    "прыжок через врага",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.Black, model.Pawn},
			want: true,
		},
		{
			name:    "белый: конечная - своя фигура",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{2, 3, model.White, model.Rook},
			want: false,
		},
		{
			name:    "белый: конечная - враг",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{2, 3, model.Black, model.Queen},
			want: true,
		},
		{
			name:    "чёрный: не бьёт свою",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{2, 3, model.Black, model.Pawn},
			want: false,
		},
		{
			name:    "чёрный: бьёт белого врага",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{2, 3, model.White, model.Queen},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := emptyBoard(8, 8)
			if tt.extraPiece != nil {
				board.SetPiece(tt.extraPiece.row, tt.extraPiece.col,
					model.NewPiece(tt.extraPiece.color, tt.extraPiece.ptype))
			}
			rules := model.NewKnightRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}
