package main

import (
	"log"
	"net"
	"os"

	"310271-sergeykhairedinov-course-1343/internal/interceptor"
	"310271-sergeykhairedinov-course-1343/internal/service" // Импорт сервиса
	pb "310271-sergeykhairedinov-course-1343/pkg/api/test"  // Импорт сгенерированного кода
	"github.com/joho/godotenv"                              // Для загрузки env файла
	"google.golang.org/grpc"                                // Работа с grpc
)

func main() {
	//подгружаем env файл для считывания порта
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error of loads env file")
	}
	GRPC_PORT := os.Getenv("GRPC_PORT")

	if GRPC_PORT == "" {
		log.Printf("error of load env file, server work on default 50051 port")
		GRPC_PORT = "50051" // значение по умолчанию если что
	}

	// Создаем gRPC-сервер и подключаем к нему интерсептор
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryLogger))

	// Создаем экземпляр нашего сервиса
	orderService := service.NewOrderService()

	// Регистрируем наш сервис в gRPC-сервере
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	// Запускаем сервер на порту

	lis, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
