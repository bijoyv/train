package reservation

import (
	"context"
	"testing"

	train "github.com/bijoyv/train/pkg/proto"
)

func TestTrainService(t *testing.T) {
	trainService := NewTrainReservationService()
	t.Run("PurchaseTicket", func(t *testing.T) {
		user := &train.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		res, err := trainService.PurchaseTicket(context.Background(), req)
		if err != nil {
			t.Fatalf("PurchaseTicket failed: %v", err)
		}

		if res.Ticket == nil {
			t.Fatal("Expected ticket in response, got nil")
		}

		if res.Ticket.User.Email != "john.doe@example.com" {
			t.Errorf("Expected user email 'john.doe@example.com', got %s", res.Ticket.User.Email)
		}

		if res.Ticket.Seat == "" {
			t.Error("Expected a seat to be assigned, got empty string")
		}
	})

	//Add testcase for user already exist
	t.Run("PurchaseTicketUserExists", func(t *testing.T) {
		user := &train.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		_, err := trainService.PurchaseTicket(context.Background(), req)
		if err == nil {
			t.Fatal("Expected error for existing user, got nil")
		}
	})

	t.Run("GetTicket", func(t *testing.T) {
		user := &train.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		_, err := trainService.PurchaseTicket(context.Background(), req)
		if err != nil {
			t.Fatalf("PurchaseTicket failed: %v", err)
		}

		getTicketReq := &train.GetTicketRequest{
			Email: "jane.doe@example.com",
		}

		res, err := trainService.GetTicket(context.Background(), getTicketReq)
		if err != nil {
			t.Fatalf("GetTicket failed: %v", err)
		}

		if res.Ticket == nil {
			t.Fatal("Expected ticket in response, got nil")
		}

		if res.Ticket.User.Email != "jane.doe@example.com" {
			t.Errorf("Expected user email 'jane.doe@example.com', got %s", res.Ticket.User.Email)
		}
	})

	t.Run("GetTicketNotFound", func(t *testing.T) {
		getTicketReq := &train.GetTicketRequest{
			Email: "nonexistent@example.com",
		}

		_, err := trainService.GetTicket(context.Background(), getTicketReq)
		if err == nil {
			t.Fatal("Expected error got get ticket, got error nil")
		}
	})

	t.Run("GetSeatsBySection", func(t *testing.T) {
		user := &train.User{
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice.smith@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		_, err := trainService.PurchaseTicket(context.Background(), req)
		if err != nil {
			t.Fatalf("PurchaseTicket failed: %v", err)
		}

		getSeatsReq := &train.GetSeatsBySectionRequest{
			Section: "A",
		}

		res, err := trainService.GetSeatsBySection(context.Background(), getSeatsReq)
		if err != nil {
			t.Fatalf("GetSeatsBySection failed: %v", err)
		}

		if len(res.Seats) == 0 {
			t.Error("Expected at least one seat in section A, got empty map")
		}
	})

	t.Run("RemoveUser", func(t *testing.T) {
		user := &train.User{
			FirstName: "Bob",
			LastName:  "Jones",
			Email:     "bob.jones@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		_, err := trainService.PurchaseTicket(context.Background(), req)
		if err != nil {
			t.Fatalf("PurchaseTicket failed: %v", err)
		}

		removeUserReq := &train.RemoveUserRequest{
			Email: "bob.jones@example.com",
		}

		res, err := trainService.RemoveUser(context.Background(), removeUserReq)
		if err != nil {
			t.Fatalf("RemoveUser failed: %v", err)
		}

		if !res.Success {
			t.Error("Expected RemoveUser to be successful, got false")
		}

		// Check if the user is removed
		_, err = trainService.GetTicket(context.Background(), &train.GetTicketRequest{
			Email: "bob.jones@example.com",
		})
		if err == nil {
			t.Error("Expected user to be removed, got ticket")
		}
	})

	t.Run("ModifySeat", func(t *testing.T) {
		user := &train.User{
			FirstName: "Carol",
			LastName:  "Williams",
			Email:     "carol.williams@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		_, err := trainService.PurchaseTicket(context.Background(), req)
		if err != nil {
			t.Fatalf("PurchaseTicket failed: %v", err)
		}

		modifySeatReq := &train.ModifySeatRequest{
			Email:   "carol.williams@example.com",
			NewSeat: "A2",
		}

		res, err := trainService.ModifySeat(context.Background(), modifySeatReq)
		if err != nil {
			t.Fatalf("ModifySeat failed: %v", err)
		}

		if !res.Success {
			t.Error("Expected ModifySeat to be successful, got false")
		}

		// Check if the seat is modified
		getTicketReq := &train.GetTicketRequest{
			Email: "carol.williams@example.com",
		}
		getTicketRes, err := trainService.GetTicket(context.Background(), getTicketReq)
		if err != nil {
			t.Fatalf("GetTicket failed: %v", err)
		}
		if getTicketRes.Ticket.Seat != "A2" {
			t.Errorf("Expected seat to be %s, got %s", "A2", getTicketRes.Ticket.Seat)
		}
	})

	t.Run("ModifySeatTaken", func(t *testing.T) {
		user := &train.User{
			FirstName: "David",
			LastName:  "Brown",
			Email:     "david.brown@example.com",
		}
		req := &train.PurchaseTicketRequest{
			From: "London",
			To:   "Paris",
			User: user,
		}

		_, err := trainService.PurchaseTicket(context.Background(), req)
		if err != nil {
			t.Fatalf("PurchaseTicket failed: %v", err)
		}
	})

}
