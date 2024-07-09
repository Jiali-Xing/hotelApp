package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	socialpb "github.com/Jiali-Xing/socialproto"
)

var (
	composePostAddr  string
	homeTimelineAddr string
	userTimelineAddr string
	socialGraphAddr  string
	numOfUsers       int
	numOfPosts       int
	numOfFollowers   int
)

func init() {
	flag.StringVar(&composePostAddr, "compose_post", "localhost:50062", "Address of the compose post service")
	flag.StringVar(&homeTimelineAddr, "home_timeline", "localhost:50059", "Address of the home timeline service")
	flag.StringVar(&userTimelineAddr, "user_timeline", "localhost:50058", "Address of the user timeline service")
	flag.StringVar(&socialGraphAddr, "social_graph", "localhost:50061", "Address of the social graph service")
	flag.IntVar(&numOfUsers, "num_of_users", 1000, "Number of users to create")
	flag.IntVar(&numOfPosts, "num_of_posts", 10, "Number of posts to create per user")
	flag.IntVar(&numOfFollowers, "num_of_followers", 10, "Number of followers per user")
	flag.Parse()
}

func getRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func main() {
	// Connect to the social graph and compose post services
	socialGraphConn, err := grpc.Dial(socialGraphAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to social graph gRPC server: %v", err)
	}
	defer socialGraphConn.Close()

	composePostConn, err := grpc.Dial(composePostAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to compose post gRPC server: %v", err)
	}
	defer composePostConn.Close()

	homeTimelineConn, err := grpc.Dial(homeTimelineAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to home timeline gRPC server: %v", err)
	}
	defer homeTimelineConn.Close()

	userTimelineConn, err := grpc.Dial(userTimelineAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user timeline gRPC server: %v", err)
	}
	defer userTimelineConn.Close()

	// Populate the services with data
	populateUsersAndFollows(socialGraphConn)
	populatePosts(composePostConn, homeTimelineConn, userTimelineConn)
}

func populateUsersAndFollows(conn *grpc.ClientConn) {
	client := socialpb.NewSocialGraphClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i := 0; i < numOfUsers; i++ {
		userId := fmt.Sprintf("user%d", i)

		// Insert user
		_, err := client.InsertUser(ctx, &socialpb.InsertUserRequest{UserId: userId})
		if err != nil {
			log.Printf("Failed to insert user %s: %v", userId, err)
		} else {
			log.Printf("Inserted user: %s", userId)
		}
	}

	for i := 0; i < numOfUsers; i++ {
		// Follow other users
		userId := fmt.Sprintf("user%d", i)
		// if i+1 < numOfUsers {
		// each user follows the last numOfFollowers users
		if i > numOfFollowers {
			for j := 0; j < numOfFollowers; j++ {
				followeeId := fmt.Sprintf("user%d", i-j)
				_, err := client.Follow(ctx, &socialpb.FollowRequest{
					FollowerId: userId,
					FolloweeId: followeeId,
				})
				if err != nil {
					log.Printf("Failed to follow user %s to user %s: %v", userId, followeeId, err)
				} else {
					log.Printf("User %s followed user %s", userId, followeeId)
				}
			}
		} else {
			for j := 0; j < i; j++ {
				followeeId := fmt.Sprintf("user%d", i-j)
				_, err := client.Follow(ctx, &socialpb.FollowRequest{
					FollowerId: userId,
					FolloweeId: followeeId,
				})
				if err != nil {
					log.Printf("Failed to follow user %s to user %s: %v", userId, followeeId, err)
				} else {
					log.Printf("User %s followed user %s", userId, followeeId)
				}
			}
		}
	}
}

func populatePosts(composePostConn, homeTimelineConn, userTimelineConn *grpc.ClientConn) {
	composeClient := socialpb.NewComposePostClient(composePostConn)
	homeClient := socialpb.NewHomeTimelineClient(homeTimelineConn)
	userClient := socialpb.NewUserTimelineClient(userTimelineConn)

	for i := 0; i < numOfUsers; i++ {
		userId := fmt.Sprintf("user%d", i)

		for j := 0; j < numOfPosts; j++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			text := getRandomString(100)
			req := &socialpb.ComposePostRequest{
				CreatorId: userId,
				Text:      text,
			}

			// Compose post
			resp, err := composeClient.ComposePost(ctx, req)
			if err != nil {
				log.Printf("Failed to compose post for user %s: %v", userId, err)
				continue
			}

			postId := resp.PostId
			log.Printf("Composed post for user %s with post ID: %s", userId, postId)

			// Read from home timeline
			homeReq := &socialpb.ReadHomeTimelineRequest{UserId: userId}
			homeResp, err := homeClient.ReadHomeTimeline(ctx, homeReq)
			if err != nil {
				log.Printf("Failed to read home timeline for user %s: %v", userId, err)
			} else {
				log.Printf("Home timeline for user %s: %v", userId, homeResp.Posts)
			}

			// Read from user timeline
			userReq := &socialpb.ReadUserTimelineRequest{UserId: userId}
			userResp, err := userClient.ReadUserTimeline(ctx, userReq)
			if err != nil {
				log.Printf("Failed to read user timeline for user %s: %v", userId, err)
			} else {
				log.Printf("User timeline for user %s: %v", userId, userResp.Posts)
			}
		}
	}
}
