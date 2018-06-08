package abrwqlvsrwmutex

import (
	"testing"
)

func BenchmarkABRWQL_LockUnlock(b *testing.B) {
	l := NewArrayBasedRWQueueLock()
	for i := 0; i < 4096; i++ {
		ticket := l.Lock()
		l.Unlock(ticket)
	}
}
