# Переменные
PROTO_DIR=./api
PROTO_FILE=$(PROTO_DIR)/order.proto
GO_OUT_DIR=./pkg/api/test
BINARY_NAME=orderservice
SERVER_DIR=./cmd/orderservice

# Генерация кода для gRPC и gRPC Gateway
generate:
	protoc -I $(PROTO_DIR) -I $(PROTO_DIR)/google \
		--go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GO_OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_FILE)

# Сборка бинарника
build: generate
	go build -o $(BINARY_NAME) $(SERVER_DIR)

# Запуск сервера
run: build
	./$(BINARY_NAME)

# Очистка сгенерированных файлов и бинарника
clean:
	rm -f $(GO_OUT_DIR)/*.pb.go $(GO_OUT_DIR)/*.gw.go
	rm -f $(BINARY_NAME)

# Команда по умолчанию (выполняется при вызове `make` без аргументов)
.DEFAULT_GOAL := run