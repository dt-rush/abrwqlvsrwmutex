package abrwqlvsrwmutex

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkRWMutex_Contention(b *testing.B) {
	l := sync.RWMutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < 512; i++ {
		wg.Add(1)
		go func() {
			l.Lock()
			time.Sleep(time.Millisecond)
			l.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}
