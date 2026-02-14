package goroutinecost

import "sync"

// SpawnUnbuffered creates a goroutine and synchronizes via an unbuffered channel.
// The sender blocks until the receiver is ready (rendezvous synchronization).
func SpawnUnbuffered() {
	ch := make(chan struct{})
	go func() { ch <- struct{}{} }()
	<-ch
}

// SpawnBuffered creates a goroutine and synchronizes via a buffered channel (cap=1).
// The sender does not block on send, reducing scheduler intervention.
func SpawnBuffered() {
	ch := make(chan struct{}, 1)
	go func() { ch <- struct{}{} }()
	<-ch
}

// SpawnWaitGroup creates a goroutine and synchronizes via sync.WaitGroup.
// No channel allocation; uses atomic counter internally.
func SpawnWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { wg.Done() }()
	wg.Wait()
}
