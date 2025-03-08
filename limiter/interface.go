package limiter

import "time"

type Store interface {
	Allow(key string, limit int, duration time.Duration) (bool, error)
}
