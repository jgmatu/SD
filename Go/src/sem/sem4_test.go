// We test the semaphore with a consumer
// and producter...
package sem

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	N          = 2
	MAXTICKETS = 10
)

var ticketSem = NewSem(0)
var holeSem = NewSem(N)
var buf = make([]*string, N) // Buffer to consumers and producters...
var wg1 = sync.WaitGroup{}
var i, j int
var mutexI = NewSem(1)
var mutexJ = NewSem(1)

func producter(pos int) {
	defer wg1.Done()
	for nticket := 0; nticket < MAXTICKETS; nticket++ {
		ticket := "Ticket... " + strconv.Itoa(nticket)
		holeSem.Down()
		mutexI.Down()
		buf[i] = &ticket
		i = (i + 1) % N
		mutexI.Up()
		ticketSem.Up()
		time.Sleep(1000 * time.Millisecond)
	}
}

func consumer(pos int) {
	defer wg1.Done()
	for nticket := 0; nticket < MAXTICKETS; nticket++ {
		ticketSem.Down()
		mutexJ.Down()
		ticket := buf[j]
		buf[j] = nil
		j = (j + 1) % N
		mutexJ.Up()
		holeSem.Up()
		time.Sleep(500 * time.Millisecond)
		fmt.Println("New Ticket!! ", *ticket)
	}
}

func TestConProd(t *testing.T) {
	for i := 0; i < 10; i++ {
		wg1.Add(2)
		go consumer(i)
		go producter(i)
	}
	wg1.Wait()
}
