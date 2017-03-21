package sem

import (
	"fmt"
	"sync"
	"testing"
)

type DownUp struct {
	c   int
	mut sync.Mutex
}

func (du DownUp) Up() {
	du.mut.Lock()
	du.c++
	du.mut.Unlock()
}

func (du DownUp) Down() {
	du.mut.Lock()
	du.c--
	du.mut.Unlock()
}

func TestSemInt(t *testing.T) {
	var duI UpDowner
	du := DownUp{0, sync.Mutex{}}
	duI = du // if compile Ok interface successfull
	duI.Up()
	duI.Up()
	duI.Up()
	duI.Down()
	fmt.Println(du.c)
}
