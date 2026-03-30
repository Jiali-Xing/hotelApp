package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/internal/social"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"

	"github.com/Jiali-Xing/plain"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	socialpb "github.com/Jiali-Xing/socialproto"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50059" // Default port if not specified
	}

	// Set up gRPC server with the appropriate interceptor
	var grpcServer *grpc.Server
	switch config.Intercept {
	case "rajomon":
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
	socialpb.RegisterHomeTimelineServer(grpcServer, &social.HomeTimelineServer{})

	ctx := context.Background()
	// Establish connections for downstream services
	postStorageConn, err := config.CreateGRPCConn(ctx, config.PostStorageAddr)
	if err != nil {
		log.Fatalf("Failed to connect to poststorage gRPC server: %v", err)
	}
	defer postStorageConn.Close()
	invoke.RegisterClient("poststorage", socialpb.NewPostStorageClient(postStorageConn))

	socialGraphConn, err := config.CreateGRPCConn(ctx, config.SocialGraphAddr)
	if err != nil {
		log.Fatalf("Failed to connect to socialgraph gRPC server: %v", err)
	}
	defer socialGraphConn.Close()
	invoke.RegisterClient("socialgraph", socialpb.NewSocialGraphClient(socialGraphConn))

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
