package grpc

import (
	"mod.go/internal/model"
	pb "mod.go/internal/proto"
)

// ============ Цвет ============

func toProtoColor(c model.Color) pb.Color {
	switch c {
	case model.White:
		return pb.Color_WHITE
	case model.Black:
		return pb.Color_BLACK
	default:
		return pb.Color_COLOR_UNSPECIFIED
	}
}

// ============ Тип фигуры ============

func toProtoPieceType(t model.PieceType) pb.PieceType {
	switch t {
	case model.Pawn:
		return pb.PieceType_PAWN
	case model.Rook:
		return pb.PieceType_ROOK
	case model.Knight:
		return pb.PieceType_KNIGHT
	case model.Bishop:
		return pb.PieceType_BISHOP
	case model.Queen:
		return pb.PieceType_QUEEN
	case model.King:
		return pb.PieceType_KING
	default:
		return pb.PieceType_PIECE_TYPE_UNSPECIFIED
	}
}

func fromProtoPieceType(t pb.PieceType) model.PieceType {
	switch t {
	case pb.PieceType_PAWN:
		return model.Pawn
	case pb.PieceType_ROOK:
		return model.Rook
	case pb.PieceType_KNIGHT:
		return model.Knight
	case pb.PieceType_BISHOP:
		return model.Bishop
	case pb.PieceType_QUEEN:
		return model.Queen
	case pb.PieceType_KING:
		return model.King
	default:
		return model.Pawn
	}
}

// ============ Фигура ============

func toProtoPiece(p *model.Piece) *pb.Piece {
	if p == nil {
		return nil
	}
	return &pb.Piece{
		Color: toProtoColor(p.Color()),
		Type:  toProtoPieceType(p.Type()),
	}
}

// ============ Доска ============

func toProtoBoard(b *model.Board) *pb.Board {
	if b == nil {
		return nil
	}
	grid := make([]*pb.Row, b.Rows())
	for r := 0; r < b.Rows(); r++ {
		cells := make([]*pb.Cell, b.Cols())
		for c := 0; c < b.Cols(); c++ {
			cells[c] = &pb.Cell{
				Piece: toProtoPiece(b.PieceAt(r, c)),
			}
		}
		grid[r] = &pb.Row{Cells: cells}
	}
	return &pb.Board{
		Rows: int32(b.Rows()),
		Cols: int32(b.Cols()),
		Grid: grid,
	}
}

// ============ Игрок ============

func toProtoPlayer(p model.Player) *pb.Player {
	return &pb.Player{
		Id:   int32(p.ID()),
		Name: p.Name(),
	}
}

// ============ Ход ============

func toProtoPromotion(m model.Move) pb.PieceType {
	if m.Promotion == model.Pawn {
		return pb.PieceType_PIECE_TYPE_UNSPECIFIED
	}
	return toProtoPieceType(m.Promotion)
}

func toProtoMove(m model.Move) *pb.Move {
	return &pb.Move{
		Id:         int32(m.ID()),
		GameId:     int32(m.GameID()),
		FromRow:    int32(m.FromRow),
		FromCol:    int32(m.FromCol),
		ToRow:      int32(m.ToRow),
		ToCol:      int32(m.ToCol),
		Captured:   toProtoPiece(m.Captured),
		Promotion:  toProtoPromotion(m),
		MovedPiece: toProtoPiece(m.MovedPiece),
		Check:      m.Check,
		Mate:       m.Mate,
	}
}

// ============ Игра ============

func toProtoGame(g model.Game) *pb.Game {
	moves := make([]*pb.Move, len(g.Moves()))
	for i, m := range g.Moves() {
		moves[i] = toProtoMove(m)
	}

	return &pb.Game{
		Id:           int32(g.ID()),
		Player1:      toProtoPlayer(g.Player1()),
		Player2:      toProtoPlayer(g.Player2()),
		Player1Color: toProtoColor(g.Player1Color()),
		Player2Color: toProtoColor(g.Player2Color()),
		Board:        toProtoBoard(g.Board()),
		Current:      toProtoColor(g.CurrentColor()),
		Moves:        moves,
	}
}
