package abrwqlvsrwmutex

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkABRWQL_Contention(b *testing.B) {
	l := NewArrayBasedRWQueueLock()
	wg := sync.WaitGroup{}
	for i := 0; i < 512; i++ {
		wg.Add(1)
		go func() {
			ticket := l.Lock()
			time.Sleep(time.Millisecond)
			l.Unlock(ticket)
			wg.Done()
		}()
	}
	wg.Wait()
}
