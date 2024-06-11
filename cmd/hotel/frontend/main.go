package main

import (
	"context"
	"github.com/Jiali-Xing/hotelApp/internal/config"

	"github.com/Jiali-Xing/hotelApp/internal/hotel"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	"log"
	"net"
	"os"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"google.golang.org/grpc"
)

type server struct {
	hotelpb.UnimplementedFrontendServiceServer
}

func (s *server) SearchHotels(ctx context.Context, req *hotelpb.SearchHotelsRequest) (*hotelpb.SearchHotelsResponse, error) {
	ctx = propagateMetadata(ctx, "frontend")
	hotels := hotel.SearchHotels(ctx, req.InDate, req.OutDate, req.Location)
	resp := &hotelpb.SearchHotelsResponse{Profiles: hotels}
	return resp, nil
}

func (s *server) StoreHotel(ctx context.Context, req *hotelpb.StoreHotelRequest) (*hotelpb.StoreHotelResponse, error) {
	ctx = propagateMetadata(ctx, "frontend")
	hotelId := hotel.StoreHotel(ctx, req.HotelId, req.Name, req.Phone, req.Location, int(req.Rate), int(req.Capacity), req.Info)
	resp := &hotelpb.StoreHotelResponse{HotelId: hotelId}
	return resp, nil
}

func (s *server) FrontendReservation(ctx context.Context, req *hotelpb.FrontendReservationRequest) (*hotelpb.FrontendReservationResponse, error) {
	ctx = propagateMetadata(ctx, "frontend")
	success := hotel.FrontendReservation(ctx, req.HotelId, req.InDate, req.OutDate, int(req.Rooms), req.Username, req.Password)
	resp := &hotelpb.FrontendReservationResponse{Success: success}
	return resp, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052" // Default port if not specified
	}
	// Establish a gRPC connection to other services
	userConn, err := createGRPCConn(config.UserAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user gRPC server: %v", err)
	}
	defer userConn.Close()
	invoke.RegisterClient("user", hotelpb.NewUserServiceClient(userConn))

	searchConn, err := createGRPCConn(config.SearchAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to search gRPC server: %v", err)
	}
	defer searchConn.Close()
	invoke.RegisterClient("search", hotelpb.NewSearchServiceClient(searchConn))

	reservationConn, err := createGRPCConn(config.ReservationAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to reservation gRPC server: %v", err)
	}
	defer reservationConn.Close()
	invoke.RegisterClient("reservation", hotelpb.NewReservationServiceClient(reservationConn))

	rateConn, err := createGRPCConn(config.RateAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to rate gRPC server: %v", err)
	}
	defer rateConn.Close()
	invoke.RegisterClient("rate", hotelpb.NewRateServiceClient(rateConn))

	profileConn, err := createGRPCConn(config.ProfileAddr, grpc.WithInsecure())
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
