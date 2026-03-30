package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/internal/hotel"
	"github.com/Jiali-Xing/plain"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50054" // Default port if not specified
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

	searchServer := &hotel.SearchServer{}
	hotelpb.RegisterSearchServiceServer(grpcServer, searchServer)

	// Listen and serve
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ctx := context.Background()
	// Establish connections for downstream services
	rateConn, err := config.CreateGRPCConn(ctx, config.RateAddr)
	if err != nil {
		log.Fatalf("Failed to connect to rate gRPC server: %v", err)
	}
	defer rateConn.Close()
	invoke.RegisterClient("rate", hotelpb.NewRateServiceClient(rateConn))

	log.Println("gRPC server listening on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
