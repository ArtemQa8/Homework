package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mod.go/internal/config"
	"mod.go/internal/model"
)

type MongoStorage struct {
	client      *mongo.Client
	db          *mongo.Database
	redisClient *redis.Client
	redisTTL    time.Duration
}

func NewMongoStorage(cfg *config.Config) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB не отвечает: %w", err)
	}

	db := mongoClient.Database(cfg.MongoDB)

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis не отвечает: %w", err)
	}

	fmt.Println("Успешно подклюено к MongoDB и Redis!")

	return &MongoStorage{
		client:      mongoClient,
		db:          db,
		redisClient: redisClient,
		redisTTL:    cfg.RedisTTL,
	}, nil
}

// Закрываем соединение Redis и MongoDB
func (s *MongoStorage) Close(ctx context.Context) error {
	if err := s.client.Disconnect(ctx); err != nil {
		return err
	}
	return s.redisClient.Close()
}

// Импорт игрока с сохранением старого ID
func (s *MongoStorage) ImportPlayer(id int, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doc := playerDoc{ID: id, Name: name}
	_, err := s.db.Collection("players").InsertOne(ctx, doc)
	return err
}

// Импорт игры с сохранением старого ID
func (s *MongoStorage) ImportGame(game model.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doc := toGameDoc(game)
	_, err := s.db.Collection("games").InsertOne(ctx, doc)
	return err
}

// Импорт хода с сохранением старого ID
func (s *MongoStorage) ImportMove(move model.Move) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doc := toMoveDoc(move)
	_, err := s.db.Collection("moves").InsertOne(ctx, doc)
	return err
}

// Установка значения счетчика (чтобы новые ID начинались после старых)
func (s *MongoStorage) SetCounter(counterName string, value int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{"$set": bson.M{"seq": value}}
	_, err := s.db.Collection("counters").UpdateOne(
		ctx,
		bson.M{"_id": counterName},
		update,
		options.Update().SetUpsert(true),
	)
	return err
}
