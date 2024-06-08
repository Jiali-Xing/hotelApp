package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	// Prompt for the address of the gRPC server
	fmt.Print("Enter the address of the gRPC server (e.g., localhost:50051): ")
	address, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	address = strings.TrimSpace(address)

	// Establish a connection to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create gRPC clients
	userClient := hotelpb.NewUserServiceClient(conn)
	searchClient := hotelpb.NewSearchServiceClient(conn)
	reservationClient := hotelpb.NewReservationServiceClient(conn)
	rateClient := hotelpb.NewRateServiceClient(conn)
	profileClient := hotelpb.NewProfileServiceClient(conn)

	// Test Login
	loginReq := &hotelpb.LoginRequest{Username: "testuser", Password: "testpassword"}
	loginResp, err := userClient.Login(context.Background(), loginReq)
	if err != nil {
		log.Printf("Login failed: %v", err)
	} else {
		log.Printf("Login response: %v", loginResp)
	}

	// Test Nearby
	nearbyReq := &hotelpb.NearbyRequest{InDate: "2024-01-01", OutDate: "2024-01-02", Location: "testlocation"}
	nearbyResp, err := searchClient.Nearby(context.Background(), nearbyReq)
	if err != nil {
		log.Printf("Nearby failed: %v", err)
	} else {
		log.Printf("Nearby response: %v", nearbyResp)
	}

	// Test MakeReservation
	makeReservationReq := &hotelpb.MakeReservationRequest{
		CustomerName: "testuser",
		HotelId:      "hotel123",
		InDate:       "2024-01-01",
		OutDate:      "2024-01-02",
		RoomNumber:   1,
	}
	makeReservationResp, err := reservationClient.MakeReservation(context.Background(), makeReservationReq)
	if err != nil {
		log.Printf("MakeReservation failed: %v", err)
	} else {
		log.Printf("MakeReservation response: %v", makeReservationResp)
	}

	// Test GetRates
	getRatesReq := &hotelpb.GetRatesRequest{HotelIds: []string{"hotel123"}}
	getRatesResp, err := rateClient.GetRates(context.Background(), getRatesReq)
	if err != nil {
		log.Printf("GetRates failed: %v", err)
	} else {
		log.Printf("GetRates response: %v", getRatesResp)
	}

	// Test GetProfiles
	getProfilesReq := &hotelpb.GetProfilesRequest{HotelIds: []string{"hotel123"}}
	getProfilesResp, err := profileClient.GetProfiles(context.Background(), getProfilesReq)
	if err != nil {
		log.Printf("GetProfiles failed: %v", err)
	} else {
		log.Printf("GetProfiles response: %v", getProfilesResp)
	}
}
