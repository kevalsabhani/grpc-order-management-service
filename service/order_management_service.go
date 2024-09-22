package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	pb "github.com/kevalsabhani/grpc-order-management-service/service/ecommerce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	orderBatchSize = 3
)

// orderManagementService implements the OrderManagementServer proto defined in ecommerce.proto
type orderManagementService struct {
	orders map[string]*pb.Order
	pb.UnimplementedOrderManagementServer
}

// GetOrder retrieves an order from the map using the order id.
// If the order id does not exist, it returns a status error with code NotFound.
// Unary RPC
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
// Unary RPC
func (s *orderManagementService) AddOrder(ctx context.Context, req *pb.Order) (*pb.OrderID, error) {
	// Generate a new UUID for the order
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	// Set the Id field of the order
	req.Id = id.String()
	// Add the order to the map
	if s.orders == nil {
		s.orders = make(map[string]*pb.Order)
	}
	s.orders[req.Id] = req
	// Return the Order ID as a proto message
	return &pb.OrderID{Value: req.Id}, status.New(codes.OK, "").Err()
}

// SearchOrders returns a stream of orders that match the search query.
// Server-streaming RPC
func (s *orderManagementService) SearchOrder(searchQuery *wrappers.StringValue, stream pb.OrderManagement_SearchOrderServer) error {
	for _, order := range s.orders {
		for _, item := range order.Items {
			if strings.Contains(item, searchQuery.Value) {
				err := stream.Send(order)
				if err != nil {
					return fmt.Errorf("error sending response: %v", err)
				}
				log.Printf("Matching order found: %s", order.Id)
				break
			}
		}
	}
	return nil
}

// UpdateOrder updates order map with orders received from client
// Client-streaming RPC
func (s *orderManagementService) UpdateOrder(stream pb.OrderManagement_UpdateOrderServer) error {
	respStr := "Updated Order IDs: "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wrappers.StringValue{Value: "Order processed:" + respStr})
		}

		s.orders[order.Id] = order
		log.Printf("Order ID: %v updated.", order.Id)
		respStr += order.Id + ", "
	}
}

// ProcessOrder process orders and combined them in different shipments based on the location
// Bidirectional-streaming RPC
func (s *orderManagementService) ProcessOrder(stream pb.OrderManagement_ProcessOrderServer) error {
	batchMaker := 1
	combShipmentMap := make(map[string]*pb.CombinedShipment)
	for {
		orderId, err := stream.Recv()
		if err != nil {
			// log.Printf("EOF: %s", orderId.Value)
			for _, shipment := range combShipmentMap {
				if err2 := stream.Send(shipment); err2 != nil {
					log.Print("error here 2")
					return err2
				}
			}
			return nil
		}

		destination := s.orders[orderId.Value].Destination
		shipment, found := combShipmentMap[destination]
		if found {
			shipment.OrderList = append(shipment.OrderList, s.orders[orderId.Value])
			combShipmentMap[destination] = shipment
		} else {
			combShipmentMap[destination] = &pb.CombinedShipment{
				Id:        "cmb - " + s.orders[orderId.Value].Destination,
				Status:    "processed",
				OrderList: []*pb.Order{s.orders[orderId.Value]},
			}
		}

		if batchMaker == orderBatchSize {
			for _, comb := range combShipmentMap {
				log.Printf("Shipping: %v -> %v", comb.Id, len(comb.OrderList))
				if err = stream.Send(comb); err != nil {
					log.Print("error here")
					return err
				}
			}
			batchMaker = 0
			combShipmentMap = make(map[string]*pb.CombinedShipment)
		} else {
			batchMaker++
		}
	}
}
