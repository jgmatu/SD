package sem

import (
	"fmt"
	"sync"
	"testing"
)

var wg = sync.WaitGroup{}
var mutexesc = NewSem(1)
var mutexnl = NewSem(1)
var torn = NewSem(1)
var nl = 0
var ctest = 0

func writer() {
	defer wg.Done()

	torn.Down()
	mutexesc.Down()
	torn.Up()
	// Region critica
	ctest++
	mutexesc.Up()
}

func reader() {
	defer wg.Done()

	torn.Down()
	torn.Up()

	mutexnl.Down()
	nl++
	if nl == 1 {
		mutexesc.Down()
	}
	mutexnl.Up()
	// Region critica
	fmt.Println("Ctest is : ", ctest)
	mutexnl.Down()

	nl--
	if nl == 0 {
		mutexesc.Up()
	}
	mutexnl.Up()
}

func TestSemTorn(t *testing.T) {
	maxCount := 100
	for i := 0; i < maxCount; i++ {
		wg.Add(1)
		go writer()
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go reader()
		}
	}
	wg.Wait()
	if ctest != maxCount {
		t.Error()
	}
}
