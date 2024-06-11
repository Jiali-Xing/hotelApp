package hotel

import (
	"bytes"
	"context"
	"crypto/sha256"
	"log"

	"github.com/Jiali-Xing/hotelApp/internal/config"

	"github.com/Jiali-Xing/hotelApp/pkg/state"

	hotelpb "github.com/Jiali-Xing/hotelproto"
	"github.com/lithammer/shortuuid"
)

type UserServer struct {
	hotelpb.UnimplementedUserServiceServer
}

func (s *UserServer) RegisterUser(ctx context.Context, req *hotelpb.RegisterUserRequest) (*hotelpb.RegisterUserResponse, error) {
	ctx = propagateMetadata(ctx, "user")
	ok := RegisterUser(ctx, req.Username, req.Password)
	resp := &hotelpb.RegisterUserResponse{Ok: ok}
	return resp, nil
}

func (s *UserServer) Login(ctx context.Context, req *hotelpb.LoginRequest) (*hotelpb.LoginResponse, error) {
	ctx = propagateMetadata(ctx, "user")
	token := Login(ctx, req.Username, req.Password)
	resp := &hotelpb.LoginResponse{Token: token}
	return resp, nil
}

func RegisterUser(ctx context.Context, username string, password string) bool {
	userId := shortuuid.New()
	salt := shortuuid.New()
	hashPass := hash(password + salt)
	user := User{
		UserId:   userId,
		Username: username,
		Password: hashPass,
		Salt:     salt,
	}

	state.SetState(ctx, username, user)
	config.DebugLog("User %s registered successfully with userId: %s and password: %s", username, userId, password)
	return true
}

func hash(str string) []byte {
	h := sha256.New()
	h.Write([]byte(str))
	val := h.Sum(nil)
	return val
}

func Login(ctx context.Context, username string, password string) string {
	config.DebugLog("Login called with username: %s and password: %s", username, password)
	user, err := state.GetState[User](ctx, username)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return "NOT-OK"
	}
	salt := user.Salt
	givenPass := hash(password + salt)
	if bytes.Equal(givenPass, user.Password) {
		return "OK"
	}
	return "NOT-OK"
}

func GetUserId(ctx context.Context, username string) string {
	user, err := state.GetState[User](ctx, username)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return ""
	}
	return user.UserId
}
