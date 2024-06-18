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

	// Start with the context
	outCtx := ctx
	// Check and append tokens
	if tokens, ok := headersIn["tokens"]; ok && len(tokens) == 1 {
		outCtx = metadata.AppendToOutgoingContext(outCtx, "tokens", tokens[0])
	}

	// Check and append user-specific metadata
	if u, uok := headersIn["u"]; uok && len(u) == 1 {
		outCtx = metadata.AppendToOutgoingContext(outCtx, "u", u[0], "b", headersIn["b"][0])
	}

	if userId, ok := headersIn["user-id"]; ok && len(userId) > 0 {
		outCtx = metadata.AppendToOutgoingContext(outCtx, "user-id", userId[0])
	}

	if reqId, ok := headersIn["request-id"]; ok && len(reqId) > 0 {
		outCtx = metadata.AppendToOutgoingContext(outCtx, "request-id", reqId[0], "timestamp", headersIn["timestamp"][0], "method", method, "name", serviceName)
	}

	// Append common metadata
	// outCtx = metadata.AppendToOutgoingContext(outCtx, "request-id", headersIn["request-id"][0], "timestamp", headersIn["timestamp"][0], "method", method, "name", serviceName)

	// Update the original context
	ctx = outCtx

	headersOut, _ := metadata.FromOutgoingContext(ctx)
	DebugLog("headersOut: %s", headersOut)

	return ctx
}
