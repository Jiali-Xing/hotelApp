package hotel

import (
	"context"
	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	hotelpb "github.com/Jiali-Xing/hotelproto"
	"log"
)

type RateServer struct {
	hotelpb.UnimplementedRateServiceServer
}

func (s *RateServer) StoreRate(ctx context.Context, req *hotelpb.StoreRateRequest) (*hotelpb.StoreRateResponse, error) {
	ctx = config.PropagateMetadata(ctx, "rate")
	hotelId, err := StoreRate(ctx, req.Rate)
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.StoreRateResponse{HotelId: hotelId}
	return resp, nil
}

func (s *RateServer) GetRates(ctx context.Context, req *hotelpb.GetRatesRequest) (*hotelpb.GetRatesResponse, error) {
	ctx = config.PropagateMetadata(ctx, "rate")
	rates, err := GetRates(ctx, req.HotelIds)
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.GetRatesResponse{Rates: rates}
	return resp, nil
}

func StoreRate(ctx context.Context, rate *hotelpb.Rate) (string, error) {
	err := state.SetState(ctx, rate.HotelId, rate)
	if err != nil {
		log.Printf("Failed to store rate for hotelId %s: %v", rate.HotelId, err)
		return "", err
	}
	return rate.HotelId, nil
}

func GetRates(ctx context.Context, hotelIds []string) ([]*hotelpb.Rate, error) {
	rates := make([]*hotelpb.Rate, len(hotelIds))
	for i, hotelId := range hotelIds {
		rate, err := state.GetState[*hotelpb.Rate](ctx, hotelId)
		if err != nil {
			log.Printf("Failed to get rate for hotelId %s: %v", hotelId, err)
			return nil, err
		}
		rates[i] = rate
	}
	return rates, nil
}
