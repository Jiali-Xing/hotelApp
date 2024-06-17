package invoke

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jiali-Xing/hotelApp/internal/config"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	socialpb "github.com/Jiali-Xing/socialproto"
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

	case "usertimeline":
		userTimelineClient, ok := client.(socialpb.UserTimelineClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "readusertimeline":
			req, ok := input.(*socialpb.ReadUserTimelineRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := userTimelineClient.ReadUserTimeline(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "writeusertimeline":
			req, ok := input.(*socialpb.WriteUserTimelineRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := userTimelineClient.WriteUserTimeline(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "hometimeline":
		homeTimelineClient, ok := client.(socialpb.HomeTimelineClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "readhometimeline":
			req, ok := input.(*socialpb.ReadHomeTimelineRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := homeTimelineClient.ReadHomeTimeline(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "writehometimeline":
			req, ok := input.(*socialpb.WriteHomeTimelineRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := homeTimelineClient.WriteHomeTimeline(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "poststorage":
		postStorageClient, ok := client.(socialpb.PostStorageClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "storepost":
			req, ok := input.(*socialpb.StorePostRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := postStorageClient.StorePost(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "storepostmulti":
			req, ok := input.(*socialpb.StorePostMultiRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := postStorageClient.StorePostMulti(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "readpost":
			req, ok := input.(*socialpb.ReadPostRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := postStorageClient.ReadPost(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "readposts":
			req, ok := input.(*socialpb.ReadPostsRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := postStorageClient.ReadPosts(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "socialgraph":
		socialGraphClient, ok := client.(socialpb.SocialGraphClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "insertuser":
			req, ok := input.(*socialpb.InsertUserRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := socialGraphClient.InsertUser(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "getfollowers":
			req, ok := input.(*socialpb.GetFollowersRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := socialGraphClient.GetFollowers(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "getfollowees":
			req, ok := input.(*socialpb.GetFolloweesRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := socialGraphClient.GetFollowees(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "follow":
			req, ok := input.(*socialpb.FollowRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := socialGraphClient.Follow(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		case "followmany":
			req, ok := input.(*socialpb.FollowManyRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := socialGraphClient.FollowMany(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	//	above is the for social network services
	//	below is the for hotel services
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
		case "storehotellocation":
			req, ok := input.(*hotelpb.StoreHotelLocationRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := searchClient.StoreHotelLocation(ctx, req)
			if err != nil {
				return res, err
			}
			res = any(resp).(T)
		default:
			return res, fmt.Errorf("unsupported method: %s", method)
		}
	case "reservation":
		reservationClient, ok := client.(hotelpb.ReservationServiceClient)
		if !ok {
			return res, fmt.Errorf("invalid client type for service: %s", app)
		}
		switch method {
		case "addhotelavailability":
			req, ok := input.(*hotelpb.AddHotelAvailabilityRequest)
			if !ok {
				return res, fmt.Errorf("invalid input type for method: %s", method)
			}
			resp, err := reservationClient.AddHotelAvailability(ctx, req)
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
