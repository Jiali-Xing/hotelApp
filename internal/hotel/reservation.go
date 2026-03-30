package hotel

import (
	"context"
	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	hotelpb "github.com/Jiali-Xing/hotelproto"
	"log"
)

type ReservationServer struct {
	hotelpb.UnimplementedReservationServiceServer
}

func (s *ReservationServer) CheckAvailability(ctx context.Context, req *hotelpb.CheckAvailabilityRequest) (*hotelpb.CheckAvailabilityResponse, error) {
	ctx = config.PropagateMetadata(ctx, "reservation")
	hotelIds, err := CheckAvailability(ctx, req.CustomerName, req.HotelIds, req.InDate, req.OutDate, int(req.RoomNumber))
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.CheckAvailabilityResponse{HotelIds: hotelIds}
	return resp, nil
}

func (s *ReservationServer) MakeReservation(ctx context.Context, req *hotelpb.MakeReservationRequest) (*hotelpb.MakeReservationResponse, error) {
	ctx = config.PropagateMetadata(ctx, "reservation")
	success, err := MakeReservation(ctx, req.CustomerName, req.HotelId, req.InDate, req.OutDate, int(req.RoomNumber))
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.MakeReservationResponse{Success: success}
	return resp, nil
}

func (s *ReservationServer) AddHotelAvailability(ctx context.Context, req *hotelpb.AddHotelAvailabilityRequest) (*hotelpb.AddHotelAvailabilityResponse, error) {
	ctx = config.PropagateMetadata(ctx, "reservation")
	hotelId, err := AddHotelAvailability(ctx, req.HotelId, int(req.Capacity))
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.AddHotelAvailabilityResponse{HotelId: hotelId}
	return resp, nil
}

func datesIntersect(inDate1 string, outDate1 string, inDate2 string, outDate2 string) bool {
	return !(inDate2 > outDate1 || inDate1 > outDate2)
}

func checkAvailability(availability hotelpb.HotelAvailability, inDate string, outDate string, numberOfRooms int) bool {
	capacity := availability.Capacity
	reservationsTheseDays := 0
	for _, reservation := range availability.Reservations {
		if datesIntersect(inDate, outDate, reservation.InDate, reservation.OutDate) {
			reservationsTheseDays++
		}
	}
	config.DebugLog("Reservations these days: %d + #Rooms: %d <> Capacity: %d", reservationsTheseDays, numberOfRooms, capacity)
	return reservationsTheseDays+numberOfRooms <= int(capacity)
}

func CheckAvailability(ctx context.Context, customerName string, hotelIds []string, inDate string, outDate string, numberOfRooms int) ([]string, error) {
	availableHotelIds := []string{}
	for _, hotelId := range hotelIds {
		availability, err := state.GetState[hotelpb.HotelAvailability](ctx, hotelId)
		if err != nil {
			log.Printf("Failed to get availability for hotelId %s: %v", hotelId, err)
			return nil, err
		}

		isAvailable := checkAvailability(availability, inDate, outDate, numberOfRooms)
		if isAvailable {
			availableHotelIds = append(availableHotelIds, hotelId)
		}
		config.DebugLog("Hotel %s is available: %v", hotelId, isAvailable)
	}
	return availableHotelIds, nil
}

func MakeReservation(ctx context.Context, customerName string, hotelId string, inDate string, outDate string, numberOfRooms int) (bool, error) {
	availability, err := state.GetState[hotelpb.HotelAvailability](ctx, hotelId)
	if err != nil {
		log.Printf("Failed to get availability for hotelId %s: %v", hotelId, err)
		return false, err
	}

	if !checkAvailability(availability, inDate, outDate, numberOfRooms) {
		return false, nil
	}

	if len(availability.Reservations) >= 10 {
		availability.Reservations = availability.Reservations[1:]
	}

	newReservation := hotelpb.Reservation{
		CustomerName: customerName,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   int32(numberOfRooms),
	}
	availability.Reservations = append(availability.Reservations, &newReservation)
	config.DebugLog("Adding reservation: %v", newReservation)
	err = state.SetState(ctx, hotelId, availability)
	if err != nil {
		log.Printf("Failed to set availability for hotelId %s: %v", hotelId, err)
		return false, err
	}
	return true, nil
}

func AddHotelAvailability(ctx context.Context, hotelId string, capacity int) (string, error) {
	availability := hotelpb.HotelAvailability{
		Reservations: []*hotelpb.Reservation{},
		Capacity:     int32(capacity),
	}
	err := state.SetState(ctx, hotelId, availability)
	if err != nil {
		log.Printf("Failed to set availability for hotelId %s: %v", hotelId, err)
		return "", err
	}
	config.DebugLog("Added availability for hotel id: %s with capacity: %d", hotelId, capacity)
	return hotelId, nil
}
