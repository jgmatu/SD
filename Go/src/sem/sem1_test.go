package sem

import (
	"fmt"
	"sync"
	"testing"
)

const (
	MAXCOUNTS = 10000
)

func TestSemMut(t *testing.T) {
	ctest := 0
	wg := sync.WaitGroup{}
	s := NewSem(1)

	for i := 0; i < MAXCOUNTS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Down()
			ctest++
			s.Up()
		}()
	}
	wg.Wait()
	if ctest != MAXCOUNTS {
		fmt.Println("Value counter : ", ctest)
		t.Error()
	}
	fmt.Println("Value of ctest : ", ctest)
}
