package main

import (
	"log"
	"net"
	"redis_test/internal/config"
	"redis_test/internal/hotel"
	"redis_test/pkg/invoke"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

func main() {
	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50055") // Listen on a port for the reservation service
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reservationServer := &hotel.ReservationServer{}
	hotelpb.RegisterReservationServiceServer(s, reservationServer)

	// Establish connections for downstream services
	searchConn, err := grpc.Dial("localhost"+config.SearchPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to search gRPC server: %v", err)
	}
	defer searchConn.Close()
	invoke.RegisterClient("search", hotelpb.NewSearchServiceClient(searchConn))

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

	log.Println("gRPC server listening on port 50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
