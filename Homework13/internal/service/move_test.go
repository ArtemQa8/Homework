package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mod.go/internal/model"
	"mod.go/internal/service"
)

func TestMoveService_Create(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, storage *mockStorage) model.Move
		wantErr bool
	}{
		{
			name: "обычный ход",
			setup: func(t *testing.T, storage *mockStorage) model.Move {
				gameID := createTestGame(t, storage)
				move := model.Move{
					FromRow: 1, FromCol: 4,
					ToRow: 3, ToCol: 4,
				}
				move.SetGameID(gameID)
				return move
			},
		},
		{
			name: "без gameID",
			setup: func(t *testing.T, storage *mockStorage) model.Move {
				return model.Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4}
			},
			wantErr: true,
		},
		{
			name: "игры не существует",
			setup: func(t *testing.T, storage *mockStorage) model.Move {
				move := model.Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4}
				move.SetGameID(999)
				return move
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewMoveService(storage)

			move := tt.setup(t, storage)
			created, err := svc.Create(move)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotZero(t, created.ID(), "ID должен быть присвоен")
			assert.Equal(t, move.FromRow, created.FromRow)
			assert.Equal(t, move.ToRow, created.ToRow)

			got, found := storage.GetMoveByID(created.ID())
			require.True(t, found)
			assert.Equal(t, created.ID(), got.ID())
		})
	}
}

func TestMoveService_Get(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewMoveService(storage)

	gameID := createTestGame(t, storage)
	move := model.Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4}
	move.SetGameID(gameID)
	created, err := svc.Create(move)
	require.NoError(t, err)

	// существующий
	got, err := svc.Get(created.ID())
	require.NoError(t, err)
	assert.Equal(t, created.ID(), got.ID())
	assert.Equal(t, 1, got.FromRow)

	// несуществующий
	_, err = svc.Get(999)
	assert.Error(t, err)
}

func TestMoveService_List(t *testing.T) {
	t.Run("пусто", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewMoveService(storage)

		assert.Empty(t, svc.List())
	})

	t.Run("несколько ходов", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewMoveService(storage)

		gameID := createTestGame(t, storage)

		m1 := model.Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4}
		m1.SetGameID(gameID)
		_, err := svc.Create(m1)
		require.NoError(t, err)

		m2 := model.Move{FromRow: 6, FromCol: 4, ToRow: 4, ToCol: 4}
		m2.SetGameID(gameID)
		_, err = svc.Create(m2)
		require.NoError(t, err)

		list := svc.List()
		require.Len(t, list, 2)
		assert.Equal(t, 1, list[0].FromRow)
		assert.Equal(t, 6, list[1].FromRow)
	})
}

func TestMoveService_Update(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewMoveService(storage)

	gameID := createTestGame(t, storage)
	move := model.Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4}
	move.SetGameID(gameID)
	created, err := svc.Create(move)
	require.NoError(t, err)

	newMove := model.Move{FromRow: 2, FromCol: 3, ToRow: 4, ToCol: 3}
	newMove.SetGameID(gameID)

	updated, err := svc.Update(created.ID(), newMove)
	require.NoError(t, err)
	assert.Equal(t, created.ID(), updated.ID())
	assert.Equal(t, 2, updated.FromRow)
	assert.Equal(t, 3, updated.FromCol)

	_, err = svc.Update(999, newMove)
	assert.Error(t, err)
}

func TestMoveService_Delete(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewMoveService(storage)

	gameID := createTestGame(t, storage)
	move := model.Move{FromRow: 1, FromCol: 4, ToRow: 3, ToCol: 4}
	move.SetGameID(gameID)
	created, err := svc.Create(move)
	require.NoError(t, err)

	err = svc.Delete(created.ID())
	require.NoError(t, err)

	_, err = svc.Get(created.ID())
	assert.Error(t, err)

	err = svc.Delete(created.ID())
	assert.Error(t, err)
}
