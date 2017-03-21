// Testing Rendezvous...
// The values are swaped...

// Order of prints...
// Hello world!
// Hello goroutine...
// Bye routine...
// Bye world...

package rendez

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var wg = sync.WaitGroup{}

var goVal = "GoRoutine Val"
var mainVal = "Main Val"

func routine(t *testing.T) {
	defer wg.Done()

	fmt.Println("Hello goroutine")
	fmt.Println("Bye goroutine")
	coupleVal := Rendezvous(1, goVal)
	if coupleVal != mainVal {
		t.Error()
	}
}

func TestRen(t *testing.T) {

	fmt.Fprintln(os.Stdout, "Hello world!!")
	wg.Add(1)
	go routine(t)
	coupleVal := Rendezvous(1, mainVal) // Wait for a goroutine...
	fmt.Println("Bye world!!!")
	if coupleVal != goVal {
		t.Error()
	}
	wg.Wait()
}
