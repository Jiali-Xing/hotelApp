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
	lis, err := net.Listen("tcp", ":50057") // Listen on a port for the profile service
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	profileServer := &hotel.ProfileServer{}
	hotelpb.RegisterProfileServiceServer(s, profileServer)

	log.Println("gRPC server listening on port 50057")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
