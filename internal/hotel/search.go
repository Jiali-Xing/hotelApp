package hotel

import (
	"context"
	"log"
	"redis_test/pkg/invoke"
	"redis_test/pkg/state"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

type SearchServer struct {
	hotelpb.UnimplementedSearchServiceServer
}

func (s *SearchServer) Nearby(ctx context.Context, req *hotelpb.NearbyRequest) (*hotelpb.NearbyResponse, error) {
	rates := Nearby(ctx, req.InDate, req.OutDate, req.Location)
	resp := &hotelpb.NearbyResponse{Rates: rates}
	return resp, nil
}

func (s *SearchServer) SearchHotels(ctx context.Context, req *hotelpb.SearchHotelsRequest) (*hotelpb.SearchHotelsResponse, error) {
	profiles := SearchHotels(ctx, req.InDate, req.OutDate, req.Location)
	resp := &hotelpb.SearchHotelsResponse{Profiles: profiles}
	return resp, nil
}

func Nearby(ctx context.Context, inDate string, outDate string, location string) []*hotelpb.Rate {
	// Find the hotel ids in that location
	hotelIds := getHotelIdsForLocation(ctx, location)

	// Get the rates for these hotels
	req := hotelpb.GetRatesRequest{HotelIds: hotelIds}
	ratesRes, err := invoke.Invoke[*hotelpb.GetRatesResponse](ctx, "rate", "GetRates", &req)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return []*hotelpb.Rate{}
	}
	return ratesRes.Rates
}

func StoreHotelLocation(ctx context.Context, hotelId string, location string) string {
	hotelIds := getHotelIdsForLocation(ctx, location)
	// Keep saved reviews bounded to 10 for consistent performance measurements
	if len(hotelIds) >= 10 {
		hotelIds = hotelIds[1:]
	}
	hotelIds = append(hotelIds, hotelId)
	state.SetState(ctx, location, hotelIds)
	return hotelId
}

func getHotelIdsForLocation(ctx context.Context, location string) []string {
	hotelIds, err := state.GetState[[]string](ctx, location)
	// If err != nil then the key does not exist
	if err != nil {
		return []string{}
	} else {
		return hotelIds
	}
}
