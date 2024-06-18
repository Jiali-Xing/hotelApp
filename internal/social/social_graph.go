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
	config.DebugLog("Error getting followers for user %s: %v", req.UserId, err)
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

// FollowMany allows a user to follow multiple users and be followed by multiple users
func (s *GraphServer) FollowMany(ctx context.Context, req *socialpb.FollowManyRequest) (*socialpb.FollowManyResponse, error) {
	sg := SGVertex{
		UserId:    req.UserId,
		Followers: req.FollowerIds,
		Followees: req.FolloweeIds,
	}
	if len(sg.Followers) >= 10 {
		sg.Followers = sg.Followers[:10]
	}
	if len(sg.Followees) >= 10 {
		sg.Followees = sg.Followees[:10]
	}
	err := state.SetState(ctx, req.UserId, sg)
	if err != nil {
		config.DebugLog("Error following many users for user %s: %v", req.UserId, err)
		return nil, err
	}
	config.DebugLog("User %s followed many users", req.UserId)
	return &socialpb.FollowManyResponse{}, nil
}

func (s *GraphServer) follow(ctx context.Context, followerId string, followeeId string) error {
	// Retrieve and update the follower's followees list
	followerSG, err := state.GetState[SGVertex](ctx, followerId)
	if err != nil && err.Error() != "key not found" {
		config.DebugLog("Error getting followees for user %s: %v", followerId, err)
		return err
	}
	if err != nil {
		followerSG = SGVertex{
			UserId:    followerId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	followerSG.Followees = append(followerSG.Followees, followeeId)
	err = state.SetState(ctx, followerId, followerSG)
	if err != nil {
		config.DebugLog("Error setting followees for user %s: %v", followerId, err)
		return err
	}

	// Retrieve and update the followee's followers list
	followeeSG, err := state.GetState[SGVertex](ctx, followeeId)
	if err != nil && err.Error() != "key not found" {
		config.DebugLog("Error getting followers for user %s: %v", followeeId, err)
		return err
	}
	if err != nil {
		followeeSG = SGVertex{
			UserId:    followeeId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	followeeSG.Followers = append(followeeSG.Followers, followerId)
	err = state.SetState(ctx, followeeId, followeeSG)
	if err != nil {
		config.DebugLog("Error setting followers for user %s: %v", followeeId, err)
		return err
	}
	config.DebugLog("User %s followed user %s", followerId, followeeId)
	return nil
}

// func (s *GraphServer) InsertUser(ctx context.Context, req *socialpb.InsertUserRequest) (*socialpb.InsertUserResponse, error) {
// 	if err := s.insertUser(ctx, req.UserId); err != nil {
// 		return nil, err
// 	}
// 	return &socialpb.InsertUserResponse{}, nil
// }

// func (s *GraphServer) GetFollowers(ctx context.Context, req *socialpb.GetFollowersRequest) (*socialpb.GetFollowersResponse, error) {
// 	followers, err := s.getFollowers(ctx, req.UserId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &socialpb.GetFollowersResponse{Followers: followers}, nil
// }

// func (s *GraphServer) GetFollowees(ctx context.Context, req *socialpb.GetFolloweesRequest) (*socialpb.GetFolloweesResponse, error) {
// 	followees, err := s.getFollowees(ctx, req.UserId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &socialpb.GetFolloweesResponse{Followees: followees}, nil
// }

// func (s *GraphServer) Follow(ctx context.Context, req *socialpb.FollowRequest) (*socialpb.FollowResponse, error) {
// 	if err := s.follow(ctx, req.FollowerId, req.FolloweeId); err != nil {
// 		return nil, err
// 	}
// 	return &socialpb.FollowResponse{}, nil
// }

// func (s *GraphServer) FollowMany(ctx context.Context, req *socialpb.FollowManyRequest) (*socialpb.FollowManyResponse, error) {
// 	if err := s.followMany(ctx, req.UserId, req.FollowerIds, req.FolloweeIds); err != nil {
// 		return nil, err
// 	}
// 	return &socialpb.FollowManyResponse{}, nil
// }

// func (s *GraphServer) getFollowers(ctx context.Context, userId string) ([]string, error) {
// 	sg, err := state.GetState[SGVertex](ctx, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sg.Followers, nil
// }

// func (s *GraphServer) getFollowees(ctx context.Context, userId string) ([]string, error) {
// 	sg, err := state.GetState[SGVertex](ctx, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sg.Followees, nil
// }

// func (s *GraphServer) follow(ctx context.Context, followerId string, followeeId string) error {
// 	// Retrieve and update the follower's followees list
// 	followerSG, err := state.GetState[SGVertex](ctx, followerId)
// 	if err != nil && err.Error() != "key not found" {
// 		return err
// 	}
// 	if err != nil {
// 		followerSG = SGVertex{
// 			UserId:    followerId,
// 			Followers: []string{},
// 			Followees: []string{},
// 		}
// 	}
// 	followerSG.Followees = append(followerSG.Followees, followeeId)
// 	if err := state.SetState(ctx, followerId, followerSG); err != nil {
// 		return err
// 	}

// 	// Retrieve and update the followee's followers list
// 	followeeSG, err := state.GetState[SGVertex](ctx, followeeId)
// 	if err != nil && err.Error() != "key not found" {
// 		return err
// 	}
// 	if err != nil {
// 		followeeSG = SGVertex{
// 			UserId:    followeeId,
// 			Followers: []string{},
// 			Followees: []string{},
// 		}
// 	}
// 	followeeSG.Followers = append(followeeSG.Followers, followerId)
// 	if err := state.SetState(ctx, followeeId, followeeSG); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *GraphServer) followMany(ctx context.Context, userId string, followerIds []string, followeeIds []string) error {
// 	sg := SGVertex{
// 		UserId:    userId,
// 		Followers: followerIds,
// 		Followees: followeeIds,
// 	}
// 	if len(sg.Followers) >= 10 {
// 		sg.Followers = sg.Followers[:10]
// 	}
// 	if len(sg.Followees) >= 10 {
// 		sg.Followees = sg.Followees[:10]
// 	}
// 	return state.SetState(ctx, userId, sg)
// }

// func (s *GraphServer) insertUser(ctx context.Context, userId string) error {
// 	sg := SGVertex{
// 		UserId:    userId,
// 		Followers: []string{},
// 		Followees: []string{},
// 	}
// 	return state.SetState(ctx, userId, sg)
// }
