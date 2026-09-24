package model

type QueenRules struct {
	Color Color
}

func NewQueenRules(color Color) *QueenRules {
	return &QueenRules{Color: color}
}

func (q *QueenRules) CanMove(move *Move, board *Board) bool {
	rowDiff := abs(move.ToRow - move.FromRow)
	colDiff := abs(move.ToCol - move.FromCol)

	if !(move.FromRow == move.ToRow || move.FromCol == move.ToCol || rowDiff == colDiff) {
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
	return enemy == nil || enemy.Color() != q.Color
}
