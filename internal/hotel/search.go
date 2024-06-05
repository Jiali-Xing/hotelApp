package hotel

import (
	"context"

	pb "github.com/Jiali-Xing/hotelproto"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}

// Implement the Nearby method
func (s *SearchService) Nearby(ctx context.Context, req *pb.NearbyRequest) (*pb.NearbyResponse, error) {
	// Dummy implementation, replace with actual logic
	rates := []*pb.Rate{
		{HotelId: "hotel1", Price: 100},
		{HotelId: "hotel2", Price: 150},
	}
	return &pb.NearbyResponse{Rates: rates}, nil
}

// Implement the SearchHotels method
func (s *SearchService) SearchHotels(ctx context.Context, req *pb.SearchHotelsRequest) (*pb.SearchHotelsResponse, error) {
	// Dummy implementation, replace with actual logic
	profiles := []*pb.HotelProfile{
		{HotelId: "hotel1", Name: "Hotel One"},
		{HotelId: "hotel2", Name: "Hotel Two"},
	}
	return &pb.SearchHotelsResponse{Profiles: profiles}, nil
}
