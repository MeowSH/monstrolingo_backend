package appenv

import (
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once

// Load reads .env once when present and keeps existing process env precedence.
func Load() {
	loadOnce.Do(func() {
		_ = godotenv.Load(".env")
	})
}
