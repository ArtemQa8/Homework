package model

type PieceWithRules interface {
	CanMove(move *Move, board *Board) bool
}

func RulesFor(piece *Piece) PieceWithRules {
	switch piece.Type() {
	case Pawn:
		return &PawnRules{Color: piece.Color()}
	case Rook:
		return &RookRules{Color: piece.Color()}
	case Knight:
		return &KnightRules{Color: piece.Color()}
	case Bishop:
		return &BishopRules{Color: piece.Color()}
	case Queen:
		return &QueenRules{Color: piece.Color()}
	case King:
		return &KingRules{Color: piece.Color()}
	default:
		return nil
	}
}
