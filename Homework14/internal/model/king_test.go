package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mod.go/internal/model"
)

func TestKingRules(t *testing.T) {
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
		{name: "вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4, want: true},
		{name: "вниз", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 4, want: true},
		{name: "влево", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 3, want: true},
		{name: "вправо", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 5, want: true},
		{name: "влево-вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 3, want: true},
		{name: "вправо-вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 5, want: true},
		{name: "влево-вниз", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 3, want: true},
		{name: "вправо-вниз", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 5, want: true},

		{name: "на 2 вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 4, want: false},
		{name: "на 2 вбок", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 6, want: false},
		{name: "буквой Г", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3, want: false},
		{name: "на 2 по диагонали", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 6, want: false},

		{name: "не двинулась", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 4, want: false},

		{
			name:    "белый: конечная - своя фигура",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.White, model.Pawn},
			want: false,
		},
		{
			name:    "белый: конечная - враг",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.Black, model.Queen},
			want: true,
		},

		{name: "чёрный: на 1 клетку", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 5, want: true},
		{
			name:    "чёрный: не бьёт свою",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.Black, model.Pawn},
			want: false,
		},
		{
			name:    "чёрный: бьёт белого врага",
			color:   model.Black,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.White, model.Queen},
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
			rules := model.NewKingRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}
