package rendez

import (
	"sync"
)

type appointment struct {
	i  interface{}
	ws sync.WaitGroup
}

type appointmentsProtected struct {
	mut sync.Mutex
	aps map[int]*appointment
}

var apsProt = appointmentsProtected{sync.Mutex{}, make(map[int]*appointment)}

func Rendezvous(tag int, val interface{}) (coupleVal interface{}) {

	apsProt.mut.Lock()
	if a, ok := apsProt.aps[tag]; !ok {
		a = &appointment{val, sync.WaitGroup{}} // New appointment!!...
		apsProt.aps[tag] = a                    // Put the appointment inside the map.
		a.ws.Add(1)                             // Before unlock the couple can't decrement...
		apsProt.mut.Unlock()
		a.ws.Wait()     // Wait to my couple
		coupleVal = a.i // get the value of my couple.
	} else {
		coupleVal, a.i = a.i, val // Get the appointment value , set appointment to couple.
		a.ws.Done()               // Wake up to my couple.
		delete(apsProt.aps, tag)  // Delete the appointment...
		apsProt.mut.Unlock()
	}
	return
}
