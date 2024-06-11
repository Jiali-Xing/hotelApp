package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

var (
	frontendAddr string
	userAddr     string
	hotelsFile   string
	numOfUsers   int
	infoSize     int
)

func init() {
	flag.StringVar(&frontendAddr, "frontend", "localhost:50052", "Address of the frontend service")
	flag.StringVar(&userAddr, "user", "localhost:50053", "Address of the user service")
	flag.StringVar(&hotelsFile, "hotels_file", "../experiments/hotel/data/hotels.json", "Path to the hotels file")
	flag.IntVar(&numOfUsers, "num_of_users", 100, "Number of users to create")
	flag.IntVar(&infoSize, "info_size", 1000, "Size of hotel info in bytes")
	flag.Parse()
}

func getRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func main() {
	test_adduser_hotel()
}

func test_adduser_hotel() {
	// Connect to the user and frontend services
	userConn, err := grpc.Dial(userAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user gRPC server: %v", err)
	}
	defer userConn.Close()

	frontendConn, err := grpc.Dial(frontendAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to frontend gRPC server: %v", err)
	}
	defer frontendConn.Close()

	// Populate the services with data
	populateHotels(frontendConn)
	populateUsers(userConn)
}

func populateHotels(conn *grpc.ClientConn) {
	client := hotelpb.NewFrontendServiceClient(conn)

	// Read hotels data from the JSON file
	data, err := ioutil.ReadFile(hotelsFile)
	if err != nil {
		log.Fatalf("Failed to read hotels file: %v", err)
	}

	var hotels []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		Address     struct {
			City string `json:"city"`
		} `json:"address"`
	}
	if err := json.Unmarshal(data, &hotels); err != nil {
		log.Fatalf("Failed to parse hotels file: %v", err)
	}

	// Add only 2 hotels to the service
	for _, hotel := range hotels[:1000] {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		req := &hotelpb.StoreHotelRequest{
			HotelId:  hotel.ID,
			Name:     hotel.Name,
			Phone:    hotel.PhoneNumber,
			Location: hotel.Address.City,
			Rate:     100,
			Capacity: 50,
			Info:     getRandomString(infoSize),
		}

		resp, err := client.StoreHotel(ctx, req)
		if err != nil {
			log.Printf("Failed to store hotel: %v", err)
		} else {
			log.Printf("Stored hotel with ID: %s", resp.HotelId)
		}
	}
}

func populateUsers(conn *grpc.ClientConn) {
	client := hotelpb.NewUserServiceClient(conn)

	for i := 0; i < numOfUsers; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		username := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("password%d", i)

		req := &hotelpb.RegisterUserRequest{
			Username: username,
			Password: password,
		}

		_, err := client.RegisterUser(ctx, req)
		if err != nil {
			log.Printf("Failed to register user %s: %v", username, err)
		} else {
			log.Printf("Registered user: %s", username)
		}

		token, err := client.Login(ctx, &hotelpb.LoginRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Printf("Failed to login user %s: %v", username, err)
		} else {
			log.Printf("Login token for user %s: %s", username, token.Token)
		}
	}
}
