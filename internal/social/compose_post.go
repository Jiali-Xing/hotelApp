package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type ComposePostServer struct {
	socialpb.UnimplementedComposePostServer
}

func (s *ComposePostServer) ComposePost(ctx context.Context, req *socialpb.ComposePostRequest) (*socialpb.ComposePostResponse, error) {
	// Invoke store_post method in poststorage service
	req1 := &socialpb.StorePostRequest{
		CreatorId: req.CreatorId,
		Text:      req.Text,
	}
	resp1, err := invoke.Invoke[socialpb.StorePostResponse](ctx, "poststorage", "StorePost", req1)
	if err != nil {
		return nil, err
	}

	postId := resp1.PostId

	// Write to user timeline
	req2 := &socialpb.WriteUserTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: []string{postId},
	}

	_, err = invoke.Invoke[string](ctx, "usertimeline", "writeusertimeline", req2)
	if err != nil {
		return nil, err
	}

	// Write to home timeline
	req3 := &socialpb.WriteHomeTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: []string{postId},
	}
	_, err = invoke.Invoke[string](ctx, "hometimeline", "writehometimeline", req3)
	if err != nil {
		return nil, err
	}

	return &socialpb.ComposePostResponse{PostId: postId}, nil
}

func (s *ComposePostServer) ComposePostMulti(ctx context.Context, req *socialpb.ComposePostMultiRequest) (*socialpb.ComposePostMultiResponse, error) {
	// Invoke store_post_multi method in poststorage service
	req1 := &socialpb.StorePostMultiRequest{
		CreatorId: req.CreatorId,
		Text:      req.Text,
		Number:    req.Number,
	}
	resp1, err := invoke.Invoke[socialpb.StorePostMultiResponse](ctx, "poststorage", "StorePostMulti", req1)
	if err != nil {
		return nil, err
	}

	postIds := resp1.PostIds

	// Write to user timeline
	req2 := &socialpb.WriteUserTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: postIds,
	}
	_, err = invoke.Invoke[string](ctx, "usertimeline", "WriteUserTimeline", req2)
	if err != nil {
		return nil, err
	}

	// Write to home timeline
	req3 := &socialpb.WriteHomeTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: postIds,
	}
	_, err = invoke.Invoke[string](ctx, "hometimeline", "WriteHomeTimeline", req3)
	if err != nil {
		return nil, err
	}

	return &socialpb.ComposePostMultiResponse{PostIds: postIds}, nil
}
