package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"google.golang.org/grpc"
// )

// // Generic gRPC client map to hold different service clients
// // This helps in reusing connections.
// var clients = make(map[string]interface{})

// // Connects to a gRPC service and returns the client
// func getClient(app string) (interface{}, error) {
// 	if client, exists := clients[app]; exists {
// 		return client, nil
// 	}

// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}

// 	var client interface{}
// 	switch app {
// 	case "UserService":
// 		client = pb.NewUserServiceClient(conn)
// 		// Add more cases for different services
// 	default:
// 		return nil, fmt.Errorf("unknown service: %s", app)
// 	}

// 	clients[app] = client
// 	return client, nil
// }

// // Generic function to perform a gRPC request
// func performGRPCRequest[T interface{}](ctx context.Context, client interface{}, method string, input interface{}) (T, error) {
// 	var response T

// 	buf, err := json.Marshal(input)
// 	if err != nil {
// 		return response, err
// 	}

// 	argBytes := bytes.NewBuffer(buf).Bytes()

// 	switch c := client.(type) {
// 	case pb.UserServiceClient:
// 		switch method {
// 		case "GetUser":
// 			req := &pb.UserRequest{}
// 			err = json.Unmarshal(argBytes, req)
// 			if err != nil {
// 				return response, err
// 			}
// 			res, err := c.GetUser(ctx, req)
// 			if err != nil {
// 				return response, err
// 			}
// 			responseBytes, err := json.Marshal(res)
// 			if err != nil {
// 				return response, err
// 			}
// 			err = json.Unmarshal(responseBytes, &response)
// 			if err != nil {
// 				return response, err
// 			}
// 		case "CreateUser":
// 			req := &pb.CreateUserRequest{}
// 			err = json.Unmarshal(argBytes, req)
// 			if err != nil {
// 				return response, err
// 			}
// 			res, err := c.CreateUser(ctx, req)
// 			if err != nil {
// 				return response, err
// 			}
// 			responseBytes, err := json.Marshal(res)
// 			if err != nil {
// 				return response, err
// 			}
// 			err = json.Unmarshal(responseBytes, &response)
// 			if err != nil {
// 				return response, err
// 			}
// 		default:
// 			return response, fmt.Errorf("unknown method: %s", method)
// 		}
// 	default:
// 		return response, fmt.Errorf("unknown client type")
// 	}

// 	return response, nil
// }

// func Invoke[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
// 	var res T
// 	client, err := getClient(app)
// 	if err != nil {
// 		log.Fatalf("could not get client: %v", err)
// 	}

// 	res, err = performGRPCRequest[T](ctx, client, method, input)
// 	if err != nil {
// 		log.Fatalf("could not perform gRPC request: %v", err)
// 	}

// 	return res
// }

// func main() {
// 	ctx := context.Background()
// 	input := &pb.UserRequest{UserId: "12345"}
// 	result := Invoke[pb.UserResponse](ctx, "UserService", "GetUser", input)
// 	fmt.Printf("Result: %+v\n", result)
// }
