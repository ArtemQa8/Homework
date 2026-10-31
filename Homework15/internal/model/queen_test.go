package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mod.go/internal/model"
)

func TestQueenRules_CanMove(t *testing.T) {
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
		{name: "прямо вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 1, toCol: 4, want: true},
		{name: "прямо вниз", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 7, toCol: 4, want: true},
		{name: "прямо влево", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 1, want: true},
		{name: "прямо вправо", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 7, want: true},

		{name: "влево-вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 1, toCol: 1, want: true},
		{name: "вправо-вверх", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 1, toCol: 7, want: true},
		{name: "влево-вниз", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 7, toCol: 1, want: true},
		{name: "вправо-вниз", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 7, toCol: 7, want: true},

		{name: "буквой Г", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3, want: false},
		{name: "не равная диагональ (2,3)", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 1, want: false},

		{name: "не двинулась", color: model.White,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 4, want: false},

		{name: "чёрный: по прямой", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 4, want: true},
		{name: "чёрный: по диагонали", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 6, want: true},
		{name: "чёрный: буквой Г", color: model.Black,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := model.NewQueenRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}

func TestQueenRules_BlockedAndCapture(t *testing.T) {
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
			name:    "блок на прямой",
			color:   model.White,
			fromRow: 4, fromCol: 1, toRow: 4, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 3, model.Black, model.Pawn},
			want: false,
		},

		{
			name:    "блок на диагонали",
			color:   model.White,
			fromRow: 2, fromCol: 2, toRow: 6, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{4, 4, model.White, model.Pawn},
			want: false,
		},

		{
			name:    "белый: конечная - своя",
			color:   model.White,
			fromRow: 2, fromCol: 2, toRow: 6, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{6, 6, model.White, model.Pawn},
			want: false,
		},
		{
			name:    "белый: конечная - враг",
			color:   model.White,
			fromRow: 2, fromCol: 2, toRow: 6, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{6, 6, model.Black, model.Queen},
			want: true,
		},
		{
			name:    "белый: конечная - пусто",
			color:   model.White,
			fromRow: 2, fromCol: 2, toRow: 6, toCol: 6,
			want: true,
		},

		{
			name:    "чёрный: не бьёт свою",
			color:   model.Black,
			fromRow: 2, fromCol: 2, toRow: 6, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{6, 6, model.Black, model.Pawn},
			want: false,
		},
		{
			name:    "чёрный: бьёт белого врага",
			color:   model.Black,
			fromRow: 2, fromCol: 2, toRow: 6, toCol: 6,
			extraPiece: &struct {
				row   int
				col   int
				color model.Color
				ptype model.PieceType
			}{6, 6, model.White, model.Queen},
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
			rules := model.NewQueenRules(tt.color)
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			assert.Equal(t, tt.want, rules.CanMove(move, board))
		})
	}
}
