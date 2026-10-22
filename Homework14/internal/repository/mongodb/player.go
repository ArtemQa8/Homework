package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mod.go/internal/model"
)

// т.к. в модели у нас поля приватные, мы добавляем такую структуру
type playerDoc struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

// Создание игрока
func (s *MongoStorage) CreatePlayer(player model.Player) model.Player {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := s.getNextSequence(ctx, "player_id")
	if err != nil {
		panic(fmt.Sprintf("ошибка получения ID игрока: %v", err))
	}
	player.SetID(id)

	doc := playerDoc{
		ID:   player.ID(),
		Name: player.Name(),
	}

	collection := s.db.Collection("players")
	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		panic(fmt.Sprintf("ошибка сохранения игрока в MongoDB: %v", err))
	}

	_ = s.logToRedis(ctx, "player", id, "CREATE", player)

	return player
}

// Получения списка всех игроков
func (s *MongoStorage) GetAllPlayers() []model.Player {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("players")
	opts := options.Find().SetSort(bson.M{"id": 1})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return []model.Player{}
	}
	defer cursor.Close(ctx)

	var docs []playerDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return []model.Player{}
	}

	players := make([]model.Player, 0, len(docs))
	for _, doc := range docs {
		p := model.NewPlayer(doc.Name)
		p.SetID(doc.ID)
		players = append(players, *p)
	}
	return players
}

// Получение игрока по ID.
func (s *MongoStorage) GetPlayerByID(id int) (model.Player, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("players")
	var doc playerDoc
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		return model.Player{}, false
	}

	p := model.NewPlayer(doc.Name)
	p.SetID(doc.ID)
	return *p, true
}

// Обновление игрока
func (s *MongoStorage) UpdatePlayer(id int, newPlayer model.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("players")
	update := bson.M{"$set": bson.M{"name": newPlayer.Name()}}

	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return fmt.Errorf("ошибка обновления игрока: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Игрок с ID %d не найден", id)
	}

	_ = s.logToRedis(ctx, "player", id, "UPDATE", newPlayer)
	return nil
}

// Удаление игрока
func (s *MongoStorage) DeletePlayer(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Collection("players")
	result, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка удаления игрока: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("Игрок с ID %d не найден", id)
	}

	_ = s.logToRedis(ctx, "player", id, "DELETE", map[string]int{"id": id})
	return nil
}
