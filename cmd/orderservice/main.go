package main

import (
	"310271-sergeykhairedinov-course-1343/internal/interceptor"
	"310271-sergeykhairedinov-course-1343/internal/service" // Импорт сервиса
	pb "310271-sergeykhairedinov-course-1343/pkg/api/test"  // Импорт сгенерированного кода
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv" // Для загрузки env файла
	"google.golang.org/grpc"   // Работа с grpc
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	//подгружаем env файл для считывания порта
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error of loads env file")
	}
	GRPC_PORT := os.Getenv("GRPC_PORT")
	HTTP_PORT := os.Getenv("HTTP_PORT")

	if GRPC_PORT == "" {
		log.Printf("error of load env file, server work on default 50051 port")
		GRPC_PORT = "50051" // значение по умолчанию если что
	}
	if HTTP_PORT == "" {
		log.Printf("error of load env file, HTTP server work on default 8080 port")
		HTTP_PORT = "8080" // значение по умолчанию для HTTP
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

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//реализация Gateway

	rt := runtime.NewServeMux()
	opts := []grpc.DialOption{(grpc.WithTransportCredentials(insecure.NewCredentials()))}
	err = pb.RegisterOrderServiceHandlerFromEndpoint(ctx, rt, "localhost:"+GRPC_PORT, opts)
	if err != nil {
		log.Fatalf("failed to serve", error(err))
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", HTTP_PORT), rt); err != nil {
			log.Printf("failed to serve", error(err))
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		log.Printf("Server Stopped")
	}
}
