package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	"github.com/Jiali-Xing/hotelApp/pkg/state"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type UserTimelineServer struct {
	socialpb.UnimplementedUserTimelineServer
}

func (s *UserTimelineServer) ReadUserTimeline(ctx context.Context, req *socialpb.ReadUserTimelineRequest) (*socialpb.ReadUserTimelineResponse, error) {
	postIds, err := state.GetState[[]string](ctx, req.UserId)
	if err != nil {
		return &socialpb.ReadUserTimelineResponse{Posts: []*socialpb.Post{}}, nil
	}

	postsReq := &socialpb.ReadPostsRequest{PostIds: postIds}
	postsResp, err := invoke.Invoke[*socialpb.ReadPostsResponse](ctx, "poststorage", "readposts", postsReq)
	if err != nil {
		return nil, err
	}

	return &socialpb.ReadUserTimelineResponse{Posts: postsResp.Posts}, nil
}

func (s *UserTimelineServer) WriteUserTimeline(ctx context.Context, req *socialpb.WriteUserTimelineRequest) (*socialpb.WriteUserTimelineResponse, error) {
	postIds, err := state.GetState[[]string](ctx, req.UserId)
	if err != nil {
		postIds = []string{}
	}
	if len(postIds) >= 10 {
		postIds = postIds[1:]
	}
	postIds = append(postIds, req.PostIds...)
	err = state.SetState(ctx, req.UserId, postIds)
	if err != nil {
		return nil, err
	}
	return &socialpb.WriteUserTimelineResponse{Success: true}, nil
}
