package model

type KingRules struct {
	Color Color
}

func NewKingRules(color Color) *KingRules {
	return &KingRules{Color: color}
}

func (k *KingRules) CanMove(move *Move, board *Board) bool {
	rowDiff := abs(move.ToRow - move.FromRow)
	colDiff := abs(move.ToCol - move.FromCol)

	if rowDiff > 1 || colDiff > 1 {
		return false
	}
	enemy := board.cells[move.ToRow][move.ToCol]
	return enemy == nil || enemy.Color() != k.Color
}
