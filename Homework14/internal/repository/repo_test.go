package repository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mod.go/internal/model"
	"mod.go/internal/repository"
)

func newTestStorage(t *testing.T) *repository.Storage {
	t.Helper()
	s := repository.NewStorage(t.TempDir())
	require.NoError(t, s.LoadFromFiles())
	return s
}

func TestStorage_LoadFromFiles_Empty(t *testing.T) {
	s := newTestStorage(t)

	assert.Empty(t, s.GetAllPlayers())
	assert.Empty(t, s.GetAllGames())
	assert.Empty(t, s.GetAllMoves())
}

func TestStorage_CreatePlayer(t *testing.T) {
	s := newTestStorage(t)

	p1 := s.CreatePlayer(*model.NewPlayer("Алексей Иванов"))
	assert.Equal(t, 1, p1.ID())
	assert.Equal(t, "Алексей Иванов", p1.Name())

	p2 := s.CreatePlayer(*model.NewPlayer("Саша Зайцев"))
	assert.Equal(t, 2, p2.ID(), "ID инкрементится")
	assert.Equal(t, "Саша Зайцев", p2.Name())
}

func TestStorage_GetPlayerByID(t *testing.T) {
	s := newTestStorage(t)
	p := s.CreatePlayer(*model.NewPlayer("Павел Недуров"))

	got, found := s.GetPlayerByID(p.ID())
	require.True(t, found)
	assert.Equal(t, p.ID(), got.ID())

	_, found = s.GetPlayerByID(999)
	assert.False(t, found)
}

func TestStorage_UpdatePlayers(t *testing.T) {
	s := newTestStorage(t)
	p := s.CreatePlayer(*model.NewPlayer("Вася Обама"))

	err := s.UpdatePlayer(p.ID(), *model.NewPlayer("Вася Пупкин"))
	require.NoError(t, err)

	got, _ := s.GetPlayerByID(p.ID())
	assert.Equal(t, "Вася Пупкин", got.Name())

	err = s.UpdatePlayer(999, *model.NewPlayer("Мистер Икс"))
	assert.Error(t, err)
}

func TestStorage_DeletePlayer(t *testing.T) {
	s := newTestStorage(t)
	p := s.CreatePlayer(*model.NewPlayer("Макс Первый"))

	require.NoError(t, s.DeletePlayer(p.ID()))
	_, found := s.GetPlayerByID(p.ID())
	assert.False(t, found)

	assert.Error(t, s.DeletePlayer(p.ID()), "повторное удаление - ошибка")
}

func TestStorage_CreateGame_NewPlayers(t *testing.T) {
	s := newTestStorage(t)

	game := model.NewGame(
		*model.NewPlayer("Иван Иванов"),
		*model.NewPlayer("Сидр Сидоров"),
		8, 8,
	)

	created, err := s.CreateGame(*game)
	require.NoError(t, err)
	assert.Equal(t, 1, created.ID())

	players := s.GetAllPlayers()
	require.Len(t, players, 2)
}

func TestStorage_CreateGame_ExistingPlayerByName(t *testing.T) {
	s := newTestStorage(t)

	// заранее создаём игрока
	existing := s.CreatePlayer(*model.NewPlayer("Пётр Николаевич"))

	game := model.NewGame(
		*model.NewPlayer("Пётр Николаевич"),
		*model.NewPlayer("Пётр Владиславович"),
		8, 8,
	)

	created, err := s.CreateGame(*game)
	require.NoError(t, err)

	assert.Equal(t, existing.ID(), created.Player1().ID())
}

func TestStorage_CreateGame_DuplicateNameError(t *testing.T) {
	s := newTestStorage(t)

	// два игрока с одинаковым именем
	s.CreatePlayer(*model.NewPlayer("Миша Герасимов"))
	s.CreatePlayer(*model.NewPlayer("Миша Герасимов"))

	game := model.NewGame(
		*model.NewPlayer("Миша Герасимов"),
		*model.NewPlayer("Миша Герасимов"),
		8, 8,
	)

	_, err := s.CreateGame(*game)
	assert.Error(t, err, "дубликат имени - ошибка")
}

func TestStorage_GetGameByID(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	got, found := s.GetGameByID(created.ID())
	require.True(t, found)
	assert.Equal(t, created.ID(), got.ID())

	_, found = s.GetGameByID(999)
	assert.False(t, found)
}

func TestStorage_DeleteGame_CleansMoves(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())
	_, err := s.CreateMove(*move)
	require.NoError(t, err)
	assert.Len(t, s.GetAllMoves(), 1)

	require.NoError(t, s.DeleteGame(created.ID()))
	assert.Empty(t, s.GetAllMoves(), "ходы должны быть очищены")
}

func TestStorage_CreateMove_RequiresGame(t *testing.T) {
	s := newTestStorage(t)

	move := model.NewMove(1, 4, 3, 4)
	_, err := s.CreateMove(*move)
	assert.Error(t, err)

	move2 := model.NewMove(1, 4, 3, 4)
	move2.SetGameID(999)
	_, err = s.CreateMove(*move2)
	assert.Error(t, err)
}

func TestStorage_CreateMove(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())

	saved, err := s.CreateMove(*move)
	require.NoError(t, err)
	assert.Equal(t, 1, saved.ID())
}

func TestStorage_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()

	// первое хранилище - создаём данные
	s1 := repository.NewStorage(dir)
	require.NoError(t, s1.LoadFromFiles())

	p := s1.CreatePlayer(*model.NewPlayer("Владимир Владимирович"))
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s1.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())
	_, err := s1.CreateMove(*move)
	require.NoError(t, err)

	// сохраняем всё на диск
	require.NoError(t, s1.SaveAll())

	// второе хранилище - читаем те же данные
	s2 := repository.NewStorage(dir)
	require.NoError(t, s2.LoadFromFiles())

	// игрок на месте
	got, found := s2.GetPlayerByID(p.ID())
	require.True(t, found)
	assert.Equal(t, "Владимир Владимирович", got.Name())

	// игра на месте
	_, found = s2.GetGameByID(created.ID())
	assert.True(t, found)

	// ход на месте
	assert.Len(t, s2.GetAllMoves(), 1)

	// счётчик ID продолжается - новый игрок получит ID = 2
	newP := s2.CreatePlayer(*model.NewPlayer("Олег Петросян"))
	assert.Equal(t, 4, newP.ID())
}

func TestStorage_UpdateGame(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	// обновляем игроков
	newGame := model.NewGame(*model.NewPlayer("C"), *model.NewPlayer("D"), 8, 8)
	err := s.UpdateGame(created.ID(), *newGame)
	require.NoError(t, err)

	got, _ := s.GetGameByID(created.ID())
	// "C" и "D" автоматически создались
	assert.NotZero(t, got.Player1().ID())
	assert.NotZero(t, got.Player2().ID())

	// несуществующий ID
	err = s.UpdateGame(999, *newGame)
	assert.Error(t, err)
}

func TestStorage_OverwriteGame(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	// сохраняем ID игроков
	origP1ID := created.Player1().ID()

	// перезаписываем "как есть" - без resolvePlayer
	updated := created
	updated.SetID(created.ID())

	err := s.OverwriteGame(created.ID(), updated)
	require.NoError(t, err)

	got, _ := s.GetGameByID(created.ID())
	assert.Equal(t, origP1ID, got.Player1().ID(), "ID игроков не пересчитывается")

	// несуществующий ID
	err = s.OverwriteGame(999, *game)
	assert.Error(t, err)
}

func TestStorage_GetMoveByID(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())
	saved, _ := s.CreateMove(*move)

	got, found := s.GetMoveByID(saved.ID())
	require.True(t, found)
	assert.Equal(t, saved.ID(), got.ID())

	_, found = s.GetMoveByID(999)
	assert.False(t, found)
}

func TestStorage_UpdateMove(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())
	saved, _ := s.CreateMove(*move)

	// обновляем координаты
	newMove := model.NewMove(2, 3, 4, 3)
	newMove.SetGameID(created.ID())

	err := s.UpdateMove(saved.ID(), *newMove)
	require.NoError(t, err)

	got, _ := s.GetMoveByID(saved.ID())
	assert.Equal(t, 2, got.FromRow)
	assert.Equal(t, 3, got.FromCol)

	// без gameID
	noGame := model.NewMove(2, 3, 4, 3)
	err = s.UpdateMove(saved.ID(), *noGame)
	assert.Error(t, err)

	// несуществующая игра
	badGame := model.NewMove(2, 3, 4, 3)
	badGame.SetGameID(999)
	err = s.UpdateMove(saved.ID(), *badGame)
	assert.Error(t, err)

	// несуществующий ход
	valid := model.NewMove(2, 3, 4, 3)
	valid.SetGameID(created.ID())
	err = s.UpdateMove(999, *valid)
	assert.Error(t, err)
}

func TestStorage_DeleteMove(t *testing.T) {
	s := newTestStorage(t)
	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())
	saved, _ := s.CreateMove(*move)

	require.NoError(t, s.DeleteMove(saved.ID()))
	_, found := s.GetMoveByID(saved.ID())
	assert.False(t, found)

	// повторное удаление
	assert.Error(t, s.DeleteMove(saved.ID()))
}

func TestStorage_ResolvePlayer_NotFoundByID(t *testing.T) {
	s := newTestStorage(t)

	p := model.NewPlayer("Иванов")
	p.SetID(999) // несуществующий

	game := model.NewGame(*p, *model.NewPlayer("Петров"), 8, 8)
	_, err := s.CreateGame(*game)
	assert.Error(t, err, "игрок с ID=999 не существует")
}

func TestStorage_ResolvePlayer_EmptyName(t *testing.T) {
	s := newTestStorage(t)

	game := model.NewGame(*model.NewPlayer(""), *model.NewPlayer("Петров"), 8, 8)
	_, err := s.CreateGame(*game)
	assert.Error(t, err, "имя игрока обязательно")
}

func TestStorage_CleanOrphanMoves(t *testing.T) {
	dir := t.TempDir()

	// Шаг 1: создаём игру с ходом, сохраняем на диск
	s1 := repository.NewStorage(dir)
	require.NoError(t, s1.LoadFromFiles())

	game := model.NewGame(*model.NewPlayer("A"), *model.NewPlayer("B"), 8, 8)
	created, _ := s1.CreateGame(*game)

	move := model.NewMove(1, 4, 3, 4)
	move.SetGameID(created.ID())
	s1.CreateMove(*move)
	require.NoError(t, s1.SaveAll())

	// Шаг 2: вручную удаляем игру из games.json, оставляя ход в moves.json
	// (симуляция "игра удалена, а ходы остались")
	gamesPath := filepath.Join(dir, "games.json")
	require.NoError(t, os.WriteFile(gamesPath, []byte("[]"), 0644))

	// Шаг 3: новое хранилище при загрузке должно почистить orphans
	s2 := repository.NewStorage(dir)
	require.NoError(t, s2.LoadFromFiles())

	// ходы должны быть очищены
	assert.Empty(t, s2.GetAllMoves(), "осиротевшие ходы должны быть удалены")
}
