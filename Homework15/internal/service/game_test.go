package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mod.go/internal/model"
	"mod.go/internal/service"
)

func TestGameService_Create(t *testing.T) {
	tests := []struct {
		name    string
		p1Name  string
		p2Name  string
		rows    int
		cols    int
		wantErr bool
	}{
		{name: "стандарт 8x8", p1Name: "Иванов", p2Name: "Петров", rows: 8, cols: 8},
		{name: "4x4 (минимум)", p1Name: "A", p2Name: "B", rows: 4, cols: 4},
		{name: "12x12", p1Name: "A", p2Name: "B", rows: 12, cols: 12},

		{name: "0x0 - 8x8 по умолчанию", p1Name: "A", p2Name: "B", rows: 0, cols: 0},
		{name: "4x0 - 4x8", p1Name: "A", p2Name: "B", rows: 4, cols: 0},
		{name: "0x8 - 8x8", p1Name: "A", p2Name: "B", rows: 0, cols: 8},

		{name: "3x3 - меньше минимума", p1Name: "A", p2Name: "B", rows: 3, cols: 3, wantErr: true},
		{name: "2x8 - меньше минимума", p1Name: "A", p2Name: "B", rows: 2, cols: 8, wantErr: true},
		{name: "пустое имя игрока 1", p1Name: "", p2Name: "B", rows: 8, cols: 8, wantErr: true},
		{name: "пустое имя игрока 2", p1Name: "A", p2Name: "", rows: 8, cols: 8, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewGameService(storage)

			p1 := *model.NewPlayer(tt.p1Name)
			p2 := *model.NewPlayer(tt.p2Name)

			game, err := svc.Create(p1, p2, tt.rows, tt.cols)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.NotZero(t, game.ID())

			assert.Equal(t, tt.p1Name, game.Player1().Name())
			assert.Equal(t, tt.p2Name, game.Player2().Name())
			assert.Equal(t, model.White, game.Player1Color())
			assert.Equal(t, model.Black, game.Player2Color())
			assert.Equal(t, model.White, game.CurrentColor())

			expectedRows := tt.rows
			expectedCols := tt.cols
			if expectedRows <= 0 {
				expectedRows = 8
			}
			if expectedCols <= 0 {
				expectedCols = 8
			}
			require.NotNil(t, game.Board())
			assert.Equal(t, expectedRows, game.Board().Rows())
			assert.Equal(t, expectedCols, game.Board().Cols())

			got, found := storage.GetGameByID(game.ID())
			require.True(t, found)
			assert.Equal(t, game.ID(), got.ID())
		})
	}
}

func TestGameService_Get(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewGameService(storage)

	created, err := svc.Create(
		*model.NewPlayer("Иванов"),
		*model.NewPlayer("Петров"),
		8, 8,
	)
	require.NoError(t, err)

	got, err := svc.Get(created.ID())
	require.NoError(t, err)
	assert.Equal(t, created.ID(), got.ID())
	assert.Equal(t, "Иванов", got.Player1().Name())

	_, err = svc.Get(999)
	assert.Error(t, err)
}

func TestGameService_List(t *testing.T) {
	t.Run("пусто", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		assert.Empty(t, svc.List())
	})

	t.Run("несколько игр", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		_, err := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
		require.NoError(t, err)
		_, err = svc.Create(*model.NewPlayer("C"), *model.NewPlayer("D"), 8, 8)
		require.NoError(t, err)

		list := svc.List()
		require.Len(t, list, 2)
		assert.Equal(t, "A", list[0].Player1().Name())
		assert.Equal(t, "C", list[1].Player1().Name())
	})
}

func TestGameService_Update(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(svc *service.GameService) int
		p1Name     string
		p2Name     string
		wantP1Name string
		wantP2Name string
		wantErr    bool
	}{
		{
			name: "обновить обоих игроков",
			setup: func(svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			p1Name:     "Сидоров",
			p2Name:     "Кузнецов",
			wantP1Name: "Сидоров",
			wantP2Name: "Кузнецов",
		},
		{
			name: "обновить только игрока 1",
			setup: func(svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			p1Name:     "Сидоров",
			p2Name:     "",
			wantP1Name: "Сидоров",
			wantP2Name: "B", //  не изменился
		},
		{
			name: "обновить только игрока 2",
			setup: func(svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			p1Name:     "",
			p2Name:     "Кузнецов",
			wantP1Name: "A", //  не изменился
			wantP2Name: "Кузнецов",
		},
		{
			name: "оба пустые - ничего не меняется",
			setup: func(svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			p1Name:     "",
			p2Name:     "",
			wantP1Name: "A",
			wantP2Name: "B",
		},
		{
			name: "несуществующий ID",
			setup: func(svc *service.GameService) int {
				return 999
			},
			p1Name:  "Сидоров",
			p2Name:  "Кузнецов",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewGameService(storage)
			id := tt.setup(svc)

			updated, err := svc.Update(id, *model.NewPlayer(tt.p1Name), *model.NewPlayer(tt.p2Name))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, id, updated.ID())
			assert.Equal(t, tt.wantP1Name, updated.Player1().Name())
			assert.Equal(t, tt.wantP2Name, updated.Player2().Name())
		})
	}
}

func TestGameService_Delete(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewGameService(storage)

	created, err := svc.Create(
		*model.NewPlayer("A"),
		*model.NewPlayer("B"),
		8, 8,
	)
	require.NoError(t, err)

	err = svc.Delete(created.ID())
	require.NoError(t, err)

	_, err = svc.Get(created.ID())
	assert.Error(t, err)

	err = svc.Delete(created.ID())
	assert.Error(t, err)
}

func TestGameService_MakeMove(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, svc *service.GameService) int
		move    *model.Move
		wantErr bool
	}{
		{
			name: "простой ход пешкой",
			setup: func(t *testing.T, svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			move: model.NewMove(1, 4, 3, 4), // e2-e4
		},
		{
			name: "несуществующая игра",
			setup: func(t *testing.T, svc *service.GameService) int {
				return 999
			},
			move:    model.NewMove(1, 4, 3, 4),
			wantErr: true,
		},
		{
			name: "нелегальный ход",
			setup: func(t *testing.T, svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			move:    model.NewMove(1, 4, 1, 5), // пешка вбок
			wantErr: true,
		},
		{
			name: "фигура чужого цвета",
			setup: func(t *testing.T, svc *service.GameService) int {
				g, _ := svc.Create(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
				return g.ID()
			},
			move:    model.NewMove(6, 4, 4, 4), // чёрная пешка при ходе белых
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewGameService(storage)
			id := tt.setup(t, svc)

			updated, err := svc.MakeMove(id, tt.move)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// ход записан в игру
			assert.Len(t, updated.Moves(), 1)

			// цвет переключился
			assert.Equal(t, model.Black, updated.CurrentColor())

			// ход сохранён в storage отдельно
			storedMoves := storage.GetAllMoves()
			require.Len(t, storedMoves, 1)
			assert.Equal(t, tt.move.FromRow, storedMoves[0].FromRow)
			assert.Equal(t, id, storedMoves[0].GameID())

			// игра перезаписана — в storage тоже есть ход
			storedGame, found := storage.GetGameByID(id)
			require.True(t, found)
			assert.Len(t, storedGame.Moves(), 1)
		})
	}
}

func TestGameService_MakeMove_Checkmate(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewGameService(storage)

	// кастомная доска: мат ладьёй на 8-й горизонтали
	board := emptyServiceBoard(8, 8)
	board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
	board.SetPiece(0, 0, model.NewPiece(model.White, model.Rook)) // a1
	board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
	board.SetPiece(6, 3, model.NewPiece(model.Black, model.Pawn)) // d7
	board.SetPiece(6, 4, model.NewPiece(model.Black, model.Pawn)) // e7
	board.SetPiece(6, 5, model.NewPiece(model.Black, model.Pawn)) // f7

	game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
	created, err := storage.CreateGame(*game)
	require.NoError(t, err)

	// Ra1-a8#
	move := model.NewMove(0, 0, 7, 0)
	updated, err := svc.MakeMove(created.ID(), move)
	require.NoError(t, err)

	moves := updated.Moves()
	require.Len(t, moves, 1)
	assert.True(t, moves[0].Mate, "должен быть мат")
	assert.True(t, moves[0].Check, "при мате Check тоже true")
}

func TestGameService_AutoMove(t *testing.T) {
	t.Run("один ход", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		created, err := svc.Create(
			*model.NewPlayer("A"),
			*model.NewPlayer("B"),
			8, 8,
		)
		require.NoError(t, err)

		assert.Empty(t, created.Moves())

		updated, err := svc.AutoMove(created.ID())
		require.NoError(t, err)

		assert.Len(t, updated.Moves(), 1)
		assert.Equal(t, model.Black, updated.CurrentColor())
		assert.Len(t, storage.GetAllMoves(), 1)
	})

	t.Run("несколько ходов подряд", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		created, err := svc.Create(
			*model.NewPlayer("A"),
			*model.NewPlayer("B"),
			8, 8,
		)
		require.NoError(t, err)

		for i := 1; i <= 5; i++ {
			updated, err := svc.AutoMove(created.ID())
			require.NoError(t, err, "ход #%d", i)
			assert.Len(t, updated.Moves(), i,
				"после %d автоходов должно быть %d ходов в игре", i, i)
		}

		assert.Len(t, storage.GetAllMoves(), 5,
			"все 5 ходов сохранены в storage")
	})

	t.Run("несуществующая игра", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		_, err := svc.AutoMove(999)
		assert.Error(t, err)
	})

	t.Run("мат - автоход не делается", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		board := emptyServiceBoard(8, 8)
		board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
		board.SetPiece(7, 0, model.NewPiece(model.White, model.Rook)) // a8 - уже бьёт по 8-й
		board.SetPiece(7, 4, model.NewPiece(model.Black, model.King)) // e8
		board.SetPiece(6, 3, model.NewPiece(model.Black, model.Pawn)) // d7
		board.SetPiece(6, 4, model.NewPiece(model.Black, model.Pawn)) // e7
		board.SetPiece(6, 5, model.NewPiece(model.Black, model.Pawn)) // f7

		game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
		created, err := storage.CreateGame(*game)
		require.NoError(t, err)

		// белый ход — переключаем ход на чёрных
		_, err = svc.MakeMove(created.ID(), model.NewMove(0, 4, 0, 3)) // Ke1-d1
		require.NoError(t, err)

		// чёрные в мате - AutoMove не должен делать ход
		updated, err := svc.AutoMove(created.ID())
		require.NoError(t, err, "мат - возврат без хода")

		assert.Len(t, updated.Moves(), 1,
			"автоход не должен добавить ход при мате")
		assert.Equal(t, model.Black, updated.CurrentColor(),
			"ход всё ещё чёрных")
		assert.Len(t, storage.GetAllMoves(), 1,
			"в storage тоже только 1 ход")
	})

	t.Run("пат - автоход не делается", func(t *testing.T) {
		storage := newMockStorage()
		svc := service.NewGameService(storage)

		// классический пат: чёрный король h8, белый ферзь f7, белый король g6
		board := emptyServiceBoard(8, 8)
		board.SetPiece(0, 7, model.NewPiece(model.Black, model.King))  // h8
		board.SetPiece(1, 5, model.NewPiece(model.White, model.Queen)) // f7
		board.SetPiece(2, 6, model.NewPiece(model.White, model.King))  // g6

		game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
		created, err := storage.CreateGame(*game)
		require.NoError(t, err)

		// белый ход: Kg6-f6 (пат сохраняется - король f6 тоже атакует g7)
		_, err = svc.MakeMove(created.ID(), model.NewMove(2, 6, 2, 5))
		require.NoError(t, err)

		// чёрные в пате
		updated, err := svc.AutoMove(created.ID())
		require.NoError(t, err, "пат - возврат без хода")

		assert.Len(t, updated.Moves(), 1)
		assert.Equal(t, model.Black, updated.CurrentColor())
	})
}

func TestGameService_MakeMove_Promotion(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, storage *mockStorage) int
		move          *model.Move
		wantErr       bool
		wantPieceType model.PieceType
	}{
		{
			name: "белая пешка - ферзь",
			setup: func(t *testing.T, storage *mockStorage) int {
				board := emptyServiceBoard(8, 8)
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn)) // e7
				game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
				created, err := storage.CreateGame(*game)
				require.NoError(t, err)
				return created.ID()
			},
			move: func() *model.Move {
				m := model.NewMove(6, 4, 7, 4)
				m.Promotion = model.Queen
				return m
			}(),
			wantPieceType: model.Queen,
		},
		{
			name: "белая пешка - конь",
			setup: func(t *testing.T, storage *mockStorage) int {
				board := emptyServiceBoard(8, 8)
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))
				game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
				created, err := storage.CreateGame(*game)
				require.NoError(t, err)
				return created.ID()
			},
			move: func() *model.Move {
				m := model.NewMove(6, 4, 7, 4)
				m.Promotion = model.Knight
				return m
			}(),
			wantPieceType: model.Knight,
		},
		{
			name: "обычный ход с Promotion - сервис сбрасывает на Pawn",
			setup: func(t *testing.T, storage *mockStorage) int {
				board := emptyServiceBoard(8, 8)
				board.SetPiece(1, 4, model.NewPiece(model.White, model.Pawn)) // e2
				game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
				created, err := storage.CreateGame(*game)
				require.NoError(t, err)
				return created.ID()
			},
			move: func() *model.Move {
				m := model.NewMove(1, 4, 3, 4) // e2-e4
				m.Promotion = model.Queen      // но ход не на последний ряд
				return m
			}(),
			wantPieceType: model.Pawn, // сервис сбросил, пешка осталась пешкой
		},
		{
			name: "пешка на последний ряд без Promotion - ошибка",
			setup: func(t *testing.T, storage *mockStorage) int {
				board := emptyServiceBoard(8, 8)
				board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn))
				game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
				created, err := storage.CreateGame(*game)
				require.NoError(t, err)
				return created.ID()
			},
			move:    model.NewMove(6, 4, 7, 4), // Promotion не задан
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			svc := service.NewGameService(storage)
			id := tt.setup(t, storage)

			updated, err := svc.MakeMove(id, tt.move)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, storage.GetAllMoves(),
					"при ошибке ход не должен быть сохранён")
				return
			}
			require.NoError(t, err)

			// фигура на целевой клетке
			piece := updated.Board().PieceAt(tt.move.ToRow, tt.move.ToCol)
			require.NotNil(t, piece)
			assert.Equal(t, tt.wantPieceType, piece.Type())
			assert.Equal(t, model.White, piece.Color())
		})
	}
}

func TestGameService_MakeMove_LeavesKingInCheck(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewGameService(storage)

	board := emptyServiceBoard(8, 8)
	board.SetPiece(0, 4, model.NewPiece(model.White, model.King)) // e1
	board.SetPiece(1, 4, model.NewPiece(model.White, model.Rook)) // e2 (заслон)
	board.SetPiece(7, 4, model.NewPiece(model.Black, model.Rook)) // e8

	game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
	created, err := storage.CreateGame(*game)
	require.NoError(t, err)

	// ход ладьёй e2-d2 открывает короля для чёрной ладьи e8
	move := model.NewMove(1, 4, 1, 3)
	_, err = svc.MakeMove(created.ID(), move)
	assert.Error(t, err, "ход должен быть отклонён - оставляет короля под шахом")

	// ничего не сохранено
	assert.Empty(t, storage.GetAllMoves(), "ход не должен быть сохранён")
	stored, _ := storage.GetGameByID(created.ID())
	assert.Empty(t, stored.Moves(), "в игре не должно быть ходов")
}

func TestGameService_AutoMove_Promotion(t *testing.T) {
	storage := newMockStorage()
	svc := service.NewGameService(storage)

	// на доске только белая пешка на e7 - единственный возможный ход e7-e8
	board := emptyServiceBoard(8, 8)
	board.SetPiece(6, 4, model.NewPiece(model.White, model.Pawn)) // e7

	game := model.NewGameWithBoard(*model.NewPlayer("A"), *model.NewPlayer("B"), board)
	created, err := storage.CreateGame(*game)
	require.NoError(t, err)

	updated, err := svc.AutoMove(created.ID())
	require.NoError(t, err)

	// ход сделан
	assert.Len(t, updated.Moves(), 1)

	// на e8 - ферзь (сервис автоматом ставит Queen)
	piece := updated.Board().PieceAt(7, 4)
	require.NotNil(t, piece)
	assert.Equal(t, model.Queen, piece.Type(), "пешка должна была превратиться в ферзя")
	assert.Equal(t, model.White, piece.Color())
}
