package invoke

import (
	"context"
	"fmt"
	"redis_test/internal/config"
	"strings"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

func Invoke[T any](ctx context.Context, app string, method string, input interface{}) (T, error) {
	var res T

	client, err := getClient(app)
	if err != nil {
		config.DebugLog("Error getting client for app %s: %v", app, err)
		return res, err
	}

	// Convert method name to lower case to make comparison case-insensitive
	method = strings.ToLower(method)

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
				config.DebugLog("Error calling Login method: %v", err)
				return res, err
			}

			if resp == nil {
				config.DebugLog("Received nil response from login gRPC call")
				return res, fmt.Errorf("received nil response from login gRPC call")
			}
			config.DebugLog("Received response from Login method: %v", resp)
			res = any(resp).(T)
		case "registeruser":
			req, ok := input.(*hotelpb.RegisterUserRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := userClient.RegisterUser(ctx, req)
			if err != nil {
				return res, err
			}

			if resp == nil {
				config.DebugLog("Received nil response from register user gRPC call")
				return res, fmt.Errorf("received nil response from register user gRPC call")
			}
			config.DebugLog("Received response from RegisterUser method: %v", resp)
			res = any(resp).(T)
			// Add more user methods here
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "search":
		searchClient, ok := client.(hotelpb.SearchServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "searchhotels":
			req, ok := input.(*hotelpb.SearchHotelsRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := searchClient.SearchHotels(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "nearby":
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
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "reservation":
		reservationClient, ok := client.(hotelpb.ReservationServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "frontendreservation":
			req, ok := input.(*hotelpb.FrontendReservationRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := reservationClient.FrontendReservation(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "makereservation":
			req, ok := input.(*hotelpb.MakeReservationRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := reservationClient.MakeReservation(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "checkavailability":
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
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "rate":
		rateClient, ok := client.(hotelpb.RateServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "getrates":
			req, ok := input.(*hotelpb.GetRatesRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := rateClient.GetRates(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "storerate":
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
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "profile":
		profileClient, ok := client.(hotelpb.ProfileServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "getprofiles":
			req, ok := input.(*hotelpb.GetProfilesRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := profileClient.GetProfiles(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "storeprofile":
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
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	// Add more services here
	default:
		return res, fmt.Errorf("unknown app: %s", app)
	}

	return res, nil
}
