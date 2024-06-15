package main

import (
	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/plain"
	socialpb "github.com/Jiali-Xing/socialproto"
	"google.golang.org/grpc"
)

func main() {
	type socialGraphServer struct {
		socialpb.UnimplementedSocialGraphServer
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051" // Default port if not specified
	}

	// Set up gRPC server with the appropriate interceptor
	var grpcServer *grpc.Server
	switch config.Intercept {
	case "charon":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.PriceTable.UnaryInterceptor))
	case "breakwaterd", "breakwater":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Breakwater.UnaryInterceptor))
	case "dagor":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Dg.UnaryInterceptorServer))
	case "plain":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(plain.UnaryInterceptor))
	default:
		grpcServer = grpc.NewServer()
	}

	// Register services
	socialpb.RegisterSocialGraphServer(grpcServer, &socialGraphServer{})

	// Listen and serve
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
