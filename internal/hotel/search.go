package hotel

import (
	"context"
	"github.com/Jiali-Xing/hotelApp/internal/config"
	"log"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	"github.com/Jiali-Xing/hotelApp/pkg/state"

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

func (s *SearchServer) StoreHotelLocation(ctx context.Context, req *hotelpb.StoreHotelLocationRequest) (*hotelpb.StoreHotelLocationResponse, error) {
	hotelId := StoreHotelLocation(ctx, req.HotelId, req.Location)
	resp := &hotelpb.StoreHotelLocationResponse{HotelId: hotelId}
	return resp, nil
}

func Nearby(ctx context.Context, inDate string, outDate string, location string) []*hotelpb.Rate {
	// Find the hotel ids in that location
	hotelIds := getHotelIdsForLocation(ctx, location)
	config.DebugLog("Found hotel ids: %v", hotelIds)
	// Get the rates for these hotels
	req := hotelpb.GetRatesRequest{HotelIds: hotelIds}
	ratesRes, err := invoke.Invoke[*hotelpb.GetRatesResponse](ctx, "rate", "GetRates", &req)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return []*hotelpb.Rate{}
	}
	config.DebugLog("Found rates: %v", ratesRes.Rates)
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
