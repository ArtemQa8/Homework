package service

import (
	"fmt"
	"math/rand"

	"mod.go/internal/model"
)

// ChooseRandomMove возвращает случайный допустимый ход для текущего игрока.
func ChooseRandomMove(game *model.Game) (*model.Move, error) {
	color := game.CurrentColor()
	rows := game.Board().Rows()
	cols := game.Board().Cols()

	var possible []*model.Move

	for fromRow := range rows {
		for fromCol := range cols {
			piece := game.Board().PieceAt(fromRow, fromCol)
			if piece == nil || piece.Color() != color {
				continue
			}
			rules := model.RulesFor(piece)
			if rules == nil {
				continue
			}
			for toRow := range rows {
				for toCol := range cols {
					if fromRow == toRow && fromCol == toCol {
						continue
					}
					move := model.NewMove(fromRow, fromCol, toRow, toCol)
					if rules.CanMove(move, game.Board()) && game.LegalMove(move, color) {
						possible = append(possible, move)
					}
				}
			}
		}
	}

	if len(possible) == 0 {
		return nil, fmt.Errorf("нет допустимых ходов")
	}

	return possible[rand.Intn(len(possible))], nil
}
