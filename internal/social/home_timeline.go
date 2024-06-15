package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type homeTimelineServer struct {
	socialpb.UnimplementedHomeTimelineServer
}

func (s *homeTimelineServer) ReadHomeTimeline(ctx context.Context, req *socialpb.ReadHomeTimelineRequest) (*socialpb.ReadHomeTimelineResponse, error) {
	postIds, err := state.GetState[[]string](ctx, req.UserId)
	if err != nil {
		return &socialpb.ReadHomeTimelineResponse{Posts: []*socialpb.Post{}}, nil
	}

	postsReq := &socialpb.ReadPostsRequest{PostIds: postIds}
	postsResp, err := invoke.Invoke[socialpb.ReadPostsResponse](ctx, "poststorage", "ro_read_posts", postsReq)
	if err != nil {
		return nil, err
	}

	return &socialpb.ReadHomeTimelineResponse{Posts: postsResp.Posts}, nil
}

func (s *homeTimelineServer) WriteHomeTimeline(ctx context.Context, req *socialpb.WriteHomeTimelineRequest) (*socialpb.WriteHomeTimelineResponse, error) {
	followersReq := &socialpb.GetFollowersRequest{UserId: req.UserId}
	followersResp, err := invoke.Invoke[socialpb.GetFollowersResponse](ctx, "socialgraph", "ro_get_followers", followersReq)
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
		state.SetState(ctx, follower, postIds)
	}

	return &socialpb.WriteHomeTimelineResponse{Success: true}, nil
}
