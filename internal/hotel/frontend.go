package hotel

import (
	"context"
	"log"

	"redis_test/pkg/invoke"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

type FrontendServer struct {
	hotelpb.UnimplementedReservationServiceServer
}

func (s *FrontendServer) FrontendReservation(ctx context.Context, req *hotelpb.FrontendReservationRequest) (*hotelpb.FrontendReservationResponse, error) {
	success := FrontendReservation(ctx, req.HotelId, req.InDate, req.OutDate, int(req.Rooms), req.Username, req.Password)
	resp := &hotelpb.FrontendReservationResponse{Success: success}
	return resp, nil
}

func FrontendReservation(ctx context.Context, hotelId, inDate, outDate string, rooms int, username, password string) bool {
	req1 := &hotelpb.LoginRequest{
		Username: username,
		Password: password,
	}
	tokenRes, err := invoke.Invoke[*hotelpb.LoginResponse](ctx, "user", "login", req1)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return false
	}

	if tokenRes.Token != "OK" {
		return false
	}

	return true
}
