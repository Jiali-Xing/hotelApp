package main

import (
	"log"
	"net"
	"redis_test/internal/hotel"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

func main() {
	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50056") // Listen on a port for the rate service
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	rateServer := &hotel.RateServer{}
	hotelpb.RegisterRateServiceServer(s, rateServer)

	// Establish connections for downstream services if needed
	// For instance, if RateServer needs to communicate with other services

	log.Println("gRPC server listening on port 50056")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
