package reservation

import (
	"context"
	"fmt"

	train "github.com/bijoyv/train/pkg/proto"
)

// TrainService implements the grpc interface using CSP.
type TrainService struct {
	ticketsChan chan map[string]*train.Ticket // Channel for ticket operations
	seatsChan   chan map[string]string        // Channel for seat operations
	train.UnimplementedTrainServiceServer
}

// helper function to assign seat
func (s *TrainService) assignSeat() string {
	seats := <-s.seatsChan
	defer func() { s.seatsChan <- seats }()

	for seat, taken := range seats {
		if taken == "" {
			seats[seat] = "taken" //mark the seat as occupied
			return seat
		}
	}
	return ""
}

// Implement grpc service methods
func (s *TrainService) PurchaseTicket(ctx context.Context, req *train.PurchaseTicketRequest) (*train.PurchaseTicketResponse, error) {
	if req.User == nil || req.User.Email == "" {
		return nil, fmt.Errorf("Invalid User Information")
	}

	seat := s.assignSeat()
	if seat == "" {
		return nil, fmt.Errorf("no seats available")
	}

	//create ticket
	ticket := &train.Ticket{
		From:  req.From,
		To:    req.To,
		User:  req.User,
		Price: 20,
		Seat:  seat,
	}

	//Communicate go routines to update tickets and seats
	tickets := <-s.ticketsChan
	seats := <-s.seatsChan

	//Check if ticket already exists
	if _, exists := tickets[req.User.Email]; exists {
		s.seatsChan <- seats
		s.ticketsChan <- tickets
		return nil, fmt.Errorf("A ticket already exists for this user")
	}
	tickets[req.User.Email] = ticket
	seats[seat] = req.User.Email

	s.seatsChan <- seats
	s.ticketsChan <- tickets

	return &train.PurchaseTicketResponse{Ticket: ticket}, nil
}
func (s *TrainService) GetTicket(ctx context.Context, req *train.GetTicketRequest) (*train.GetTicketResponse, error) {
	tickets := <-s.ticketsChan
	defer func() {
		s.ticketsChan <- tickets
	}()
	ticket, exists := tickets[req.Email]

	if !exists {
		return nil, fmt.Errorf("Ticket not found for user %s", req.Email)
	}
	return &train.GetTicketResponse{Ticket: ticket}, nil

}
func (s *TrainService) GetSeatsBySection(ctx context.Context, req *train.GetSeatsBySectionRequest) (*train.GetSeatsBySectionResponse, error) {
	seats := <-s.seatsChan
	defer func() { s.seatsChan <- seats }()

	result := make(map[string]string)

	for seat, email := range seats {
		if string(seat[0]) == req.Section {
			result[seat] = email
		}
	}
	return &train.GetSeatsBySectionResponse{Seats: result}, nil
}

func (s *TrainService) RemoveUser(ctx context.Context, req *train.RemoveUserRequest) (*train.RemoveUserResponse, error) {
	tickets := <-s.ticketsChan
	seats := <-s.seatsChan

	ticket, exists := tickets[req.Email]
	if !exists {
		return nil, fmt.Errorf("User %s not found with ticket", req.Email)
	}

	//remove seat
	delete(seats, ticket.Seat)

	//remove ticket
	delete(tickets, req.Email)

	s.seatsChan <- seats
	s.ticketsChan <- tickets

	return &train.RemoveUserResponse{Success: true}, nil
}

func (s *TrainService) ModifySeat(ctx context.Context, req *train.ModifySeatRequest) (*train.ModifySeatResponse, error) {

	tickets := <-s.ticketsChan
	seats := <-s.seatsChan

	ticket, exists := tickets[req.Email]

	if !exists {
		s.ticketsChan <- tickets
		s.seatsChan <- seats
		return nil, fmt.Errorf("user %s not found with ticket", req.Email)
	}
	//check if new seat is available
	if seat, taken := seats[req.NewSeat]; taken && seat != "" {

		s.ticketsChan <- tickets
		s.seatsChan <- seats
		return nil, fmt.Errorf("seat %s already in use", seat)
	}

	//update seat details
	delete(seats, ticket.Seat)
	seats[req.NewSeat] = req.Email

	//update ticket
	ticket.Seat = req.NewSeat

	s.ticketsChan <- tickets
	s.seatsChan <- seats
	return &train.ModifySeatResponse{Success: true}, nil

}

// initialize the service
func NewTrainReservationService() *TrainService {
	tickets := make(map[string]*train.Ticket)
	seats := make(map[string]string)

	//initialize for twenty seats now for each section
	for i := 1; i <= 20; i++ {

		seat := fmt.Sprintf("A%d", i)
		seats[seat] = ""
		seat = fmt.Sprintf("B%d", i)
		seats[seat] = ""
	}

	ts := &TrainService{
		ticketsChan: make(chan map[string]*train.Ticket, 1),
		seatsChan:   make(chan map[string]string, 1),
	}
	ts.ticketsChan <- tickets
	ts.seatsChan <- seats

	return ts
}
