package model

type RookRules struct {
	Color Color
}

func NewRookRules(color Color) *RookRules {
	return &RookRules{Color: color}
}

func (r *RookRules) CanMove(move *Move, board *Board) bool {
	if move.FromRow != move.ToRow && move.FromCol != move.ToCol {
		return false
	}
	rowStep := sign(move.ToRow - move.FromRow)
	colStep := sign(move.ToCol - move.FromCol)

	rr, cc := move.FromRow+rowStep, move.FromCol+colStep
	for rr != move.ToRow || cc != move.ToCol {
		if board.cells[rr][cc] != nil {
			return false
		}
		rr += rowStep
		cc += colStep
	}
	enemy := board.cells[move.ToRow][move.ToCol]
	return enemy == nil || enemy.Color() != r.Color
}
