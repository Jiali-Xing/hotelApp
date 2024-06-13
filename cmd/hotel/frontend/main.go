package main

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/internal/config"

	"log"
	"net"
	"os"

	"github.com/Jiali-Xing/hotelApp/internal/hotel"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

type server struct {
	hotelpb.UnimplementedFrontendServiceServer
}

func (s *server) SearchHotels(ctx context.Context, req *hotelpb.SearchHotelsRequest) (*hotelpb.SearchHotelsResponse, error) {
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
	ctx = config.PropagateMetadata(ctx, "frontend")
	success, err := hotel.FrontendReservation(ctx, req.HotelId, req.InDate, req.OutDate, int(req.Rooms), req.Username, req.Password)
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

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	hotelServer := &server{}
	hotelpb.RegisterFrontendServiceServer(s, hotelServer)

	log.Println("gRPC server listening on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
