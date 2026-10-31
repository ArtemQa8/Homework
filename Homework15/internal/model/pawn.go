package model

type PawnRules struct {
	Color Color
}

func NewPawnRules(color Color) *PawnRules {
	return &PawnRules{Color: color}
}

func (p *PawnRules) CanMove(move *Move, board *Board) bool {
	shift := 1
	startRow := 1
	if p.Color == Black {
		shift = -1
		startRow = board.Rows() - 2
	}

	if move.FromCol == move.ToCol {
		if move.ToRow == move.FromRow+shift {
			return board.cells[move.ToRow][move.ToCol] == nil
		}
		if move.ToRow == move.FromRow+2*shift && move.FromRow == startRow {
			return board.cells[move.FromRow+shift][move.FromCol] == nil &&
				board.cells[move.ToRow][move.ToCol] == nil
		}
	}
	if (move.ToCol == move.FromCol+1 || move.ToCol == move.FromCol-1) &&
		move.ToRow == move.FromRow+shift {
		enemy := board.cells[move.ToRow][move.ToCol]
		return enemy != nil && enemy.Color() != p.Color
	}
	return false
}
