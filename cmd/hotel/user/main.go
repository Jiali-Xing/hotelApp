package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/internal/hotel"
	"github.com/Jiali-Xing/plain"

	hotelpb "github.com/Jiali-Xing/hotelproto"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053" // Default port if not specified
	}
	// Set up gRPC server

	// Set up gRPC server with the appropriate interceptor
	var grpcServer *grpc.Server
	switch config.Intercept {
	case "charon":
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

	userServer := &hotel.UserServer{}
	hotelpb.RegisterUserServiceServer(grpcServer, userServer)

	// Listen and serve
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port " + port)

	// create a couple of users to start with
	go func() {
		ctx := context.Background()
		hotel.RegisterUser(ctx, "user1", "password1")
		hotel.RegisterUser(ctx, "user2", "password2")
		hotel.RegisterUser(ctx, "user3", "password3")

		token := hotel.Login(ctx, "user1", "password2")
		println(token)
		token = hotel.Login(ctx, "user2", "password3")
		println(token)
		token = hotel.Login(ctx, "user1", "password1")
		println(token)
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
