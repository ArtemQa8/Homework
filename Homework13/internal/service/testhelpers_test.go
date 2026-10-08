package service_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"mod.go/internal/model"
)

// createTestGame создаёт игру в моке и возвращает её ID.
func createTestGame(t *testing.T, storage *mockStorage) int {
	t.Helper()
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, err := storage.CreateGame(*game)
	require.NoError(t, err)
	return created.ID()
}

// emptyServiceBoard создаёт доску rows×cols без фигур.
func emptyServiceBoard(rows, cols int) *model.Board {
	b := model.NewBoard(rows, cols)
	for row := range rows {
		for col := range cols {
			b.SetPiece(row, col, nil)
		}
	}
	return b
}
