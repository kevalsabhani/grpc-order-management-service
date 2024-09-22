package main

import (
	"context"
	"io"
	"log"

	"github.com/golang/protobuf/ptypes/wrappers"
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

	// Add order 1
	order := pb.Order{
		Items:       []string{"Google Pixel 3A", "Mac Book Pro", "Apple Watch S4"},
		Destination: "San Jose, CA",
		Price:       2300.00,
	}
	addResp1, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 1 added: %s", addResp1.Value)
	log.Println()

	// Add order 2
	order = pb.Order{
		Items:       []string{"Google Pixel 5A", "Mac Book M3", "Apple Watch S7"},
		Destination: "San Jose, CA",
		Price:       2300.00,
	}
	addResp2, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 2 added: %s", addResp2.Value)
	log.Println()

	// Get order
	getResp, err := c.GetOrder(context.Background(), &pb.OrderID{Value: addResp1.Value})
	if err != nil {
		log.Fatalf("Could not get order: %v", err)
	}
	log.Printf("Order retrieved: %s", getResp.String())
	log.Println()

	// Search orders
	searchStream, err := c.SearchOrder(context.Background(), &wrappers.StringValue{Value: "Google"})
	if err != nil {
		log.Fatalf("Could not search order: %v", err)
	}
	for {
		searchResp, err := searchStream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("Order found: %s", searchResp.String())
	}
	log.Println()

	// Update orders
	updateStream, err := c.UpdateOrder(context.Background())
	if err != nil {
		log.Fatalf("could not update orders: %v", err)
	}

	// Updating order 1
	updateOrder1 := pb.Order{
		Id:          addResp1.Value,
		Items:       []string{"Motorola edge 50 fusion", "iPhone 16"},
		Destination: "San Jose, CA",
		Price:       2700.00,
	}
	err = updateStream.Send(&updateOrder1)
	if err != nil {
		log.Fatalf("failed to send Order 1: %v", err)
	}
	log.Println()

	// Updating order 2
	updateOrder2 := pb.Order{
		Id:          addResp2.Value,
		Items:       []string{"Google pixel 9"},
		Destination: "San Jose, CA",
		Price:       500.00,
	}
	err = updateStream.Send(&updateOrder2)
	if err != nil {
		log.Fatalf("failed to send Order 2: %v", err)
	}
	log.Println()

	updateResp, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to fetch resp: %v", err)
	}
	log.Println(updateResp.Value)
}
