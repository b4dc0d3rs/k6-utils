package k6utils

import (
	"time"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/k6utils", new(K6Utils))
}

type K6Utils struct {
	csvRecords [][]string
	header []string
}

func (c *K6Utils) SleepMilliseconds(sleepMilliseconds int) {
	time.Sleep(time.Duration(sleepMilliseconds) * time.Millisecond)
}
