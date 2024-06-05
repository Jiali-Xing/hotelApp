package invoke

import (
	"fmt"
	"sync"

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
	})
	return conn, err
}

func registerClient(service string, client interface{}) {
	mu.Lock()
	defer mu.Unlock()
	clients[service] = client
}

func getClient(service string) (interface{}, error) {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[service]
	if !exists {
		return nil, fmt.Errorf("client not registered for service: %s", service)
	}
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
	}
}
