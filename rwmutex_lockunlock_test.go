package abrwqlvsrwmutex

import (
	"sync"
	"testing"
)

func BenchmarkRWMutex_LockUnlock(b *testing.B) {
	l := sync.RWMutex{}
	for i := 0; i < 4096; i++ {
		l.Lock()
		l.Unlock()
	}
}
