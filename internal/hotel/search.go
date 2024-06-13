package hotel

import (
	"context"
	"log"

	"github.com/Jiali-Xing/hotelApp/internal/config"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	"github.com/Jiali-Xing/hotelApp/pkg/state"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

type SearchServer struct {
	hotelpb.UnimplementedSearchServiceServer
}

func (s *SearchServer) Nearby(ctx context.Context, req *hotelpb.NearbyRequest) (*hotelpb.NearbyResponse, error) {
	ctx = config.PropagateMetadata(ctx, "search")
	rates, err := Nearby(ctx, req.InDate, req.OutDate, req.Location)
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.NearbyResponse{Rates: rates}
	return resp, nil
}

func (s *SearchServer) StoreHotelLocation(ctx context.Context, req *hotelpb.StoreHotelLocationRequest) (*hotelpb.StoreHotelLocationResponse, error) {
	ctx = config.PropagateMetadata(ctx, "search")
	hotelId, err := StoreHotelLocation(ctx, req.HotelId, req.Location)
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.StoreHotelLocationResponse{HotelId: hotelId}
	return resp, nil
}

func Nearby(ctx context.Context, inDate string, outDate string, location string) ([]*hotelpb.Rate, error) {
	// Find the hotel ids in that location
	hotelIds, err := getHotelIdsForLocation(ctx, location)
	if err != nil {
		return nil, err
	}
	config.DebugLog("Found hotel ids: %v", hotelIds)
	// Get the rates for these hotels
	req := hotelpb.GetRatesRequest{HotelIds: hotelIds}
	ratesRes, err := invoke.Invoke[*hotelpb.GetRatesResponse](ctx, "rate", "GetRates", &req)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return nil, err
	}
	config.DebugLog("Found rates: %v", ratesRes.Rates)
	return ratesRes.Rates, nil
}

func StoreHotelLocation(ctx context.Context, hotelId string, location string) (string, error) {
	hotelIds, err := getHotelIdsForLocation(ctx, location)
	if err != nil {
		return "", err
	}
	// Keep saved reviews bounded to 10 for consistent performance measurements
	if len(hotelIds) >= 10 {
		hotelIds = hotelIds[1:]
	}
	hotelIds = append(hotelIds, hotelId)
	state.SetState(ctx, location, hotelIds)
	return hotelId, nil
}

func getHotelIdsForLocation(ctx context.Context, location string) ([]string, error) {
	hotelIds, err := state.GetState[[]string](ctx, location)
	// If err != nil then the key does not exist
	if err != nil {
		return []string{}, err
	}
	return hotelIds, nil
}
