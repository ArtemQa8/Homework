package model_test

import "mod.go/internal/model"

// emptyBoard создаёт доску rows×cols без фигур.
// Используется в тестах правил фигур, где нужна чистая доска.
func emptyBoard(rows, cols int) *model.Board {
	b := model.NewBoard(rows, cols)
	for row := range rows {
		for col := range cols {
			b.SetPiece(row, col, nil)
		}
	}
	return b
}
