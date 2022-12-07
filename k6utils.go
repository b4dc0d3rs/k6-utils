package k6utils

import (
	"time"

	"github.com/patrickmn/go-cache"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/k6utils", new(K6Utils))
}

type K6Utils struct {
	csvRecords [][]string
	header []string
	cache *cache.Cache
}

func (c *K6Utils) SleepMilliseconds(sleepMilliseconds int) {
	time.Sleep(time.Duration(sleepMilliseconds) * time.Millisecond)
}
