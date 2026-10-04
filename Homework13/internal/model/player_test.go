package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/model"
)

func TestNewPlayer(t *testing.T) {
	p := model.NewPlayer("Иванов")

	assert.Equal(t, "Иванов", p.Name())
	assert.Equal(t, 0, p.ID())

	p.SetID(42)
	assert.Equal(t, 42, p.ID())

	assert.Equal(t, "игрок", p.ObjectType())
}

func TestPlayer_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		player model.Player
		want   string
	}{
		{
			name:   "id=0",
			player: *model.NewPlayer("Иванов"),
			want:   `{"id": 0, "имя": "Иванов"}`,
		},
		{
			name: "с id",
			player: func() model.Player {
				p := model.NewPlayer("Петров")
				p.SetID(7)
				return *p
			}(),
			want: `{"id": 7, "имя": "Петров"}`,
		},
		{
			name:   "пустое имя",
			player: *model.NewPlayer(""),
			want:   `{"id": 0, "имя": ""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.player)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestPlayer_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  int
		wantNam string
		wantErr bool
	}{
		{name: "обычный", input: `{"id": 42, "имя": "Иванов"}`, wantID: 42, wantNam: "Иванов"},
		{name: "пустое имя", input: `{"id": 7, "имя": ""}`, wantID: 7, wantNam: ""},
		{name: "только имя", input: `{"имя": "Петров"}`, wantID: 0, wantNam: "Петров"},
		{name: "только id", input: `{"id": 5}`, wantID: 5, wantNam: ""},
		{name: "пустой объект", input: `{}`, wantID: 0, wantNam: ""},
		{name: "сломанный JSON", input: `{`, wantErr: true},
		{name: "неверный тип id", input: `{"id": "abc", "имя": "Иванов"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p model.Player
			err := json.Unmarshal([]byte(tt.input), &p)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, p.ID())
			assert.Equal(t, tt.wantNam, p.Name())
		})
	}
}

func TestPlayer_RoundTrip(t *testing.T) {
	original := model.NewPlayer("Иванов")
	original.SetID(42)

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var restored model.Player
	err = json.Unmarshal(data, &restored)
	require.NoError(t, err)

	assert.Equal(t, original.ID(), restored.ID())
	assert.Equal(t, original.Name(), restored.Name())
}
