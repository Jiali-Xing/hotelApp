package social

import (
	"context"
	"log"

	"github.com/Jiali-Xing/hotelApp/internal/config"
	"github.com/Jiali-Xing/hotelApp/pkg/invoke"
	socialpb "github.com/Jiali-Xing/socialproto"
)

type NginxServer struct {
	socialpb.UnimplementedNginxServiceServer
}

func (s *NginxServer) ComposePost(ctx context.Context, req *socialpb.ComposePostRequest) (*socialpb.ComposePostResponse, error) {
	ctx = config.PropagateMetadata(ctx, "nginx")
	resp, err := invoke.Invoke[*socialpb.ComposePostResponse](ctx, "compose", "ComposePost", req)
	if err != nil {
		log.Printf("Error forwarding compose post request: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *NginxServer) ReadUserTimeline(ctx context.Context, req *socialpb.ReadUserTimelineRequest) (*socialpb.ReadUserTimelineResponse, error) {
	ctx = config.PropagateMetadata(ctx, "nginx")
	resp, err := invoke.Invoke[*socialpb.ReadUserTimelineResponse](ctx, "usertimeline", "ReadUserTimeline", req)
	if err != nil {
		log.Printf("Error forwarding read user timeline request: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *NginxServer) ReadHomeTimeline(ctx context.Context, req *socialpb.ReadHomeTimelineRequest) (*socialpb.ReadHomeTimelineResponse, error) {
	ctx = config.PropagateMetadata(ctx, "nginx")
	resp, err := invoke.Invoke[*socialpb.ReadHomeTimelineResponse](ctx, "hometimeline", "ReadHomeTimeline", req)
	if err != nil {
		log.Printf("Error forwarding read home timeline request: %v", err)
		return nil, err
	}
	return resp, nil
}
