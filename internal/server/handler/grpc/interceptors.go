// Package grpc provides middleware and interceptors for gRPC services,
// including logging and monitoring capabilities.
package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

// GrpcInterceptorLogger creates a gRPC server interceptor that logs each call,
// outputting information about the duration of execution and the response status.
func GrpcInterceptorLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Call the handler to complete the normal execution of a unary RPC
		resp, err := handler(ctx, req)

		duration := time.Since(start)

		// Log the request details
		slog.Info("Request", slog.String("method", info.FullMethod), slog.Duration("duration", duration))

		// Get the status code of the response
		statusCode := status.Code(err)

		// Log the response details
		slog.Info("Response", slog.String("status", statusCode.String()), slog.Duration("duration", duration))

		return resp, err
	}
}
