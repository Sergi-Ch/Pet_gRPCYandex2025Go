package main

import (
	"log"
	"net"

	"310271-sergeykhairedinov-course-1343/internal/service" // Импорт твоего сервиса
	pb "310271-sergeykhairedinov-course-1343/pkg/api/test"  // Импорт сгенерированного кода
	"google.golang.org/grpc"
)

func main() {
	// Создаем gRPC-сервер
	grpcServer := grpc.NewServer()

	// Создаем экземпляр нашего сервиса
	orderService := service.NewOrderService()

	// Регистрируем наш сервис в gRPC-сервере
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	// Запускаем сервер на порту 2051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
