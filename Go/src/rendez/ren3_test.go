// In this test we check the code before the Rendezvous is executed...
// with two counters and as well as check Rendezvous with the same tag.

// we don't who take the appointment first then we don't know what values
// of couple will be...

package rendez

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	TWEETS = 10
	SLEEPS = 10
)

var wg3 = sync.WaitGroup{}
var nTweet = 0
var nSleep = 0

func countersTest(coupleVal interface{}, t *testing.T) {
	if coupleVal == "Tweet" && nTweet != TWEETS {
		t.Error()
	}
	if coupleVal == "Sleep" && nSleep != SLEEPS {
		t.Error()
	}
}

func Sleep(tag int, val interface{}, t *testing.T) {
	defer wg3.Done()
	for i := 0; i < SLEEPS; i++ {
		time.Sleep(time.Second)
		nSleep++
	}
	coupleVal := Rendezvous(tag, val)
	countersTest(coupleVal, t)
}

func PrintTweet(tag int, val interface{}, t *testing.T) {
	defer wg3.Done()
	for i := 0; i < TWEETS; i++ {
		fmt.Println("Tweet!")
		nTweet++
	}
	coupleVal := Rendezvous(tag, val)
	countersTest(coupleVal, t)
}

func IWaitAndTest(tag int, val interface{}, t *testing.T) {
	defer wg3.Done()
	fmt.Println("waiting...")
	coupleVal := Rendezvous(tag, val)
	countersTest(coupleVal, t)
}

func TestRen3(t *testing.T) {
	wg3.Add(4)

	go IWaitAndTest(1, "Wait", t)
	go Sleep(1, "Sleep", t)
	go IWaitAndTest(1, "Wait", t)
	go PrintTweet(1, "Tweet", t)

	wg3.Wait()
}
