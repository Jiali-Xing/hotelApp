package social

import (
	"context"
	"log"

	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type NginxServer struct {
	socialpb.UnimplementedNginxServiceServer
}

func ComposePost(ctx context.Context, req *socialpb.ComposePostRequest) (*socialpb.ComposePostResponse, error) {
	resp, err := invoke.Invoke[*socialpb.ComposePostResponse](ctx, "compose", "ComposePost", req)
	if err != nil {
		log.Printf("Error forwarding compose post request: %v", err)
		return nil, err
	}
	return resp, nil
}

func ReadUserTimeline(ctx context.Context, req *socialpb.ReadUserTimelineRequest) (*socialpb.ReadUserTimelineResponse, error) {
	resp, err := invoke.Invoke[*socialpb.ReadUserTimelineResponse](ctx, "usertimeline", "ReadUserTimeline", req)
	if err != nil {
		log.Printf("Error forwarding read user timeline request: %v", err)
		return nil, err
	}
	return resp, nil
}

func ReadHomeTimeline(ctx context.Context, req *socialpb.ReadHomeTimelineRequest) (*socialpb.ReadHomeTimelineResponse, error) {
	resp, err := invoke.Invoke[*socialpb.ReadHomeTimelineResponse](ctx, "hometimeline", "ReadHomeTimeline", req)
	if err != nil {
		log.Printf("Error forwarding read home timeline request: %v", err)
		return nil, err
	}
	return resp, nil
}
