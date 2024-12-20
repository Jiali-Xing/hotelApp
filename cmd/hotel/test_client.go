package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

var (
	skipTest bool
)

func main() {
	flag.Parse()

	// Establish gRPC connections to all services
	frontendConn, err := grpc.Dial(config.FrontendAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to frontend gRPC server: %v", err)
	}
	defer frontendConn.Close()
	frontendClient := hotelpb.NewFrontendServiceClient(frontendConn)

	userConn, err := grpc.Dial(config.UserAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user gRPC server: %v", err)
	}
	defer userConn.Close()
	userClient := hotelpb.NewUserServiceClient(userConn)

	searchConn, err := grpc.Dial(config.SearchAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to search gRPC server: %v", err)
	}
	defer searchConn.Close()
	searchClient := hotelpb.NewSearchServiceClient(searchConn)

	reservationConn, err := grpc.Dial(config.ReservationAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to reservation gRPC server: %v", err)
	}
	defer reservationConn.Close()
	reservationClient := hotelpb.NewReservationServiceClient(reservationConn)

	rateConn, err := grpc.Dial(config.RateAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to rate gRPC server: %v", err)
	}
	defer rateConn.Close()
	rateClient := hotelpb.NewRateServiceClient(rateConn)

	profileConn, err := grpc.Dial(config.ProfileAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to profile gRPC server: %v", err)
	}
	defer profileConn.Close()
	profileClient := hotelpb.NewProfileServiceClient(profileConn)

	// Run tests
	skipTest = true

	if config.RunLocally {

		// Debugging the store hotel request

		// Test Frontend Service
		testSearchHotels(frontendClient)
		testMakeReservation(frontendClient)

		if !skipTest {
			// Test User Service
			testRegisterUser(userClient)
			testLogin(userClient)

			// Test Search Service
			testStoreHotelLocation(searchClient)

			// Test Reservation Service
			testCheckAvailability(reservationClient)

			// Test Rate Service
			testGetRates(rateClient)

			// Test Profile Service
			testGetProfiles(profileClient)
		}
	} else {
		//	only test search and reservation of the frontend service
		testSearchHotels(frontendClient)
		testMakeReservation(frontendClient)
	}
}

func testRegisterUser(client hotelpb.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.RegisterUserRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	resp, err := client.RegisterUser(ctx, req)
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}
	log.Printf("RegisterUser response: %v", resp)
}

func testLogin(client hotelpb.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	log.Printf("Login response: %v", resp)
}

func testSearchHotels(client hotelpb.FrontendServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.SearchHotelsRequest{
		InDate:   "2023-04-17",
		OutDate:  "2023-04-19",
		Location: "new-york-city-ny-0",
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "tokens", "99999", "request-id", "12345", "method", "search-hotel", "u", "1", "b", "1", "timestamp", "12345", "name", "frontend", "method", "search-hotel")
	resp, err := client.SearchHotels(ctx, req)
	if err != nil {
		log.Fatalf("Failed to search hotels: %v", err)
	}
	log.Printf("SearchHotels response: %v", resp)
}

func testStoreHotelLocation(client hotelpb.SearchServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.StoreHotelLocationRequest{
		HotelId:  "H001",
		Location: "Test City",
	}

	resp, err := client.StoreHotelLocation(ctx, req)
	if err != nil {
		log.Fatalf("Failed to store hotel location: %v", err)
	}
	log.Printf("StoreHotelLocation response: %v", resp)
}

func testCheckAvailability(client hotelpb.ReservationServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.CheckAvailabilityRequest{
		CustomerName: "user1",
		HotelIds:     []string{"1"},
		InDate:       "2024-06-01",
		OutDate:      "2024-06-10",
		RoomNumber:   1,
	}

	resp, err := client.CheckAvailability(ctx, req)
	if err != nil {
		log.Fatalf("Failed to check availability: %v", err)
	}
	log.Printf("CheckAvailability response: %v", resp)
}

func testMakeReservation(client hotelpb.FrontendServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.FrontendReservationRequest{
		HotelId:  "1",
		InDate:   "2023-04-17",
		OutDate:  "2023-04-19",
		Rooms:    1,
		Username: "user1",
		Password: "password1",
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "tokens", "99999", "request-id", "12345", "method", "search-hotel", "u", "1", "b", "1", "timestamp", "12345", "name", "frontend", "method", "search-hotel", "id", "1", "demand", "3")
	resp, err := client.FrontendReservation(ctx, req)
	if err != nil {
		log.Fatalf("Failed to make reservation: %v", err)
	}
	log.Printf("MakeReservation response: %v", resp)
}

func testGetRates(client hotelpb.RateServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.GetRatesRequest{
		HotelIds: []string{"H001"},
	}

	resp, err := client.GetRates(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get rates: %v", err)
	}
	log.Printf("GetRates response: %v", resp)
}

func testGetProfiles(client hotelpb.ProfileServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &hotelpb.GetProfilesRequest{
		HotelIds: []string{"H001"},
	}

	resp, err := client.GetProfiles(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}
	log.Printf("GetProfiles response: %v", resp)
}
