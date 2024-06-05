package hotel

import (
	"context"
	"redis_test/pkg/state"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

type RateServer struct {
	hotelpb.UnimplementedRateServiceServer
}

func (s *RateServer) StoreRate(ctx context.Context, req *hotelpb.StoreRateRequest) (*hotelpb.StoreRateResponse, error) {
	hotelId := StoreRate(ctx, req.Rate)
	resp := &hotelpb.StoreRateResponse{HotelId: hotelId}
	return resp, nil
}

func (s *RateServer) GetRates(ctx context.Context, req *hotelpb.GetRatesRequest) (*hotelpb.GetRatesResponse, error) {
	rates := GetRates(ctx, req.HotelIds)
	resp := &hotelpb.GetRatesResponse{Rates: rates}
	return resp, nil
}

func StoreRate(ctx context.Context, rate *hotelpb.Rate) string {
	state.SetState(ctx, rate.HotelId, rate)
	return rate.HotelId
}

func GetRates(ctx context.Context, hotelIds []string) []*hotelpb.Rate {
	rates := make([]*hotelpb.Rate, len(hotelIds))
	for i, hotelId := range hotelIds {
		rate, err := state.GetState[*hotelpb.Rate](ctx, hotelId)
		if err != nil {
			panic(err)
		}
		rates[i] = rate
	}
	return rates
}
