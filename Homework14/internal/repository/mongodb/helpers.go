package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoStorage) getNextSequence(ctx context.Context, counterName string) (int, error) {
	collection := s.db.Collection("counters")
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var counter struct {
		Seq int `bson:"seq"`
	}

	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": counterName}, update, opts).Decode(&counter)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения счетчика %s: %w", counterName, err)
	}
	return counter.Seq, nil
}

func (s *MongoStorage) logToRedis(ctx context.Context, entityType string, entityID int, action string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга лога: %w", err)
	}

	key := fmt.Sprintf("log:%s:%d:%d", entityType, entityID, time.Now().UnixNano())
	value := fmt.Sprintf("[%s] %s", action, string(data))

	err = s.redisClient.Set(ctx, key, value, s.redisTTL).Err()
	if err != nil {
		return fmt.Errorf("ошибка записи лога в Redis: %w", err)
	}
	return nil
}
