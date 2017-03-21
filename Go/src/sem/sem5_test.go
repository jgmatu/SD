package sem

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	MAXVAL = 6
	MAXGO  = 10 * MAXVAL
)

var sem = NewSem(MAXVAL)
var count = 0
var wg5 = sync.WaitGroup{}
var mut = NewSem(1)

func gocount(t *testing.T) {
	defer wg5.Done()

	// As much 6 go routines at same time.
	sem.Down()

	mut.Down()
	count++
	if count > MAXVAL {
		t.Error()
	}
	fmt.Println("Count : ", count)
	mut.Up()

	// Time to wait MAXVAL go routines inside critic region.
	time.Sleep(time.Millisecond)

	mut.Down()
	count--
	mut.Up()

	sem.Up()
}

func TestSem(t *testing.T) {
	for i := 0; i < MAXGO; i++ {
		wg5.Add(1)
		go gocount(t)
	}
	wg5.Wait()

	semNeg := NewSem(-MAXVAL)
	if semNeg != nil {
		t.Error()
	}
}
