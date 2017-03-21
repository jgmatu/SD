package sem

import (
	"sync"
)

// To use a semaphore... the methods Up() and Down()
// are used to increment and decrement the semaphore
// you must implement interface UpDowner to get a semaphore.
type UpDowner interface {
	Up()
	Down()
}

type Sem struct {
	mutex  sync.Mutex
	tokens int
	cond   *sync.Cond
}

func NewSem(ntok int) *Sem {
	if ntok < 0 {
		return nil
	}
	s := new(Sem)
	s.mutex = sync.Mutex{}
	s.tokens = ntok
	s.cond = sync.NewCond(&s.mutex)
	return s
}

func (s *Sem) Up() {
	s.mutex.Lock()
	s.tokens++
	s.cond.Signal()
	s.mutex.Unlock()
}

func (s *Sem) Down() {
	s.mutex.Lock()
	for s.tokens == 0 {
		s.cond.Wait()
	}
	s.tokens--
	s.mutex.Unlock()
}
