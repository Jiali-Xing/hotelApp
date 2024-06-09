package main

import (
	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/internal/hotel"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50056" // Default port if not specified
	}
	// Set up gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	rateServer := &hotel.RateServer{}
	hotelpb.RegisterRateServiceServer(s, rateServer)

	// Establish connections for downstream services if needed
	// For instance, if RateServer needs to communicate with other services

	log.Println("gRPC server listening on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
