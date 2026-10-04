package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/model"
)

func TestNewMove(t *testing.T) {
	move := model.NewMove(1, 2, 3, 4)

	assert.Equal(t, 1, move.FromRow)
	assert.Equal(t, 2, move.FromCol)
	assert.Equal(t, 3, move.ToRow)
	assert.Equal(t, 4, move.ToCol)

	assert.Equal(t, 0, move.ID())
	assert.Equal(t, 0, move.GameID())

	move.SetID(42)
	move.SetGameID(7)
	assert.Equal(t, 42, move.ID())
	assert.Equal(t, 7, move.GameID())

	assert.Equal(t, "ход", move.ObjectType())
}

func TestMove_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		move *model.Move
		want string
	}{
		{
			name: "минимальный - только координаты",
			move: model.NewMove(1, 2, 3, 4),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4
			}`,
		},
		{
			name: "с Promotion = Pawn - поле пропускается",
			move: func() *model.Move {
				m := model.NewMove(1, 2, 3, 4)
				m.Promotion = model.Pawn
				return m
			}(),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4
			}`,
		},
		{
			name: "с Promotion = Queen",
			move: func() *model.Move {
				m := model.NewMove(6, 4, 7, 4)
				m.Promotion = model.Queen
				return m
			}(),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 6, "отСтолбец": 4,
				"вСтрока": 7, "вСтолбец": 4,
				"превращение": "Ферзь"
			}`,
		},
		{
			name: "со съеденной фигурой",
			move: func() *model.Move {
				m := model.NewMove(1, 2, 3, 4)
				m.Captured = model.NewPiece(model.Black, model.Rook)
				return m
			}(),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4,
				"съедена": {"цвет": "Чёрные", "тип": "Ладья"}
			}`,
		},
		{
			name: "с ходившей фигурой",
			move: func() *model.Move {
				m := model.NewMove(1, 2, 3, 4)
				m.MovedPiece = model.NewPiece(model.White, model.Knight)
				return m
			}(),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4,
				"ходившаяФигура": {"цвет": "Белые", "тип": "Конь"}
			}`,
		},
		{
			name: "с Check",
			move: func() *model.Move {
				m := model.NewMove(1, 2, 3, 4)
				m.Check = true
				return m
			}(),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4,
				"шах": true
			}`,
		},
		{
			name: "с Mate",
			move: func() *model.Move {
				m := model.NewMove(1, 2, 3, 4)
				m.Mate = true
				return m
			}(),
			want: `{
				"играID": 0, "id": 0,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4,
				"мат": true
			}`,
		},
		{
			name: "полный - всё сразу",
			move: func() *model.Move {
				m := model.NewMove(6, 4, 7, 4)
				m.SetID(42)
				m.SetGameID(7)
				m.Promotion = model.Queen
				m.Captured = model.NewPiece(model.Black, model.Rook)
				m.MovedPiece = model.NewPiece(model.White, model.Pawn)
				m.Check = true
				m.Mate = true
				return m
			}(),
			want: `{
				"играID": 7, "id": 42,
				"отСтрока": 6, "отСтолбец": 4,
				"вСтрока": 7, "вСтолбец": 4,
				"превращение": "Ферзь",
				"съедена": {"цвет": "Чёрные", "тип": "Ладья"},
				"ходившаяФигура": {"цвет": "Белые", "тип": "Пешка"},
				"шах": true, "мат": true
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.move)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestMove_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, m model.Move)
		wantErr bool
	}{
		{
			name:  "минимальный",
			input: `{"отСтрока": 1, "отСтолбец": 2, "вСтрока": 3, "вСтолбец": 4}`,
			check: func(t *testing.T, m model.Move) {
				assert.Equal(t, 1, m.FromRow)
				assert.Equal(t, 2, m.FromCol)
				assert.Equal(t, 3, m.ToRow)
				assert.Equal(t, 4, m.ToCol)
				assert.Equal(t, 0, m.ID())
				assert.Equal(t, 0, m.GameID())
				assert.Nil(t, m.Captured)
				assert.Nil(t, m.MovedPiece)
				assert.Equal(t, model.Pawn, m.Promotion)
				assert.False(t, m.Check)
				assert.False(t, m.Mate)
			},
		},
		{
			name: "с promotion",
			input: `{
				"отСтрока": 6, "отСтолбец": 4,
				"вСтрока": 7, "вСтолбец": 4,
				"превращение": "Ферзь"
			}`,
			check: func(t *testing.T, m model.Move) {
				assert.Equal(t, model.Queen, m.Promotion)
			},
		},
		{
			name: "с ID и GameID",
			input: `{
				"играID": 7, "id": 42,
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4
			}`,
			check: func(t *testing.T, m model.Move) {
				assert.Equal(t, 7, m.GameID())
				assert.Equal(t, 42, m.ID())
			},
		},
		{
			name: "с Captured",
			input: `{
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4,
				"съедена": {"цвет": "Чёрные", "тип": "Ладья"}
			}`,
			check: func(t *testing.T, m model.Move) {
				require.NotNil(t, m.Captured)
				assert.Equal(t, model.Black, m.Captured.Color())
				assert.Equal(t, model.Rook, m.Captured.Type())
			},
		},
		{
			name: "с Check и Mate",
			input: `{
				"отСтрока": 1, "отСтолбец": 2,
				"вСтрока": 3, "вСтолбец": 4,
				"шах": true, "мат": true
			}`,
			check: func(t *testing.T, m model.Move) {
				assert.True(t, m.Check)
				assert.True(t, m.Mate)
			},
		},
		{
			name:    "сломанный JSON",
			input:   `{`,
			wantErr: true,
		},
		{
			name:    "неверный тип - строка вместо числа",
			input:   `{"отСтрока": "abc", "отСтолбец": 2, "вСтрока": 3, "вСтолбец": 4}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m model.Move
			err := json.Unmarshal([]byte(tt.input), &m)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, m)
		})
	}
}

func TestMove_RoundTrip(t *testing.T) {
	original := model.NewMove(6, 4, 7, 4)
	original.SetID(42)
	original.SetGameID(7)
	original.Promotion = model.Queen
	original.Captured = model.NewPiece(model.Black, model.Rook)
	original.MovedPiece = model.NewPiece(model.White, model.Pawn)
	original.Check = true
	original.Mate = true

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var restored model.Move
	err = json.Unmarshal(data, &restored)
	require.NoError(t, err)

	assert.Equal(t, original.FromRow, restored.FromRow)
	assert.Equal(t, original.FromCol, restored.FromCol)
	assert.Equal(t, original.ToRow, restored.ToRow)
	assert.Equal(t, original.ToCol, restored.ToCol)
	assert.Equal(t, original.ID(), restored.ID())
	assert.Equal(t, original.GameID(), restored.GameID())
	assert.Equal(t, original.Promotion, restored.Promotion)
	assert.Equal(t, original.Check, restored.Check)
	assert.Equal(t, original.Mate, restored.Mate)

	require.NotNil(t, restored.Captured)
	assert.Equal(t, original.Captured.Color(), restored.Captured.Color())
	assert.Equal(t, original.Captured.Type(), restored.Captured.Type())

	require.NotNil(t, restored.MovedPiece)
	assert.Equal(t, original.MovedPiece.Color(), restored.MovedPiece.Color())
	assert.Equal(t, original.MovedPiece.Type(), restored.MovedPiece.Type())
}
