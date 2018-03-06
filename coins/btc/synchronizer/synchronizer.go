package synchronizer

import "github.com/shiftdevices/godbb/util/locker"

// Synchronizer keeps track of a reference counter. It is useful to keep track of outstanding tasks
// that run in goroutines.
type Synchronizer struct {
	requestsCounter int32
	onSyncStarted   func()
	onSyncFinished  func()
	wait            chan struct{}
	waitLock        locker.Locker
}

// NewSynchronizer creates a new Synchronizer. onSyncStarted is called when the counter is first
// incremented. onSyncFinished is called when the counter is last decremented.
func NewSynchronizer(onSyncStarted func(), onSyncFinished func()) *Synchronizer {
	synchronizer := &Synchronizer{
		requestsCounter: 0,
		onSyncStarted:   onSyncStarted,
		onSyncFinished:  onSyncFinished,
		wait:            nil,
	}
	return synchronizer
}

// IncRequestsCounter increments the counter, and returns a function to decrement it which must be
// called after the task has finished.
func (synchronizer *Synchronizer) IncRequestsCounter() func() {
	defer synchronizer.waitLock.Lock()()
	synchronizer.requestsCounter++
	if synchronizer.requestsCounter == 1 {
		synchronizer.onSyncStarted()
		synchronizer.wait = make(chan struct{})
	}
	return synchronizer.decRequestsCounter
}

func (synchronizer *Synchronizer) decRequestsCounter() {
	defer synchronizer.waitLock.Lock()()
	synchronizer.requestsCounter--
	if synchronizer.requestsCounter == 0 {
		synchronizer.onSyncFinished()
		// Everyone waiting will be notified by this.
		close(synchronizer.wait)
		synchronizer.wait = nil
	} else if synchronizer.requestsCounter < 0 {
		panic("request counter cannot be negative")
	}
}

// WaitSynchronized blocks until all pending synchronization tasks are finished.
func (synchronizer *Synchronizer) WaitSynchronized() {
	if func() int32 {
		defer synchronizer.waitLock.RLock()()
		return synchronizer.requestsCounter
	}() == 0 {
		return
	}
	<-synchronizer.wait
}