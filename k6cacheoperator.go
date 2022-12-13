package k6utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func (k6utils *K6Utils) CreateCacheWithExpiryInSeconds(durationInSeconds int) {
	duration := time.Duration(durationInSeconds) * time.Second
	k6utils.cache = cache.New(duration, duration)
}

func (k6utils *K6Utils) PutToCache(key string, value string) {
	k6utils.cache.Set(key, value, cache.DefaultExpiration)
}

func (k6utils *K6Utils) GetFromCache(key string) interface{} {
	stringValue, _ := k6utils.cache.Get(key)
	return stringValue
}

func (k6utils *K6Utils) RemoveFromCache(key string) {
	k6utils.cache.Delete(key)
}