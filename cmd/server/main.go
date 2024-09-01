package main

import (
	"log"
	"net"

	train "github.com/bijoyv/train/pkg/proto"
	reservation "github.com/bijoyv/train/pkg/train"
	"google.golang.org/grpc"
)

func main() {
	// Create a TrainService instance
	trainService := reservation.NewTrainReservationService()

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
