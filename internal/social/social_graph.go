package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/pkg/state"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type socialGraphServer struct {
	socialpb.UnimplementedSocialGraphServer
}

func (s *socialGraphServer) InsertUser(ctx context.Context, req *socialpb.InsertUserRequest) (*socialpb.InsertUserResponse, error) {
	s.insertUser(ctx, req.UserId)
	return &socialpb.InsertUserResponse{}, nil
}

func (s *socialGraphServer) GetFollowers(ctx context.Context, req *socialpb.GetFollowersRequest) (*socialpb.GetFollowersResponse, error) {
	followers := s.getFollowers(ctx, req.UserId)
	return &socialpb.GetFollowersResponse{Followers: followers}, nil
}

func (s *socialGraphServer) GetFollowees(ctx context.Context, req *socialpb.GetFolloweesRequest) (*socialpb.GetFolloweesResponse, error) {
	followees := s.getFollowees(ctx, req.UserId)
	return &socialpb.GetFolloweesResponse{Followees: followees}, nil
}

func (s *socialGraphServer) Follow(ctx context.Context, req *socialpb.FollowRequest) (*socialpb.FollowResponse, error) {
	s.follow(ctx, req.FollowerId, req.FolloweeId)
	return &socialpb.FollowResponse{}, nil
}

func (s *socialGraphServer) FollowMany(ctx context.Context, req *socialpb.FollowManyRequest) (*socialpb.FollowManyResponse, error) {
	s.followMany(ctx, req.UserId, req.FollowerIds, req.FolloweeIds)
	return &socialpb.FollowManyResponse{}, nil
}

func (s *socialGraphServer) getFollowers(ctx context.Context, userId string) []string {
	sg, err := state.GetState[SGVertex](ctx, userId)
	if err != nil {
		panic(err)
	}
	return sg.Followers
}

func (s *socialGraphServer) getFollowees(ctx context.Context, userId string) []string {
	sg, err := state.GetState[SGVertex](ctx, userId)
	if err != nil {
		panic(err)
	}
	return sg.Followees
}

func (s *socialGraphServer) follow(ctx context.Context, followerId string, followeeId string) {
	sg, err := state.GetState[SGVertex](ctx, followerId)
	if err != nil {
		sg = SGVertex{
			UserId:    followerId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	sg.Followees = append(sg.Followees, followeeId)
	err = state.SetState(ctx, followerId, sg)
	if err != nil {
		panic(err)
	}

	sg, err = state.GetState[SGVertex](ctx, followeeId)
	if err != nil {
		sg = SGVertex{
			UserId:    followeeId,
			Followers: []string{},
			Followees: []string{},
		}
	}
	sg.Followers = append(sg.Followers, followerId)
	err = state.SetState(ctx, followeeId, sg)
	if err != nil {
		panic(err)
	}
}

func (s *socialGraphServer) followMany(ctx context.Context, userId string, followerIds []string, followeeIds []string) {
	sg := SGVertex{
		UserId:    userId,
		Followers: followerIds,
		Followees: followeeIds,
	}
	if len(sg.Followers) >= 10 {
		sg.Followers = sg.Followers[:10]
	}
	if len(sg.Followees) >= 10 {
		sg.Followees = sg.Followees[:10]
	}
	err := state.SetState(ctx, userId, sg)
	if err != nil {
		panic(err)
	}
}

func (s *socialGraphServer) insertUser(ctx context.Context, userId string) {
	sg := SGVertex{
		Followers: []string{},
		Followees: []string{},
	}
	err := state.SetState(ctx, userId, sg)
	if err != nil {
		panic(err)
	}
}
