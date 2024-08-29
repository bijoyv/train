package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	train "github.com/bijoyv/train/proto"
	"google.golang.org/grpc"
)

// Define the gRPC service and message structures
type TrainService struct {
	mu      sync.Mutex
	Tickets map[string]*train.Ticket // Store tickets by user email
	Seats   map[string]string        // Store seat assignments: "A1": "user@example.com"
	train.UnimplementedTrainServiceServer
}

// Implement the gRPC service methods
func (s *TrainService) PurchaseTicket(ctx context.Context, req *train.PurchaseTicketRequest) (*train.PurchaseTicketResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Basic validation
	if req.User == nil || req.User.Email == "" {
		return nil, fmt.Errorf("invalid user information")
	}

	// Check if a ticket already exists for this user
	if _, exists := s.Tickets[req.User.Email]; exists {
		return nil, fmt.Errorf("a ticket already exists for this user")
	}

	// Assign a seat (simple logic for demo)
	seat := s.assignSeat()
	if seat == "" {
		return nil, fmt.Errorf("no seats available")
	}

	// Create the ticket
	ticket := &train.Ticket{
		From:  req.From,
		To:    req.To,
		User:  req.User,
		Price: 20, // Fixed price for now
		Seat:  seat,
	}

	// Store the ticket and seat assignment
	s.Tickets[req.User.Email] = ticket
	s.Seats[seat] = req.User.Email

	return &train.PurchaseTicketResponse{Ticket: ticket}, nil
}

func (s *TrainService) GetTicket(ctx context.Context, req *train.GetTicketRequest) (*train.GetTicketResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.Tickets[req.Email]
	if !exists {
		return nil, fmt.Errorf("ticket not found for user: %s", req.Email)
	}
	return &train.GetTicketResponse{Ticket: ticket}, nil
}

func (s *TrainService) GetSeatsBySection(ctx context.Context, req *train.GetSeatsBySectionRequest) (*train.GetSeatsBySectionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seats := make(map[string]string)
	for seat, email := range s.Seats {
		if string(seat[0]) == req.Section {
			seats[seat] = email
		}
	}

	return &train.GetSeatsBySectionResponse{Seats: seats}, nil
}

func (s *TrainService) RemoveUser(ctx context.Context, req *train.RemoveUserRequest) (*train.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.Tickets[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", req.Email)
	}

	// Remove seat assignment
	delete(s.Seats, ticket.Seat)

	// Remove ticket
	delete(s.Tickets, req.Email)

	return &train.RemoveUserResponse{Success: true}, nil
}

func (s *TrainService) ModifySeat(ctx context.Context, req *train.ModifySeatRequest) (*train.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.Tickets[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", req.Email)
	}

	// Check if the new seat is available
	if _, taken := s.Seats[req.NewSeat]; taken {
		return nil, fmt.Errorf("seat %s is already taken", req.NewSeat)
	}

	// Update seat assignment
	delete(s.Seats, ticket.Seat)
	s.Seats[req.NewSeat] = req.Email

	// Update ticket
	ticket.Seat = req.NewSeat

	return &train.ModifySeatResponse{Success: true}, nil
}

// Helper function for seat assignment (replace with your logic)
func (s *TrainService) assignSeat() string {
	for seat, taken := range s.Seats {
		if taken != "taken" {
			s.Seats[seat] = "taken" // Mark the seat as taken
			return seat
		}
	}
	return ""
}

// Initialize the Tickets and Seats
func initialize(trainService *TrainService) {
	//Initialize for twenty seats
	for i := 1; i <= 20; i++ {
		seat := fmt.Sprintf("A%d", i)
		trainService.Seats[seat] = ""
	}

}

func main() {
	// Create a TrainService instance
	trainService := &TrainService{
		Tickets: make(map[string]*train.Ticket),
		Seats:   make(map[string]string),
	}

	initialize(trainService)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":50051") // Choose your port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	//Register the TrainService
	train.RegisterTrainServiceServer(grpcServer, trainService)

	log.Println("Server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
