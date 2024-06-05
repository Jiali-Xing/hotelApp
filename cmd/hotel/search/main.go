package main

import (
	"log"
	"net"
	"os"
	"redis_test/internal/config"
	"redis_test/internal/hotel"
	"redis_test/pkg/invoke"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50054" // Default port if not specified
	}

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	searchServer := &hotel.SearchServer{}
	hotelpb.RegisterSearchServiceServer(s, searchServer)

	// Establish connections for downstream services
	rateConn, err := grpc.Dial("localhost"+config.RatePort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to rate gRPC server: %v", err)
	}
	defer rateConn.Close()
	invoke.RegisterClient("rate", hotelpb.NewRateServiceClient(rateConn))

	profileConn, err := grpc.Dial("localhost"+config.ProfilePort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to profile gRPC server: %v", err)
	}
	defer profileConn.Close()
	invoke.RegisterClient("profile", hotelpb.NewProfileServiceClient(profileConn))

	log.Println("gRPC server listening on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
