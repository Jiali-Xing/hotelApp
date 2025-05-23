package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/internal/hotel"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	hotelpb "github.com/Jiali-Xing/hotelproto"
	"github.com/Jiali-Xing/plain"
	"github.com/valyala/fastrand"
	"google.golang.org/grpc"
)

type server struct {
	hotelpb.UnimplementedFrontendServiceServer
}

// Helper function to generate a random username
// func generateRandomUserAndPassword() (string, string) {
// 	randomNum := fastrand.Uint32n(100)
// 	username := fmt.Sprintf("user%d", randomNum)
// 	password := fmt.Sprintf("password%d", randomNum)
// 	return username, password
// }

// Helper function to generate a random hotel ID
func generateRandomHotelID() string {
	randomNum := fastrand.Uint32n(1000) + 1
	return fmt.Sprintf("%d", randomNum)
}

// Helper function to generate random dates within a specified range
func generateRandomDates() (string, string) {
	start := time.Now().AddDate(0, 0, int(fastrand.Uint32n(30))) // Randomize start date
	end := start.AddDate(0, 0, int(fastrand.Uint32n(30)))        // Randomize end date
	inDate := start.Format("2006-01-02")
	outDate := end.Format("2006-01-02")
	return inDate, outDate
}

func generateRandomCity() string {
	// return a random city from the list houston-tx-0 new-york-city-ny-0 los-angeles-ca-0 chicago-il-0
	cities := []string{"houston-tx-0", "new-york-city-ny-0", "los-angeles-ca-0", "chicago-il-0"}
	randomNum := fastrand.Uint32n(4)
	return cities[randomNum]
}

func (s *server) SearchHotels(ctx context.Context, req *hotelpb.SearchHotelsRequest) (*hotelpb.SearchHotelsResponse, error) {
	// Randomize InDate and OutDate
	req.InDate, req.OutDate = generateRandomDates()
	// Randomize Location
	req.Location = generateRandomCity()

	ctx = config.PropagateMetadata(ctx, "frontend")
	hotels, err := hotel.SearchHotels(ctx, req.InDate, req.OutDate, req.Location)
	if err != nil {
		log.Printf("Error searching hotels: %v", err)
		return nil, err
	}
	resp := &hotelpb.SearchHotelsResponse{Profiles: hotels}
	return resp, nil
}

func (s *server) StoreHotel(ctx context.Context, req *hotelpb.StoreHotelRequest) (*hotelpb.StoreHotelResponse, error) {
	ctx = config.PropagateMetadata(ctx, "frontend")
	hotelId, err := hotel.StoreHotel(ctx, req.HotelId, req.Name, req.Phone, req.Location, int(req.Rate), int(req.Capacity), req.Info)
	if err != nil {
		log.Printf("Error storing hotel: %v", err)
		return nil, err
	}
	resp := &hotelpb.StoreHotelResponse{HotelId: hotelId}
	return resp, nil
}

func (s *server) FrontendReservation(ctx context.Context, req *hotelpb.FrontendReservationRequest) (*hotelpb.FrontendReservationResponse, error) {
	// Randomize the Username, Password, HotelId, InDate, and OutDate
	// if username is user0, password is password0, if username is user99, password is password99, etc.
	username, password, err := config.GenerateUserAndPassword(ctx)
	if err != nil {
		log.Printf("Error generating username and password: %v", err)
		return nil, err
	}
	req.HotelId = generateRandomHotelID()
	req.InDate, req.OutDate = generateRandomDates()

	ctx = config.PropagateMetadata(ctx, "frontend")
	success, err := hotel.FrontendReservation(ctx, req.HotelId, req.InDate, req.OutDate, int(req.Rooms), username, password)
	if err != nil {
		log.Printf("Error making reservation: %v", err)
		return nil, err
	}
	resp := &hotelpb.FrontendReservationResponse{Success: success}
	return resp, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052" // Default port if not specified
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
	case "topdown":
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(config.Topdown.UnaryInterceptor))
		go config.Topdown.StartServer(8082)
		config.DebugLog(fmt.Sprintf("Topdown gRPC RL agent server listening on port %d", 8082))
	default:
		grpcServer = grpc.NewServer()
	}

	// Register the frontend service
	hotelServer := &server{}
	hotelpb.RegisterFrontendServiceServer(grpcServer, hotelServer)

	// Listen and serve
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ctx := context.Background()
	// Establish a gRPC connection to other services
	userConn, err := config.CreateGRPCConn(ctx, config.UserAddr)
	if err != nil {
		log.Fatalf("Failed to connect to user gRPC server: %v", err)
	}
	defer userConn.Close()
	invoke.RegisterClient("user", hotelpb.NewUserServiceClient(userConn))

	searchConn, err := config.CreateGRPCConn(ctx, config.SearchAddr)
	if err != nil {
		log.Fatalf("Failed to connect to search gRPC server: %v", err)
	}
	defer searchConn.Close()
	invoke.RegisterClient("search", hotelpb.NewSearchServiceClient(searchConn))

	reservationConn, err := config.CreateGRPCConn(ctx, config.ReservationAddr)
	if err != nil {
		log.Fatalf("Failed to connect to reservation gRPC server: %v", err)
	}
	defer reservationConn.Close()
	invoke.RegisterClient("reservation", hotelpb.NewReservationServiceClient(reservationConn))

	rateConn, err := config.CreateGRPCConn(ctx, config.RateAddr)
	if err != nil {
		log.Fatalf("Failed to connect to rate gRPC server: %v", err)
	}
	defer rateConn.Close()
	invoke.RegisterClient("rate", hotelpb.NewRateServiceClient(rateConn))

	profileConn, err := config.CreateGRPCConn(ctx, config.ProfileAddr)
	if err != nil {
		log.Fatalf("Failed to connect to profile gRPC server: %v", err)
	}
	defer profileConn.Close()
	invoke.RegisterClient("profile", hotelpb.NewProfileServiceClient(profileConn))

	log.Println("gRPC server listening on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
