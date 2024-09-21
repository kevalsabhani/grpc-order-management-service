package main

import (
	"context"
	"log"

	pb "github.com/kevalsabhani/grpc-order-management-service/client/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50052"
)

func main() {
	// Create a gRPC client connection to the server
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOrderManagementClient(conn)

	// Add order
	order := pb.Order{
		Items:       []string{"Google Pixel 3A", "Mac Book Pro", "Apple Watch S4"},
		Destination: "San Jose, CA",
		Price:       2300.00,
	}
	addResp, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order added: %s", addResp.Value)

	// Get order
	getResp, err := c.GetOrder(context.Background(), &pb.OrderID{Value: addResp.Value})
	if err != nil {
		log.Fatalf("Could not get order: %v", err)
	}
	log.Printf("Order retrieved: %s", getResp.String())
}
