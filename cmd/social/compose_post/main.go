package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"

	"github.com/Jiali-Xing/plain"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/internal/social"
	socialpb "github.com/Jiali-Xing/socialproto"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50062" // Default port if not specified
	}

	// Set up gRPC server with the appropriate interceptor
	var grpcServer *grpc.Server
	switch config.Intercept {
	case "rajomon":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.PriceTable.UnaryInterceptor))
	case "breakwater", "breakwaterd":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Breakwater.UnaryInterceptor))
	case "dagor":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Dg.UnaryInterceptorServer))
	case "plain":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(plain.UnaryInterceptor))
	default:
		grpcServer = grpc.NewServer()
	}

	// Register services
	socialpb.RegisterComposePostServer(grpcServer, &social.ComposePostServer{})

	ctx := context.Background()
	// Establish connections for downstream services
	userTimelineConn, err := config.CreateGRPCConn(ctx, config.UserTimelineAddr)
	if err != nil {
		log.Fatalf("Failed to connect to usertimeline gRPC server: %v", err)
	}
	defer userTimelineConn.Close()
	invoke.RegisterClient("usertimeline", socialpb.NewUserTimelineClient(userTimelineConn))

	homeTimelineConn, err := config.CreateGRPCConn(ctx, config.HomeTimelineAddr)
	if err != nil {
		log.Fatalf("Failed to connect to hometimeline gRPC server: %v", err)
	}
	defer homeTimelineConn.Close()
	invoke.RegisterClient("hometimeline", socialpb.NewHomeTimelineClient(homeTimelineConn))

	postStorageConn, err := config.CreateGRPCConn(ctx, config.PostStorageAddr)
	if err != nil {
		log.Fatalf("Failed to connect to poststorage gRPC server: %v", err)
	}
	defer postStorageConn.Close()
	invoke.RegisterClient("poststorage", socialpb.NewPostStorageClient(postStorageConn))

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
