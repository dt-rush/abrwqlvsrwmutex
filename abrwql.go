// Array-Based Read-Write Queueing Lock
// implementation of a circular array-based queueing lock which is
// safe from overflow at the expense of relatively-busy checking on a single
// atomic value (dequeueCount) in the case of overflow
package abrwqlvsrwmutex

import (
	"sync/atomic"
	"time"
)

const LOCK_QUEUE_SZ = 256

const LOCK_SLEEP = 5 * time.Microsecond

type ArrayBasedRWQueueLock struct {
	arr          []int
	ticket       uint32
	dequeueCount uint32
	nReaders     uint32
}

func NewArrayBasedRWQueueLock() *ArrayBasedRWQueueLock {
	abql := ArrayBasedRWQueueLock{
		arr:          make([]int, LOCK_QUEUE_SZ),
		ticket:       0,
		dequeueCount: 0,
		nReaders:     0}
	abql.arr[0] = 1
	return &abql
}

func (l *ArrayBasedRWQueueLock) RLock() {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(LOCK_QUEUE_SZ) {
		time.Sleep(LOCK_SLEEP)
	}
	for l.arr[ticket%uint32(LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(LOCK_SLEEP)
	}
	// increment nReaders
	atomic.AddUint32(&l.nReaders, 1)
	// move the queue forward after incrementing nReaders, so that
	// either another call to RLock() can get the queue head or else
	// a call to Lock() can get the queue head but wait for nReaders = 0
	l.arr[int(ticket)%LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
	return
}

func (l *ArrayBasedRWQueueLock) Lock() uint32 {
	ticket := atomic.AddUint32(&l.ticket, 1) - 1
	for ticket-l.dequeueCount >= uint32(LOCK_QUEUE_SZ) {
		time.Sleep(LOCK_SLEEP)
	}
	for l.arr[ticket%uint32(LOCK_QUEUE_SZ)] != 1 {
		time.Sleep(LOCK_SLEEP)
	}
	// wait here if our turn in the queue came because the prior lock was
	// an RLock releasing itself in hopes of triggering another RLock() instance
	for l.nReaders > 0 {
		time.Sleep(LOCK_SLEEP)
	}
	return ticket
}

func (l *ArrayBasedRWQueueLock) RUnlock() {
	// decrement nReaders
	atomic.AddUint32(&l.nReaders, ^uint32(0))
}

func (l *ArrayBasedRWQueueLock) Unlock(ticket uint32) {
	l.arr[int(ticket)%LOCK_QUEUE_SZ] = 0
	l.arr[int(ticket+1)%LOCK_QUEUE_SZ] = 1
	atomic.AddUint32(&l.dequeueCount, 1)
}
