package cache

import "time"

type Memory struct {
	caches []Cache
}

func NewMemory(caches []Cache) *Memory {
	return &Memory{caches: caches}
}

func (receiver *Memory) Refresh() {
	for _, cache := range receiver.caches {
		err := cache.Load()
		if err != nil {
			panic(err)
		}
		go func(c Cache) {
			for {
				time.Sleep(c.GetRefreshGap())
				err := c.Load()
				if err != nil {
					panic(err)
				}
			}
		}(cache)
	}
}
