package hotel

import (
	"context"
	"log"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	hotelpb "github.com/Jiali-Xing/hotelproto"
)

type ProfileServer struct {
	hotelpb.UnimplementedProfileServiceServer
}

func (s *ProfileServer) StoreProfile(ctx context.Context, req *hotelpb.StoreProfileRequest) (*hotelpb.StoreProfileResponse, error) {
	ctx = config.PropagateMetadata(ctx, "profile")
	hotelId, err := StoreProfile(ctx, req.Profile)
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.StoreProfileResponse{HotelId: hotelId}
	return resp, nil
}

func (s *ProfileServer) GetProfiles(ctx context.Context, req *hotelpb.GetProfilesRequest) (*hotelpb.GetProfilesResponse, error) {
	ctx = config.PropagateMetadata(ctx, "profile")
	profiles, err := GetProfiles(ctx, req.HotelIds)
	if err != nil {
		return nil, err
	}
	resp := &hotelpb.GetProfilesResponse{Profiles: profiles}
	return resp, nil
}

func StoreProfile(ctx context.Context, profile *hotelpb.HotelProfile) (string, error) {
	err := state.SetState(ctx, profile.HotelId, profile)
	if err != nil {
		log.Printf("Failed to store profile for hotelId %s: %v", profile.HotelId, err)
		return "", err
	}
	return profile.HotelId, nil
}

func GetProfiles(ctx context.Context, hotelIds []string) ([]*hotelpb.HotelProfile, error) {
	// Bulk
	var profiles []*hotelpb.HotelProfile
	var err error
	if len(hotelIds) > 0 {
		profiles, err = state.GetBulkStateDefault[*hotelpb.HotelProfile](ctx, hotelIds, &hotelpb.HotelProfile{})
		if err != nil {
			log.Printf("Failed to get profiles for hotelIds %v: %v", hotelIds, err)
			return nil, err
		}
	} else {
		profiles = make([]*hotelpb.HotelProfile, len(hotelIds))
	}
	config.DebugLog("Found profiles: %v for hotel ids: %v", profiles, hotelIds)
	return profiles, nil
}
