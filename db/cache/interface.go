package cache

import "time"

type Cache interface {
	Load() error
	GetRefreshGap() time.Duration
}
