package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func UnaryLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	log.Printf("Received request: %s  payload: %v", info.FullMethod, req)

	resp, err := handler(ctx, req)

	log.Printf("Sent response: %s, duration: %v , response: %v , error: %v \n\n", info.FullMethod, time.Since(start), resp, err)
	return resp, err
}
