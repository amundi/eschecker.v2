package worker

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type testWorker struct {
	i int
	sync.Mutex
}

func collectorTest(t *testWorker) {
	G_WorkQueue <- t
}

func (t *testWorker) DoRequest() {
	t.Lock()
	t.i++
	t.Unlock()
}

func TestStartDispatcher(t *testing.T) {
	StartDispatcher(32)
	test := new(testWorker)
	test.i = 0
	for i := 0; i < 32; i++ {
		collectorTest(test)
	}
	time.Sleep(1 * time.Millisecond)
	StopAllWorkers(32)
	assert.Equal(t, 32, test.i)
}
