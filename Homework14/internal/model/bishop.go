package model

type BishopRules struct {
	Color Color
}

func NewBishopRules(color Color) *BishopRules {
	return &BishopRules{Color: color}
}

func (b *BishopRules) CanMove(move *Move, board *Board) bool {
	rowDiff := abs(move.ToRow - move.FromRow)
	colDiff := abs(move.ToCol - move.FromCol)

	if rowDiff != colDiff || rowDiff == 0 {
		return false
	}

	rowStep := sign(move.ToRow - move.FromRow)
	colStep := sign(move.ToCol - move.FromCol)

	r, c := move.FromRow+rowStep, move.FromCol+colStep
	for r != move.ToRow || c != move.ToCol {
		if board.cells[r][c] != nil {
			return false
		}
		r += rowStep
		c += colStep
	}

	enemy := board.cells[move.ToRow][move.ToCol]
	return enemy == nil || enemy.Color() != b.Color
}
