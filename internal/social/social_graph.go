package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type GraphServer struct {
	socialpb.UnimplementedSocialGraphServer
}

type SGVertex struct {
	UserId    string   `json:"userId"`
	Followers []string `json:"followers"`
	Followees []string `json:"followees"` // Fixed to match the JSON key
}

// InsertUser inserts a user with an empty social graph
func (s *GraphServer) InsertUser(ctx context.Context, req *socialpb.InsertUserRequest) (*socialpb.InsertUserResponse, error) {
	sg := SGVertex{
		UserId:    req.UserId,
		Followers: []string{},
		Followees: []string{},
	}
	err := state.SetState(ctx, req.UserId, sg)
	if err != nil {
		config.DebugLog("Error inserting user %s: %v", req.UserId, err)
		return nil, err
	}
	config.DebugLog("Inserted user: %s", req.UserId)
	return &socialpb.InsertUserResponse{}, nil
}

// GetFollowers retrieves the list of followers for a given user
func (s *GraphServer) GetFollowers(ctx context.Context, req *socialpb.GetFollowersRequest) (*socialpb.GetFollowersResponse, error) {
	sg, err := state.GetState[SGVertex](ctx, req.UserId)
	if err != nil {
		config.DebugLog("Error getting followers for user %s: %v", req.UserId, err)
		return nil, err
	}
	config.DebugLog("Retrieved followers for user %s: %v", req.UserId, sg.Followers)
	return &socialpb.GetFollowersResponse{Followers: sg.Followers}, nil
}

// GetFollowees retrieves the list of followees for a given user
func (s *GraphServer) GetFollowees(ctx context.Context, req *socialpb.GetFolloweesRequest) (*socialpb.GetFolloweesResponse, error) {
	sg, err := state.GetState[SGVertex](ctx, req.UserId)
	if err != nil {
		config.DebugLog("Error getting followees for user %s: %v", req.UserId, err)
		return nil, err
	}
	config.DebugLog("Retrieved followees for user %s: %v", req.UserId, sg.Followees)
	return &socialpb.GetFolloweesResponse{Followees: sg.Followees}, nil
}

// Follow allows a user to follow another user
func (s *GraphServer) Follow(ctx context.Context, req *socialpb.FollowRequest) (*socialpb.FollowResponse, error) {
	err := s.follow(ctx, req.FollowerId, req.FolloweeId)
	if err != nil {
		config.DebugLog("Error following user %s to user %s: %v", req.FollowerId, req.FolloweeId, err)
		return nil, err
	}
	config.DebugLog("User %s followed user %s", req.FollowerId, req.FolloweeId)
	return &socialpb.FollowResponse{}, nil
}

// follow is a helper function to handle the following logic
func (s *GraphServer) follow(ctx context.Context, followerId string, followeeId string) error {
	// Retrieve the follower's state
	sgFollower, err := state.GetState[SGVertex](ctx, followerId)
	if err != nil {
		config.DebugLog("Error getting state for follower %s: %v", followerId, err)
		sgFollower = SGVertex{
			UserId:    followerId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	config.DebugLog("Before following: follower %s has followees: %v", followerId, sgFollower.Followees)
	// Add the followee to the follower's followees
	sgFollower.Followees = append(sgFollower.Followees, followeeId)
	err = state.SetState(ctx, followerId, sgFollower)
	if err != nil {
		config.DebugLog("Error setting state for follower %s: %v", followerId, err)
		return err
	}
	config.DebugLog("After following: follower %s has followees: %v", followerId, sgFollower.Followees)

	// Retrieve the followee's state
	sgFollowee, err := state.GetState[SGVertex](ctx, followeeId)
	if err != nil {
		config.DebugLog("Error getting state for followee %s: %v", followeeId, err)
		sgFollowee = SGVertex{
			UserId:    followeeId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	config.DebugLog("Before following: followee %s has followers: %v", followeeId, sgFollowee.Followers)
	// Add the follower to the followee's followers
	sgFollowee.Followers = append(sgFollowee.Followers, followerId)
	err = state.SetState(ctx, followeeId, sgFollowee)
	if err != nil {
		config.DebugLog("Error setting state for followee %s: %v", followeeId, err)
		return err
	}
	config.DebugLog("After following: followee %s has followers: %v", followeeId, sgFollowee.Followers)
	return nil
}
