# Переменные
PROTO_DIR=proto
PROTO_FILE=$(PROTO_DIR)/order.proto
OUT_DIR=pkg/api/test
SERVER_BINARY=bin/orderserver

# Генерация кода для gRPC
generate:
	protoc --go_out=. --go-grpc_out=. $(PROTO_FILE)

# Сборка бинарника
build: generate
	go build -o $(SERVER_BINARY) cmd/orderservice/main.go

# Запуск сервера
run: build
	./$(SERVER_BINARY)

# Очистка сгенерированных файлов и бинарника
clean:
	rm -rf $(OUT_DIR)/*.pb.go $(OUT_DIR)/*.grpc.pb.go $(SERVER_BINARY)

# Псевдоним для запуска всех команд по порядку
all: generate build run