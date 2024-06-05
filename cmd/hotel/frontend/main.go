package main

import (
	"flag"
	"log"
	"net"
	"redis_test/pkg/invoke"

	"redis_test/internal/hotel"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

var (
	useSeparatePorts bool
)

func init() {
	flag.BoolVar(&useSeparatePorts, "useSeparatePorts", false, "Use separate ports for each downstream service")
	flag.Parse()
}

func main() {
	var userPort, searchPort, reservationPort, ratePort, profilePort string

	if useSeparatePorts {
		// Define different ports for each service
		userPort = ":50053"
		searchPort = ":50054"
		reservationPort = ":50055"
		ratePort = ":50056"
		profilePort = ":50057"
	} else {
		// Use the same port but different URLs (Kubernetes-like environment)
		userPort = ":50051"
		searchPort = ":50051"
		reservationPort = ":50051"
		ratePort = ":50051"
		profilePort = ":50051"
	}

	// Establish gRPC connections
	userConn, err := grpc.Dial("localhost"+userPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user gRPC server: %v", err)
	}
	defer userConn.Close()

	searchConn, err := grpc.Dial("localhost"+searchPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to search gRPC server: %v", err)
	}
	defer searchConn.Close()

	reservationConn, err := grpc.Dial("localhost"+reservationPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to reservation gRPC server: %v", err)
	}
	defer reservationConn.Close()

	rateConn, err := grpc.Dial("localhost"+ratePort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to rate gRPC server: %v", err)
	}
	defer rateConn.Close()

	profileConn, err := grpc.Dial("localhost"+profilePort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to profile gRPC server: %v", err)
	}
	defer profileConn.Close()

	// Register gRPC clients
	invoke.RegisterClient("user", hotelpb.NewUserServiceClient(userConn))
	invoke.RegisterClient("search", hotelpb.NewSearchServiceClient(searchConn))
	invoke.RegisterClient("reservation", hotelpb.NewReservationServiceClient(reservationConn))
	invoke.RegisterClient("rate", hotelpb.NewRateServiceClient(rateConn))
	invoke.RegisterClient("profile", hotelpb.NewProfileServiceClient(profileConn))

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50052") // Listen on a different port for the frontend server
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	hotelServer := &hotel.FrontendServer{}
	hotelpb.RegisterReservationServiceServer(s, hotelServer)

	log.Println("gRPC server listening on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
