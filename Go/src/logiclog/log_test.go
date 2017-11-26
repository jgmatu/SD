package logiclog

import (
	"sync"
	"testing"
)

func isBadOrder(less bool, err error) bool {
	return err == nil && less
}

func checkOrder(t *testing.T) {
      lines := readlines("order.txt")
	prev := Line{}

	for i, line := range lines {
		if i > 0 {
			less, err := line.mark.IsLess(prev.mark)
			if isBadOrder(less, err) {
				t.Error()
			}
		}
		prev = line
	}
}

func orderLogs() {
	files := make([]string, 3)
	output := true // Test not output
	files[0] = "A.txt"
	files[1] = "B.txt"
	files[2] = "C.txt"
	Order("order.txt", files, !output)
}

func TestMsg(t *testing.T) {
	wg := sync.WaitGroup{}
	ch1 := make(chan Msg)
	ch2 := make(chan Msg)
	ch3 := make(chan Msg)

	wg.Add(3)
	go func() {
		defer wg.Done()
		logf := NewLog("a", "A.txt")

		logf.Event("A: Hi!")
		logf.Event("A: I am bored...")
		logf.Send("A: Shall we go to the cinema.", ch1)
		logf.Receive(ch3)
		logf.Close()
	}()
	go func() {
		defer wg.Done()
		logf := NewLog("b", "B.txt")

		logf.Event("B: I am playing computer...")
		logf.Receive(ch1)
		logf.Send("B: We go to the cinema!", ch2)
		logf.Close()
	}()
	go func() {
		defer wg.Done()
		logf := NewLog("c", "C.txt")

		logf.Event("C: I am having dinner...")
		logf.Receive(ch2)
		logf.Send("C: I can go!!", ch3)
		logf.Close()
	}()
	wg.Wait()
	orderLogs()
	checkOrder(t)
}
