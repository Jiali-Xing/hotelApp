package invoke

import (
	"context"
	"fmt"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

func Invoke[T any](ctx context.Context, app string, method string, input interface{}) (T, error) {
	var res T

	client, err := getClient(app)
	if err != nil {
		return res, err
	}

	switch app {
	case "user":
		userClient, ok := client.(hotelpb.UserServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "login":
			req, ok := input.(*hotelpb.LoginRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := userClient.Login(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "registerUser":
			req, ok := input.(*hotelpb.RegisterUserRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := userClient.RegisterUser(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
			// Add more user methods here
		}
	case "search":
		searchClient, ok := client.(hotelpb.SearchServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "SearchHotels":
			req, ok := input.(*hotelpb.SearchHotelsRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := searchClient.SearchHotels(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "Nearby":
			req, ok := input.(*hotelpb.NearbyRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := searchClient.Nearby(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
			// Add more search methods here
		}
	case "reservation":
		reservationClient, ok := client.(hotelpb.ReservationServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "FrontendReservation":
			req, ok := input.(*hotelpb.FrontendReservationRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := reservationClient.FrontendReservation(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "MakeReservation":
			req, ok := input.(*hotelpb.MakeReservationRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := reservationClient.MakeReservation(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "CheckAvailability":
			req, ok := input.(*hotelpb.CheckAvailabilityRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := reservationClient.CheckAvailability(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
			// Add more reservation methods here
		}
	case "rate":
		rateClient, ok := client.(hotelpb.RateServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "GetRates":
			req, ok := input.(*hotelpb.GetRatesRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := rateClient.GetRates(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "StoreRate":
			req, ok := input.(*hotelpb.StoreRateRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := rateClient.StoreRate(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
			// Add more rate methods here
		}
	case "profile":
		profileClient, ok := client.(hotelpb.ProfileServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "GetProfiles":
			req, ok := input.(*hotelpb.GetProfilesRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := profileClient.GetProfiles(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "StoreProfile":
			req, ok := input.(*hotelpb.StoreProfileRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := profileClient.StoreProfile(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
			// Add more profile methods here
		}
	// Add more services here
	default:
		return res, fmt.Errorf("unknown app: %s", app)
	}

	return res, nil
}
