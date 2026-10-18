package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/service"
)

func TestPlayerService_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "обычное имя", input: "Иванов"},
		{name: "имя с пробелом", input: "Иван Иванов"},
		{name: "пробелы вокруг", input: "  Иванов  "},
		{name: "пустое имя", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewPlayerService(storage)

			player, err := svc.Create(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.input, player.Name())
			assert.NotZero(t, player.ID(), "ID должен быть присвоен")

			got, found := storage.GetPlayerByID(player.ID())
			require.True(t, found)
			assert.Equal(t, player.Name(), got.Name())
		})
	}
}

func TestPlayerService_Get(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewPlayerService(storage)

	// создаём игрока
	created, err := svc.Create("Иванов")
	require.NoError(t, err)

	// существующий ID
	got, err := svc.Get(created.ID())
	require.NoError(t, err)
	assert.Equal(t, created.ID(), got.ID())
	assert.Equal(t, "Иванов", got.Name())

	// несуществующий ID
	_, err = svc.Get(999)
	assert.Error(t, err)
}

func TestPlayerService_List(t *testing.T) {
	t.Run("пусто", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewPlayerService(storage)

		assert.Empty(t, svc.List())
	})

	t.Run("несколько игроков", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewPlayerService(storage)

		_, err := svc.Create("Иванов")
		require.NoError(t, err)
		_, err = svc.Create("Петров")
		require.NoError(t, err)

		list := svc.List()
		require.Len(t, list, 2)
		assert.Equal(t, "Иванов", list[0].Name())
		assert.Equal(t, "Петров", list[1].Name())
	})
}

func TestPlayerService_Update(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(svc *service.PlayerService) int
		newName string
		wantErr bool
	}{
		{
			name: "существующий игрок",
			setup: func(svc *service.PlayerService) int {
				p, _ := svc.Create("Иванов")
				return p.ID()
			},
			newName: "Сидоров",
			wantErr: false,
		},
		{
			name: "несуществующий ID",
			setup: func(svc *service.PlayerService) int {
				return 999
			},
			newName: "Сидоров",
			wantErr: true,
		},
		{
			name: "пустое имя",
			setup: func(svc *service.PlayerService) int {
				p, _ := svc.Create("Иванов")
				return p.ID()
			},
			newName: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewPlayerService(storage)
			id := tt.setup(svc)

			updated, err := svc.Update(id, tt.newName)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, id, updated.ID())
			assert.Equal(t, tt.newName, updated.Name())
		})
	}
}

func TestPlayerService_Delete(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewPlayerService(storage)

	p, err := svc.Create("Иванов")
	require.NoError(t, err)

	err = svc.Delete(p.ID())
	require.NoError(t, err)

	_, err = svc.Get(p.ID())
	assert.Error(t, err)

	err = svc.Delete(p.ID())
	assert.Error(t, err)
}
