package hotel

import (
	"context"
	"log"

	"redis_test/pkg/invoke"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

func SearchHotels(ctx context.Context, inDate string, outDate string, location string) []*hotelpb.HotelProfile {
	req1 := &hotelpb.NearbyRequest{InDate: inDate, OutDate: outDate, Location: location}
	hotelIdsRes, err := invoke.Invoke[*hotelpb.NearbyResponse](ctx, "search", "Nearby", req1)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return nil
	}
	rates := hotelIdsRes.Rates

	hotelIds := make([]string, len(rates))
	for i, rate := range rates {
		hotelIds[i] = rate.HotelId
	}

	req2 := &hotelpb.CheckAvailabilityRequest{
		CustomerName: "",
		HotelIds:     hotelIds,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   1,
	}
	availableHotelIdsRes, err := invoke.Invoke[*hotelpb.CheckAvailabilityResponse](ctx, "reservation", "CheckAvailability", req2)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return nil
	}

	req3 := &hotelpb.GetProfilesRequest{HotelIds: availableHotelIdsRes.HotelIds}
	profilesRes, err := invoke.Invoke[*hotelpb.GetProfilesResponse](ctx, "profile", "GetProfiles", req3)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return nil
	}
	return profilesRes.Profiles
}

func StoreHotel(ctx context.Context, hotelId string, name string, phone string, location string, rate int, capacity int, info string) string {
	req1 := &hotelpb.StoreHotelLocationRequest{Location: location, HotelId: hotelId}
	_, err := invoke.Invoke[*hotelpb.StoreHotelLocationResponse](ctx, "search", "StoreHotelLocation", req1)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return ""
	}

	req2 := &hotelpb.StoreRateRequest{Rate: &hotelpb.Rate{HotelId: hotelId, Price: int32(rate)}}
	_, err = invoke.Invoke[*hotelpb.StoreRateResponse](ctx, "rate", "StoreRate", req2)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return ""
	}

	req3 := &hotelpb.AddHotelAvailabilityRequest{
		HotelId:  hotelId,
		Capacity: int32(capacity),
	}
	_, err = invoke.Invoke[*hotelpb.AddHotelAvailabilityResponse](ctx, "reservation", "AddHotelAvailability", req3)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return ""
	}

	hotelProfile := &hotelpb.HotelProfile{
		HotelId: hotelId,
		Name:    name,
		Phone:   phone,
		Info:    info,
	}
	req4 := &hotelpb.StoreProfileRequest{Profile: hotelProfile}
	_, err = invoke.Invoke[*hotelpb.StoreProfileResponse](ctx, "profile", "StoreProfile", req4)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return ""
	}
	return hotelId
}

func FrontendReservation(ctx context.Context, hotelId string, inDate string, outDate string, rooms int, username string, password string) bool {
	req1 := &hotelpb.LoginRequest{
		Username: username,
		Password: password,
	}
	tokenRes, err := invoke.Invoke[*hotelpb.LoginResponse](ctx, "user", "Login", req1)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return false
	}
	if tokenRes.Token != "OK" {
		return false
	}

	req2 := &hotelpb.MakeReservationRequest{
		CustomerName: username,
		HotelId:      hotelId,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   int32(rooms),
	}
	successRes, err := invoke.Invoke[*hotelpb.MakeReservationResponse](ctx, "reservation", "MakeReservation", req2)
	if err != nil {
		log.Printf("Error invoking gRPC method: %v", err)
		return false
	}
	return successRes.Success
}
