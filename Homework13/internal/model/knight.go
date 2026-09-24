package model

type KnightRules struct {
	Color Color
}

func NewKnightRules(color Color) *KnightRules {
	return &KnightRules{Color: color}
}

func (k *KnightRules) CanMove(move *Move, board *Board) bool {
	rowDiff := abs(move.ToRow - move.FromRow)
	colDiff := abs(move.ToCol - move.FromCol)

	if !((rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2)) {
		return false
	}
	enemy := board.cells[move.ToRow][move.ToCol]
	return enemy == nil || enemy.Color() != k.Color
}
