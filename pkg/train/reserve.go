package reservation

import (
	"context"
	"fmt"

	train "github.com/bijoyv/train/pkg/proto"
)

// TrainService implements the grpc interface using CSP.
type TrainService struct {
	ops chan func(map[string]*train.Ticket, map[string]string)
	train.UnimplementedTrainServiceServer
}

// This is for running a go routine to make the data local for synchronization
func (s *TrainService) Run() {
	tickets := make(map[string]*train.Ticket)
	seats := initializeSeats()

	for op := range s.ops {
		op(tickets, seats)
	}
}

// helper function to initialize seats for section A and B
func initializeSeats() map[string]string {
	seats := make(map[string]string)
	for i := 1; i <= 20; i++ {
		seats[fmt.Sprintf("A%d", i)] = ""
		seats[fmt.Sprintf("B%d", i)] = ""
	}
	return seats
}

// helper function to assign seat
func assignSeat(seats map[string]string) string {

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
	result := make(chan error, 1)
	resTicket := make(chan *train.Ticket, 1)
	s.ops <- func(tickets map[string]*train.Ticket, seats map[string]string) {
		seat := assignSeat(seats)
		if seat == "" {
			result <- fmt.Errorf("no seats available")
			return
		}
		if _, exists := tickets[req.User.Email]; exists {
			seats[seat] = "" //reset the taken as we are not purchasing
			result <- fmt.Errorf("ticket already exist for this user")
			return
		}

		ticket := &train.Ticket{
			From:  req.From,
			To:    req.To,
			User:  req.User,
			Price: 20,
			Seat:  seat,
		}
		tickets[req.User.Email] = ticket
		seats[seat] = req.User.Email
		resTicket <- ticket
	}
	select {
	case er := <-result:
		return nil, er
	case tkt := <-resTicket:

		return &train.PurchaseTicketResponse{Ticket: tkt}, nil
	}
}
func (s *TrainService) GetTicket(ctx context.Context, req *train.GetTicketRequest) (*train.GetTicketResponse, error) {
	er := make(chan error, 1)
	resTicket := make(chan *train.Ticket, 1)
	s.ops <- func(tickets map[string]*train.Ticket, seats map[string]string) {
		ticket, exists := tickets[req.Email]
		if !exists {
			er <- fmt.Errorf("Ticket not found for user %s", req.Email)
			return
		}
		resTicket <- ticket

	}

	select {
	case e := <-er:
		return nil, e
	case t := <-resTicket:

		return &train.GetTicketResponse{Ticket: t}, nil
	}
}
func (s *TrainService) GetSeatsBySection(ctx context.Context, req *train.GetSeatsBySectionRequest) (*train.GetSeatsBySectionResponse, error) {

	seatc := make(chan map[string]string, 1)
	s.ops <- func(tickets map[string]*train.Ticket, seats map[string]string) {
		result := make(map[string]string)
		for seat, email := range seats {
			if string(seat[0]) == req.Section {
				result[seat] = email
			}

		}
		seatc <- result
	}

	result := <-seatc
	return &train.GetSeatsBySectionResponse{Seats: result}, nil
}

func (s *TrainService) RemoveUser(ctx context.Context, req *train.RemoveUserRequest) (*train.RemoveUserResponse, error) {
	er := make(chan error, 1)
	result := make(chan bool, 1)

	s.ops <- func(tickets map[string]*train.Ticket, seats map[string]string) {

		ticket, exists := tickets[req.Email]
		if !exists {
			er <- fmt.Errorf("User %s not found with ticket", req.Email)
			return
		}
		//		delete(seats, ticket.Seat)
		seats[ticket.Seat] = ""
		delete(tickets, req.Email)

		result <- true
	}
	select {
	case e := <-er:
		return nil, e
	case <-result:

		return &train.RemoveUserResponse{Success: true}, nil
	}
}

func (s *TrainService) ModifySeat(ctx context.Context, req *train.ModifySeatRequest) (*train.ModifySeatResponse, error) {
	er := make(chan error, 1)
	result := make(chan bool, 1)

	s.ops <- func(tickets map[string]*train.Ticket, seats map[string]string) {

		ticket, exists := tickets[req.Email]
		if !exists {
			er <- fmt.Errorf("user %s not found with ticket", req.Email)
			return
		}
		if seat, taken := seats[req.NewSeat]; taken && seat != "" {
			er <- fmt.Errorf("seat %s already in use", seat)
			return
		}
		delete(seats, ticket.Seat)
		seats[req.NewSeat] = req.Email
		ticket.Seat = req.NewSeat
		result <- true
	}
	select {
	case e := <-er:
		return nil, e
	case <-result:
		return &train.ModifySeatResponse{Success: true}, nil
	}

}

// initialize the service
func NewTrainReservationService() *TrainService {
	ts := &TrainService{
		ops: make(chan func(map[string]*train.Ticket, map[string]string)),
	}
	go ts.Run()
	return ts
}
