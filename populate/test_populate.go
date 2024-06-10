package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

func main() {
	// Replace "localhost:50052" with the appropriate address of your frontend service
	userAddress := "localhost:50053"
	userConn, err := grpc.Dial(userAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to frontend gRPC server: %v", err)
	}
	defer userConn.Close()

	frontendAddress := "localhost:50052"
	frontendConn, err := grpc.Dial(frontendAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to frontend gRPC server: %v", err)
	}
	defer frontendConn.Close()

	// Example calls to the respective service clients
	storeHotel(frontendConn)
	addUsers(userConn)
}

func storeHotel(conn *grpc.ClientConn) {
	client := hotelpb.NewFrontendServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.StoreHotelRequest{
		HotelId:  "H001",
		Name:     "Test Hotel",
		Phone:    "1234567890",
		Location: "Test City",
		Rate:     100,
		Capacity: 50,
		Info:     "This is a test hotel.",
	}

	resp, err := client.StoreHotel(ctx, req)
	if err != nil {
		log.Fatalf("Failed to store hotel: %v", err)
	}

	log.Printf("Stored hotel with ID: %s", resp.HotelId)
}

func addUsers(conn *grpc.ClientConn) {
	client := hotelpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.RegisterUserRequest{
		Username: "user11",
		Password: "11",
	}

	client.RegisterUser(ctx, req)

	token, _ := client.Login(ctx, &hotelpb.LoginRequest{
		Username: "user11",
		Password: "11",
	})

	log.Printf("Token: %s", token)
}
