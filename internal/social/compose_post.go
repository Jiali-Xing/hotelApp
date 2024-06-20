package social

import (
	"context"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type ComposePostServer struct {
	socialpb.UnimplementedComposePostServer
}

func (s *ComposePostServer) ComposePost(ctx context.Context, req *socialpb.ComposePostRequest) (*socialpb.ComposePostResponse, error) {
	ctx = config.PropagateMetadata(ctx, "composepost")
	// Invoke store_post method in poststorage service
	req1 := &socialpb.StorePostRequest{
		CreatorId: req.CreatorId,
		Text:      req.Text,
	}
	config.DebugLog("Composing post: %+v", req1)
	resp1, err := invoke.Invoke[*socialpb.StorePostResponse](ctx, "poststorage", "StorePost", req1)
	if err != nil {
		config.DebugLog("Error storing post: %v", err)
		return nil, err
	}
	config.DebugLog("Post stored successfully: %+v", resp1)

	postId := resp1.PostId

	// Write to user timeline
	req2 := &socialpb.WriteUserTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: []string{postId},
	}
	config.DebugLog("Writing to user timeline: %+v", req2)
	_, err = invoke.Invoke[*socialpb.WriteUserTimelineResponse](ctx, "usertimeline", "WriteUserTimeline", req2)
	if err != nil {
		config.DebugLog("Error writing to user timeline: %v", err)
		return nil, err
	}
	config.DebugLog("User timeline updated successfully")

	// Write to home timeline
	req3 := &socialpb.WriteHomeTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: []string{postId},
	}
	config.DebugLog("Writing to home timeline: %+v", req3)
	_, err = invoke.Invoke[*socialpb.WriteHomeTimelineResponse](ctx, "hometimeline", "WriteHomeTimeline", req3)
	if err != nil {
		config.DebugLog("Error writing to home timeline: %v", err)
		return nil, err
	}
	config.DebugLog("Home timeline updated successfully")

	return &socialpb.ComposePostResponse{PostId: postId}, nil
}

func (s *ComposePostServer) ComposePostMulti(ctx context.Context, req *socialpb.ComposePostMultiRequest) (*socialpb.ComposePostMultiResponse, error) {
	ctx = config.PropagateMetadata(ctx, "composepost")
	// Invoke store_post_multi method in poststorage service
	req1 := &socialpb.StorePostMultiRequest{
		CreatorId: req.CreatorId,
		Text:      req.Text,
		Number:    req.Number,
	}
	config.DebugLog("Composing multiple posts: %+v", req1)
	resp1, err := invoke.Invoke[*socialpb.StorePostMultiResponse](ctx, "poststorage", "StorePostMulti", req1)
	if err != nil {
		config.DebugLog("Error storing multiple posts: %v", err)
		return nil, err
	}
	config.DebugLog("Multiple posts stored successfully: %+v", resp1)

	postIds := resp1.PostIds

	// Write to user timeline
	req2 := &socialpb.WriteUserTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: postIds,
	}
	config.DebugLog("Writing to user timeline: %+v", req2)
	_, err = invoke.Invoke[*socialpb.WriteUserTimelineResponse](ctx, "usertimeline", "WriteUserTimeline", req2)
	if err != nil {
		config.DebugLog("Error writing to user timeline: %v", err)
		return nil, err
	}
	config.DebugLog("User timeline updated successfully")

	// Write to home timeline
	req3 := &socialpb.WriteHomeTimelineRequest{
		UserId:  req.CreatorId,
		PostIds: postIds,
	}
	config.DebugLog("Writing to home timeline: %+v", req3)
	_, err = invoke.Invoke[*socialpb.WriteHomeTimelineResponse](ctx, "hometimeline", "WriteHomeTimeline", req3)
	if err != nil {
		config.DebugLog("Error writing to home timeline: %v", err)
		return nil, err
	}
	config.DebugLog("Home timeline updated successfully")

	return &socialpb.ComposePostMultiResponse{PostIds: postIds}, nil
}
