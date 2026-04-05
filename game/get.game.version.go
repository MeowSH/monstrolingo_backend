package game

import "context"

// GetGameVersion returns the currently supported game version.
//
//encore:api public method=GET path=/game/version
func GetGameVersion(ctx context.Context) (*gameVersionResponse, error) {
	_ = ctx
	return newGameVersionResponse(), nil
}
