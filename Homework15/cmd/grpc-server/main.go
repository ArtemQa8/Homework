package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	grpchandler "mod.go/internal/api/grpc"
	"mod.go/internal/config"
	pb "mod.go/internal/proto"
	"mod.go/internal/repository/mongodb"
	"mod.go/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	storage, err := mongodb.NewMongoStorage(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("не удалось открыть порт 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

	playerService := service.NewPlayerService(storage)
	gameService := service.NewGameService(storage)
	moveService := service.NewMoveService(storage)

	pb.RegisterPlayerServiceServer(grpcServer, grpchandler.NewPlayerServer(playerService))
	pb.RegisterGameServiceServer(grpcServer, grpchandler.NewGameServer(gameService))
	pb.RegisterMoveServiceServer(grpcServer, grpchandler.NewMoveServer(moveService))

	go func() {
		fmt.Println("gRPC-сервер запущен на :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	fmt.Println("\nПолучен сигнал остановки, завершаем работу...")
	grpcServer.GracefulStop()

	if err := storage.Close(context.Background()); err != nil {
		log.Printf("Ошибка при закрытии БД: %v", err)
	}

	fmt.Println("gRPC-сервер остановлен")
}
