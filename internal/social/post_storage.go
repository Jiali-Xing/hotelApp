package social

import (
	"context"
	"time"

	"github.com/Jiali-Xing/hotelApp/pkg/state"
	"github.com/lithammer/shortuuid"

	socialpb "github.com/Jiali-Xing/socialproto"
)

type PostStorageServer struct {
	socialpb.UnimplementedPostStorageServer
}

func (s *PostStorageServer) StorePost(ctx context.Context, req *socialpb.StorePostRequest) (*socialpb.StorePostResponse, error) {
	postId := s.storePost(ctx, req.CreatorId, req.Text)
	return &socialpb.StorePostResponse{PostId: postId}, nil
}

func (s *PostStorageServer) StorePostMulti(ctx context.Context, req *socialpb.StorePostMultiRequest) (*socialpb.StorePostMultiResponse, error) {
	postIds := s.storePostMulti(ctx, req.CreatorId, req.Text, int(req.Number))
	return &socialpb.StorePostMultiResponse{PostIds: postIds}, nil
}

func (s *PostStorageServer) ReadPost(ctx context.Context, req *socialpb.ReadPostRequest) (*socialpb.ReadPostResponse, error) {
	post, err := state.GetState[socialpb.Post](ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	return &socialpb.ReadPostResponse{Post: &post}, nil
}

func (s *PostStorageServer) ReadPosts(ctx context.Context, req *socialpb.ReadPostsRequest) (*socialpb.ReadPostsResponse, error) {
	retPosts, err := state.GetBulkState[socialpb.Post](ctx, req.PostIds)
	if err != nil {
		return nil, err
	}
	posts := make([]*socialpb.Post, len(retPosts))
	for i, post := range retPosts {
		posts[i] = &post
	}
	return &socialpb.ReadPostsResponse{Posts: posts}, nil
}

func (s *PostStorageServer) storePost(ctx context.Context, creatorId string, text string) string {
	postIds := s.storePostMulti(ctx, creatorId, text, 1)
	return postIds[0]
}

func (s *PostStorageServer) storePostMulti(ctx context.Context, creatorId string, text string, number int) []string {
	posts := make(map[string]interface{}, number)
	postIds := make([]string, number)
	for i := 0; i < number; i++ {
		postId := shortuuid.New()
		timestamp := time.Now().Unix()
		posts[postId] = socialpb.Post{
			PostId:    postId,
			CreatorId: creatorId,
			Text:      text,
			Timestamp: timestamp,
		}
		postIds[i] = postId
	}
	state.SetBulkState(ctx, posts)
	return postIds
}
