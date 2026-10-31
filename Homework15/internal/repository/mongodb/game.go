package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mod.go/internal/model"
)

type boardDoc struct {
	Rows  int           `bson:"rows"`
	Cols  int           `bson:"cols"`
	Cells [][]*pieceDoc `bson:"cells"`
}

type gameDoc struct {
	ID           int         `bson:"id"`
	Player1      playerDoc   `bson:"player1"`
	Player2      playerDoc   `bson:"player2"`
	Player1Color model.Color `bson:"player1_color"`
	Player2Color model.Color `bson:"player2_color"`
	Board        *boardDoc   `bson:"board"`
	Current      model.Color `bson:"current"`
	Moves        []moveDoc   `bson:"moves"`
}

// Создать игру
func (s *MongoStorage) CreateGame(game model.Game) (model.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p1, err := s.resolvePlayer(ctx, game.Player1())
	if err != nil {
		return model.Game{}, err
	}
	p2, err := s.resolvePlayer(ctx, game.Player2())
	if err != nil {
		return model.Game{}, err
	}
	game.SetPlayer1(p1)
	game.SetPlayer2(p2)

	id, err := s.getNextSequence(ctx, "game_id")
	if err != nil {
		return model.Game{}, fmt.Errorf("ошибка получения ID игры: %w", err)
	}
	game.SetID(id)

	doc := toGameDoc(game)

	collection := s.db.Collection("games")
	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return model.Game{}, fmt.Errorf("ошибка сохранения игры: %w", err)
	}

	_ = s.logToRedis(ctx, "game", id, "CREATE", game)
	return game, nil
}

// Получить список всех игр.
func (s *MongoStorage) GetAllGames() []model.Game {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("games")
	opts := options.Find().SetSort(bson.M{"id": 1})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return []model.Game{}
	}
	defer cursor.Close(ctx)

	var docs []gameDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return []model.Game{}
	}

	games := make([]model.Game, 0, len(docs))
	for _, doc := range docs {
		games = append(games, fromGameDoc(doc))
	}
	return games
}

// Получить игру по ID
func (s *MongoStorage) GetGameByID(id int) (model.Game, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("games")
	var doc gameDoc
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		return model.Game{}, false
	}

	return fromGameDoc(doc), true
}

// Обновить игру
func (s *MongoStorage) UpdateGame(id int, newGame model.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p1, err := s.resolvePlayer(ctx, newGame.Player1())
	if err != nil {
		return err
	}
	p2, err := s.resolvePlayer(ctx, newGame.Player2())
	if err != nil {
		return err
	}
	newGame.SetPlayer1(p1)
	newGame.SetPlayer2(p2)

	doc := toGameDoc(newGame)
	collection := s.db.Collection("games")

	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("ошибка обновления игры: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Игра с ID %d не найдена", id)
	}

	_ = s.logToRedis(ctx, "game", id, "UPDATE", newGame)
	return nil
}

// OverwriteGame перезаписывает игру целиком (используется после MakeMove).
func (s *MongoStorage) OverwriteGame(id int, newGame model.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := toGameDoc(newGame)

	collection := s.db.Collection("games")
	filter := bson.M{"id": id}
	update := bson.M{"$set": doc}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("ошибка обновления игры: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Игра с ID %d не найдена", id)
	}

	_ = s.logToRedis(ctx, "game", id, "UPDATE", newGame)
	return nil
}

// Удалить игру
func (s *MongoStorage) DeleteGame(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("games")
	result, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка удаления игры: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("Игра с ID %d не найдена", id)
	}

	_ = s.logToRedis(ctx, "game", id, "DELETE", map[string]int{"id": id})
	return nil
}

// Конвертеры
func toGameDoc(g model.Game) gameDoc {
	p1 := playerDoc{ID: g.Player1().ID(), Name: g.Player1().Name()}
	p2 := playerDoc{ID: g.Player2().ID(), Name: g.Player2().Name()}

	var board *boardDoc
	if g.Board() != nil {
		b := g.Board()
		cells := make([][]*pieceDoc, b.Rows())
		for r := 0; r < b.Rows(); r++ {
			cells[r] = make([]*pieceDoc, b.Cols())
			for c := 0; c < b.Cols(); c++ {
				cells[r][c] = toPieceDoc(b.PieceAt(r, c))
			}
		}
		board = &boardDoc{Rows: b.Rows(), Cols: b.Cols(), Cells: cells}
	}

	moves := make([]moveDoc, len(g.Moves()))
	for i, m := range g.Moves() {
		moves[i] = toMoveDoc(m)
	}

	return gameDoc{
		ID:           g.ID(),
		Player1:      p1,
		Player2:      p2,
		Player1Color: g.Player1Color(),
		Player2Color: g.Player2Color(),
		Board:        board,
		Current:      g.CurrentColor(),
		Moves:        moves,
	}
}

func fromGameDoc(doc gameDoc) model.Game {
	p1 := model.NewPlayer(doc.Player1.Name)
	p1.SetID(doc.Player1.ID)
	p2 := model.NewPlayer(doc.Player2.Name)
	p2.SetID(doc.Player2.ID)

	var board *model.Board
	if doc.Board != nil && doc.Board.Rows > 0 && doc.Board.Cols > 0 {
		board = model.NewBoard(doc.Board.Rows, doc.Board.Cols)
		for r := 0; r < doc.Board.Rows; r++ {
			for c := 0; c < doc.Board.Cols; c++ {
				if r < len(doc.Board.Cells) && c < len(doc.Board.Cells[r]) {
					board.SetPiece(r, c, fromPieceDoc(doc.Board.Cells[r][c]))
				}
			}
		}
	}

	game := model.NewGameWithBoard(*p1, *p2, board)
	game.SetID(doc.ID)
	game.SetPlayer1Color(doc.Player1Color)
	game.SetPlayer2Color(doc.Player2Color)
	game.SetCurrentColor(doc.Current)

	moves := make([]model.Move, len(doc.Moves))
	for i, m := range doc.Moves {
		moves[i] = fromMoveDoc(m)
	}
	game.SetMoves(moves)

	return *game
}

func (s *MongoStorage) resolvePlayer(ctx context.Context, player model.Player) (model.Player, error) {
	if player.ID() != 0 {
		p, found := s.GetPlayerByID(player.ID())
		if !found {
			return model.Player{}, fmt.Errorf("Игрок с ID %d не найден", player.ID())
		}
		return p, nil
	}

	if player.Name() == "" {
		return model.Player{}, errors.New("Имя игрока обязательно")
	}

	collection := s.db.Collection("players")
	var doc playerDoc
	err := collection.FindOne(ctx, bson.M{"name": player.Name()}).Decode(&doc)

	if err == nil {
		p := model.NewPlayer(doc.Name)
		p.SetID(doc.ID)
		return *p, nil
	}

	if err != mongo.ErrNoDocuments {
		return model.Player{}, fmt.Errorf("ошибка поиска игрока: %w", err)
	}

	return s.CreatePlayer(player), nil
}
