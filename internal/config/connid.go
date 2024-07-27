package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

// Extract yyy from conn-id metadata
func getConnectionID(ctx context.Context) (uint32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("missing metadata")
	}

	connIDs := md.Get("conn-id")
	if len(connIDs) == 0 {
		return 0, fmt.Errorf("missing conn-id")
	}

	connID := connIDs[0]
	parts := strings.Split(connID, "c")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid conn-id format")
	}

	conn, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid conn-id number: %v", err)
	}

	return uint32(conn), nil
}

// Helper function to generate a user and password based on yyy
func GenerateUserAndPassword(ctx context.Context) (string, string, error) {
	connID, err := getConnectionID(ctx)
	if err != nil {
		return "", "", err
	}

	username := fmt.Sprintf("user%d", connID)
	password := fmt.Sprintf("password%d", connID)
	return username, password, nil
}
