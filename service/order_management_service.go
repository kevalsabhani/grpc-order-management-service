package main

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/kevalsabhani/grpc-order-management-service/service/ecommerce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// orderManagementService implements the OrderManagementServer proto defined in ecommerce.proto
type orderManagementService struct {
	orders map[string]*pb.Order
	pb.UnimplementedOrderManagementServer
}

// GetOrder retrieves an order from the map using the order id.
// If the order id does not exist, it returns a status error with code NotFound.
func (s *orderManagementService) GetOrder(ctx context.Context, req *pb.OrderID) (*pb.Order, error) {
	order, exists := s.orders[req.Value]
	if exists {
		return order, nil
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist: %s", req.Value)
}

// AddOrder adds a new order to the map.
// It generates a new UUID for the order and sets the Id field of the order,
// then adds the order to the map.
// The method returns an OrderID proto message with the value set to the
// newly generated UUID.
func (s *orderManagementService) AddOrder(ctx context.Context, req *pb.Order) (*pb.OrderID, error) {
	// Generate a new UUID for the order
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	// Set the Id field of the order
	req.Id = id.String()
	// Add the order to the map
	s.orders[req.Id] = req
	// Return the Order ID as a proto message
	return &pb.OrderID{Value: req.Id}, status.New(codes.OK, "").Err()
}
