package service

import (
	pb "310271-sergeykhairedinov-course-1343/pkg/api/test"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type OrderServes struct {
	mu     sync.Mutex
	orders map[string]*pb.Order
	pb.UnimplementedOrderServiceServer
}

func NewOrderService() *OrderServes {
	return &OrderServes{orders: make(map[string]*pb.Order)}
}

func (s *OrderServes) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := req.GetItem()
	quantity := req.GetQuantity()
	id := uuid.New().String()

	order := &pb.Order{
		Id:       id,
		Item:     item,
		Quantity: quantity,
	}
	s.orders[id] = order
	return &pb.CreateOrderResponse{Id: id}, nil

}

func (s *OrderServes) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	id := req.GetId()
	s.mu.Lock()
	defer s.mu.Unlock()
	order, err := s.orders[id]
	if !err {
		return nil, status.Errorf(codes.NotFound, "order whith id %s not found", id)
	}
	return &pb.GetOrderResponse{Order: order}, nil
}

func (s *OrderServes) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := req.GetItem()
	quantity := req.GetQuantity()
	id := req.GetId()

	order := &pb.Order{
		Id:       id,
		Item:     item,
		Quantity: quantity,
	}
	s.orders[id] = order
	return &pb.UpdateOrderResponse{Order: order}, nil

}

func (s *OrderServes) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	id := req.Id
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.orders[id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order with id %s not found", id)
	}
	delete(s.orders, id)
	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderServes) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	SliceOfOrders := make([]*pb.Order, 0, len(s.orders))

	for _, value := range s.orders {
		SliceOfOrders = append(SliceOfOrders, value)
	}
	return &pb.ListOrdersResponse{Orders: SliceOfOrders}, nil
}
