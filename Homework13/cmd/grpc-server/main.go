package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	grpchandler "mod.go/internal/api/grpc"
	pb "mod.go/internal/proto"
	"mod.go/internal/repository"
	"mod.go/internal/service"
)

func main() {
	storage := repository.NewStorage()
	if err := storage.LoadFromFiles(); err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
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

	if err := storage.SaveAll(); err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
	}

	fmt.Println("gRPC-сервер остановлен, данные сохранены.")
}
