package main

import (
	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/internal/hotel"
	"github.com/Jiali-Xing/plain"
	"log"
	"net"
	"os"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50055" // Default port if not specified
	}

	// Set up gRPC server with the appropriate interceptor
	var grpcServer *grpc.Server
	switch config.Intercept {
	case "rajomon":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.PriceTable.UnaryInterceptor))
	case "breakwaterd":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Breakwater.UnaryInterceptor))
	case "dagor":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Dg.UnaryInterceptorServer))
	case "plain":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(plain.UnaryInterceptor))
	default:
		grpcServer = grpc.NewServer()
	}

	reservationServer := &hotel.ReservationServer{}
	hotelpb.RegisterReservationServiceServer(grpcServer, reservationServer)

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
