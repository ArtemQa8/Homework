package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"mod.go/internal/model"
)

// ChooseRandomMove возвращает случайный допустимый ход для текущего игрока.
func ChooseRandomMove(game *model.Game) (*model.Move, error) {
	color := game.CurrentColor()
	rows := game.Board().Rows()
	cols := game.Board().Cols()

	var possible []*model.Move

	for fromRow := 0; fromRow < rows; fromRow++ {
		for fromCol := 0; fromCol < cols; fromCol++ {
			piece := game.Board().PieceAt(fromRow, fromCol)
			if piece == nil || piece.Color() != color {
				continue
			}
			rules := model.RulesFor(piece)
			if rules == nil {
				continue
			}
			for toRow := 0; toRow < rows; toRow++ {
				for toCol := 0; toCol < cols; toCol++ {
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

// MakeAutoMove выбирает случайный ход и делает его, с задержкой (для старого main).
func MakeAutoMove(game *model.Game) (*model.Move, error) {
	move, err := ChooseRandomMove(game)
	if err != nil {
		return nil, err
	}

	delay := time.Duration(2000+rand.Intn(2000)) * time.Millisecond
	time.Sleep(delay)

	return move, nil
}

func SavePlayers(game *model.Game, bus chan<- model.StorageObject) {
	bus <- game.Player1()
	bus <- game.Player2()
}

func SaveGame(game *model.Game, bus chan<- model.StorageObject) {
	bus <- game
}

func SaveMove(game *model.Game, bus chan<- model.StorageObject) {
	moves := game.Moves()
	if len(moves) > 0 {
		last := moves[len(moves)-1]
		bus <- last
	}
}

func SimulateGame(ctx context.Context, state *model.GameState, bus chan<- model.StorageObject) {
	defer func() {
		state.Finished = true
	}()

	for {
		select {
		case <-ctx.Done():
			state.LastMove = "Симуляция остановлена"
			return
		default:
		}

		if state.Game.Checkmate(state.Game.CurrentColor()) {
			winner := "Белые"
			if state.Game.CurrentColor() == model.White {
				winner = "Чёрные"
			}
			state.LastMove = fmt.Sprintf("Мат! Победили %s", winner)
			return
		}

		color := state.Game.CurrentColor()
		var name string
		if color == model.White {
			name = state.Game.Player1().Name()
		} else {
			name = state.Game.Player2().Name()
		}

		start := time.Now()

		move, err := MakeAutoMove(state.Game)
		if err != nil {
			state.LastMove = fmt.Sprintf("Ошибка: %v", err)
			return
		}

		if piece := state.Game.Board().PieceAt(move.FromRow, move.FromCol); piece != nil && piece.Type() == model.Pawn {
			lastRow := 0
			if piece.Color() == model.White {
				lastRow = state.Game.Board().Rows() - 1
			}
			if move.ToRow == lastRow {
				move.Promotion = model.Queen
			}
		}

		if err := state.Game.MakeMove(move); err != nil {
			state.LastMove = fmt.Sprintf("Ошибка: %v", err)
			return
		}
		SaveMove(state.Game, bus)

		elapsed := time.Since(start)
		state.LastMove = fmt.Sprintf("%s (%s)", name, model.FormatMove(*move))
		state.MoveTime = elapsed

		select {
		case <-ctx.Done():
			state.LastMove = "Симуляция остановлена"
			return
		case <-time.After(1 * time.Second):
		}
	}
}
