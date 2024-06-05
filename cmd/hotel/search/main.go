package main

import (
	"fmt"
	"log"
	"net"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"redis_test/internal/hotel"

	pb "github.com/Jiali-Xing/hotelproto"
)

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))

	// Create a gRPC server
	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &hotel.SearchService{})
	reflection.Register(server)

	// Listen on a port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("gRPC server listening on port 50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
