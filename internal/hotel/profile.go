package hotel

import (
	"context"
	"github.com/Jiali-Xing/hotelApp/internal/config"

	"github.com/Jiali-Xing/hotelApp/pkg/state"

	hotelpb "github.com/Jiali-Xing/hotelproto"
)

type ProfileServer struct {
	hotelpb.UnimplementedProfileServiceServer
}

func (s *ProfileServer) StoreProfile(ctx context.Context, req *hotelpb.StoreProfileRequest) (*hotelpb.StoreProfileResponse, error) {
	ctx = propagateMetadata(ctx, "profile")
	hotelId := StoreProfile(ctx, req.Profile)
	resp := &hotelpb.StoreProfileResponse{HotelId: hotelId}
	return resp, nil
}

func (s *ProfileServer) GetProfiles(ctx context.Context, req *hotelpb.GetProfilesRequest) (*hotelpb.GetProfilesResponse, error) {
	ctx = propagateMetadata(ctx, "profile")
	profiles := GetProfiles(ctx, req.HotelIds)
	resp := &hotelpb.GetProfilesResponse{Profiles: profiles}
	return resp, nil
}

func StoreProfile(ctx context.Context, profile *hotelpb.HotelProfile) string {
	state.SetState(ctx, profile.HotelId, profile)
	return profile.HotelId
}

func GetProfiles(ctx context.Context, hotelIds []string) []*hotelpb.HotelProfile {
	// Bulk
	var profiles []*hotelpb.HotelProfile
	if len(hotelIds) > 0 {
		profiles = state.GetBulkStateDefault[*hotelpb.HotelProfile](ctx, hotelIds, &hotelpb.HotelProfile{})
	} else {
		profiles = make([]*hotelpb.HotelProfile, len(hotelIds))
	}
	config.DebugLog("Found profiles: %v for hotel ids: %v", profiles, hotelIds)
	return profiles
}
