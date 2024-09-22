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

	// Add order 3
	order = pb.Order{
		Items:       []string{"Google Pixel 3A", "Mac Book Pro"},
		Destination: "Mountain View, CA",
		Price:       1800.00,
	}
	addResp3, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 3 added: %s", addResp3.Value)

	// Add order 4
	order = pb.Order{
		Items:       []string{"Apple Watch S4"},
		Destination: "San Jose, CA",
		Price:       400.00,
	}
	addResp4, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 4 added: %s", addResp4.Value)

	// Add order 5
	order = pb.Order{
		Items:       []string{"Google Home Mini", "Google Nest Hub"},
		Destination: "Mountain View, CA",
		Price:       400.00,
	}
	addResp5, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 5 added: %s", addResp5.Value)

	// Add order 6
	order = pb.Order{
		Items:       []string{"Amazon Echo"},
		Destination: "San Jose, CA",
		Price:       30.00,
	}
	addResp6, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 6 added: %s", addResp6.Value)

	// Add order 7
	order = pb.Order{
		Items:       []string{"Amazon Echo", "Apple iPhone XS"},
		Destination: "Mountain View, CA",
		Price:       300.00,
	}
	addResp7, err := c.AddOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Could not add order: %v", err)
	}
	log.Printf("Order 7 added: %s", addResp7.Value)
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
	log.Println()

	// Process orders
	processStream, err := c.ProcessOrder(context.Background())
	if err != nil {
		log.Fatalf("could not process orders: %v", err)
	}

	// process order 1
	if err = processStream.Send(addResp1); err != nil {
		log.Fatalf("could not send order %v", addResp1.Value)
	}

	// process order 2
	if err = processStream.Send(addResp2); err != nil {
		log.Fatalf("could not send order %v", addResp2.Value)
	}

	ch := make(chan struct{})
	go fetchStreamingResp(processStream, ch)

	// process order 3
	if err = processStream.Send(addResp3); err != nil {
		log.Fatalf("could not send order %v", addResp3.Value)
	}

	// process order 4
	if err = processStream.Send(addResp4); err != nil {
		log.Fatalf("could not send order %v", addResp4.Value)
	}

	// process order 5
	if err = processStream.Send(addResp5); err != nil {
		log.Fatalf("could not send order %v", addResp5.Value)
	}

	// process order 6
	if err = processStream.Send(addResp6); err != nil {
		log.Fatalf("could not send order %v", addResp6.Value)
	}

	// process order 7
	if err = processStream.Send(addResp7); err != nil {
		log.Fatalf("could not send order %v", addResp7.Value)
	}

	if err = processStream.CloseSend(); err != nil {
		log.Fatal(err)
	}

	<-ch
}

func fetchStreamingResp(processStream pb.OrderManagement_ProcessOrderClient, ch chan<- struct{}) {
	for {
		shipment, err := processStream.Recv()
		if err == io.EOF || err != nil {
			break
		}

		log.Printf("Shipment orders: %v", shipment)
	}
	ch <- struct{}{}
}
