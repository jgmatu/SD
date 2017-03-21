// Comprobamos que Rendezvous no tiene condiciones de carrera
// poniendo contadores e incrementandolos primero en una goroutine
// y despues en la otra solo se asegura que el codigo de antes del rendez
// se ha ejecutado por eso last llama a rendez antes de incrementar el contador
// no incrementara hasta que su pareja no lo haya incrementado primero.

package rendez

import (
	"fmt"
	"sync"
	"testing"
)

const (
	MAXROUTINES = 10
	MAXCOUNT    = 10000
)

var counters = make([]int, MAXROUTINES)
var wg4 = sync.WaitGroup{}

func Before(tag int, val interface{}, t *testing.T) {
	defer wg4.Done()

	for i := 0; i < MAXCOUNT; i++ {
		counters[tag]++
	}
	coupleVal := Rendezvous(tag, val)

	if coupleVal != "Later" {
		t.Error()
	}
}

func Later(tag int, val interface{}, t *testing.T) {
	defer wg4.Done()

	coupleVal := Rendezvous(tag, val)
	for i := 0; i < MAXCOUNT; i++ {
		counters[tag]++
	}

	if coupleVal != "Before" {
		t.Error()
	}
}

func TestRunningConditions(t *testing.T) {
	wg4.Add(2 * MAXROUTINES)
	for i := 0; i < MAXROUTINES; i++ {
		go Before(i, "Before", t)
		go Later(i, "Later", t)
	}
	wg4.Wait()

	for i := 0; i < MAXROUTINES; i++ {
		if counters[i] != 2*MAXCOUNT {
			fmt.Println("%d", counters[i])
			t.Error()
		}
	}
}
