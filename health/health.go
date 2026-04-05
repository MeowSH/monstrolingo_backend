package health

import "context"

type StatusResponse struct {
	Status string `json:"status"`
}

// Ping exposes a basic health endpoint for startup checks.
//
//encore:api public method=GET path=/health
func Ping(ctx context.Context) (*StatusResponse, error) {
	_ = ctx
	return &StatusResponse{Status: "ok"}, nil
}
