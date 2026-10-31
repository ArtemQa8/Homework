package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mod.go/internal/model"
)

type pieceDoc struct {
	Color model.Color     `bson:"color"`
	Type  model.PieceType `bson:"type"`
}

func toPieceDoc(p *model.Piece) *pieceDoc {
	if p == nil {
		return nil
	}
	return &pieceDoc{Color: p.Color(), Type: p.Type()}
}

func fromPieceDoc(doc *pieceDoc) *model.Piece {
	if doc == nil {
		return nil
	}
	return model.NewPiece(doc.Color, doc.Type)
}

type moveDoc struct {
	ID         int             `bson:"id"`
	GameID     int             `bson:"game_id"`
	FromRow    int             `bson:"from_row"`
	FromCol    int             `bson:"from_col"`
	ToRow      int             `bson:"to_row"`
	ToCol      int             `bson:"to_col"`
	Captured   *pieceDoc       `bson:"captured,omitempty"`
	Promotion  model.PieceType `bson:"promotion"`
	MovedPiece *pieceDoc       `bson:"moved_piece,omitempty"`
	Check      bool            `bson:"check"`
	Mate       bool            `bson:"mate"`
}

// Создать ход.
func (s *MongoStorage) CreateMove(move model.Move) (model.Move, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, found := s.GetGameByID(move.GameID()); !found {
		return model.Move{}, fmt.Errorf("Игра с ID %d не найдена", move.GameID())
	}

	id, err := s.getNextSequence(ctx, "move_id")
	if err != nil {
		return model.Move{}, fmt.Errorf("ошибка получения ID хода: %w", err)
	}
	move.SetID(id)

	doc := toMoveDoc(move)

	collection := s.db.Collection("moves")
	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return model.Move{}, fmt.Errorf("ошибка сохранения хода: %w", err)
	}

	_ = s.logToRedis(ctx, "move", id, "CREATE", move)
	return move, nil
}

// Получить список всех ходов.
func (s *MongoStorage) GetAllMoves() []model.Move {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("moves")
	opts := options.Find().SetSort(bson.M{"id": 1})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return []model.Move{}
	}
	defer cursor.Close(ctx)

	var docs []moveDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return []model.Move{}
	}

	moves := make([]model.Move, 0, len(docs))
	for _, doc := range docs {
		moves = append(moves, fromMoveDoc(doc))
	}
	return moves
}

func (s *MongoStorage) GetMoveByID(id int) (model.Move, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("moves")
	var doc moveDoc
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		return model.Move{}, false
	}

	return fromMoveDoc(doc), true
}

func (s *MongoStorage) UpdateMove(id int, newMove model.Move) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := toMoveDoc(newMove)
	collection := s.db.Collection("moves")

	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("ошибка обновления хода: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Ход с ID %d не найден", id)
	}

	_ = s.logToRedis(ctx, "move", id, "UPDATE", newMove)
	return nil
}

func (s *MongoStorage) DeleteMove(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("moves")
	result, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка удаления хода: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("Ход с ID %d не найден", id)
	}

	_ = s.logToRedis(ctx, "move", id, "DELETE", map[string]int{"id": id})
	return nil
}

func toMoveDoc(m model.Move) moveDoc {
	return moveDoc{
		ID:         m.ID(),
		GameID:     m.GameID(),
		FromRow:    m.FromRow,
		FromCol:    m.FromCol,
		ToRow:      m.ToRow,
		ToCol:      m.ToCol,
		Captured:   toPieceDoc(m.Captured),
		Promotion:  m.Promotion,
		MovedPiece: toPieceDoc(m.MovedPiece),
		Check:      m.Check,
		Mate:       m.Mate,
	}
}

func fromMoveDoc(doc moveDoc) model.Move {
	m := model.Move{
		FromRow:    doc.FromRow,
		FromCol:    doc.FromCol,
		ToRow:      doc.ToRow,
		ToCol:      doc.ToCol,
		Captured:   fromPieceDoc(doc.Captured),
		Promotion:  doc.Promotion,
		MovedPiece: fromPieceDoc(doc.MovedPiece),
		Check:      doc.Check,
		Mate:       doc.Mate,
	}
	m.SetID(doc.ID)
	m.SetGameID(doc.GameID)
	return m
}
