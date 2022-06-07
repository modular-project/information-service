package middleware

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(i interface{}) error {
	return status.Errorf(codes.Unknown, "panic triggered: %v", i)
}
