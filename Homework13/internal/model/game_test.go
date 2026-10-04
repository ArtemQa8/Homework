package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/model"
)

func TestNewGame(t *testing.T) {
	game := model.NewGame("Иван Иванов", "Пётр Петров", 8, 8)

	assert.Equal(t, "Иван Иванов", game.Player1().Name())
	assert.Equal(t, "Пётр Петров", game.Player2().Name())
	assert.Equal(t, model.White, game.Player1Color())
	assert.Equal(t, model.Black, game.Player2Color())
	assert.Equal(t, model.White, game.CurrentColor())
	assert.Equal(t, 0, game.ID())

	require.NotNil(t, game.Board())
	assert.Equal(t, 8, game.Board().Rows())
	assert.Equal(t, 8, game.Board().Cols())

	assert.Empty(t, game.Moves())
}

func TestGame_Methods(t *testing.T) {
	game := model.NewGame("Иван Иванов", "Пётр Петров", 8, 8)

	game.SetID(12)
	assert.Equal(t, 12, game.ID())

	newP1 := model.NewPlayer("Олег Сидоров")
	newP2 := model.NewPlayer("Дмитрий Кузнецов")
	game.SetPlayer1(*newP1)
	game.SetPlayer2(*newP2)
	assert.Equal(t, "Олег Сидоров", game.Player1().Name())
	assert.Equal(t, "Дмитрий Кузнецов", game.Player2().Name())

	game.SetPlayer1Color(model.Black)
	game.SetPlayer2Color(model.White)
	assert.Equal(t, model.Black, game.Player1Color())
	assert.Equal(t, model.White, game.Player2Color())

	newBoard := model.NewBoard(4, 4)
	game.SetBoard(newBoard)
	assert.Equal(t, 4, game.Board().Rows())
	assert.Equal(t, 4, game.Board().Cols())
}

func TestGame_MakeMove_StartMoves(t *testing.T) {
	tests := []struct {
		name      string
		fromRow   int
		fromCol   int
		toRow     int
		toCol     int
		wantColor model.Color
	}{
		{name: "пешка: e2-e4 (на 2)", fromRow: 1, fromCol: 4, toRow: 3, toCol: 4, wantColor: model.Black},
		{name: "пешка: e2-e3 (на 1)", fromRow: 1, fromCol: 4, toRow: 2, toCol: 4, wantColor: model.Black},
		{name: "пешка: d2-d4", fromRow: 1, fromCol: 3, toRow: 3, toCol: 3, wantColor: model.Black},
		{name: "пешка: a2-a3", fromRow: 1, fromCol: 0, toRow: 2, toCol: 0, wantColor: model.Black},
		{name: "конь: g1-f3", fromRow: 0, fromCol: 6, toRow: 2, toCol: 5, wantColor: model.Black},
		{name: "конь: b1-c3", fromRow: 0, fromCol: 1, toRow: 2, toCol: 2, wantColor: model.Black},
		{name: "конь: b1-a3", fromRow: 0, fromCol: 1, toRow: 2, toCol: 0, wantColor: model.Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := model.NewGame("A", "B", 8, 8)

			require.Empty(t, game.Moves())
			require.Equal(t, model.White, game.CurrentColor())

			require.NotNil(t, game.Board().PieceAt(tt.fromRow, tt.fromCol))
			require.Nil(t, game.Board().PieceAt(tt.toRow, tt.toCol))

			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			err := game.MakeMove(move)
			require.NoError(t, err)

			moves := game.Moves()
			require.Len(t, moves, 1)
			last := moves[0]
			assert.Equal(t, tt.fromRow, last.FromRow)
			assert.Equal(t, tt.fromCol, last.FromCol)
			assert.Equal(t, tt.toRow, last.ToRow)
			assert.Equal(t, tt.toCol, last.ToCol)
			assert.NotNil(t, last.MovedPiece)
			assert.False(t, last.Check)
			assert.False(t, last.Mate)

			assert.Equal(t, tt.wantColor, game.CurrentColor())

			assert.Nil(t, game.Board().PieceAt(tt.fromRow, tt.fromCol))
			moved := game.Board().PieceAt(tt.toRow, tt.toCol)
			require.NotNil(t, moved)
		})
	}
}

func TestGame_MakeMove_AllPieces(t *testing.T) {
	tests := []struct {
		name      string
		pieceType model.PieceType
		fromRow   int
		fromCol   int
		toRow     int
		toCol     int
	}{
		{name: "ладья: горизонталь", pieceType: model.Rook,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 7},
		{name: "ладья: вертикаль", pieceType: model.Rook,
			fromRow: 4, fromCol: 4, toRow: 1, toCol: 4},
		{name: "слон: диагональ вниз", pieceType: model.Bishop,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 6},
		{name: "слон: диагональ вверх", pieceType: model.Bishop,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 2},
		{name: "ферзь: вертикаль", pieceType: model.Queen,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 4},
		{name: "ферзь: диагональ", pieceType: model.Queen,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 6},
		{name: "король: 1 вперёд", pieceType: model.King,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 4},
		{name: "король: 1 по диагонали", pieceType: model.King,
			fromRow: 4, fromCol: 4, toRow: 5, toCol: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			board := emptyBoard(8, 8)
			board.SetPiece(tt.fromRow, tt.fromCol,
				model.NewPiece(model.White, tt.pieceType))

			game := model.NewGame("A", "B", 8, 8)
			game.SetBoard(board)

			require.Equal(t, model.White, game.CurrentColor())

			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			err := game.MakeMove(move)
			require.NoError(t, err)

			assert.Nil(t, board.PieceAt(tt.fromRow, tt.fromCol))
			moved := board.PieceAt(tt.toRow, tt.toCol)
			require.NotNil(t, moved)
			assert.Equal(t, tt.pieceType, moved.Type())
			assert.Equal(t, model.White, moved.Color())

			assert.Equal(t, model.Black, game.CurrentColor())

			moves := game.Moves()
			require.Len(t, moves, 1)
			assert.Equal(t, tt.pieceType, moves[0].MovedPiece.Type())
		})
	}
}

func TestGame_MakeMove_BlackStartMoves(t *testing.T) {
	tests := []struct {
		name    string
		fromRow int
		fromCol int
		toRow   int
		toCol   int
	}{
		{name: "пешка: e7-e5 (на 2)", fromRow: 6, fromCol: 4, toRow: 4, toCol: 4},
		{name: "пешка: e7-e6 (на 1)", fromRow: 6, fromCol: 4, toRow: 5, toCol: 4},
		{name: "пешка: d7-d5", fromRow: 6, fromCol: 3, toRow: 4, toCol: 3},
		{name: "пешка: a7-a6", fromRow: 6, fromCol: 0, toRow: 5, toCol: 0},
		{name: "конь: g8-f6", fromRow: 7, fromCol: 6, toRow: 5, toCol: 5},
		{name: "конь: b8-c6", fromRow: 7, fromCol: 1, toRow: 5, toCol: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := model.NewGame("A", "B", 8, 8)

			whiteMove := model.NewMove(1, 4, 3, 4) // e2-e4
			require.NoError(t, game.MakeMove(whiteMove))
			require.Equal(t, model.Black, game.CurrentColor())

			require.NotNil(t, game.Board().PieceAt(tt.fromRow, tt.fromCol))
			require.Nil(t, game.Board().PieceAt(tt.toRow, tt.toCol))

			blackMove := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			err := game.MakeMove(blackMove)
			require.NoError(t, err)

			assert.Equal(t, model.White, game.CurrentColor())

			moves := game.Moves()
			require.Len(t, moves, 2)

			last := moves[1]
			assert.Equal(t, tt.fromRow, last.FromRow)
			assert.Equal(t, tt.fromCol, last.FromCol)
			assert.Equal(t, tt.toRow, last.ToRow)
			assert.Equal(t, tt.toCol, last.ToCol)
			require.NotNil(t, last.MovedPiece)
			assert.Equal(t, model.Black, last.MovedPiece.Color())

			assert.Nil(t, game.Board().PieceAt(tt.fromRow, tt.fromCol))
			moved := game.Board().PieceAt(tt.toRow, tt.toCol)
			require.NotNil(t, moved)
			assert.Equal(t, model.Black, moved.Color())
		})
	}
}

func TestGame_MakeMove_FewMoves(t *testing.T) {
	game := model.NewGame("A", "B", 8, 8)

	sequence := []struct {
		name      string
		fromRow   int
		fromCol   int
		toRow     int
		toCol     int
		wantColor model.Color
	}{
		{name: "1. e2-e4 (белые)", fromRow: 1, fromCol: 4, toRow: 3, toCol: 4, wantColor: model.Black},
		{name: "2. e7-e5 (чёрные)", fromRow: 6, fromCol: 4, toRow: 4, toCol: 4, wantColor: model.White},
		{name: "3. g1-f3 (белые)", fromRow: 0, fromCol: 6, toRow: 2, toCol: 5, wantColor: model.Black},
		{name: "4. b8-c6 (чёрные)", fromRow: 7, fromCol: 1, toRow: 5, toCol: 2, wantColor: model.White},
		{name: "5. f1-c4 (белые)", fromRow: 0, fromCol: 5, toRow: 3, toCol: 2, wantColor: model.Black},
		{name: "6. g8-f6 (чёрные)", fromRow: 7, fromCol: 6, toRow: 5, toCol: 5, wantColor: model.White},
	}

	for i, tt := range sequence {
		t.Run(tt.name, func(t *testing.T) {
			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			err := game.MakeMove(move)
			require.NoError(t, err)

			assert.Equal(t, tt.wantColor, game.CurrentColor())

			moves := game.Moves()
			require.Len(t, moves, i+1)

			assert.Nil(t, game.Board().PieceAt(tt.fromRow, tt.fromCol))

			moved := game.Board().PieceAt(tt.toRow, tt.toCol)
			require.NotNil(t, moved)
		})
	}

	assert.Len(t, game.Moves(), 6)
	assert.Equal(t, model.White, game.CurrentColor())
}

func TestGame_MakeMove_Errors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *model.Game
		move    *model.Move
		wantErr bool
	}{
		// стандартная доска - не важно, где фигура
		{
			name: "выход за границы (toRow=99)",
			setup: func() *model.Game {
				return model.NewGame("A", "B", 8, 8)
			},
			move:    model.NewMove(1, 4, 99, 4),
			wantErr: true,
		},
		{
			name: "выход за границы (fromCol=-1)",
			setup: func() *model.Game {
				return model.NewGame("A", "B", 8, 8)
			},
			move:    model.NewMove(1, -1, 3, 4),
			wantErr: true,
		},
		{
			name: "стартовая клетка пуста",
			setup: func() *model.Game {
				return model.NewGame("A", "B", 8, 8)
			},
			move:    model.NewMove(4, 4, 5, 4), // e5 пусто
			wantErr: true,
		},
		{
			name: "фигура чужого цвета при ходе белых",
			setup: func() *model.Game {
				return model.NewGame("A", "B", 8, 8)
			},
			move:    model.NewMove(6, 4, 4, 4), // e7-e5 чёрной пешкой
			wantErr: true,
		},
		// кастомная доска
		{
			name: "пешка назад",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Pawn))
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:    model.NewMove(4, 4, 3, 4), // пешка назад
			wantErr: true,
		},
		{
			name: "ладья по диагонали",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Rook))
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:    model.NewMove(4, 4, 2, 2), // ладья по диагонали
			wantErr: true,
		},
		{
			name: "конь прямо",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Knight))
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:    model.NewMove(4, 4, 6, 4), // конь прямо
			wantErr: true,
		},
		{
			name: "ход оставляет короля под шахом - откат",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(1, 4, model.NewPiece(model.White, model.Rook)) // заслон
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.Rook)) // бьёт по вертикали

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:    model.NewMove(1, 4, 1, 3), // ладья уходит с заслона
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()

			movesBefore := len(game.Moves())
			colorBefore := game.CurrentColor()

			err := game.MakeMove(tt.move)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, game.Moves(), movesBefore,
				"история ходов не должна пополняться при ошибке")
			assert.Equal(t, colorBefore, game.CurrentColor(),
				"цвет не должен переключаться при ошибке")
		})
	}
}

func TestGame_MakeMove_Captures(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() *model.Game
		move         *model.Move
		wantCaptured *model.Piece
		wantErr      bool
	}{
		{
			name: "пешка бьёт врага по диагонали",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(3, 4, model.NewPiece(model.White, model.Pawn)) // e4
				board.SetPiece(4, 5, model.NewPiece(model.Black, model.Pawn)) // f5
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:         model.NewMove(3, 4, 4, 5), // e4xf5
			wantCaptured: model.NewPiece(model.Black, model.Pawn),
		},

		{
			name: "ладья бьёт врага по горизонтали",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 0, model.NewPiece(model.White, model.Rook))  // a5
				board.SetPiece(4, 4, model.NewPiece(model.Black, model.Queen)) // e5
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:         model.NewMove(4, 0, 4, 4), // a5xe5
			wantCaptured: model.NewPiece(model.Black, model.Queen),
		},

		{
			name: "не бьёт свою фигуру",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 0, model.NewPiece(model.White, model.Rook))
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Pawn)) // своя
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:    model.NewMove(4, 0, 4, 4),
			wantErr: true,
		},

		{
			name: "конь бьёт через фигуры",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Knight)) // e5
				board.SetPiece(4, 5, model.NewPiece(model.Black, model.Pawn))   // f5 — блок справа
				board.SetPiece(3, 4, model.NewPiece(model.Black, model.Pawn))   // e4 — блок сверху
				board.SetPiece(2, 5, model.NewPiece(model.Black, model.Queen))  // f6 — цель взятия
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:         model.NewMove(4, 4, 2, 5), // e5xf6, через блоки
			wantCaptured: model.NewPiece(model.Black, model.Queen),
		},

		{
			name: "пешка не бьёт прямо (враг впереди)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(3, 4, model.NewPiece(model.White, model.Pawn)) // e4
				board.SetPiece(4, 4, model.NewPiece(model.Black, model.Pawn)) // e5 — прямо впереди
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			move:    model.NewMove(3, 4, 4, 4),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()

			err := game.MakeMove(tt.move)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			moved := game.Board().PieceAt(tt.move.ToRow, tt.move.ToCol)
			require.NotNil(t, moved)

			moves := game.Moves()
			require.Len(t, moves, 1)
			last := moves[0]

			require.NotNil(t, last.Captured, "Captured должен быть заполнен")
			assert.Equal(t, tt.wantCaptured.Type(), last.Captured.Type())
			assert.Equal(t, tt.wantCaptured.Color(), last.Captured.Color())
		})
	}
}

func TestGame_MakeMove_Castling(t *testing.T) {
	tests := []struct {
		name        string
		kingToCol   int
		rookFromCol int
		rookToCol   int
	}{
		{name: "короткая O-O", kingToCol: 6, rookFromCol: 7, rookToCol: 5},
		{name: "длинная O-O-O", kingToCol: 2, rookFromCol: 0, rookToCol: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := emptyBoard(8, 8)
			board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
			board.SetPiece(0, tt.rookFromCol, model.NewPiece(model.White, model.Rook))

			game := model.NewGame("A", "B", 8, 8)
			game.SetBoard(board)

			move := model.NewMove(0, 4, 0, tt.kingToCol)
			err := game.MakeMove(move)
			require.NoError(t, err)

			king := board.PieceAt(0, tt.kingToCol)
			require.NotNil(t, king)
			assert.Equal(t, model.King, king.Type())

			rook := board.PieceAt(0, tt.rookToCol)
			require.NotNil(t, rook)
			assert.Equal(t, model.Rook, rook.Type())

			assert.Nil(t, board.PieceAt(0, 4))              // e1
			assert.Nil(t, board.PieceAt(0, tt.rookFromCol)) // где была ладья

			moves := game.Moves()
			require.Len(t, moves, 1)
			require.NotNil(t, moves[0].MovedPiece)
			assert.Equal(t, model.King, moves[0].MovedPiece.Type())
		})
	}
}

func TestGame_MakeMove_Castling_Black(t *testing.T) {
	tests := []struct {
		name        string
		kingToCol   int
		rookFromCol int
		rookToCol   int
	}{
		{name: "чёрная короткая O-O", kingToCol: 6, rookFromCol: 7, rookToCol: 5},
		{name: "чёрная длинная O-O-O", kingToCol: 2, rookFromCol: 0, rookToCol: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := emptyBoard(8, 8)

			board.SetPiece(4, 0, model.NewPiece(model.White, model.Rook)) // a5

			board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
			board.SetPiece(7, tt.rookFromCol, model.NewPiece(model.Black, model.Rook))

			game := model.NewGame("A", "B", 8, 8)
			game.SetBoard(board)

			whiteMove := model.NewMove(4, 0, 3, 0)
			require.NoError(t, game.MakeMove(whiteMove))
			require.Equal(t, model.Black, game.CurrentColor())

			move := model.NewMove(7, 4, 7, tt.kingToCol)
			err := game.MakeMove(move)
			require.NoError(t, err)

			king := board.PieceAt(7, tt.kingToCol)
			require.NotNil(t, king)
			assert.Equal(t, model.King, king.Type())
			assert.Equal(t, model.Black, king.Color())

			rook := board.PieceAt(7, tt.rookToCol)
			require.NotNil(t, rook)
			assert.Equal(t, model.Rook, rook.Type())

			assert.Nil(t, board.PieceAt(7, 4)) // e8
			assert.Nil(t, board.PieceAt(7, tt.rookFromCol))

			moves := game.Moves()
			require.Len(t, moves, 2)
			require.NotNil(t, moves[1].MovedPiece)
			assert.Equal(t, model.King, moves[1].MovedPiece.Type())
		})
	}
}

func TestGame_MakeMove_Castling_Restrictions(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *model.Game
		fromCol int
		toCol   int
	}{
		{
			name: "путь заблокирован (g1)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))   // e1
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))   // h1
				board.SetPiece(0, 6, model.NewPiece(model.White, model.Bishop)) // g1 — блок
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "путь заблокирован (f1)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))
				board.SetPiece(0, 5, model.NewPiece(model.White, model.Bishop)) // f1 — блок
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "нет ладьи на h1",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				// ладью НЕ ставим
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "король под шахом",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.Rook)) // e8 — бьёт e1
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "промежуточная f1 под боем",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))
				board.SetPiece(7, 5, model.NewPiece(model.Black, model.Rook)) // f8 — бьёт f1
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "конечное поле g1 под боем",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))
				// чёрная ладья бьёт g1, но НЕ f1 и НЕ e1
				board.SetPiece(7, 6, model.NewPiece(model.Black, model.Rook)) // g8 — бьёт g1
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "король ходил (e1-f1-e1)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // чтобы чёрные могли ходить
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// туда-обратно
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 5))) // Ke1-f1
				require.NoError(t, game.MakeMove(model.NewMove(7, 4, 7, 3))) // Ke8-d8
				require.NoError(t, game.MakeMove(model.NewMove(0, 5, 0, 4))) // Kf1-e1
				require.NoError(t, game.MakeMove(model.NewMove(7, 3, 7, 4))) // Kd8-e8
				return game
			},
			fromCol: 4, toCol: 6,
		},
		{
			name: "ладья ходила (h1-h2-h1)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook))
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				require.NoError(t, game.MakeMove(model.NewMove(0, 7, 1, 7))) // Rh1-h2
				require.NoError(t, game.MakeMove(model.NewMove(7, 4, 7, 3))) // Ke8-d8
				require.NoError(t, game.MakeMove(model.NewMove(1, 7, 0, 7))) // Rh2-h1
				require.NoError(t, game.MakeMove(model.NewMove(7, 3, 7, 4))) // Kd8-e8
				return game
			},
			fromCol: 4, toCol: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()

			movesBefore := len(game.Moves())
			colorBefore := game.CurrentColor()

			move := model.NewMove(0, tt.fromCol, 0, tt.toCol)
			err := game.MakeMove(move)
			assert.Error(t, err, "рокировка должна быть отклонена")

			// состояние не изменилось
			assert.Len(t, game.Moves(), movesBefore,
				"история ходов не должна пополняться")
			assert.Equal(t, colorBefore, game.CurrentColor(),
				"цвет не должен переключаться")
		})
	}
}

func TestGame_MakeMove_Castling_Restrictions_Black(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *model.Game
		toCol int
	}{
		{
			name: "чёрный король ходил (e8-d8-e8)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.Rook)) // h8
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3))) // Ke1-d1
				require.NoError(t, game.MakeMove(model.NewMove(7, 4, 7, 3))) // Ke8-d8
				require.NoError(t, game.MakeMove(model.NewMove(0, 3, 0, 4))) // Kd1-e1
				require.NoError(t, game.MakeMove(model.NewMove(7, 3, 7, 4))) // Kd8-e8
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3))) // Ke1-d1
				return game
			},
			toCol: 6,
		},
		{
			name: "чёрная ладья H ходила (h8-h7-h8)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.Rook)) // h8
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3))) // Ke1-d1
				require.NoError(t, game.MakeMove(model.NewMove(7, 7, 6, 7))) // Rh8-h7
				require.NoError(t, game.MakeMove(model.NewMove(0, 3, 0, 4))) // Kd1-e1
				require.NoError(t, game.MakeMove(model.NewMove(6, 7, 7, 7))) // Rh7-h8
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3))) // Ke1-d1
				return game
			},
			toCol: 6,
		},
		{
			name: "чёрная ладья A ходила (a8-a7-a8)",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
				board.SetPiece(7, 0, model.NewPiece(model.Black, model.Rook)) // a8
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3))) // Ke1-d1
				require.NoError(t, game.MakeMove(model.NewMove(7, 0, 6, 0))) // Ra8-a7
				require.NoError(t, game.MakeMove(model.NewMove(0, 3, 0, 4))) // Kd1-e1
				require.NoError(t, game.MakeMove(model.NewMove(6, 0, 7, 0))) // Ra7-a8
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3))) // Ke1-d1
				return game
			},
			toCol: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()

			movesBefore := len(game.Moves())
			colorBefore := game.CurrentColor()

			move := model.NewMove(7, 4, 7, tt.toCol)
			err := game.MakeMove(move)
			assert.Error(t, err, "рокировка должна быть отклонена")

			assert.Len(t, game.Moves(), movesBefore,
				"история ходов не должна пополняться")
			assert.Equal(t, colorBefore, game.CurrentColor(),
				"цвет не должен переключаться")
		})
	}
}

func TestGame_MakeMove_Castling_SmallBoards(t *testing.T) {
	tests := []struct {
		name    string
		rows    int
		cols    int
		kingCol int
		wantErr bool
	}{
		{name: "4x4 - король уходит за доску", rows: 4, cols: 4, kingCol: 2, wantErr: true},
		{name: "5x5", rows: 5, cols: 5, kingCol: 2, wantErr: false},
		{name: "6x6", rows: 6, cols: 6, kingCol: 3, wantErr: false},
		{name: "7x7", rows: 7, cols: 7, kingCol: 3, wantErr: false},
		{name: "8x8 (стандарт)", rows: 8, cols: 8, kingCol: 4, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := emptyBoard(tt.rows, tt.cols)
			board.SetPiece(0, tt.kingCol, model.NewPiece(model.White, model.King))
			board.SetPiece(0, tt.cols-1, model.NewPiece(model.White, model.Rook))

			game := model.NewGame("A", "B", tt.rows, tt.cols)
			game.SetBoard(board)

			toCol := tt.kingCol + 2
			move := model.NewMove(0, tt.kingCol, 0, toCol)
			err := game.MakeMove(move)

			if tt.wantErr {
				assert.Error(t, err, "ожидаем ошибку на %dx%d", tt.rows, tt.cols)
				return
			}

			require.NoError(t, err, "рокировка должна пройти на %dx%d", tt.rows, tt.cols)

			newKing := board.PieceAt(0, toCol)
			require.NotNil(t, newKing)
			assert.Equal(t, model.King, newKing.Type())

			newRook := board.PieceAt(0, toCol-1)
			require.NotNil(t, newRook)
			assert.Equal(t, model.Rook, newRook.Type())

			assert.Nil(t, board.PieceAt(0, tt.kingCol))
		})
	}
}

func TestGame_MakeMove_EnPassant(t *testing.T) {
	type checkFn func(t *testing.T, board *model.Board, move *model.Move)

	tests := []struct {
		name    string
		setup   func() (*model.Game, *model.Move)
		wantErr bool
		check   checkFn
	}{
		{
			name: "белый en passant - классический",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Pawn)) // e5
				board.SetPiece(6, 3, model.NewPiece(model.Black, model.Pawn)) // d7
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// 1. Белые: Ke1-d1 (нейтральный)
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3)))
				// 2. Чёрные: d7-d5 (двойной)
				require.NoError(t, game.MakeMove(model.NewMove(6, 3, 4, 3)))

				// 3. Белые: e5xd6 (en passant)
				return game, model.NewMove(4, 4, 5, 3)
			},
			check: func(t *testing.T, board *model.Board, move *model.Move) {
				ourPawn := board.PieceAt(5, 3)
				require.NotNil(t, ourPawn)
				assert.Equal(t, model.Pawn, ourPawn.Type())
				assert.Equal(t, model.White, ourPawn.Color())

				assert.Nil(t, board.PieceAt(4, 4))
				assert.Nil(t, board.PieceAt(4, 3))

				require.NotNil(t, move.Captured)
				assert.Equal(t, model.Pawn, move.Captured.Type())
				assert.Equal(t, model.Black, move.Captured.Color())
			},
		},
		{
			name: "чёрный en passant - зеркально",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(3, 3, model.NewPiece(model.Black, model.Pawn)) // d4
				board.SetPiece(1, 4, model.NewPiece(model.White, model.Pawn)) // e2

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// 1. Белые: e2-e4 (двойной)
				require.NoError(t, game.MakeMove(model.NewMove(1, 4, 3, 4)))

				// 2. Чёрные: d4xe3 (en passant)
				return game, model.NewMove(3, 3, 2, 4)
			},
			check: func(t *testing.T, board *model.Board, move *model.Move) {
				ourPawn := board.PieceAt(2, 4)
				require.NotNil(t, ourPawn)
				assert.Equal(t, model.Pawn, ourPawn.Type())
				assert.Equal(t, model.Black, ourPawn.Color())

				assert.Nil(t, board.PieceAt(3, 3))
				assert.Nil(t, board.PieceAt(3, 4))

				require.NotNil(t, move.Captured)
				assert.Equal(t, model.Pawn, move.Captured.Type())
				assert.Equal(t, model.White, move.Captured.Color())
			},
		},
		{
			name: "враг ходил на 1 клетку - не en passant",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Pawn)) // e5
				board.SetPiece(5, 3, model.NewPiece(model.Black, model.Pawn)) // d6
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// 1. Белые: Ke1-d1
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3)))
				// 2. Чёрные: d6-d5 (одиночный)
				require.NoError(t, game.MakeMove(model.NewMove(5, 3, 4, 3)))

				// 3. Белые пробуют e5xd6 - не en passant
				return game, model.NewMove(4, 4, 5, 3)
			},
			wantErr: true,
		},
		{
			name: "враг не в соседнем столбце - не en passant",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Pawn)) // e5
				board.SetPiece(6, 2, model.NewPiece(model.Black, model.Pawn)) // c7
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// 1. Белые: Ke1-d1
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3)))
				// 2. Чёрные: c7-c5 (двойной, но НЕ на соседнем столбце)
				require.NoError(t, game.MakeMove(model.NewMove(6, 2, 4, 2)))

				// 3. Белые пробуют e5xd6 - не en passant (враг на c5, не d5)
				return game, model.NewMove(4, 4, 5, 3)
			},
			wantErr: true,
		},
		{
			name: "не сразу - en passant просрочен",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(4, 4, model.NewPiece(model.White, model.Pawn)) // e5
				board.SetPiece(6, 3, model.NewPiece(model.Black, model.Pawn)) // d7
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// 1. Белые: Ke1-d1 (нейтральный)
				require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 3)))
				// 2. Чёрные: d7-d5 (двойной)
				require.NoError(t, game.MakeMove(model.NewMove(6, 3, 4, 3)))
				// 3. Белые: Kd1-e1 (нейтральный)
				require.NoError(t, game.MakeMove(model.NewMove(0, 3, 0, 4)))
				// 4. Чёрные: Ke8-d8 (нейтральный)
				require.NoError(t, game.MakeMove(model.NewMove(7, 4, 7, 3)))

				// 5. Белые: e5xd6 - ПРОСРОЧЕН
				return game, model.NewMove(4, 4, 5, 3)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, move := tt.setup()

			movesBefore := len(game.Moves())
			colorBefore := game.CurrentColor()

			err := game.MakeMove(move)

			if tt.wantErr {
				assert.Error(t, err, "ожидаем ошибку")
				assert.Len(t, game.Moves(), movesBefore,
					"история ходов не должна пополняться при ошибке")
				assert.Equal(t, colorBefore, game.CurrentColor(),
					"цвет не должен переключаться при ошибке")
				return
			}

			require.NoError(t, err)
			assert.Len(t, game.Moves(), movesBefore+1)

			moves := game.Moves()
			last := moves[len(moves)-1]

			assert.False(t, last.Check)
			assert.False(t, last.Mate)

			if tt.check != nil {
				tt.check(t, game.Board(), &last)
			}
		})
	}
}

func TestGame_MakeMove_CheckAndMate(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (*model.Game, *model.Move)
		wantCheck bool
		wantMate  bool
	}{
		{
			name: "обычный ход - ни шаха, ни мата",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))  // e1
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))  // e8
				board.SetPiece(0, 3, model.NewPiece(model.White, model.Queen)) // d1
				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game, model.NewMove(0, 3, 1, 3) // Qd1-d2
			},
			wantCheck: false,
			wantMate:  false,
		},
		{
			name: "шах ферзём по вертикали",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))  // e1
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))  // e8
				board.SetPiece(3, 3, model.NewPiece(model.White, model.Queen)) // d4

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// Qd4-e4+ - шах по вертикали e
				return game, model.NewMove(3, 3, 3, 4)
			},
			wantCheck: true,
			wantMate:  false,
		},
		{
			name: "мат ладьёй на последней горизонтали",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
				board.SetPiece(0, 0, model.NewPiece(model.White, model.Rook)) // a1

				// чёрный король зажат своими пешками
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
				board.SetPiece(6, 3, model.NewPiece(model.Black, model.Pawn)) // d7
				board.SetPiece(6, 4, model.NewPiece(model.Black, model.Pawn)) // e7
				board.SetPiece(6, 5, model.NewPiece(model.Black, model.Pawn)) // f7

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// Ra1-a8# - мат
				return game, model.NewMove(0, 0, 7, 0)
			},
			wantCheck: true,
			wantMate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, move := tt.setup()

			err := game.MakeMove(move)
			require.NoError(t, err)

			moves := game.Moves()
			require.Len(t, moves, 1)
			last := moves[0]

			assert.Equal(t, tt.wantCheck, last.Check, "флаг Check")
			assert.Equal(t, tt.wantMate, last.Mate, "флаг Mate")

			if tt.wantMate {
				assert.True(t, last.Check, "при мате Check тоже должен быть true")
			}
		})
	}
}

func TestGame_MakeMove_Promotion(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*model.Game, *model.Move)
		wantErr       bool
		wantPieceType model.PieceType
		wantColor     model.Color
	}{
		{
			name: "белая e7-e8=Q",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(6, 4, 7, 4)
				move.Promotion = model.Queen
				return game, move
			},
			wantPieceType: model.Queen,
			wantColor:     model.White,
		},
		{
			name: "белая e7-e8=R",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(6, 4, 7, 4)
				move.Promotion = model.Rook
				return game, move
			},
			wantPieceType: model.Rook,
			wantColor:     model.White,
		},
		{
			name: "белая e7-e8=B",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(6, 4, 7, 4)
				move.Promotion = model.Bishop
				return game, move
			},
			wantPieceType: model.Bishop,
			wantColor:     model.White,
		},
		{
			name: "белая e7-e8=N",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(6, 4, 7, 4)
				move.Promotion = model.Knight
				return game, move
			},
			wantPieceType: model.Knight,
			wantColor:     model.White,
		},
		{
			name: "чёрная e2-e1=Q",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(1, 4, model.NewPiece(model.Black, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				// первый ход белых — переключаем на чёрных
				require.NoError(t, game.MakeMove(model.NewMove(0, 0, 0, 1))) // Ka1-b1

				move := model.NewMove(1, 4, 0, 4)
				move.Promotion = model.Queen
				return game, move
			},
			wantPieceType: model.Queen,
			wantColor:     model.Black,
		},
		{
			name: "обычный ход пешки - Promotion игнорируется",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(1, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(1, 4, 3, 4)
				move.Promotion = model.Queen
				return game, move
			},
			wantPieceType: model.Pawn,
			wantColor:     model.White,
		},
		{
			name: "пешка на последний ряд без Promotion",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				return game, model.NewMove(6, 4, 7, 4)
			},
			wantErr: true,
		},
		{
			name: "пешка на последний ряд с Promotion=Pawn",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(6, 4, 7, 4)
				move.Promotion = model.Pawn
				return game, move
			},
			wantErr: true,
		},
		{
			name: "пешка на последний ряд с Promotion=King",
			setup: func() (*model.Game, *model.Move) {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 0, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 7, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)

				move := model.NewMove(6, 4, 7, 4)
				move.Promotion = model.King
				return game, move
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, move := tt.setup()

			movesBefore := len(game.Moves())

			err := game.MakeMove(move)

			if tt.wantErr {
				assert.Error(t, err, "ожидаем ошибку")
				assert.Len(t, game.Moves(), movesBefore,
					"история не должна пополняться при ошибке")
				return
			}

			require.NoError(t, err)

			promoted := game.Board().PieceAt(move.ToRow, move.ToCol)
			require.NotNil(t, promoted, "фигура должна быть на целевой клетке")
			assert.Equal(t, tt.wantPieceType, promoted.Type(),
				"тип фигуры на целевой клетке")
			assert.Equal(t, tt.wantColor, promoted.Color(),
				"цвет фигуры на целевой клетке")

			assert.Nil(t, game.Board().PieceAt(move.FromRow, move.FromCol))

			moves := game.Moves()
			last := moves[len(moves)-1]
			require.NotNil(t, last.MovedPiece)
			assert.Equal(t, model.Pawn, last.MovedPiece.Type(),
				"в истории — пешка (до превращения)")
		})
	}
}

func TestGame_CheckmateStalemate(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *model.Game
		color         model.Color
		wantCheckmate bool
		wantStalemate bool
	}{
		{
			name: "мат",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 0, model.NewPiece(model.White, model.Rook))

				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))
				board.SetPiece(6, 3, model.NewPiece(model.Black, model.Pawn))
				board.SetPiece(6, 4, model.NewPiece(model.Black, model.Pawn))
				board.SetPiece(6, 5, model.NewPiece(model.Black, model.Pawn))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			color:         model.Black,
			wantCheckmate: true,
			wantStalemate: false,
		},
		{
			name: "пат - король h8, ферзь f7, белый король g6",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(2, 6, model.NewPiece(model.White, model.King))  // g6
				board.SetPiece(1, 5, model.NewPiece(model.White, model.Queen)) // f7

				board.SetPiece(0, 7, model.NewPiece(model.Black, model.King)) // h8

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			color:         model.Black,
			wantCheckmate: false,
			wantStalemate: true,
		},
		{
			name: "обычная позиция - ни мата, ни пата",
			setup: func() *model.Game {
				board := emptyBoard(8, 8)
				board.SetPiece(0, 4, model.NewPiece(model.White, model.King))
				board.SetPiece(7, 4, model.NewPiece(model.Black, model.King))

				game := model.NewGame("A", "B", 8, 8)
				game.SetBoard(board)
				return game
			},
			color:         model.Black,
			wantCheckmate: false,
			wantStalemate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()

			assert.Equal(t, tt.wantCheckmate, game.Checkmate(tt.color),
				"Checkmate(%v)", tt.color)
			assert.Equal(t, tt.wantStalemate, game.Stalemate(tt.color),
				"Stalemate(%v)", tt.color)
		})
	}
}

func TestGame_MovesIsCopy(t *testing.T) {
	game := model.NewGame("A", "B", 8, 8)

	// делаем пару ходов
	require.NoError(t, game.MakeMove(model.NewMove(1, 4, 3, 4))) // e2-e4
	require.NoError(t, game.MakeMove(model.NewMove(6, 4, 4, 4))) // e7-e5

	moves1 := game.Moves()
	require.Len(t, moves1, 2)

	moves1[0] = model.Move{FromRow: 99, FromCol: 99}

	moves2 := game.Moves()
	require.Len(t, moves2, 2)
	assert.Equal(t, 1, moves2[0].FromRow, "внутренний ход не должен быть испорчен")
	assert.Equal(t, 4, moves2[0].FromCol)

	moves2 = append(moves2, model.Move{FromRow: 42, FromCol: 42})

	moves3 := game.Moves()
	assert.Len(t, moves3, 2, "append к копии не влияет на оригинал")
}

func TestGame_MarshalJSON(t *testing.T) {
	game := model.NewGame("Иванов", "Петров", 8, 8)
	game.SetID(42)

	data, err := json.Marshal(game)
	require.NoError(t, err)

	var parsed struct {
		ID      int `json:"id"`
		Player1 struct {
			ID    int    `json:"id"`
			Name  string `json:"имя"`
			Color string `json:"цвет"`
		} `json:"игрок1"`
		Player2 struct {
			ID    int    `json:"id"`
			Name  string `json:"имя"`
			Color string `json:"цвет"`
		} `json:"игрок2"`
		Board   json.RawMessage `json:"доска"`
		Current string          `json:"текущий"`
		Moves   []any           `json:"ходы"`
	}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, 42, parsed.ID)
	assert.Equal(t, "Иванов", parsed.Player1.Name)
	assert.Equal(t, "Белые", parsed.Player1.Color)
	assert.Equal(t, "Петров", parsed.Player2.Name)
	assert.Equal(t, "Чёрные", parsed.Player2.Color)
	assert.NotEmpty(t, parsed.Board)
	assert.Equal(t, "Белые", parsed.Current)
	assert.Empty(t, parsed.Moves)
}

func TestGame_UnmarshalJSON(t *testing.T) {
	t.Run("полный JSON", func(t *testing.T) {
		input := `{
			"id": 42,
			"игрок1": {"id": 1, "имя": "Иванов", "цвет": "Белые"},
			"игрок2": {"id": 2, "имя": "Петров", "цвет": "Чёрные"},
			"доска": {
				"строки": 8,
				"столбцы": 8,
				"клетки": [
					[{"цвет":"Белые","тип":"Ладья"},null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null],
					[null,null,null,null,null,null,null,null]
				]
			},
			"текущий": "Чёрные",
			"ходы": []
		}`

		var game model.Game
		err := json.Unmarshal([]byte(input), &game)
		require.NoError(t, err)

		assert.Equal(t, 42, game.ID())
		assert.Equal(t, "Иванов", game.Player1().Name())
		assert.Equal(t, 1, game.Player1().ID())
		assert.Equal(t, model.White, game.Player1Color())
		assert.Equal(t, "Петров", game.Player2().Name())
		assert.Equal(t, 2, game.Player2().ID())
		assert.Equal(t, model.Black, game.Player2Color())
		assert.Equal(t, model.Black, game.CurrentColor())

		require.NotNil(t, game.Board())
		assert.Equal(t, 8, game.Board().Rows())

		rook := game.Board().PieceAt(0, 0)
		require.NotNil(t, rook)
		assert.Equal(t, model.Rook, rook.Type())
		assert.Equal(t, model.White, rook.Color())

		assert.Empty(t, game.Moves())
	})

	t.Run("с ходами", func(t *testing.T) {
		input := `{
			"id": 1,
			"игрок1": {"id": 1, "имя": "A", "цвет": "Белые"},
			"игрок2": {"id": 2, "имя": "B", "цвет": "Чёрные"},
			"доска": {"строки": 8, "столбцы": 8, "клетки": [
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null]
			]},
			"текущий": "Белые",
			"ходы": [
				{"играID": 1, "id": 10, "отСтрока": 1, "отСтолбец": 4, "вСтрока": 3, "вСтолбец": 4},
				{"играID": 1, "id": 11, "отСтрока": 6, "отСтолбец": 4, "вСтрока": 4, "вСтолбец": 4}
			]
		}`

		var game model.Game
		err := json.Unmarshal([]byte(input), &game)
		require.NoError(t, err)

		moves := game.Moves()
		require.Len(t, moves, 2)

		assert.Equal(t, 10, moves[0].ID())
		assert.Equal(t, 1, moves[0].FromRow)
		assert.Equal(t, 4, moves[0].FromCol)
		assert.Equal(t, 3, moves[0].ToRow)
		assert.Equal(t, 4, moves[0].ToCol)

		assert.Equal(t, 11, moves[1].ID())
	})

	t.Run("с флагами рокировки", func(t *testing.T) {
		input := `{
			"id": 1,
			"игрок1": {"id": 1, "имя": "A", "цвет": "Белые"},
			"игрок2": {"id": 2, "имя": "B", "цвет": "Чёрные"},
			"доска": {"строки": 8, "столбцы": 8, "клетки": [
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null],
				[null,null,null,null,null,null,null,null]
			]},
			"текущий": "Белые",
			"ходы": [],
			"белыйКорольДвигался": true,
			"белаяЛадьяAДвигалась": true,
			"белаяЛадьяHДвигалась": true
		}`

		var game model.Game
		err := json.Unmarshal([]byte(input), &game)
		require.NoError(t, err)

		board := game.Board()
		board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
		board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook)) // h1

		movesBefore := len(game.Moves())
		err = game.MakeMove(model.NewMove(0, 4, 0, 6)) // рокировка
		assert.Error(t, err, "рокировка должна падать — флаги установлены")
		assert.Len(t, game.Moves(), movesBefore)
	})

	t.Run("сломанный JSON", func(t *testing.T) {
		var game model.Game
		err := json.Unmarshal([]byte(`{`), &game)
		assert.Error(t, err)
	})

	t.Run("неверный тип поля", func(t *testing.T) {
		var game model.Game
		err := json.Unmarshal([]byte(`{"id": "abc"}`), &game)
		assert.Error(t, err)
	})
}

func TestGame_RoundTrip(t *testing.T) {
	game := model.NewGame("Иванов", "Петров", 8, 8)
	game.SetID(42)

	// делаем несколько ходов
	require.NoError(t, game.MakeMove(model.NewMove(1, 4, 3, 4))) // e2-e4
	require.NoError(t, game.MakeMove(model.NewMove(6, 4, 4, 4))) // e7-e5
	require.NoError(t, game.MakeMove(model.NewMove(0, 6, 2, 5))) // Ng1-f3
	require.NoError(t, game.MakeMove(model.NewMove(7, 1, 5, 2))) // Nb8-c6

	// туда
	data, err := json.Marshal(game)
	require.NoError(t, err)

	// обратно
	var restored model.Game
	err = json.Unmarshal(data, &restored)
	require.NoError(t, err)

	// ID
	assert.Equal(t, game.ID(), restored.ID())

	// игроки
	assert.Equal(t, game.Player1().Name(), restored.Player1().Name())
	assert.Equal(t, game.Player2().Name(), restored.Player2().Name())

	// цвета
	assert.Equal(t, game.Player1Color(), restored.Player1Color())
	assert.Equal(t, game.Player2Color(), restored.Player2Color())

	// текущий ход
	assert.Equal(t, game.CurrentColor(), restored.CurrentColor())

	// доска — размеры
	require.NotNil(t, restored.Board())
	assert.Equal(t, game.Board().Rows(), restored.Board().Rows())
	assert.Equal(t, game.Board().Cols(), restored.Board().Cols())

	// доска — каждая клетка
	for row := range game.Board().Rows() {
		for col := range game.Board().Cols() {
			orig := game.Board().PieceAt(row, col)
			rest := restored.Board().PieceAt(row, col)
			if orig == nil {
				assert.Nil(t, rest, "(%d,%d) должно быть пусто", row, col)
				continue
			}
			require.NotNil(t, rest, "(%d,%d) не должно быть пусто", row, col)
			assert.Equal(t, orig.Type(), rest.Type(), "тип (%d,%d)", row, col)
			assert.Equal(t, orig.Color(), rest.Color(), "цвет (%d,%d)", row, col)
		}
	}

	// ходы
	origMoves := game.Moves()
	restMoves := restored.Moves()
	require.Len(t, restMoves, len(origMoves))
	for i := range origMoves {
		assert.Equal(t, origMoves[i].FromRow, restMoves[i].FromRow)
		assert.Equal(t, origMoves[i].FromCol, restMoves[i].FromCol)
		assert.Equal(t, origMoves[i].ToRow, restMoves[i].ToRow)
		assert.Equal(t, origMoves[i].ToCol, restMoves[i].ToCol)
	}
}

func TestGame_RoundTrip_Flags(t *testing.T) {
	// Сценарий: король ходит туда-обратно, флаг "король ходил" должен сохраниться

	board := emptyBoard(8, 8)
	board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
	board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook)) // h1
	board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
	board.SetPiece(7, 0, model.NewPiece(model.Black, model.Rook)) // a8

	game := model.NewGame("A", "B", 8, 8)
	game.SetBoard(board)

	// белые: Ke1-f1; чёрные: Ke8-d8; белые: Kf1-e1; чёрные: Kd8-e8
	require.NoError(t, game.MakeMove(model.NewMove(0, 4, 0, 5)))
	require.NoError(t, game.MakeMove(model.NewMove(7, 4, 7, 3)))
	require.NoError(t, game.MakeMove(model.NewMove(0, 5, 0, 4)))
	require.NoError(t, game.MakeMove(model.NewMove(7, 3, 7, 4)))

	// round-trip
	data, err := json.Marshal(game)
	require.NoError(t, err)

	var restored model.Game
	err = json.Unmarshal(data, &restored)
	require.NoError(t, err)

	// пытаемся рокироваться
	movesBefore := len(restored.Moves())
	err = restored.MakeMove(model.NewMove(0, 4, 0, 6))
	assert.Error(t, err, "рокировка должна быть отклонена — король уже ходил")
	assert.Len(t, restored.Moves(), movesBefore, "история не должна пополняться")
}

func TestGame_MakeMove_Errors_More(t *testing.T) {
	tests := []struct {
		name    string
		piece   model.PieceType
		fromRow int
		fromCol int
		toRow   int
		toCol   int
	}{
		{
			name:    "ферзь по букве Г",
			piece:   model.Queen,
			fromRow: 4, fromCol: 4, toRow: 2, toCol: 3,
		},
		{
			name:    "король на 2 клетки",
			piece:   model.King,
			fromRow: 4, fromCol: 4, toRow: 6, toCol: 4,
		},
		{
			name:    "слон по прямой",
			piece:   model.Bishop,
			fromRow: 4, fromCol: 4, toRow: 4, toCol: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := emptyBoard(8, 8)
			board.SetPiece(tt.fromRow, tt.fromCol,
				model.NewPiece(model.White, tt.piece))

			game := model.NewGame("A", "B", 8, 8)
			game.SetBoard(board)

			movesBefore := len(game.Moves())
			colorBefore := game.CurrentColor()

			move := model.NewMove(tt.fromRow, tt.fromCol, tt.toRow, tt.toCol)
			err := game.MakeMove(move)
			assert.Error(t, err, "ожидаем ошибку")

			assert.Len(t, game.Moves(), movesBefore,
				"история не должна пополняться")
			assert.Equal(t, colorBefore, game.CurrentColor(),
				"цвет не должен переключаться")
		})
	}
}

func TestGame_MakeMove_Castling_Check(t *testing.T) {
	// Рокировка, после которой белая ладья на f1 ставит шах чёрному королю на f8
	board := emptyBoard(8, 8)
	board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
	board.SetPiece(0, 7, model.NewPiece(model.White, model.Rook)) // h1
	board.SetPiece(7, 5, model.NewPiece(model.Black, model.King)) // f8

	game := model.NewGame("A", "B", 8, 8)
	game.SetBoard(board)

	// короткая рокировка: Ke1-g1, ладья h1-f1
	move := model.NewMove(0, 4, 0, 6)
	err := game.MakeMove(move)
	require.NoError(t, err)

	moves := game.Moves()
	require.Len(t, moves, 1)
	last := moves[0]

	// после рокировки ладья на f1 бьёт чёрного короля на f8 по вертикали f
	assert.True(t, last.Check, "должен быть шах")
	assert.False(t, last.Mate, "не мат — у чёрного короля есть куда уйти")

	// король и ладья на новых местах
	king := game.Board().PieceAt(0, 6)
	require.NotNil(t, king)
	assert.Equal(t, model.King, king.Type())

	rook := game.Board().PieceAt(0, 5)
	require.NotNil(t, rook)
	assert.Equal(t, model.Rook, rook.Type())
}

func TestGame_SimpleMethods(t *testing.T) {
	game := model.NewGame("A", "B", 8, 8)

	// ObjectType
	assert.Equal(t, "игра", game.ObjectType())

	// InCheck - на старте никто не под шахом
	assert.False(t, game.InCheck(model.White))
	assert.False(t, game.InCheck(model.Black))

	// SetLastMoveID - сначала без ходов, потом с ходом
	game.SetLastMoveID(999) // ходов нет - метод молча выходит

	require.NoError(t, game.MakeMove(model.NewMove(1, 4, 3, 4))) // e2-e4
	game.SetLastMoveID(42)

	moves := game.Moves()
	require.Len(t, moves, 1)
	assert.Equal(t, 42, moves[0].ID(), "ID последнего хода должен обновиться")
}
