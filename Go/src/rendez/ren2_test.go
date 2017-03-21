package rendez

import (
	"testing"
	"time"
)

var routineVal = "ready to follow..."

func sleepR(t *testing.T) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	Rendezvous(0, routineVal)
}

func TestRenS(t *testing.T) {
	go sleepR(t)
	coupleVal := Rendezvous(0, "wait routine...")

	if coupleVal != routineVal {
		t.Error()
	}
}
