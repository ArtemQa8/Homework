package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mod.go/internal/model"
)

func TestRookRules_SimpleMoves(t *testing.T) {
	board := emptyBoard(8, 8)

	tests := []struct {
		name    string
		color   model.Color
		fromRow int
		fromCol int
		toRow   int
		toCol   int
		want    bool
	}{
		{name: "вверх на 3", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 1, toCol: 4, want: true},
		{name: "вниз на 3", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 7, toCol: 4, want: true},
		{name: "вверх на 1", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 3, toCol: 4, want: true},

		{name: "вправо на 3", color: model.White,
			fromRow: 4, fromCol: 2, toRow: 4, toCol: 5, want: true},
		{name: "влево на 3", color: model.White,
			fromRow: 4, fromCol: 5, toRow: 4, toCol: 2, want: true},

		{name: "по диагонали", color: model.White,
			fromRow: 3, fromCol: 3, toRow: 5, toCol: 5, want: false},
		{name: "буквой Г", color: model.White,
			fromRow: 3, fromCol: 3, toRow: 5, toCol: 4, want: false},

		{name: "не двинулась", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 4, want: false},

		{name: "чёрная: вверх на 2", color: model.Black,
			fromRow: 5, fromCol: 2, toRow: 3, toCol: 2, want: true},
		{name: "чёрная: по диагонали", color: model.Black,
			fromRow: 3, fromCol: 3, toRow: 4, toCol: 4, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := model.NewRookRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}

func TestRookRules_BlockedAndCapture(t *testing.T) {
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
		{
			name:    "белая: вверх, блок на пути",
			color:   model.White,
			fromRow: 6, fromCol: 3, toRow: 2, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 3, model.Black, model.Pawn},
			want: false,
		},
		{
			name:    "белая: вниз, блок на пути",
			color:   model.White,
			fromRow: 2, fromCol: 3, toRow: 6, toCol: 3,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 3, model.Black, model.Knight},
			want: false,
		},
		{
			name:    "белая: вправо, блок на пути",
			color:   model.White,
			fromRow: 3, fromCol: 1, toRow: 3, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 3, model.Black, model.Bishop},
			want: false,
		},
		{
			name:    "белая: влево, блок на пути",
			color:   model.White,
			fromRow: 3, fromCol: 6, toRow: 3, toCol: 1,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 4, model.White, model.Pawn},
			want: false,
		},
		{
			name:    "белая: конечная - своя фигура",
			color:   model.White,
			fromRow: 3, fromCol: 1, toRow: 3, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 5, model.White, model.Pawn},
			want: false,
		},
		{
			name:    "белая: конечная - враг, можно бить",
			color:   model.White,
			fromRow: 3, fromCol: 1, toRow: 3, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 5, model.Black, model.Queen},
			want: true,
		},
		{
			name:    "белая: конечная - пусто",
			color:   model.White,
			fromRow: 3, fromCol: 1, toRow: 3, toCol: 5,
			want: true,
		},
		{
			name:    "белая: блок сразу за стартом",
			color:   model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 5, model.Black, model.Pawn},
			want: false,
		},
		{
			name:    "чёрная: не бьёт свою",
			color:   model.Black,
			fromRow: 3, fromCol: 1, toRow: 3, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 5, model.Black, model.Pawn},
			want: false,
		},
		{
			name:    "чёрная: бьёт белого врага",
			color:   model.Black,
			fromRow: 3, fromCol: 1, toRow: 3, toCol: 5,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{3, 5, model.White, model.Queen},
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
			rules := model.NewRookRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}
