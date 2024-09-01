package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	train "github.com/bijoyv/train/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientCommands defines the available commands for the client
type ClientCommands struct {
	Command string
	From    string
	To      string
	Email   string
	Section string
	NewSeat string
}

func main() {
	// Define command-line flags
	cmd := flag.String("cmd", "", "Command to execute: purchase, getticket, getseats, removeuser, modifyseat")
	from := flag.String("from", "", "Origin station (required for purchase)")
	to := flag.String("to", "", "Destination station (required for purchase)")
	email := flag.String("email", "", "User email (required for purchase, getticket, removeuser, modifyseat)")
	section := flag.String("section", "", "Seat section (required for getseats)")
	newSeat := flag.String("newseat", "", "New seat (required for modifyseat)")

	flag.Parse()

	clientCommands := ClientCommands{
		Command: *cmd,
		From:    *from,
		To:      *to,
		Email:   *email,
		Section: *section,
		NewSeat: *newSeat,
	}

	// Validate input
	if err := validateInput(clientCommands); err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// Correct way to create a gRPC connection:
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC connection: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := train.NewTrainServiceClient(conn)

	// Execute the command
	switch clientCommands.Command {
	case "purchase":
		executePurchase(client, clientCommands.From, clientCommands.To, clientCommands.Email)
	case "getticket":
		executeGetTicket(client, clientCommands.Email)
	case "getseats":
		executeGetSeats(client, clientCommands.Section)
	case "removeuser":
		executeRemoveUser(client, clientCommands.Email)
	case "modifyseat":
		executeModifySeat(client, clientCommands.Email, clientCommands.NewSeat)
	default:
		flag.Usage()
		os.Exit(1)
	}

}

// validateInput validates the command-line input
func validateInput(cmd ClientCommands) error {
	switch cmd.Command {
	case "purchase":
		if cmd.From == "" || cmd.To == "" || cmd.Email == "" {
			return fmt.Errorf("purchase requires --from, --to, and --email")
		}
	case "getticket", "removeuser", "modifyseat":
		if cmd.Email == "" {
			return fmt.Errorf("%s requires --email", cmd.Command)
		}
	case "getseats":
		if cmd.Section == "" {
			return fmt.Errorf("getseats requires --section")
		}
	default:
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
	return nil
}

// executePurchase handles the purchase command
func executePurchase(client train.TrainServiceClient, from, to, email string) {
	user := &train.User{
		Email: email,
	}
	purchaseRequest := &train.PurchaseTicketRequest{
		From: from,
		To:   to,
		User: user,
	}
	purchaseResponse, err := client.PurchaseTicket(context.Background(), purchaseRequest)
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	fmt.Println("Ticket purchased:", purchaseResponse.Ticket)
}

// executeGetTicket handles the getticket command
func executeGetTicket(client train.TrainServiceClient, email string) {
	getTicketRequest := &train.GetTicketRequest{
		Email: email,
	}
	getTicketResponse, err := client.GetTicket(context.Background(), getTicketRequest)
	if err != nil {
		log.Fatalf("could not get ticket: %v", err)
	}
	fmt.Println("Ticket retrieved:", getTicketResponse.Ticket)
}

// executeGetSeats handles the getseats command
func executeGetSeats(client train.TrainServiceClient, section string) {
	getSeatsBySectionRequest := &train.GetSeatsBySectionRequest{
		Section: section,
	}
	getSeatsBySectionResponse, err := client.GetSeatsBySection(context.Background(), getSeatsBySectionRequest)
	if err != nil {
		log.Fatalf("could not get seats by section: %v", err)
	}
	fmt.Println("Seats in section", section, ":", getSeatsBySectionResponse.Seats)
}

// executeRemoveUser handles the removeuser command
func executeRemoveUser(client train.TrainServiceClient, email string) {
	removeUserRequest := &train.RemoveUserRequest{
		Email: email,
	}
	removeUserResponse, err := client.RemoveUser(context.Background(), removeUserRequest)
	if err != nil {
		log.Fatalf("could not remove user: %v", err)
	}
	fmt.Println("User removed successfully:", removeUserResponse.Success)
}

// executeModifySeat handles the modifyseat command
func executeModifySeat(client train.TrainServiceClient, email, newSeat string) {
	modifySeatRequest := &train.ModifySeatRequest{
		Email:   email,
		NewSeat: newSeat,
	}
	modifySeatResponse, err := client.ModifySeat(context.Background(), modifySeatRequest)
	if err != nil {
		log.Fatalf("could not modify seat: %v", err)
	}
	fmt.Println("Seat modified successfully:", modifySeatResponse.Success)
}
