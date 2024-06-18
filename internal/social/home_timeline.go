package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type HomeTimelineServer struct {
	socialpb.UnimplementedHomeTimelineServer
}

func (s *HomeTimelineServer) ReadHomeTimeline(ctx context.Context, req *socialpb.ReadHomeTimelineRequest) (*socialpb.ReadHomeTimelineResponse, error) {
	ctx = config.PropagateMetadata(ctx, "hometimeline")
	postIds, err := state.GetState[[]string](ctx, req.UserId)
	if err != nil {
		return &socialpb.ReadHomeTimelineResponse{Posts: []*socialpb.Post{}}, nil
	}

	postsReq := &socialpb.ReadPostsRequest{PostIds: postIds}
	postsResp, err := invoke.Invoke[*socialpb.ReadPostsResponse](ctx, "poststorage", "readposts", postsReq)
	if err != nil {
		return nil, err
	}

	return &socialpb.ReadHomeTimelineResponse{Posts: postsResp.Posts}, nil
}

func (s *HomeTimelineServer) WriteHomeTimeline(ctx context.Context, req *socialpb.WriteHomeTimelineRequest) (*socialpb.WriteHomeTimelineResponse, error) {
	ctx = config.PropagateMetadata(ctx, "hometimeline")
	followersReq := &socialpb.GetFollowersRequest{UserId: req.UserId}
	followersResp, err := invoke.Invoke[*socialpb.GetFollowersResponse](ctx, "socialgraph", "getfollowers", followersReq)
	if err != nil {
		return nil, err
	}

	for _, follower := range followersResp.Followers {
		postIds, err := state.GetState[[]string](ctx, follower)
		if err != nil {
			postIds = []string{}
		}
		if len(postIds) >= 10 {
			postIds = postIds[1:]
		}
		postIds = append(postIds, req.PostIds...)
		err = state.SetState(ctx, follower, postIds)
		if err != nil {
			return nil, err
		}
	}

	return &socialpb.WriteHomeTimelineResponse{Success: true}, nil
}
