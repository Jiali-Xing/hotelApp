package config

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
)

// propagateMetadata extracts metadata from the incoming context and appends it to the outgoing context.
func PropagateMetadata(ctx context.Context, serviceName string) context.Context {
	headersIn, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No incoming metadata found")
		return ctx
	}

	method := ""
	if methods, found := headersIn["method"]; found && len(methods) > 0 {
		method = methods[0]
	}
	DebugLog("headersIn: %s", headersIn)

	// Append the request-id and timestamp of headersIn to the outgoing context
	if tokens, ok := headersIn["tokens"]; ok && len(tokens) == 1 {
		ctx = metadata.AppendToOutgoingContext(ctx, "tokens", tokens[0], "request-id", headersIn["request-id"][0], "timestamp", headersIn["timestamp"][0], "method", method, "name", serviceName)
	} else if u, uok := headersIn["u"]; uok && len(u) == 1 {
		ctx = metadata.AppendToOutgoingContext(ctx, "u", u[0], "b", headersIn["b"][0], "request-id", headersIn["request-id"][0], "timestamp", headersIn["timestamp"][0], "method", method, "name", serviceName)
	} else {
		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", headersIn["request-id"][0], "timestamp", headersIn["timestamp"][0], "method", method, "name", serviceName)
	}

	headersOut, _ := metadata.FromOutgoingContext(ctx)
	DebugLog("headersOut: %s", headersOut)

	return ctx
}
