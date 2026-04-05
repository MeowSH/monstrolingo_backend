package game

import "monstrolingo_backend/internal/gameversion"

type gameVersionResponse struct {
	GameVersion string `json:"game_version"`
}

func newGameVersionResponse() *gameVersionResponse {
	return &gameVersionResponse{
		GameVersion: gameversion.Value(),
	}
}
