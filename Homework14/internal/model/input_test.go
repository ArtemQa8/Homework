package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/model"
)

func TestParseCell(t *testing.T) {
	boardDefault := model.NewBoard(8, 8)
	boardSmall := model.NewBoard(4, 4)
	boardBig := model.NewBoard(30, 30)

	tests := []struct {
		name     string
		board    *model.Board
		notation string
		wantRow  int
		wantCol  int
		wantErr  bool
	}{
		// wantErr если не пишем, то там false
		{name: "a1 8x8", board: boardDefault, notation: "a1", wantRow: 0, wantCol: 0},
		{name: "a1 4x4", board: boardSmall, notation: "a1", wantRow: 0, wantCol: 0},
		{name: "a1 30x30", board: boardBig, notation: "a1", wantRow: 0, wantCol: 0},

		// wantCol,wantRow не пишем, так как хотим ошибку и до них не дойдет
		{name: "e4 8x8", board: boardDefault, notation: "e4", wantRow: 3, wantCol: 4},
		{name: "e4 4x4", board: boardSmall, notation: "e4", wantErr: true},
		{name: "e4 30x30", board: boardBig, notation: "e4", wantRow: 3, wantCol: 4},

		{name: "c8 8x8", board: boardDefault, notation: "c8", wantRow: 7, wantCol: 2},
		{name: "c8 4x4", board: boardSmall, notation: "c8", wantErr: true},
		{name: "c8 30x30", board: boardBig, notation: "c8", wantRow: 7, wantCol: 2},

		{name: "aa1 8x8", board: boardDefault, notation: "aa1", wantErr: true},
		{name: "aa1 4x4", board: boardSmall, notation: "aa1", wantErr: true},
		{name: "aa1 30x30", board: boardBig, notation: "aa1", wantRow: 0, wantCol: 26},

		// Не зависящие от размера доски, проверяем на стандартной
		{name: "пустая строка", board: boardDefault, notation: "", wantErr: true},
		{name: "только буква", board: boardDefault, notation: "e", wantErr: true},
		{name: "только цифра", board: boardDefault, notation: "5", wantErr: true},
		{name: "недопустимый символ", board: boardDefault, notation: "a-4", wantErr: true},
		{name: "ноль", board: boardDefault, notation: "b0", wantErr: true},
		{name: "цифра спереди", board: boardDefault, notation: "2c", wantErr: true},
		// Этот выявил ошибку, тесты сделаны не зря!
		// А еще убрали мертвый код, ну тесты!
		{name: "две клетки", board: boardDefault, notation: "b3b4", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, err := model.ParseCell(tt.notation, tt.board)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantRow, row)
			assert.Equal(t, tt.wantCol, col)
		})
	}
}

func TestLettersToIndex(t *testing.T) {
	tests := []struct {
		name    string
		letters string
		want    int
	}{
		{name: "А заглавная", letters: "A", want: 0},
		{name: "а строчная", letters: "a", want: 0},
		{name: "B", letters: "B", want: 1},
		{name: "Z", letters: "Z", want: 25},
		{name: "АA", letters: "AA", want: 26},
		{name: "AB", letters: "AB", want: 27},
		{name: "BA", letters: "BA", want: 52},
		{name: "ZZ", letters: "ZZ", want: 701},
		{name: "AAA", letters: "AAA", want: 702},
		{name: "ABC", letters: "ABC", want: 730},

		{name: "пустая строка", letters: "", want: -1},
		{name: "цифра", letters: "4", want: -1},
		{name: "буква с цифрой", letters: "C5", want: -1},
		{name: "спецсимвол", letters: "!", want: -1},
		{name: "пробел", letters: " ", want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, model.LettersToIndex(tt.letters))
		})
	}
}
