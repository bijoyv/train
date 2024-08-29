package main

import (
	"context"
	"fmt"
	"log"

	train "github.com/bijoyv/train/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := train.NewTrainServiceClient(conn)

	// Example usage: Purchase a ticket
	user := &train.User{
		Email: "user@example.com",
	}
	purchaseRequest := &train.PurchaseTicketRequest{
		From: "London",
		To:   "Paris",
		User: user,
	}
	purchaseResponse, err := client.PurchaseTicket(context.Background(), purchaseRequest)
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	fmt.Println("Ticket purchased:", purchaseResponse.Ticket)

	// Example usage: Get a ticket
	getTicketRequest := &train.GetTicketRequest{
		Email: "user@example.com",
	}
	getTicketResponse, err := client.GetTicket(context.Background(), getTicketRequest)
	if err != nil {
		log.Fatalf("could not get ticket: %v", err)
	}
	fmt.Println("Ticket retrieved:", getTicketResponse.Ticket)

	// Example usage: Get seats by section
	getSeatsBySectionRequest := &train.GetSeatsBySectionRequest{
		Section: "A",
	}
	getSeatsBySectionResponse, err := client.GetSeatsBySection(context.Background(), getSeatsBySectionRequest)
	if err != nil {
		log.Fatalf("could not get seats by section: %v", err)
	}
	fmt.Println("Seats in section A:", getSeatsBySectionResponse.Seats)

}
