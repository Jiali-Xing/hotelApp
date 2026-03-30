package invoke

import (
	"fmt"
	"sync"

	"github.com/Jiali-Xing/hotelApp/internal/config"

	"google.golang.org/grpc"
)

var (
	conn     *grpc.ClientConn
	clients  = make(map[string]interface{})
	mu       sync.Mutex
	initOnce sync.Once
)

func getClientConn(address string) (*grpc.ClientConn, error) {
	var err error
	initOnce.Do(func() {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		config.DebugLog("Established connection to gRPC server at address: %s", address)
	})
	config.DebugLog("Returning existing connection to gRPC server at address: %s", address)
	return conn, err
}

func registerClient(service string, client interface{}) {
	mu.Lock()
	defer mu.Unlock()
	clients[service] = client
	config.DebugLog("Registered client for service: %s", service)
}

func getClient(service string) (interface{}, error) {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[service]
	if !exists {
		return nil, fmt.Errorf("client not registered for service: %s", service)
	}
	config.DebugLog("Retrieved client for service: %s", service)
	return client, nil
}

func RegisterClient(service string, client interface{}) {
	registerClient(service, client)
}

func CloseConnection() {
	mu.Lock()
	defer mu.Unlock()
	if conn != nil {
		conn.Close()
		config.DebugLog("Closed connection to gRPC server")
	}
}
